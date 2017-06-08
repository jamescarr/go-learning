package main

import (
	"bufio"
	"bytes"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func failOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func TestReadingBytes(t *testing.T) {
	// because I'm a noob with go!
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	reader := bytes.NewReader([]byte("Here is a string...."))

	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n == 0 || err != nil {
			break
		}
		writer.Write(buf[:n])
		writer.Flush()
	}
	assert.Equal(t, "Here is a string....", b.String())

}

func TestStoreNumberOfChunks(t *testing.T) {
	reader := bytes.NewReader([]byte("Here is a string...."))
	mc := memcache.New("localhost:11211")

	chunker := NewChunker(Config{host: "localhost:11211"})

	chunker.store("foo", reader)

	it, _ := mc.Get("foo")

	assert.Equal(t, `{"chunks":1}`, string(it.Value))
}

func TestSpecifyingChunkSize(t *testing.T) {
	reader := bytes.NewReader([]byte("abcdef"))
	mc := memcache.New("localhost:11211")

	chunker := NewChunker(Config{host: "localhost:11211", chunkSize: 2})

	chunker.store("foo", reader)

	it, _ := mc.Get("foo")

	assert.Equal(t, `{"chunks":3}`, string(it.Value))
}

func TestWhatGoesInIsWhatComesOut(t *testing.T) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	reader := bytes.NewReader([]byte("abcdef"))

	chunker := NewChunker(Config{host: "localhost:11211", chunkSize: 2})
	chunker.store("foo", reader)
	chunker.get("foo", writer)

	writer.Flush()
	assert.Equal(t, "abcdef", b.String())
}
