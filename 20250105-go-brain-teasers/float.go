package main

import (
	"fmt"
)

func main() {
	n := 1.1
	fmt.Println(n * n)

	fmt.Println(1.1 * 1.1)
	fmt.Println(1.1 * n)

	a := 1.1
	b := 100.0
	fmt.Println(a*b == 110)
	fmt.Println(a * b)
}
