package cache

type getReq struct {
	key   string
	reply chan getResp
}

type getResp struct {
	value interface{}
	ok    bool
}

type setReq struct {
	key   string
	value interface{}
}

// ChannelCache is an in-memory key-value cache managed via channels.
type ChannelCache struct {
	getCh chan getReq
	setCh chan setReq
	done  chan struct{}
}

func NewChannelCache() *ChannelCache {
	c := &ChannelCache{
		getCh: make(chan getReq, 256),
		setCh: make(chan setReq, 256),
		done:  make(chan struct{}),
	}
	go c.run()
	return c
}

func (c *ChannelCache) run() {
	data := make(map[string]interface{})
	for {
		select {
		case req := <-c.getCh:
			v, ok := data[req.key]
			req.reply <- getResp{v, ok}
		case req := <-c.setCh:
			data[req.key] = req.value
		case <-c.done:
			return
		}
	}
}

func (c *ChannelCache) Get(key string) (interface{}, bool) {
	reply := make(chan getResp, 1)
	c.getCh <- getReq{key: key, reply: reply}
	r := <-reply
	return r.value, r.ok
}

func (c *ChannelCache) Set(key string, value interface{}) {
	c.setCh <- setReq{key: key, value: value}
}

func (c *ChannelCache) Close() {
	close(c.done)
}
