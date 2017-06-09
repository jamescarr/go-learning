package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
	"testing"
)

const MC_HOST = "localhost:11211"

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
	mc := memcache.New(MC_HOST)

	chunker := NewChunker(Config{host: MC_HOST})

	chunker.store("foo", reader)

	it, _ := mc.Get("foo")

	assert.Equal(t, 1, read(it.Value).Chunks)
}

func TestSpecifyingChunkSize(t *testing.T) {
	reader := bytes.NewReader([]byte("abcdef"))
	mc := memcache.New(MC_HOST)

	chunker := NewChunker(Config{host: MC_HOST, chunkSize: 2})

	chunker.store("foo", reader)

	it, _ := mc.Get("foo")

	assert.Equal(t, 3, read(it.Value).Chunks)
}

func TestWhatGoesInIsWhatComesOut(t *testing.T) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	reader := bytes.NewReader([]byte("abcdef"))

	chunker := NewChunker(Config{host: MC_HOST, chunkSize: 2})
	chunker.store("foo", reader)
	chunker.get("foo", writer)

	writer.Flush()
	assert.Equal(t, "abcdef", b.String())
}

func TestChecksumIsComputed(t *testing.T) {
	reader := bytes.NewReader([]byte("abcdef"))
	hash := sha256.New()
	hash.Write([]byte("abcdef"))
	expectedChecksum := fmt.Sprintf("%x", hash.Sum(nil))
	mc := memcache.New(MC_HOST)

	chunker := NewChunker(Config{host: MC_HOST, chunkSize: 2})

	chunker.store("foo", reader)

	it, _ := mc.Get("foo")

	assert.Equal(t, expectedChecksum, read(it.Value).CheckSum)

}

func read(str []byte) LargeValue {
	val := LargeValue{}
	json.Unmarshal(str, &val)
	return val
}
