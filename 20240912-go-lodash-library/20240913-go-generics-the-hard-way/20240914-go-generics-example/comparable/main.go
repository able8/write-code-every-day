package main

// NOTE This won't work since operator == not defined for T
// invalid operation: e == v (incomparable types in type set
func findFunc1[T any](a []T, v T) int {
	for i, e := range a {
		if e == v {
			return i
		}
	}
	return -1
}

func findFunc[T comparable](a []T, v T) int {
	for i, e := range a {
		if e == v {
			return i
		}
	}
	return -1
}

func main() {
	print(findFunc([]int{1, 2, 3, 4, 5, 6}, 5))
}
