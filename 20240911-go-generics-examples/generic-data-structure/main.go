// https://www.atatus.com/blog/generics-in-golang/

package main

import (
	"fmt"
	"reflect"
)

// a generic data structure, a stack, that uses a type parameter:
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}
func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		panic("stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func main() {
	// The type parameter "T" is replaced with the actual types of the stacks create
	intStack := &Stack[int]{}
	stringStack := &Stack[string]{}
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)
	stringStack.Push("apple")
	stringStack.Push("banana")
	stringStack.Push("cherry")
	fmt.Println(intStack.Pop())    // prints 3
	fmt.Println(stringStack.Pop()) // prints cherry
}

// Generic functions:
// This function takes a slice of any type T and a value of type T and returns the index of the value in the slice.
// The any keyword in the type parameter specifies that any type can be used.
func findIndex[T any](slice []T, value T) int {
	for i, v := range slice {
		if reflect.DeepEqual(v, value) {
			return i
		}
	}
	return -1
}

type Equatable interface {
	Equals(other interface{}) bool
}

// Constraints on type parameters:
// This defines a type constraint on the type parameter T that requires it to implement the Equatable interface.
// This allows the findIndex function to use the Equals method to compare values of type T.
func findIndex2[T Equatable](slice []T, value T) int {
	for i, v := range slice {
		if v.Equals(value) {
			return i
		}
	}
	return -1
}
