package main

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Println("Caught error", err)
	}
}
func runResponder(address string) {
	responder, err := zmq.NewSocket(zmq.REP)
	handleError(err)
	defer responder.Close()
	responder.Bind(address)

	for {
		log.Println("Wait for next request from client ...")
		request, err := responder.Recv(0)
		handleError(err)
		log.Printf("Received request: [%s]\n", request)

		//  Do some 'work'
		time.Sleep(time.Second)

		//  Send reply back to client
		responder.Send("World", 0)
	}
}

func runRequester(address string, done chan bool) {
	requester, err := zmq.NewSocket(zmq.REQ)
	handleError(err)
	defer requester.Close()
	requester.Connect(address)

	for request := 0; request < 10; request++ {
		log.Println("Sending Requests")
		requester.Send("Hello", 0)
		reply, err := requester.Recv(0)
		handleError(err)
		log.Printf("Received reply %d [%s]\n", request, reply)
	}
	done <- true
}
func main() {
	done := make(chan bool, 1)

	go runRequester("tcp://localhost:5560", done)
	go runResponder("tcp://*:5560")

	<-done
}
