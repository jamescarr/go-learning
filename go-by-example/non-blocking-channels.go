package main

import "fmt"

func main() {
  messages := make(chan string)
  signals := make(chan bool)

  select {
  case msg := <-messages:
    fmt.Println("Recevied message ", msg)
  default:
    fmt.Println("No message received")
  }
  
  msg := "hi"

  select {
  case messages <- msg:
    fmt.Println("Sent message", msg)
  default:
    fmt.Println("No message sent")
  }

  select {
  case msg := <-messages:
      fmt.Println("received message", msg)
  case sig := <-signals:
      fmt.Println("received signal", sig)
  default:
      fmt.Println("no activity")
  }

}
