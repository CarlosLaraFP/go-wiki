package interfaces

import (
	"fmt"
	"io"
	"os"
)

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

type File interface {
	Open(name string) (*os.File, error)
}

// PrintFile takes in an interface with (read text file -> print to terminal)
func PrintFile(file File) error {
	if len(os.Args) != 2 {
		return fmt.Errorf("filename must be provided as command line argument: %v", os.Args)
	}

	f, err := file.Open(os.Args[1])
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", os.Args[0], err)
	}
	defer f.Close()

	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		return fmt.Errorf("error reading & printing file %s: %v", os.Args[0], err)
	}

	return nil
}
