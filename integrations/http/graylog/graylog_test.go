package graylog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTestCase(t *testing.T) {
	t.Log("setup test case")
	client := NewGraylog("http://localhost:9000/api", "admin", "admin")
	streams, _, _ := client.Streams.List()
	for _, stream := range streams.Streams {
		client.Streams.Delete(stream.ID)
	}
}

func TestGetStreams(t *testing.T) {
	setupTestCase(t)
	client := NewGraylog("http://localhost:9000/api", "admin", "admin")

	streams, _, _ := client.Streams.List()

	assert.Equal(t, 0, streams.Total)
}

func TestAddNewStream(t *testing.T) {
	setupTestCase(t)
	client := NewGraylog("http://localhost:9000/api", "admin", "admin")
	stream := Stream{
		Title:       "test",
		Description: "This is a test",
	}

	client.Streams.Create(stream)

}
