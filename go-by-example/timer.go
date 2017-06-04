package main

import (
  "fmt"
  "time"
)

func stop(timer *time.Timer, done chan bool) {
  stop := timer.Stop()

  if stop {
    fmt.Println("Timer 2 stopped")
    done <- true
  }
}

func wait(timer *time.Timer, done chan bool) {
  fmt.Println("Waiting...")
  <-timer.C
  fmt.Println("Timer 2 expired")
  done <- true
}

func main() {
  timer1 := time.NewTimer(time.Second * 2)
  done := make(chan bool)

  <-timer1.C
  fmt.Println("Timer 1 expired")

  timer2 := time.NewTimer(time.Second)
  go wait(timer2, done)
  go stop(timer2, done)
  <-done
}
