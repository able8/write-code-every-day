package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3}
	b := append(a[:1], 10)
	fmt.Printf("a=%v, b=%v\n", a, b)
}
