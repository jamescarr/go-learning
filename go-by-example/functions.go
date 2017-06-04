package main

import (
  "fmt"
  "time"
)

func plus(a int, b int) int {
  return a + b
}

func sum_and_product(a int, b int) (int, int) {
  return a + b, a * b
}

func sum(nums ...int) int {
  total := 0
  for _, num := range nums {
    total += num
  }
  return total
}

func intSeq() func() int {
  n := 0
  return func() int {
      n += 1
      return n
  }
}

func fact(n int) int {
  if n == 0 {
    return 1
  }
  return n * fact(n-1)
}
func main() {
  res := plus(1, 2)

  fmt.Println("1 + 2 = ", res)

  res, product := sum_and_product(2, 4)

  fmt.Println("2 + 4 = ", res, " | 2 * 4 = ", product)

  fmt.Println("Sum is ", sum(1, 2, 3, 4))

  numbers := []int{2, 4, 6, 8}
  fmt.Println("Sum is ", sum(numbers...))

  seq := intSeq()
  n := 0
  for n < 10 {
    n += seq()
    fmt.Println("seq is ", n)
    time.Sleep(300 * time.Millisecond)
  }

  fmt.Println(fact(8))

}
