package main

import (
    "fmt"
    "gucache"
    "log"
)

var db = map[string]string{
    "Tom":  "sb",
    "Jack": "123",
    "Sam":  "regrdg",
}

func Load(key string) ([]byte, error){
    log.Printf("[SlowDB] search key %s", key)
    if v, ok := db[key]; ok {
        return []byte(v), nil
    }
    return nil, fmt.Errorf("%s not exist", key)
}

func main() {
	g := gucache.NewGroup("main", 2<<20, gucache.GetterFunc(Load))
	fmt.Println(g.Get("Tom"))
    fmt.Println(g.Get("Tom"))
	g.Set("test", "123")
	fmt.Println(g.Get("test"))
}
