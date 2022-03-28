package sudogo

import (
	"math/rand"
	"time"
)

// Generates puzzles of the given kind using a particular random number generator.
type Generator struct {
	kind   *Kind
	solver Solver
	random *rand.Rand
}

func NewGenerator(kind *Kind) Generator {
	return NewSeededGenerator(kind, time.Now().Unix())
}

func NewSeededGenerator(kind *Kind, seed int64) Generator {
	return NewRandomGenerator(kind, rand.New(rand.NewSource(seed)))
}

func generatorSolver(kind *Kind) Solver {
	solver := NewSolver(kind.Empty())
	solver.steps = GenerateSolveSteps
	return solver
}

func NewRandomGenerator(kind *Kind, random *rand.Rand) Generator {
	return Generator{kind, generatorSolver(kind), random}
}

func (gen *Generator) Reset() {
	gen.solver = generatorSolver(gen.kind)
}

func (gen *Generator) Puzzle() *Puzzle {
	return &gen.solver.puzzle
}

func (gen *Generator) Solver() *Solver {
	return &gen.solver
}

func (gen *Generator) IsComplete() bool {
	return gen.solver.Solved()
}

func (gen *Generator) GetUnsolved() *CellGroups {
	return pointerAt(gen.solver.unsolved, 0)
}

func (gen *Generator) GetRandomUnsolved() *CellGroups {
	return randomPointer(gen.random, gen.solver.unsolved)
}

func (gen *Generator) GetRandom(match func(other *Cell) bool) *CellGroups {
	matches := 0
	for _, group := range gen.solver.cells {
		if match(group.cell) {
			matches++
		}
	}
	if matches == 0 {
		return nil
	}
	chosen := gen.random.Intn(matches)
	for _, group := range gen.solver.cells {
		if match(group.cell) {
			chosen--
			if chosen < 0 {
				return &group
			}
		}
	}
	return nil
}

// The smallest number of candidates in a cell without a value. 0=solved, 1=naked
func (gen *Generator) GetPressure() int {
	pressure := 0
	for _, group := range gen.solver.unsolved {
		if pressure == 0 || pressure > group.cell.candidates.Count {
			pressure = group.cell.candidates.Count
		}
	}
	return pressure
}

func (gen *Generator) GetRandomPressured() *CellGroups {
	minCount := gen.GetPressure()

	return gen.GetRandom(func(other *Cell) bool {
		return other.Empty() && other.candidates.Count == minCount
	})
}

func (gen *Generator) Attempt() *Puzzle {
	for !gen.IsComplete() {
		gen.solver.Solve()

		if gen.IsComplete() {
			break
		}

		if gen.GetPressure() == 0 {
			return nil
		}

		randomGroup := gen.GetRandomUnsolved()
		randomValue := randomElement[int](gen.random, randomGroup.cell.Candidates(), 0)

		gen.solver.SetGroup(randomGroup, randomValue)
	}
	return gen.Puzzle()
}

func (gen *Generator) Attempts(tries int) (*Puzzle, int) {
	var generated *Puzzle = nil
	for i := 0; i < tries; i++ {
		generated = gen.Attempt()
		if generated != nil {
			return generated, i + 1
		} else {
			gen.Reset()
		}
	}
	return nil, tries
}

func (gen *Generator) Generate() (*Puzzle, int) {
	return gen.Attempts(1 << 14)
}
