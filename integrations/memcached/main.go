package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	mc := memcache.New("127.0.0.1:11211")
	mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

	it, _ := mc.Get("foo")

	fmt.Println(string(it.Value))
}
