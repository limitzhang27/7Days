package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64                    // 允许最大内存
	nBytes   int64                    // 当前使用内存
	ll       *list.List               // 双向链表
	cache    map[string]*list.Element // 哈希，key是字符串， value是链表指针
	// optional and executed when an entry is purged
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可为nil
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

func New(maxBytes int64, ovEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: ovEvicted,
	}
}

// Get 查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除元素
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // 返回双向链表中最后一个元素，如果链表为空会返回nil
	if ele != nil {
		c.ll.Remove(ele)                                       // 链表删除元素
		kv := ele.Value.(*entry)                               // 元素转成 entry 类型
		delete(c.cache, kv.key)                                // 删除哈希表中的值
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 重置当前使用内存大小
		// 存在回调函数则执行
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 已存在
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len() - kv.value.Len())
	} else {
		ele = c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.nBytes = c.nBytes + int64(value.Len()) + int64(len(key))
	}
	// 等于0的不判断
	if c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
