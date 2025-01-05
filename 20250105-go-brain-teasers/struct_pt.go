package main

import (
	"fmt"
)

// Point is a 2D point
type Point struct {
	X int
	Y int
}

func main() {
	p := Point{1, 2}
	fmt.Printf("%v\n", p) // {1 2}
}
