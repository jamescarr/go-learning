package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	data, err := ioutil.ReadFile("file-reading.go")
	failOnError(err)
	fmt.Println(string(data))

	f, err := os.Open("file-reading.go")
	defer f.Close()
	failOnError(err)

	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	failOnError(err)
	fmt.Printf("%d bytes: %s\n", n1, string(b1))

	o2, err := f.Seek(6, 0)
	failOnError(err)
	b2 := make([]byte, 2)
	n2, err := f.Read(b2)
	failOnError(err)
	fmt.Printf("%d bytes @ %d: %s\n", n2, o2, string(b2))

	o3, err := f.Seek(6, 0)
	failOnError(err)
	b3 := make([]byte, 2)
	n3, err := io.ReadAtLeast(f, b3, 2)
	failOnError(err)
	fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))

	_, err = f.Seek(0, 0)
	failOnError(err)

	r5 := bufio.NewReader(f)
	b4, err := r5.Peek(5)
	failOnError(err)
	fmt.Printf("5 bytes: %s\n", string(b4))

	fmt.Println("read the whole file")

}
