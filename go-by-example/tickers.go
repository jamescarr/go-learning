package main

import (
  "fmt"
  "time"
  "reflect"
)

func main() {
  ticker := time.NewTicker(time.Millisecond * 500)
  go func() {
    for t := range ticker.C {
      fmt.Println("Tick at", t.UTC())
      fmt.Println("Tick is of type", reflect.TypeOf(t))
    }
  }()

  time.Sleep(time.Second * 5)
  ticker.Stop()
  fmt.Println("Ticker has been stopped")
}
