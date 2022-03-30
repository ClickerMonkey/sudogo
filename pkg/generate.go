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

func (gen *Generator) GetRandomPressured() *CellGroups {
	minCount := gen.solver.GetMinCandidateCount()

	return gen.GetRandom(func(other *Cell) bool {
		return other.Empty() && other.candidates.Count == minCount
	})
}

func (gen *Generator) Attempt() *Puzzle {
	for !gen.IsComplete() {
		gen.solver.Solve(SolverLimit{})

		if gen.IsComplete() {
			break
		}

		if gen.solver.GetMinCandidateCount() == 0 {
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

type ClearLimits struct {
	SolverLimit
	symmetric bool
	maxStates int
}

func (limits ClearLimits) Extend(extend ClearLimits) ClearLimits {
	out := limits
	if extend.symmetric && !out.symmetric {
		out.symmetric = true
	}
	if extend.maxBatches > 0 {
		out.maxBatches = extend.maxBatches
	}
	if extend.maxCost > 0 {
		out.maxCost = extend.maxCost
	}
	if extend.minCost > 0 {
		out.minCost = extend.minCost
	}
	if extend.maxLogs > 0 {
		out.maxLogs = extend.maxLogs
	}
	if extend.maxPlacements > 0 {
		out.maxPlacements = extend.maxPlacements
	}
	if extend.maxStates > 0 {
		out.maxStates = extend.maxStates
	}
	return out
}

var DifficultyBeginner = ClearLimits{
	SolverLimit: SolverLimit{minCost: 3600, maxCost: 4500},
	symmetric:   true,
}
var DifficultyEasy = ClearLimits{
	SolverLimit: SolverLimit{minCost: 4300, maxCost: 5500},
	symmetric:   true,
}
var DifficultyMedium = ClearLimits{
	SolverLimit: SolverLimit{minCost: 5300, maxCost: 6900},
	symmetric:   true,
}
var DifficultyTricky = ClearLimits{
	SolverLimit: SolverLimit{minCost: 6500, maxCost: 9300},
	symmetric:   true,
}
var DifficultyFiendish = ClearLimits{
	SolverLimit: SolverLimit{minCost: 8300, maxCost: 14000},
	symmetric:   false,
}
var DifficultyDiabolical = ClearLimits{
	SolverLimit: SolverLimit{minCost: 11000, maxCost: 25000},
	symmetric:   false,
}

func (gen *Generator) ClearCells(puzzle *Puzzle, limits ClearLimits) (*Puzzle, int) {
	if puzzle == nil || (limits.maxBatches == 0 && limits.maxCost == 0 && limits.maxLogs == 0 && limits.maxPlacements == 0 && limits.maxStates == 0) {
		return nil, 0
	}

	states := 0

	type AttemptState struct {
		puzzle    Puzzle
		solver    Solver
		available []*Cell
	}

	attempts := NewStack[AttemptState](limits.maxPlacements)

	initial := puzzle.Clone()
	attempts.Push(AttemptState{
		puzzle: initial,
		solver: initial.Solver(),
		available: pointersWhere(initial.cells, func(cell *Cell) bool {
			return cell.HasValue()
		}),
	})

	for !attempts.Empty() {
		last := attempts.Peek()

		if len(last.available) == 0 {
			attempts.Pop()
			continue
		}

		next := last.puzzle.Clone()
		nextSolver := next.Solver()

		cell := randomPointer(gen.random, last.available)
		cellSymmetric := last.puzzle.GetSymmetric(cell)

		doSymmetric := limits.symmetric && cellSymmetric.HasValue()

		next.Remove(cell.col, cell.row)
		if doSymmetric {
			next.Remove(cellSymmetric.col, cellSymmetric.row)
		}

		last.available = removeValue(last.available, cell)
		if doSymmetric {
			last.available = removeValue(last.available, cellSymmetric)
		}

		if len(last.available) == 0 {
			attempts.Pop()
		}

		nextSolutions := next.GetSolutions(SolutionLimit{
			SolverLimit:  limits.SolverLimit,
			maxSolutions: 2,
		})

		if len(nextSolutions) == 1 {
			uniqueSolution := nextSolutions[0]
			states++

			if !uniqueSolution.canContinue(limits.SolverLimit, 0) {
				return &next, states
			}

			if limits.maxStates > 0 && states >= limits.maxStates {
				break
			}

			attempts.Push(AttemptState{
				puzzle:    next,
				solver:    nextSolver,
				available: sliceClone(last.available),
			})
		}
	}

	return nil, states
}
