package metrics

type snapshotResp struct {
	count int64
	sum   float64
}

// ChannelAggregator collects metrics via channels (fan-in pattern).
type ChannelAggregator struct {
	recordCh   chan float64
	snapshotCh chan chan snapshotResp
	done       chan struct{}
}

func NewChannelAggregator() *ChannelAggregator {
	a := &ChannelAggregator{
		recordCh:   make(chan float64, 256),
		snapshotCh: make(chan chan snapshotResp),
		done:       make(chan struct{}),
	}
	go a.run()
	return a
}

func (a *ChannelAggregator) run() {
	var count int64
	var sum float64
	for {
		select {
		case v := <-a.recordCh:
			count++
			sum += v
		case reply := <-a.snapshotCh:
			reply <- snapshotResp{count, sum}
		case <-a.done:
			return
		}
	}
}

func (a *ChannelAggregator) Record(value float64) {
	a.recordCh <- value
}

func (a *ChannelAggregator) Snapshot() (count int64, sum float64) {
	reply := make(chan snapshotResp, 1)
	a.snapshotCh <- reply
	r := <-reply
	return r.count, r.sum
}

func (a *ChannelAggregator) Close() {
	close(a.done)
}
