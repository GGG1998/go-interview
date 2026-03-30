package counter

// ChannelCounter is a shared counter managed via channels.
type ChannelCounter struct {
	incCh chan struct{}
	getCh chan chan int64
	done  chan struct{}
}

func NewChannelCounter() *ChannelCounter {
	c := &ChannelCounter{
		incCh: make(chan struct{}, 256),
		getCh: make(chan chan int64),
		done:  make(chan struct{}),
	}
	go c.run()
	return c
}

func (c *ChannelCounter) run() {
	var value int64
	for {
		select {
		case <-c.incCh:
			value++
		case reply := <-c.getCh:
			reply <- value
		case <-c.done:
			return
		}
	}
}

func (c *ChannelCounter) Increment() {
	c.incCh <- struct{}{}
}

func (c *ChannelCounter) Get() int64 {
	reply := make(chan int64, 1)
	c.getCh <- reply
	return <-reply
}

func (c *ChannelCounter) Close() {
	close(c.done)
}
