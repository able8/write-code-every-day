package main

import "fmt"

// Numeric expresses a type constraint satisfied by any numeric type.
type Numeric interface {
	uint | uint8 | uint16 | uint32 | uint64 |
		int | int8 | int16 | int32 | ~int64 |
		float32 | float64 |
		complex64 | complex128
}

func NewT[T any]() *T {

	// NewT[T] *T returns a *T with a nil value since there was no instance of a T allocated.

	// var t *T
	// return t

	// the output will be a memory address that references a new instance of an int value.
	return new(T)
}

func Sum[T Numeric](args ...T) T {
	sum := new(T)
	for i := 0; i < len(args); i++ {
		*sum += args[i]
	}
	return *sum
}

func main() {
	fmt.Println(NewT[int]())
}
