package tinycache

import (
	"sync"

	"./lru"
)

type coreCache struct {
	mu         	sync.Mutex
	lru        	*lru.Cache
	cap 		int64
}

func (c *coreCache) set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cap)
	}
	c.lru.Set(key, tinyString(value))
}

func (c *coreCache) get(key string) (value string, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return string(v.(tinyString)), ok
	}

	return
}