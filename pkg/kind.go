package sudogo

import (
	"math"
)

type Constraint func(puzzle *Puzzle) bool

type Kind struct {
	Boxs        Size
	BoxSize     Size
	Constraints []Constraint
}

func (kind *Kind) Width() int {
	return kind.BoxSize.Width * kind.Boxs.Width
}

func (kind *Kind) Height() int {
	return kind.BoxSize.Height * kind.Boxs.Height
}

func (kind *Kind) Digits() int {
	return kind.BoxSize.Area()
}

func (kind *Kind) Area() int {
	return kind.Width() * kind.Height()
}

func (kind *Kind) DigitsSize() int {
	return int(math.Floor(math.Log10(float64(kind.Digits())))) + 1
}

func (kind *Kind) Empty() Puzzle {
	return New(kind)
}

func (kind *Kind) Create(values [][]int) Puzzle {
	instance := New(kind)
	instance.SetAll(values)
	return instance
}

func (kind *Kind) Generator() Generator {
	return NewGenerator(kind)
}

var Classic = &Kind{
	Boxs:    Size{3, 3},
	BoxSize: Size{3, 3},
}
