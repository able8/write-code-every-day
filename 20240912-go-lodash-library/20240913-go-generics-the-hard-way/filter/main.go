package main

import (
	"fmt"
)

func filterFunc[T any](a []T, f func(T) bool) []T {
	var n []T
	for _, e := range a {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func main() {
	vi := []int{1, 2, 3, 4, 5, 6}
	vi = filterFunc[int](vi, func(v int) bool {
		return v < 4
	})
	fmt.Println(vi)

	vi2 := []string{"1", "12", "123"}
	vi2 = filterFunc[string](vi2, func(v string) bool {
		return len(v) > 1
	})
	fmt.Println(vi2)

}
