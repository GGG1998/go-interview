package workerpool

import "sync"

// MutexPool is a worker pool with a mutex-protected task queue.
type MutexPool struct {
	mu    sync.Mutex
	cond  *sync.Cond
	queue []func()
	stop  bool
	wg    sync.WaitGroup
}

func NewMutexPool(workers int) *MutexPool {
	p := &MutexPool{}
	p.cond = sync.NewCond(&p.mu)
	p.wg.Add(workers)
	for range workers {
		go p.worker()
	}
	return p
}

func (p *MutexPool) worker() {
	defer p.wg.Done()
	for {
		p.mu.Lock()
		for len(p.queue) == 0 && !p.stop {
			p.cond.Wait()
		}
		if p.stop && len(p.queue) == 0 {
			p.mu.Unlock()
			return
		}
		task := p.queue[0]
		p.queue = p.queue[1:]
		p.mu.Unlock()
		task()
	}
}

func (p *MutexPool) Submit(task func()) {
	p.mu.Lock()
	p.queue = append(p.queue, task)
	p.mu.Unlock()
	p.cond.Signal()
}

func (p *MutexPool) Close() {
	p.mu.Lock()
	p.stop = true
	p.mu.Unlock()
	p.cond.Broadcast()
	p.wg.Wait()
}
