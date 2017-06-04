package main

import (
  "time"
  "fmt"
)

func doSomething(c chan<- string, n int) {
  time.Sleep(2 * time.Second)
  c <- fmt.Sprintf("result %d", n)
}

/*
This is interesting... apparently literals can be boxed into different 
types in go. In this case, passing a literal int into the method call
here turns it into a time.Duration struct.
*/
func receiveWithTimeout(c <-chan string, timeoutSeconds time.Duration) {
  select {
  case res := <-c:
    fmt.Println("The Result is", res)
  case <-time.After(time.Second * timeoutSeconds):
    fmt.Println("Timeout")
  }
}

func main() {
  // assume we're making an external call that returns on a channel
  // after two seconds.
  c1 := make(chan string, 1)

  go doSomething(c1, 1)
  receiveWithTimeout(c1, 1)

  c2 := make(chan string, 1)

  go doSomething(c2, 2)
  receiveWithTimeout(c2, 3)
}
