package main

import (
	b64 "encoding/base64"
	"fmt"
)

func main() {
	x := "foo"
	encoded := b64.StdEncoding.EncodeToString([]byte(x))

	fmt.Println(x, "encoded is", encoded)

	decoded, _ := b64.StdEncoding.DecodeString(encoded)
	fmt.Println(encoded, "decoded is", string(decoded))

	uEnc := b64.URLEncoding.EncodeToString([]byte(x))
	fmt.Println(uEnc)
	uDec, _ := b64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))

}
