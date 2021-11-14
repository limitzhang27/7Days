package main

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	geecache.NewGroup("source", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Printf("[Slow DB] search key %s\n", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exit", key)
		}))

	s := geecache.NewHTTPPool("source")
	_ = http.ListenAndServe(":9527", s)
}
