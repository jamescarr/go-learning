package main

import (
	"flag"
	"github.com/golang/protobuf/proto"
	"github.com/jamescarr/go-learning/integrations/communication/protocolbuffers/example"
	"log"
	"net"
	"strconv"
)

func main() {
	dest := flag.String("d", "127.0.0.1:2110", "Enter the destnation socket address")

	test := &example.Test{
		Label: proto.String("hello"),
		Type:  proto.Int32(17),
		Reps:  []int64{1, 2, 3},
		Optionalgroup: &example.Test_OptionalGroup{
			RequiredField: proto.String("good bye"),
		},
	}
	data, err := proto.Marshal(test)
	checkError(err)
	sendDataToDest(data, dest)
}

func sendDataToDest(data []byte, dst *string) {
	conn, err := net.Dial("tcp", *dst)
	checkError(err)
	n, err := conn.Write(data)
	checkError(err)
	log.Println("Sent " + strconv.Itoa(n) + " bytes")
}

func checkError(err error) {
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
}
