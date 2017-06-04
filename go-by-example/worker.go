package main

import "fmt"
import "time"

func worker(id int, jobs <-chan int, results chan<- int) {
  for j := range jobs {
    fmt.Println("worker", id, "started  job", j)
    time.Sleep(time.Second)
    fmt.Println("worker", id, "finished job", j)
    results <- j * 2
  }
}

func main(){
  jobs := make(chan int, 100)
  results := make(chan int, 100)

  for w :=0; w < 3; w++ {
    go worker(w, jobs, results)
  }

  for i := 0; i <= 5; i++ {
    jobs <- i
  }
  close(jobs)

  for a := 0; a <= 5; a++ {
    r := <-results
    fmt.Println("Result is", r)
  }
}
