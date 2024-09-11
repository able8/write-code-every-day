package main

import "fmt"

// Stack is a generic stack data structure.
type Stack[T any] struct {
	items []T
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the element from the top of the stack.
func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zeroValue T
		return zeroValue, fmt.Errorf("stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

func main() {
	// Example with integers
	intStack := Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)

	fmt.Println("Int Stack Size:", intStack.Size())

	poppedInt, err := intStack.Pop()
	if err == nil {
		fmt.Println("Popped Int:", poppedInt)
	}

	// Example with strings
	stringStack := Stack[string]{}
	stringStack.Push("apple")
	stringStack.Push("banana")
	stringStack.Push("cherry")

	fmt.Println("String Stack Size:", stringStack.Size())

	poppedString, err := stringStack.Pop()
	if err == nil {
		fmt.Println("Popped String:", poppedString)
	}
}
