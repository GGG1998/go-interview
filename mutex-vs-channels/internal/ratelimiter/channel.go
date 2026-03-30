package ratelimiter

import "time"

type allowReq struct {
	clientID string
	reply    chan bool
}

// ChannelLimiter is a token-bucket rate limiter managed via channels.
type ChannelLimiter struct {
	allowCh chan allowReq
	resetCh chan struct{}
	done    chan struct{}
}

func NewChannelLimiter(limit int, window time.Duration) *ChannelLimiter {
	l := &ChannelLimiter{
		allowCh: make(chan allowReq, 256),
		resetCh: make(chan struct{}),
		done:    make(chan struct{}),
	}
	go l.run(limit)
	return l
}

func (l *ChannelLimiter) run(limit int) {
	tokens := make(map[string]int)
	for {
		select {
		case req := <-l.allowCh:
			if tokens[req.clientID] >= limit {
				req.reply <- false
			} else {
				tokens[req.clientID]++
				req.reply <- true
			}
		case <-l.resetCh:
			tokens = make(map[string]int)
		case <-l.done:
			return
		}
	}
}

func (l *ChannelLimiter) Allow(clientID string) bool {
	reply := make(chan bool, 1)
	l.allowCh <- allowReq{clientID: clientID, reply: reply}
	return <-reply
}

func (l *ChannelLimiter) Reset() {
	l.resetCh <- struct{}{}
}

func (l *ChannelLimiter) Close() {
	close(l.done)
}
