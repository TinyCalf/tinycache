package tinycache

import (
	"sync"

	"./lru"
)

type cache struct {
	mu         sync.Mutex //lru的读操作也有修改链表，所以不选用读写锁
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes)
	}
	c.lru.Set(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, exists := c.lru.Get(key); exists {
		return v.(ByteView), ok
	}

	return
}
