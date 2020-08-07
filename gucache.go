package gucache

import (
	"fmt"
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
}

var (
	m      sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxBytes int64, getter Getter, storageType ...int64) *Group {
	if getter == nil {
		panic("getter func can't be nil")
	}
	m.Lock()
	defer m.Unlock()
	var sType int64
	if len(storageType) == 1 {
		sType = storageType[0]
	} else {
		sType = 0
	}
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: maxBytes, storageType: sType},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	m.RLock()
	defer m.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("empty key")
	}
	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	}
	return g.loadByGetter(key)
}

func (g *Group) Set(key string, value string) {
	v := ByteView{b: []byte(value)}
	g.mainCache.set(key, v)
}

func (g *Group) loadByGetter(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.mainCache.set(key, value)
	return value, nil
}
