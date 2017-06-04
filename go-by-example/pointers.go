package main

import (
  "fmt"
)

func zeroval(n int) {
  n = 0
}

func zeroptr(n *int) {
  *n = 0
}

func main() {
  i := 1

  zeroval(i)
  fmt.Println("zeroval: ", i)

  zeroptr(&i)
  fmt.Println("zeroptr: ", i)

  fmt.Println("pointer: ", &i)
}
