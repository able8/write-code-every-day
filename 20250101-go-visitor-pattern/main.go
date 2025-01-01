package main

import "fmt"

// Visitor is a behavioral design pattern
// that allows adding new behaviors to existing class hierarchy
//  without altering any existing code.

type Visitor interface {
	visitForSquare(*Square)
	visitForCircle(*Circle)
}

type Shape interface {
	getType() string
	accept(Visitor)
}

type Square struct {
	side int
}

func (s *Square) accept(v Visitor) {
	v.visitForSquare(s)
}

func (s *Square) getType() string {
	return "Square"
}

type Circle struct {
	radius int
}

func (c *Circle) accept(v Visitor) {
	v.visitForCircle(c)
}

func (c *Circle) getType() string {
	return "Circle"
}

type AreaCalculator struct {
	area int
}

func (a *AreaCalculator) visitForSquare(s *Square) {
	// Calculate area for square.
	// Then assign in to the area instance variable.
	fmt.Println("Calculating area for square")
}

func (a *AreaCalculator) visitForCircle(s *Circle) {
	fmt.Println("Calculating area for circle")

}

type MiddleCoordinates struct {
	x int
	y int
}

func (a *MiddleCoordinates) visitForSquare(s *Square) {
	// Calculate middle point coordinates for square.
	// Then assign in to the x and y instance variable.
	fmt.Println("Calculating middle point coordinates for square")
}

func (a *MiddleCoordinates) visitForCircle(c *Circle) {
	fmt.Println("Calculating middle point coordinates for circle")
}

func main() {
	square := &Square{side: 2}
	circle := &Circle{radius: 3}

	areaCalculator := &AreaCalculator{}

	square.accept(areaCalculator)
	circle.accept(areaCalculator)

	fmt.Println()
	middleCoordinates := &MiddleCoordinates{}
	square.accept(middleCoordinates)
	circle.accept(middleCoordinates)
}
