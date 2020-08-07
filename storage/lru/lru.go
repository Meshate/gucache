package lru

import (
	"container/list"
	. "github.com/Meshate/gucache/storage"
	"unsafe"
)

type lru struct {
	maxBytes         int64
	guBytes          int64
	data             *list.List
	cache            map[string]*list.Element
	OnElementRemoved func(key string, value Value)
}

type element struct {
	key   string
	value Value
}

func New(maxBytes int64, onElementRemoved func(key string, value Value)) *lru {
	return &lru{
		maxBytes:         maxBytes,
		guBytes:          0,
		data:             list.New(),
		cache:            make(map[string]*list.Element),
		OnElementRemoved: onElementRemoved,
	}
}

func (l *lru) Set(key string, value Value) {
	if item, ok := l.cache[key]; ok {
		l.data.MoveToFront(item)
		ele := item.Value.(*element)
		l.guBytes += int64(unsafe.Sizeof(value) - unsafe.Sizeof(ele.value))
		ele.value = value
	} else {
		ele := &element{
			key:   key,
			value: value,
		}
		v := l.data.PushFront(ele)
		l.cache[key] = v
		l.guBytes += int64(unsafe.Sizeof(*ele))
	}
	for l.guBytes != 0 && l.maxBytes < l.guBytes {
		l.removeOne()
	}
}

func (l *lru) Get(key string) (value Value, ok bool) {
	if item, ok := l.cache[key]; ok {
		l.data.MoveToFront(item)
		ele := item.Value.(*element)
		return ele.value, true
	}
	return nil, false
}

func (l *lru) removeOne() {
	last := l.data.Back()
	if last != nil {
		l.data.Remove(last)
		ele := last.Value.(*element)
		delete(l.cache, ele.key)
		l.guBytes -= int64(unsafe.Sizeof(*ele))
		if l.OnElementRemoved != nil {
			l.OnElementRemoved(ele.key, ele.value)
		}
	}
}
