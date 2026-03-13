package interfaces

import (
	"fmt"
	"math"
)

// Interface definition
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Struct 1
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Struct 2
type Square struct {
	Side float64
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

func (s Square) Perimeter() float64 {
	return 4 * s.Side
}

// Function using interface
func PrintShapeInfo(s Shape) {
	fmt.Printf("Type: %T\n", s)
	fmt.Printf("Area: %.2f\n", s.Area())
	fmt.Printf("Perimeter: %.2f\n\n", s.Perimeter())
}

func Run() {
	c := Circle{Radius: 5}
	s := Square{Side: 4}

	PrintShapeInfo(c)
	PrintShapeInfo(s)
}
