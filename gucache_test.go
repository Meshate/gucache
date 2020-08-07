package gucache

import (
	"sync"
	"testing"
)

var g = NewGroup("main", 2<<20, GetterFunc(func(key string) ([]byte, error) {
	return nil, nil
}))

func BenchmarkGroup_Set(b *testing.B) {
	b.ResetTimer()
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for i := 0; i < b.N; i++ {
			g.Set("test", []byte("long bench in"))
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		for i := 0; i < b.N; i++ {
			g.Set("test2", []byte("long bench in"))
		}
		wg.Done()
	}()
	wg.Wait()
}

func BenchmarkGroup_Get(b *testing.B) {
	b.ResetTimer()
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for i := 0; i < b.N; i++ {
			g.Get("test")
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		for i := 0; i < b.N; i++ {
			g.Get("test2")
		}
		wg.Done()
	}()
	wg.Wait()
}
