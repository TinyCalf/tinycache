package lru

import "container/list"

//Value 使用Len来获取其大小
type Value interface {
	Len() int
}

//Cache 是一个LRU缓存 不要直接构建它
type Cache struct {
	maxBytes int64
	nBytes   int64
	cache    map[string]*list.Element
	ll       *list.List
}

//New 返回一个新建的Cache实例
func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes,
		0,
		make(map[string]*list.Element),
		list.New(),
	}
}

//Len 获取已缓存元素长度
func (c *Cache) Len() int {
	return c.ll.Len()
}

//Set 设置
func (c *Cache) Set(key string, value Value) {
	if ele, exists := c.cache[key]; exists {
		c.ll.MoveToFront(ele)
		obj := ele.Value.(*cacheObject)
		c.nBytes += int64(value.Len()) - int64(obj.value.Len())
		obj.value = value
	} else {
		obj := &cacheObject{key, value}
		ele := c.ll.PushFront(obj)
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.removeOldest()
	}
}

//Get 获取
func (c *Cache) Get(key string) (value Value, exists bool) {
	if ele, exists := c.cache[key]; exists {
		c.ll.MoveToFront(ele)
		obj := ele.Value.(*cacheObject)
		return obj.value, true
	}
	return
}

func (c *Cache) removeOldest() {
	if ele := c.ll.Back(); ele != nil {
		c.ll.Remove(ele)
		obj := ele.Value.(*cacheObject)
		delete(c.cache, obj.key)
		c.nBytes -= int64(len(obj.key)) + int64(obj.value.Len())
	}
}

type cacheObject struct {
	key   string
	value Value
}
