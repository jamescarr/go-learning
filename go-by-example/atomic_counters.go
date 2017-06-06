package main

import (
  "fmt"
  "time"
  "sync/atomic"
)


func main() {
  var ops uint64 = 0
  var non_atomic_ops uint64 = 0

  for i := 0; i < 50; i++ {
    go func() {
      atomic.AddUint64(&ops, 1)
      non_atomic_ops += 1
      time.Sleep(time.Millisecond)
    }()
  }

  time.Sleep(time.Second)

  opsFinal := atomic.LoadUint64(&ops)
  fmt.Println("ops:", opsFinal)
  fmt.Println("non-atomic ops:", non_atomic_ops)
}

