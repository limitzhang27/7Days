package lru

import (
	"reflect"
	"testing"
)

// 新建一个 Value 对象做测试
type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"

	l := len(k1 + k2 + v1 + v2)
	lru := New(int64(l), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"

	keys := make([]string, 0)
	callback := func(k string, v Value) {
		keys = append(keys, k)
	}

	l := len(k1 + k2 + v1 + v2)
	lru := New(int64(l), callback)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldest key1 failed")
	}
	expect := []string{"key1"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("call Onevicted failed, expect keys equals to %s", expect)
	}
}
