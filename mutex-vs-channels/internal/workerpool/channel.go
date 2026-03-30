package workerpool

import "sync"

// ChannelPool is a worker pool using a buffered channel as the task queue.
type ChannelPool struct {
	taskCh chan func()
	wg     sync.WaitGroup
}

func NewChannelPool(workers, queueSize int) *ChannelPool {
	p := &ChannelPool{
		taskCh: make(chan func(), queueSize),
	}
	p.wg.Add(workers)
	for range workers {
		go p.worker()
	}
	return p
}

func (p *ChannelPool) worker() {
	defer p.wg.Done()
	for task := range p.taskCh {
		task()
	}
}

func (p *ChannelPool) Submit(task func()) {
	p.taskCh <- task
}

func (p *ChannelPool) Close() {
	close(p.taskCh)
	p.wg.Wait()
}
