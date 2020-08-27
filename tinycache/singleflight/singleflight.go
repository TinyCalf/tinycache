package singleflight

import (
	"sync"
)

type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}

//Group 是一组访问
type Group struct {
	mu sync.Mutex
	m map[string]*call
}

// Do 代理fn的执行
// 传入的func基本是一次IO调用，时间一般比较长，可以看作100ms;
// Do保证的就是这100ms内相同的并发请求能只消耗一次IO
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock() // 该锁保证的是相同call不被并发创建，所以调用fn之前要及时解锁，以免变成单线程请求
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() //找到了相同的请求正在进行中，等待直到请求完成
		return c.val, c.err
	}
	//没找到call 说明是当前并发下的第一个请求
	//那创建一个新请求，并增加一个等待锁的队列,使得并发进来的其他请求等待
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock() //删除的时候还是需要Lock防止后来的请求发现这个call，但实际已经完成了
	delete(g.m, key) //delete删除的是map中的call，上面的c已经持有这个指针了，所以不会对上面有影响
	g.mu.Unlock()

	return c.val, c.err
}