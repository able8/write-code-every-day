package main

import (
	"fmt"
	"log"
)

type PrintInterface interface {
	print()
}

type PrintNum int

func (p PrintNum) print() {
	log.Println("Num:", p)
}

type PrintText string

func (p PrintText) print() {
	log.Println("Text:", p)
}

// PrintSlice prints elements of a slice of any type
func PrintSlice[T PrintInterface](s []T) {
	for _, value := range s {
		value.print()
	}
}

func main() {
	// Example with PrintSlice
	stringSlice := []PrintText{"apple", "banana", "orange"}
	intSlice := []PrintNum{5, 2, 9, 1, 7}

	fmt.Println("String slice:")
	PrintSlice(stringSlice)

	fmt.Println("\nInteger slice:")
	PrintSlice(intSlice)
}
