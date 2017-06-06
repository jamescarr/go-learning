package main

import (
  "fmt"
  "math/rand"
  "sync"
  "sync/atomic"
  "time"
)

func writer(mutex *sync.Mutex, state map[int]int, writeOps *uint64) {
  for {
    key := rand.Intn(5)
    val := rand.Intn(100)
    withLock(mutex, func(){
      state[key] = val
    })
    atomic.AddUint64(writeOps, 1)
    time.Sleep(time.Millisecond)
  }
}

func reader(mutex *sync.Mutex, state map[int]int, ops *uint64) {
  total := 0
  for {
    key := rand.Intn(5)
    withLock(mutex, func(){
      total += state[key]
    })
    atomic.AddUint64(ops, 1)
    time.Sleep(time.Millisecond)
  }
}

func withLock(mutex *sync.Mutex, fn func()) {
  mutex.Lock()
  fn()
  mutex.Unlock()
}

func main() {
  var state = make(map[int]int)
  var readOps uint64 = 0
  var writeOps uint64 = 0
  var mutex = &sync.Mutex{}
  for r := 0; r < 100; r++ {
    go reader(mutex, state, &readOps)
  }
  for w := 0; w < 10; w++ {
    go writer(mutex, state, &writeOps)
  }

  time.Sleep(time.Second)

  readOpsFinal := atomic.LoadUint64(&readOps)
  fmt.Println("readOps:", readOpsFinal)
  writeOpsFinal := atomic.LoadUint64(&writeOps)
  fmt.Println("writeOps:", writeOpsFinal)

  withLock(mutex, func() {
    fmt.Println("state:", state)
  })
}
