package geecache

import (
	"fmt"
	"geecache/pb"
	"geecache/singleflight"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker

	// 使用 singlefight.Group 保证
	// 每一个key只会请求一次
	loader *singleflight.Group
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheByte int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheByte: cacheByte},
		loader:    &singleflight.Group{},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	return groups[name]
}

func (g *Group) RegisterPeer(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if value, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return value, nil
	}
	return g.load(key)
}

// 获取数据，先从远程结点中获取，如果为设置远程结点就从本地获取
func (g *Group) load(key string) (value ByteView, err error) {
	res, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			// 根据 key 获取远程结点，通过远程结点获取数据
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
		return g.getLocally(key) // 从本地获取
	})
	return res.(ByteView), err
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
