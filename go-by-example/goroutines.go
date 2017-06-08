package main

import "fmt"

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ": ", i)
	}
}

func main() {
	f("direct")

	go f("goroutine")
	go f("goroutine 2")
	go f("goroutine 3")

	go func(msg string) {
		fmt.Println(msg)
	}("going")

	var input string
	fmt.Scanln(&input) // needs to be done so program doesn't exit
	fmt.Println("done")
}
