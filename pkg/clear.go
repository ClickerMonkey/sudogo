package sudogo

type ClearLimits struct {
	SolverLimit
	Symmetric bool
	MaxStates int
}

func (limits ClearLimits) Extend(extend ClearLimits) ClearLimits {
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

var DifficultyBeginner = ClearLimits{
	SolverLimit: SolverLimit{MinCost: 3600, MaxCost: 4500},
	Symmetric:   true,
}
var DifficultyEasy = ClearLimits{
	SolverLimit: SolverLimit{MinCost: 4300, MaxCost: 5500},
	Symmetric:   true,
}
var DifficultyMedium = ClearLimits{
	SolverLimit: SolverLimit{MinCost: 5300, MaxCost: 6900},
	Symmetric:   true,
}
var DifficultyHard = ClearLimits{
	SolverLimit: SolverLimit{MinCost: 6000, MaxCost: 7200},
	Symmetric:   false,
}
var DifficultyTricky = ClearLimits{
	SolverLimit: SolverLimit{MinCost: 6500, MaxCost: 9300},
	Symmetric:   false,
}
var DifficultyFiendish = ClearLimits{
	SolverLimit: SolverLimit{MinCost: 8300, MaxCost: 14000},
	Symmetric:   false,
}
var DifficultyDiabolical = ClearLimits{
	SolverLimit: SolverLimit{MinCost: 11000, MaxCost: 25000},
	Symmetric:   false,
}

func (gen *Generator) ClearCells(puzzle *Puzzle, limits ClearLimits) (*Puzzle, int) {
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

		nextSolutions := next.GetSolutions(SolutionLimit{
			SolverLimit:  limits.SolverLimit,
			MaxSolutions: 2,
		})

		if len(nextSolutions) == 1 {
			uniqueSolution := nextSolutions[0]
			states++

			if !uniqueSolution.CanContinue(limits.SolverLimit, 0) {
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
