package sudogo

import (
	"math"
	"math/rand"
	"time"
)

type Constraint func(puzzle *PuzzleInstance) bool

type Puzzle struct {
	Boxs        Size
	BoxSize     Size
	Constraints []Constraint
}

func (puzzle *Puzzle) Width() int {
	return puzzle.BoxSize.Width * puzzle.Boxs.Width
}

func (puzzle *Puzzle) Height() int {
	return puzzle.BoxSize.Height * puzzle.Boxs.Height
}

func (puzzle *Puzzle) Digits() int {
	return puzzle.BoxSize.Area()
}

func (puzzle *Puzzle) DigitsSize() int {
	return int(math.Floor(math.Log10(float64(puzzle.Digits())))) + 1
}

func (puzzle *Puzzle) DefaultCandidates() uint64 {
	return uint64((1 << (puzzle.Digits() + 1)) - 1)
}

func (puzzle *Puzzle) Empty() PuzzleInstance {
	return NewInstance(puzzle)
}

func (puzzle *Puzzle) Create(values [][]int) PuzzleInstance {
	instance := NewInstance(puzzle)
	instance.SetAll(values)
	return instance
}

func (puzzle *Puzzle) Generate() PuzzleInstance {
	instance := NewInstance(puzzle)
	instance.Generate(rand.New(rand.NewSource(time.Now().Unix())))
	return instance
}

var Classic = &Puzzle{
	Boxs:    Size{3, 3},
	BoxSize: Size{3, 3},
}
