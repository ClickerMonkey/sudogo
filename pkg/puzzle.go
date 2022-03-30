package sudogo

import (
	"fmt"
	"strings"
)

type Puzzle struct {
	Kind  *Kind
	Cells []Cell
}

func New(kind *Kind) Puzzle {
	boxsWide, boxsHigh, _, boxHeight, size := kind.GetDimensions()
	cellCount := size * size
	cells := make([]Cell, cellCount)

	for i := 0; i < cellCount; i++ {
		cell := &cells[i]
		cell.Id = i
		cell.Value = 0
		cell.Row = i / size
		cell.Col = i % size
		cell.Box = ((cell.Row / boxsHigh) * boxHeight) + (cell.Col / boxsWide)
		cell.candidates.Fill(size)
	}

	return Puzzle{kind, cells}
}

func (puzzle *Puzzle) Clone() Puzzle {
	kind := puzzle.Kind
	n := len(puzzle.Cells)
	cells := make([]Cell, n)

	copy(cells, puzzle.Cells)

	return Puzzle{kind, cells}
}

func (puzzle *Puzzle) Clear() {
	candidates := puzzle.Kind.Digits()
	for i := range puzzle.Cells {
		c := &puzzle.Cells[i]
		c.Value = 0
		c.candidates.Fill(candidates)
	}
}

func (puzzle *Puzzle) Solver() Solver {
	return NewSolver(*puzzle)
}

func (puzzle *Puzzle) Get(col int, row int) *Cell {
	return &puzzle.Cells[row*puzzle.Kind.Size()+col]
}

func (puzzle *Puzzle) Set(col int, row int, value int) bool {
	return puzzle.SetCell(puzzle.Get(col, row), value)
}

func (puzzle *Puzzle) SetCell(cell *Cell, value int) bool {
	set := cell.SetValue(value)

	if set {
		for i := range puzzle.Cells {
			c := &puzzle.Cells[i]
			if c.Empty() && c.InGroup(cell) {
				c.RemoveCandidate(value)
			}
		}
	}

	return set
}

func (puzzle *Puzzle) Remove(col int, row int) bool {
	return puzzle.RemoveCell(&puzzle.Cells[row*puzzle.Kind.Size()+col])
}

func (puzzle *Puzzle) RemoveCell(cell *Cell) bool {
	removed := cell.HasValue()
	if removed {
		cell.Value = 0
		cell.candidates = puzzle.GetCandidatesFor(cell)

		for i := range puzzle.Cells {
			other := &puzzle.Cells[i]
			if cell.InGroup(other) {
				if other.Empty() {
					other.candidates = puzzle.GetCandidatesFor(other)
				}
			}
		}
	}
	return removed
}

func (puzzle *Puzzle) SetAll(values [][]int) int {
	size := puzzle.Kind.Size()
	sets := 0

	puzzle.Clear()

	for y, row := range values {
		if y >= size {
			break
		}
		for x, value := range row {
			if x >= size {
				break
			}
			if value > 0 && puzzle.Set(x, y, value) {
				sets++
			}
		}
	}

	return sets
}

func (puzzle *Puzzle) GetAll() [][]int {
	s := puzzle.Kind.Size()
	all := make([][]int, s)

	for y := 0; y < s; y++ {
		all[y] = make([]int, s)
	}

	for _, c := range puzzle.Cells {
		all[c.Row][c.Col] = c.Value
	}

	return all
}

func (puzzle *Puzzle) GetRow(rowIndex int) []int {
	row := make([]int, puzzle.Kind.Size())

	for _, c := range puzzle.Cells {
		if c.Row == rowIndex {
			row[c.Col] = c.Value
		}
	}

	return row
}

func (puzzle *Puzzle) GetRowCells(rowIndex int) []*Cell {
	row := make([]*Cell, puzzle.Kind.Size())

	for i := range puzzle.Cells {
		c := &puzzle.Cells[i]
		if c.Row == rowIndex {
			row[c.Col] = c
		}
	}

	return row
}

func (puzzle *Puzzle) GetColumn(columnIndex int) []int {
	column := make([]int, puzzle.Kind.Size())

	for _, c := range puzzle.Cells {
		if c.Col == columnIndex {
			column[c.Row] = c.Value
		}
	}

	return column
}

func (puzzle *Puzzle) GetColumnCells(columnIndex int) []*Cell {
	column := make([]*Cell, puzzle.Kind.Size())

	for i := range puzzle.Cells {
		c := &puzzle.Cells[i]
		if c.Col == columnIndex {
			column[c.Row] = c
		}
	}

	return column
}

func (puzzle *Puzzle) GetSymmetric(cell *Cell) *Cell {
	n := puzzle.Kind.Size() - 1

	return puzzle.Get(n-cell.Col, n-cell.Row)
}

func (puzzle *Puzzle) IsSolved() bool {
	size := puzzle.Kind.Size()
	rows := make([]Candidates, size)
	cols := make([]Candidates, size)
	boxs := make([]Candidates, size)

	complete := Candidates{}
	complete.Fill(size)

	for i := range puzzle.Cells {
		cell := &puzzle.Cells[i]
		if cell.Empty() {
			return false
		}

		rows[cell.Row].Set(cell.Value, true)
		cols[cell.Col].Set(cell.Value, true)
		boxs[cell.Box].Set(cell.Value, true)
	}

	for i := 0; i < size; i++ {
		if rows[i].Value != complete.Value {
			return false
		}
		if cols[i].Value != complete.Value {
			return false
		}
		if boxs[i].Value != complete.Value {
			return false
		}
	}

	return true
}

func (puzzle *Puzzle) GetCandidatesFor(cell *Cell) Candidates {
	candidates := Candidates{}
	candidates.Fill(puzzle.Kind.Size())

	for k := range puzzle.Cells {
		other := &puzzle.Cells[k]
		if cell.InGroup(other) && other.HasValue() {
			candidates.Set(other.Value, false)
		}
	}

	return candidates
}

func (puzzle *Puzzle) IsValid() bool {
	size := puzzle.Kind.Size()
	rows := make([]Candidates, size)
	cols := make([]Candidates, size)
	boxs := make([]Candidates, size)

	for i := range puzzle.Cells {
		cell := &puzzle.Cells[i]

		if cell.Value > 0 {
			if rows[cell.Row].Has(cell.Value) {
				return false
			}
			if cols[cell.Col].Has(cell.Value) {
				return false
			}
			if boxs[cell.Box].Has(cell.Value) {
				return false
			}

			rows[cell.Row].Set(cell.Value, true)
			cols[cell.Col].Set(cell.Value, true)
			boxs[cell.Box].Set(cell.Value, true)
		}

		candidates := puzzle.GetCandidatesFor(cell)

		if cell.Value > 0 {
			if !candidates.Has(cell.Value) {
				return false
			}
		} else {
			if candidates.Count == 0 {
				return false
			}
			if (candidates.Value & cell.candidates.Value) != cell.candidates.Value {
				return false
			}
		}
	}
	return true
}

func (puzzle *Puzzle) UniqueId() string {
	sb := strings.Builder{}
	for k := range puzzle.Cells {
		other := &puzzle.Cells[k]
		if other.Empty() {
			sb.WriteString(".")
		} else {
			sb.WriteString(fmt.Sprint(other.Value))
		}
	}
	return sb.String()
}

func (puzzle *Puzzle) HasUniqueSolution() bool {
	return len(puzzle.GetSolutions(SolutionLimit{maxSolutions: 2})) == 1
}

type SolutionLimit struct {
	SolverLimit
	maxSolutions int
}

func (puzzle *Puzzle) GetSolutions(limits SolutionLimit) []*Solver {
	solutions := make([]*Solver, 0, max(0, limits.maxSolutions))
	unique := map[string]bool{}

	solvers := NewQueue[Solver]()
	solvers.Offer(puzzle.Solver())

	for !solvers.Empty() {
		solver := solvers.Poll()
		solution, solved := solver.Solve(limits.SolverLimit)

		if !solved {
			min := solver.GetMinCandidateCount()
			minCell := solver.GetGroupWhere(func(group *CellGroups) bool {
				return group.Cell.candidates.Count == min
			})
			if minCell != nil {
				for _, candidate := range minCell.Cell.Candidates() {
					newSolver := solution.Solver()
					if newSolver.Set(minCell.Cell.Col, minCell.Cell.Row, candidate) {
						solvers.Offer(newSolver)
					}
				}
			}
		} else {
			id := solution.UniqueId()
			if !unique[id] {
				solutions = append(solutions, solver)
				if limits.maxSolutions > 0 && len(solutions) == limits.maxSolutions {
					break
				}
				unique[id] = true
			}
		}
	}

	return solutions
}
