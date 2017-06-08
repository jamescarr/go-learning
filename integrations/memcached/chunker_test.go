package main

import (
	"bufio"
	"bytes"
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
