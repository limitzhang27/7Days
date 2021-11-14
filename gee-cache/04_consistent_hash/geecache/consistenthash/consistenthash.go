package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 定义 hash 函数类型，允许用户自定义hash算法
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // hash 函数
	replicas int            // 虚拟结点数量
	keys     []int          // hash 环
	hashMap  map[int]string // 虚拟结点对应真实节点
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  map[int]string{},
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	l := len(m.keys)
	// 查询hash的位置
	idx := sort.Search(l, func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%l]]
}
