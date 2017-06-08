package main

import (
  "fmt"
  "math/rand"
  "sync/atomic"
  "time"
)


type readOp struct {
  key int
  resp chan int
}

type writeOp struct
  key int
  val int
  resp chan bool
}

func main() {
  var readOps uint64 = 0
  var writeOps uint64 = 0

  reads := make(chan *readOps)
  writes := make(chan *writeOps)

  go func() {
    var state = make(map[int]int)
    for {
      select {
      case read := <-reads
        read.resp = state[read.key]
      }
    }
  }()
}
