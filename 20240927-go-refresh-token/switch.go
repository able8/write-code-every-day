package main

import "fmt"

func f() int { return 1 }

func main() {
	a := 2

	fmt.Println("*** value switch ***")

	switch a {
	case 1:
		fmt.Println("1")
	case 2:
		fmt.Println("2")
	default:
		fmt.Println("default")
	}

	fmt.Println("\n*** multi-value case ***")
	switch a {
	case 1:
		fmt.Println("1")
	case 2, 3, 4:
		fmt.Println("2, 3, or 4")
		//	case 1, 2:
		//		fmt.Println("1 or 2")
	}

	fmt.Println("\n*** value switch with initializer ***")
	switch a = f(); a {
	case 1:
		fmt.Println("1")
	case 2:
		fmt.Println("2")
	default:
		fmt.Println("default")
	}

	fmt.Println("\n*** switch with case expressions ***")
	switch {
	case a == 1:
		fmt.Println("1")
	case a >= 2 && a <= 4:
		fmt.Println("2")
	case a <= 5:
		fmt.Println("3")
	}

	fmt.Println("\n*** type switch ***")

	var b any
	b = "hello"
	switch v := b.(type) {
	case int:
		fmt.Println("b is an int:", v)
	case string, []byte:
		fmt.Println("b is a string:", v)
	}

	fmt.Println("\n*** type parameter switch ***")
	do[bool](a)
	do[bool](true)
	do[int]([]int{1, 2, 3})
}

func do[T comparable](a any) {
	switch v := a.(type) {
	case int:
		fmt.Println("a is an int:", v)
	case T:
		fmt.Printf("a is of type %T: %v\n", v, v)
	case []T:
		fmt.Printf("a is a slice of %T: %v\n", v, v)
	case []byte:
		fmt.Println("a is a byte slice:", v)
	}
}
