package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
)

func hashIt(algo hash.Hash, s string) string {
	algo.Write([]byte(s))
	return fmt.Sprintf("%x", algo.Sum(nil))
}

func main() {
	s := "Hello World"

	fmt.Println(s)
	fmt.Println(hashIt(md5.New(), s))
	fmt.Println(hashIt(sha1.New(), s))
	fmt.Println(hashIt(sha256.New(), s))
}
