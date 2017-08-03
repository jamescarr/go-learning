package main

import (
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/jamescarr/go-learning/integrations/communication/protocolbuffers/example"
)

func checkError(err error) {
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
}

func handleProtoClient(conn net.Conn, c chan *example.Test) {
	log.Println("Connection established")
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)
	checkError(err)
	log.Println("Decoding Protobuf message")
	//Create an struct pointer of type example.Test struct
	protodata := new(example.Test)
	//Convert all the data retrieved into the example.Test struct type
	err = proto.Unmarshal(data[0:n], protodata)
	checkError(err)
	//Push the protobuf message into a channel
	c <- protodata
}

func main() {
	c := make(chan *example.Test)

	go func() {
		for {
			message := <-c
			log.Println("Received message", message)

		}
	}()

	listener, err := net.Listen("tcp", "127.0.0.1:2110")
	checkError(err)
	log.Println("Server listening on 2110")
	for {
		if conn, err := listener.Accept(); err == nil {
			//If err is nil then that means that data is available for us so we take up this data and pass it to a new goroutine
			go handleProtoClient(conn, c)
		} else {
			continue
		}
	}

}
