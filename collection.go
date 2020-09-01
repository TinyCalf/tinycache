package tinycache

import (
	"fmt"
)

// Loader 用来获取源数据
type Loader func(key string) (string, error)

var collections = make(map[string]*Collection)

//Collection 是一个集合
//表示独立的缓存命名空间，有独立的获取源数据的方法
type Collection struct {
	name string
	cache coreCache
	loader Loader
}

// CreateCollection 创建一个新的集合
// cap 用于指定该集合所占用的内存，超出内存部分会执行lru
func CreateCollection(name string, cap int64, loader Loader) *Collection {
	c := &Collection{
		name: name,
		cache: coreCache{cap: cap},
		loader: loader,
	}
	collections[name] = c
	return c
}

// GetCollection 获取一个创建过的集合，不存在返回nil
func GetCollection(name string) (*Collection) {
	return collections[name]
}

// Get 从集合中获取键值
func (c *Collection) Get(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key is required")
	}

	peer, isSelf := globalPeerManager.selectPeer(c.name + key);
	if peer != nil && !isSelf  {
		value, err := peer.get(c.name, key)
		if err == nil {
			return value, nil
		}
	}
	
	if v, ok := c.cache.get(key); ok {
		return v, nil
	}

	v, e := c.loadFromLoader(key)
	return v, e
}

func (c *Collection) loadFromLoader(key string) (string, error) {
	value, err := c.loader(key)
	if err != nil{
		return "", err
	}
	c.cache.set(key, value)
	return value, nil
}