package sudogo

import (
	"math/rand"
	"time"
)

// Generates puzzles of the given kind using a particular random number generator.
type Generator struct {
	Kind   *Kind
	solver Solver
	Random *rand.Rand
}

func NewGenerator(kind *Kind) Generator {
	return NewRandomGenerator(kind, Random())
}

func NewSeededGenerator(kind *Kind, seed int64) Generator {
	return NewRandomGenerator(kind, RandomSeeded(seed))
}

func RandomSeed() int64 {
	return time.Now().UnixNano()
}

func Random() *rand.Rand {
	return RandomSeeded(RandomSeed())
}

func RandomSeeded(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func generatorSolver(kind *Kind) Solver {
	solver := NewSolver(kind.Empty())
	solver.Steps = GenerateSolveSteps
	return solver
}

func NewRandomGenerator(kind *Kind, random *rand.Rand) Generator {
	return Generator{kind, generatorSolver(kind), random}
}

func (gen *Generator) Reset() {
	gen.solver = generatorSolver(gen.Kind)
}

func (gen *Generator) Puzzle() *Puzzle {
	return &gen.solver.Puzzle
}

func (gen *Generator) Solver() *Solver {
	return &gen.solver
}

func (gen *Generator) IsComplete() bool {
	return gen.solver.Solved()
}

func (gen *Generator) GetUnsolved() *Cell {
	return pointerAt(gen.solver.Unsolved, 0)
}

func (gen *Generator) GetRandomUnsolved() *Cell {
	return randomPointer(gen.Random, gen.solver.Unsolved)
}

func (gen *Generator) Attempt() *Puzzle {
	for !gen.IsComplete() {
		gen.solver.Solve(SolveLimit{})

		if gen.IsComplete() {
			break
		}

		if gen.solver.GetMinCandidateCount() == 0 {
			return nil
		}

		randomCell := gen.GetRandomUnsolved()
		randomValue := randomElement(gen.Random, randomCell.Candidates(), 0)

		gen.solver.SetCell(randomCell, randomValue)
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
