package main

import (
  "fmt"
  "time"
)

func rangeOverChannel() {
  queue := make(chan string, 2)
  queue <- "one"
  queue <- "two"
  close(queue)

  for elem := range(queue) {
    fmt.Println(elem)
  }
}

func main() {
  jobs := make(chan int, 5)
  done := make(chan bool)

  go func() {
    for {
      j, more := <-jobs
      if more {
        fmt.Println("Received job", j)
        time.Sleep(time.Second)
      } else {
        fmt.Println("Received all jobs")
        done <- true
        return
      }
    }
  }()

  for j := 0; j <=3; j++ {
    jobs <- j
    fmt.Println("Job sent")
  }

  close(jobs)

  fmt.Println("All Jobs sent!")

  <-done

  rangeOverChannel()
}
