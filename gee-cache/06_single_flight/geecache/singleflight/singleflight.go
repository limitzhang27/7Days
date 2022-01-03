package singleflight

import "sync"

// call 代表正在进行中或者已经结束的请求。
// 使用 sync.WaitGroup 避免重复进入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 是singleflight 的主数据结构，管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// Do 方法，
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock() // 先解锁，减少锁阻塞的时间，

	c.val, c.err = fn()
	c.wg.Done()

	// 请求完就删了，只是防止同一时刻重复请求
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
