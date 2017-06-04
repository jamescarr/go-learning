package main

import (
  "fmt"
  "time"
)

func sender(ch chan<- string) func() {
  return func() {
    ch <- "pong"
  }
}

func read(ch <-chan string) {
  msg := <-ch
  fmt.Println(msg)
}

func bufferedChannels() {
  messages := make(chan string, 2)

  messages <- "buffered"
  messages <- "channel"

  fmt.Println(<-messages)
  fmt.Println(<-messages)
}

func worker(done chan bool) {
  fmt.Print("working...")
  time.Sleep(2 * time.Second)
  fmt.Println("done")
  done <- true
}

func main() {
  messages := make(chan string)

  go func() { messages <- "ping" }()
  go sender(messages)()

  go read(messages)
  go read(messages)

  bufferedChannels()

  done := make(chan bool, 1)
  go worker(done)
  <-done
}
