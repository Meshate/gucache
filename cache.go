package gucache

import (
	"github.com/Meshate/gucache/storage"
	"github.com/Meshate/gucache/storage/lru"
	"sync"
)

const(
	StorageLru int64 = iota
)

type dataCache interface {
	Set(key string, value storage.Value)
	Get(key string) (value storage.Value, ok bool)
}

type cache struct {
	m           sync.RWMutex
	cacheBytes  int64
	storageType int64
	cache       dataCache
}

func (c *cache) set(key string, value ByteView) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.cache == nil {
		switch c.storageType {
		case StorageLru:
			c.cache = lru.New(c.cacheBytes, nil)
		default:
			panic("need a storageType")
		}
	}
	c.cache.Set(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	if c.cache == nil {
		return
	}
	if v, ok := c.cache.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
