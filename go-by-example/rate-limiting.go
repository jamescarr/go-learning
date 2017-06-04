package main

import (
  "fmt"
  "time"
)

func main() {
  requests := generateRequests(5)

  limiter := time.Tick(time.Millisecond * 200)

  readChannelWithLimiter(requests, limiter)

  burstyLimiter := make(chan time.Time, 3)

  for i := 0; i < 3; i++ {
    burstyLimiter <- time.Now()
  }

  go func() {
    for t := range time.Tick(time.Millisecond * 200) {
        burstyLimiter <- t
    }
  }()

  burstyRequests := generateRequests(5)

  readChannelWithLimiter(burstyRequests, burstyLimiter)
}

func readChannelWithLimiter(requests chan int, limiter <-chan time.Time) {
  for req := range requests {
    <-limiter
    fmt.Println("request", req, time.Now())
  }

}
func generateRequests(size int) chan int {
  requests := make(chan int, size)
  for i := 1; i <= size; i++ {
      requests <- i
  }
  close(requests)
  return requests
}
