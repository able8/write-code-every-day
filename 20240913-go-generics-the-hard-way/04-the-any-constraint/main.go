package main

import (
	"fmt"
)

// Sum returns the sum of the provided arguments.
// func Sum[T any](args ...T) T {
func Sum[T int | int8 | int32 | int64](args ...T) T {
	var sum T
	for i := 0; i < len(args); i++ {
		sum += args[i]
	}
	return sum
}

func main() {
	fmt.Println(Sum([]int{1, 2, 3}...))
	fmt.Println(Sum([]int8{1, 2, 3}...))
}
