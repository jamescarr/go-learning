package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"io"
)

type LargeValue struct {
	Chunks   int    `json:"chunks"`
	CheckSum string `json:"checksum"`
}

type Config struct {
	host      string
	chunkSize int
}

type Chunker struct {
	client *memcache.Client
	config Config
}

func failOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func (c *Chunker) get(key string, out io.Writer) {
	item, _ := c.client.Get(key)
	val := LargeValue{}

	json.Unmarshal(item.Value, &val)

	for i := 0; i < val.Chunks; i++ {
		chunkKey := fmt.Sprintf("%s-%d", key, i)
		chunk, err := c.client.Get(chunkKey)
		failOnError(err)
		out.Write(chunk.Value)
	}
}

func (c *Chunker) store(key string, input io.Reader) {
	val := LargeValue{Chunks: 0}
	hash := sha256.New()

	buf := make([]byte, c.config.chunkSize)
	for {
		n, err := input.Read(buf)
		if n == 0 || err != nil {
			break
		}
		hash.Write(buf[:n])
		chunkKey := fmt.Sprintf("%s-%d", key, val.Chunks)
		c.client.Set(&memcache.Item{Key: chunkKey, Value: buf[:n]})

		// write each chunk out here
		val.Chunks += 1
	}

	val.CheckSum = fmt.Sprintf("%x", hash.Sum(nil))
	serialized, err := json.Marshal(val)
	failOnError(err)
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
