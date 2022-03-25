package sudogo

import (
	"math"
)

type Constraint func(puzzle *Puzzle) bool

type PuzzleKind struct {
	Boxs        Size
	BoxSize     Size
	Constraints []Constraint
}

func (kind *PuzzleKind) Width() int {
	return kind.BoxSize.Width * kind.Boxs.Width
}

func (kind *PuzzleKind) Height() int {
	return kind.BoxSize.Height * kind.Boxs.Height
}

func (kind *PuzzleKind) Digits() int {
	return kind.BoxSize.Area()
}

func (kind *PuzzleKind) Area() int {
	return kind.Width() * kind.Height()
}

func (kind *PuzzleKind) DigitsSize() int {
	return int(math.Floor(math.Log10(float64(kind.Digits())))) + 1
}

func (kind *PuzzleKind) DefaultCandidates() uint64 {
	return uint64((1 << (kind.Digits() + 1)) - 1)
}

func (kind *PuzzleKind) Empty() Puzzle {
	return New(kind)
}

func (kind *PuzzleKind) Create(values [][]int) Puzzle {
	instance := New(kind)
	instance.SetAll(values)
	return instance
}

func (kind *PuzzleKind) Generator() Generator {
	return NewGenerator(kind)
}

var Classic = &PuzzleKind{
	Boxs:    Size{3, 3},
	BoxSize: Size{3, 3},
}
