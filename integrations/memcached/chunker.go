package main

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"io"
)

type LargeValue struct {
	Chunks int `json:"chunks"`
}

type Config struct {
	host      string
	chunkSize int
}

type Chunker struct {
	client *memcache.Client
	config Config
}

func (c *Chunker) get(key string, out io.Writer) {
	item, _ := c.client.Get(key)
	val := LargeValue{}
	json.Unmarshal(item.Value, &val)

	for i := 0; i < val.Chunks; i++ {
		chunkKey := fmt.Sprintf("%s-%d", key, i)
		chunk, _ := c.client.Get(chunkKey)
		out.Write(chunk.Value)
	}
}

func (c *Chunker) store(key string, input io.Reader) {
	val := LargeValue{Chunks: 0}
	buf := make([]byte, c.config.chunkSize)
	for {
		n, err := input.Read(buf)
		if n == 0 || err != nil {
			break
		}
		chunkKey := fmt.Sprintf("%s-%d", key, val.Chunks)
		c.client.Set(&memcache.Item{Key: chunkKey, Value: buf})

		// write each chunk out here
		val.Chunks += 1
	}

	serialized, _ := json.Marshal(val)
	fmt.Println(string(serialized))
	c.client.Set(&memcache.Item{Key: key, Value: serialized})
}

func NewChunker(config Config) *Chunker {
	if config.chunkSize == 0 {
		config.chunkSize = 1024
	}

	return &Chunker{
		client: memcache.New(config.host),
		config: config,
	}
}
