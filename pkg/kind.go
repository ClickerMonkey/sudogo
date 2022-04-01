package sudogo

import (
	"math"
)

// The classic 9x9 puzzle with 9 boxes of 3x3 and digits 1-9.
var Classic = &Kind{
	BoxSize: Size{3, 3},
}

// A 4x4 puzzle with 4 boxes of 2x2 and digits 1-4.
var Kind2x2 = &Kind{
	BoxSize: Size{2, 2},
}

// A 6x6 puzzle with 6 boxes of 3x2 and digits 1-6.
var Kind3x2 = &Kind{
	BoxSize: Size{3, 2},
}

// A 12x12 puzzle with 12 boxes of 4x3 and digits 1-12.
var Kind4x3 = &Kind{
	BoxSize: Size{4, 3},
}

// A 16x16 puzzle with 16 boxes of 4x4 and digits 1-16.
var Kind4x4 = &Kind{
	BoxSize: Size{4, 4},
}

type Kind struct {
	BoxSize     Size
	Constraints []Constraint
}

func NewKind(boxWidth int, boxHeight int) *Kind {
	return &Kind{
		BoxSize: Size{
			Width:  boxWidth,
			Height: boxHeight,
		},
	}
}

// The width, height, and number of digits in this puzzle kind.
func (kind *Kind) Size() int {
	return kind.BoxSize.Width * kind.BoxSize.Height
}

// The number of unique digits in this puzzle kind.
func (kind *Kind) Digits() int {
	return kind.Size()
}

// How many boxes wide the puzzle would be.
func (kind *Kind) BoxesWide() int {
	return kind.BoxSize.Height
}

// How many boxes high the puzzle would be.
func (kind *Kind) BoxesHigh() int {
	return kind.BoxSize.Width
}

// How many cells would be in the puzzle.
func (kind *Kind) Area() int {
	return kind.Size() * kind.Size()
}

// How many characters it could take to print out the largest digit of a value in this puzzle kind.
func (kind *Kind) DigitsSize() int {
	return int(math.Floor(math.Log10(float64(kind.Digits())))) + 1
}

// Returns the dimensions of a puzzle of this kind.
func (kind *Kind) GetDimensions() (boxsWide int, boxsHigh int, boxWidth int, boxHeight int, size int) {
	w := kind.BoxSize.Width
	h := kind.BoxSize.Height
	return h, w, w, h, w * h
}

// Creates an empty puzzle of this kind.
func (kind *Kind) Empty() Puzzle {
	return New(kind)
}

// Creates a puzzle with an initial set of values of this kind.
func (kind *Kind) Create(values [][]int) Puzzle {
	instance := New(kind)
	instance.SetAll(values)
	return instance
}

// Creates a generator for puzzles of this kind.
func (kind *Kind) Generator() Generator {
	return NewGenerator(kind)
}
