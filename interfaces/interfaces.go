package interfaces

import "fmt"

type Shape interface {
	GetArea() float64
	PrintArea()
}

type Square struct {
	Length float64
}

func (s *Square) GetArea() float64 {
	return s.Length * s.Length
}

func (s *Square) PrintArea() {
	fmt.Printf("Square area: %f\n", s.GetArea())
}

type Triangle struct {
	Base   float64
	Height float64
}

func (t *Triangle) GetArea() float64 {
	return t.Base * t.Height * 0.5
}

func (t *Triangle) PrintArea() {
	fmt.Printf("Triangle area: %f\n", t.GetArea())
}

func ShowShape(s Shape) {
	s.PrintArea()
}
