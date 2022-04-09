package sudogo

type ClearLimit struct {
	SolveLimit
	Symmetric bool
	MaxStates int
	Fast      bool
}

func (limits ClearLimit) Extend(extend ClearLimit) ClearLimit {
	out := limits
	if extend.Symmetric && !out.Symmetric {
		out.Symmetric = true
	}
	if extend.MaxBatches > 0 {
		out.MaxBatches = extend.MaxBatches
	}
	if extend.MaxCost > 0 {
		out.MaxCost = extend.MaxCost
	}
	if extend.MinCost > 0 {
		out.MinCost = extend.MinCost
	}
	if extend.MaxLogs > 0 {
		out.MaxLogs = extend.MaxLogs
	}
	if extend.MaxPlacements > 0 {
		out.MaxPlacements = extend.MaxPlacements
	}
	if extend.MaxStates > 0 {
		out.MaxStates = extend.MaxStates
	}
	return out
}

var DifficultyBeginner = ClearLimit{
	SolveLimit: SolveLimit{MinCost: 3600, MaxCost: 4500},
	Symmetric:  true,
}
var DifficultyEasy = ClearLimit{
	SolveLimit: SolveLimit{MinCost: 4300, MaxCost: 5500},
	Symmetric:  true,
}
var DifficultyMedium = ClearLimit{
	SolveLimit: SolveLimit{MinCost: 5300, MaxCost: 6900},
	Symmetric:  true,
}
var DifficultyHard = ClearLimit{
	SolveLimit: SolveLimit{MinCost: 6000, MaxCost: 7200},
	Symmetric:  false,
}
var DifficultyTricky = ClearLimit{
	SolveLimit: SolveLimit{MinCost: 6500, MaxCost: 9300},
	Symmetric:  false,
}
var DifficultyFiendish = ClearLimit{
	SolveLimit: SolveLimit{MinCost: 8300, MaxCost: 14000},
	Symmetric:  false,
}
var DifficultyDiabolical = ClearLimit{
	SolveLimit: SolveLimit{MinCost: 11000, MaxCost: 25000},
	Symmetric:  false,
}

func (gen *Generator) ClearCells(puzzle *Puzzle, limits ClearLimit) (*Puzzle, int) {
	if puzzle == nil || (limits.MaxBatches == 0 && limits.MaxCost == 0 && limits.MaxLogs == 0 && limits.MaxPlacements == 0 && limits.MaxStates == 0) {
		return nil, 0
	}

	states := 0

	type AttemptState struct {
		puzzle    Puzzle
		available []*Cell
	}

	attempts := NewStack[AttemptState](limits.MaxPlacements)

	initial := puzzle.Clone()
	attempts.Push(AttemptState{
		puzzle: initial,
		available: pointersWhere(initial.Cells, func(cell *Cell) bool {
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

		cell := randomPointer(gen.Random, last.available)
		cellSymmetric := last.puzzle.GetSymmetric(cell)

		doSymmetric := limits.Symmetric && cellSymmetric.HasValue()

		next.Remove(cell.Col, cell.Row)
		if doSymmetric {
			next.Remove(cellSymmetric.Col, cellSymmetric.Row)
		}

		last.available = removeValue(last.available, cell)
		if doSymmetric {
			last.available = removeValue(last.available, cellSymmetric)
		}

		if len(last.available) == 0 {
			attempts.Pop()
		}

		var solver *Solver = nil

		if limits.Fast {
			nextSolver := next.Solver()
			nextSolution, _ := nextSolver.Solve(limits.SolveLimit)
			if nextSolution != nil && nextSolution.IsSolved() {
				solver = &nextSolver
			}
		} else {
			nextSolutions := next.GetSolutions(SolutionsLimit{
				SolveLimit:   limits.SolveLimit,
				MaxSolutions: 2,
			})

			if len(nextSolutions) == 1 {
				solver = nextSolutions[0]
			}
		}

		if solver != nil {
			states++

			if !solver.CanContinue(limits.SolveLimit, 0) {
				return &next, states
			}

			if limits.MaxStates > 0 && states >= limits.MaxStates {
				break
			}

			attempts.Push(AttemptState{
				puzzle:    next,
				available: sliceClone(last.available),
			})
		}
	}

	return nil, states
}
