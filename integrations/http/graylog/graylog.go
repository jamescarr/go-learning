package graylog

import (
	"fmt"
	"github.com/dghubble/sling"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"log"
	"net/http"
)

type Stream struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
type StreamsList struct {
	Total   int      `json:"total"`
	Streams []Stream `json:"streams"`
}

type StreamCreated struct {
	StreamID string `json:"stream_id"`
}

type StreamsService struct {
	sling *sling.Sling
}

type GraylogError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *GraylogError) Error() string {
	return e.Message
}

func failOnError(e error) {
	if e != nil {
		panic(e)
	}
}
func newStreamsService(sling *sling.Sling) *StreamsService {
	return &StreamsService{
		sling: sling.Path("api/streams"),
	}
}

func (s *StreamsService) List() (StreamsList, *http.Response, error) {
	list := StreamsList{}
	graylogError := new(GraylogError)
	resp, err := s.sling.New().Get("").Receive(&list, &graylogError)
	failOnError(err)
	return list, resp, err

}

func (s *StreamsService) Create(stream Stream) {
	success := StreamCreated{}
	err := new(GraylogError)
	_, _ = s.sling.New().Post("").BodyJSON(stream).Receive(success, err)

}

func (s *StreamsService) Delete(streamId string) {
	success := GraylogError{}
	err := new(GraylogError)
	req, _ := s.sling.New().Delete(fmt.Sprintf("/api/streams/%s", streamId)).Request()
	s.sling.Do(req, &success, &err)
	failOnError(err)
	log.Printf("Deleting stream with id", streamId)

}

type Graylog struct {
	sling   *sling.Sling
	Streams *StreamsService
}

func NewGraylog(address string, username string, password string) *Graylog {
	client := cleanhttp.DefaultPooledClient()
	base := sling.New().Client(client).Base(address).SetBasicAuth(username, password)

	return &Graylog{
		sling:   base,
		Streams: newStreamsService(base.New()),
	}
}
