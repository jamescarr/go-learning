package main

import (
  "fmt"
)

type Person struct {
  Name string
  Age int
}

type Rectangle struct {
  Width, Height int
}

func (r Rectangle) area() int {
  return r.Width * r.Height
}

func (r Rectangle) perim() int {
  return 2*r.Width + 2*r.Height
}

type GeometricShape interface {
  area() int
  perim() int
}

func PrintShape(g GeometricShape) {
  fmt.Println("area: ", g.area())
  fmt.Println("perim: ", g.perim())
}

func main() {
  fmt.Println(Person{"Bob", 20})

  fmt.Println(Person{Name:"Alice", Age:30})

  fmt.Println(Person{Name: "Fred"})

  fmt.Println(&Person{Name: "Ann", Age: 40})

  s := Person{Name: "Sean", Age: 50}
  sp := &s
  fmt.Println(sp.Age)

  sp.Age = 51
  fmt.Println(s)

  fmt.Println("\nMethods....\n\n")
  r := Rectangle{Width: 10, Height: 5}

  rp := &r
  rv := r
  r.Width = 12

  PrintShape(r)
  PrintShape(rp)
  PrintShape(rv)
}
