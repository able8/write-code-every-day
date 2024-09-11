package main

// https://github.com/akutz/go-generics-the-hard-way/blob/main/03-getting-started/02-syntax.md

import "fmt"

// Sum returns the sum of the provided arguments.
func Sum(args ...int) int {
	var sum int
	for i := 0; i < len(args); i++ {
		sum += args[i]
	}
	return sum
}

// Sum returns the sum of the provided arguments.
func SumGen[T int | int64](args ...T) T {
	var sum T
	for i := 0; i < len(args); i++ {
		sum += args[i]
	}
	return sum
}

func main() {
	fmt.Println(SumGen([]int64{1, 2, 3}...))
}
