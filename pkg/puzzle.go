package sudogo

import (
	"fmt"
	"strings"
)

type Puzzle struct {
	kind  *Kind
	cells []Cell
}

func New(kind *Kind) Puzzle {
	boxsWide, boxsHigh, _, boxHeight, size := kind.GetDimensions()
	cellCount := size * size
	cells := make([]Cell, cellCount)

	for i := 0; i < cellCount; i++ {
		cell := &cells[i]
		cell.id = i
		cell.value = 0
		cell.row = i / size
		cell.col = i % size
		cell.box = ((cell.row / boxsHigh) * boxHeight) + (cell.col / boxsWide)
		cell.candidates.Fill(size)
	}

	return Puzzle{kind, cells}
}

func (puzzle *Puzzle) Clone() Puzzle {
	kind := puzzle.kind
	n := len(puzzle.cells)
	cells := make([]Cell, n)

	copy(cells, puzzle.cells)

	return Puzzle{kind, cells}
}

func (puzzle *Puzzle) Clear() {
	candidates := puzzle.kind.Digits()
	for i := range puzzle.cells {
		c := &puzzle.cells[i]
		c.value = 0
		c.candidates.Fill(candidates)
	}
}

func (puzzle *Puzzle) Solver() Solver {
	return NewSolver(*puzzle)
}

func (puzzle *Puzzle) Get(col int, row int) *Cell {
	return &puzzle.cells[row*puzzle.kind.Size()+col]
}

func (puzzle *Puzzle) Set(col int, row int, value int) bool {
	return puzzle.SetCell(puzzle.Get(col, row), value)
}

func (puzzle *Puzzle) SetCell(cell *Cell, value int) bool {
	set := cell.SetValue(value)

	if set {
		for i := range puzzle.cells {
			c := &puzzle.cells[i]
			if c.Empty() && c.InGroup(cell) {
				c.RemoveCandidate(value)
			}
		}
	}

	return set
}

func (puzzle *Puzzle) Remove(col int, row int) bool {
	return puzzle.RemoveCell(&puzzle.cells[row*puzzle.kind.Size()+col])
}

func (puzzle *Puzzle) RemoveCell(cell *Cell) bool {
	removed := cell.HasValue()
	if removed {
		cell.value = 0
		cell.candidates = puzzle.GetCandidatesFor(cell)

		for i := range puzzle.cells {
			other := &puzzle.cells[i]
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
	size := puzzle.kind.Size()
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
	s := puzzle.kind.Size()
	all := make([][]int, s)

	for y := 0; y < s; y++ {
		all[y] = make([]int, s)
	}

	for _, c := range puzzle.cells {
		all[c.row][c.col] = c.value
	}

	return all
}

func (puzzle *Puzzle) GetRow(rowIndex int) []int {
	row := make([]int, puzzle.kind.Size())

	for _, c := range puzzle.cells {
		if c.row == rowIndex {
			row[c.col] = c.value
		}
	}

	return row
}

func (puzzle *Puzzle) GetRowCells(rowIndex int) []*Cell {
	row := make([]*Cell, puzzle.kind.Size())

	for i := range puzzle.cells {
		c := &puzzle.cells[i]
		if c.row == rowIndex {
			row[c.col] = c
		}
	}

	return row
}

func (puzzle *Puzzle) GetColumn(columnIndex int) []int {
	column := make([]int, puzzle.kind.Size())

	for _, c := range puzzle.cells {
		if c.col == columnIndex {
			column[c.row] = c.value
		}
	}

	return column
}

func (puzzle *Puzzle) GetColumnCells(columnIndex int) []*Cell {
	column := make([]*Cell, puzzle.kind.Size())

	for i := range puzzle.cells {
		c := &puzzle.cells[i]
		if c.col == columnIndex {
			column[c.row] = c
		}
	}

	return column
}

func (puzzle *Puzzle) GetSymmetric(cell *Cell) *Cell {
	n := puzzle.kind.Size() - 1

	return puzzle.Get(n-cell.col, n-cell.row)
}

func (puzzle *Puzzle) IsSolved() bool {
	size := puzzle.kind.Size()
	rows := make([]Candidates, size)
	cols := make([]Candidates, size)
	boxs := make([]Candidates, size)

	complete := Candidates{}
	complete.Fill(size)

	for i := range puzzle.cells {
		cell := &puzzle.cells[i]
		if cell.Empty() {
			return false
		}

		rows[cell.row].Set(cell.value, true)
		cols[cell.col].Set(cell.value, true)
		boxs[cell.box].Set(cell.value, true)
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
	candidates.Fill(puzzle.kind.Size())

	for k := range puzzle.cells {
		other := &puzzle.cells[k]
		if cell.InGroup(other) && other.HasValue() {
			candidates.Set(other.value, false)
		}
	}

	return candidates
}

func (puzzle *Puzzle) IsValid() bool {
	size := puzzle.kind.Size()
	rows := make([]Candidates, size)
	cols := make([]Candidates, size)
	boxs := make([]Candidates, size)

	for i := range puzzle.cells {
		cell := &puzzle.cells[i]

		if cell.value > 0 {
			if rows[cell.row].Has(cell.value) {
				return false
			}
			if cols[cell.col].Has(cell.value) {
				return false
			}
			if boxs[cell.box].Has(cell.value) {
				return false
			}

			rows[cell.row].Set(cell.value, true)
			cols[cell.col].Set(cell.value, true)
			boxs[cell.box].Set(cell.value, true)
		}

		candidates := puzzle.GetCandidatesFor(cell)

		if cell.value > 0 {
			if !candidates.Has(cell.value) {
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
	for k := range puzzle.cells {
		other := &puzzle.cells[k]
		if other.Empty() {
			sb.WriteString(".")
		} else {
			sb.WriteString(fmt.Sprint(other.value))
		}
	}
	return sb.String()
}

func (puzzle *Puzzle) HasUniqueSolution() bool {
	return len(puzzle.GetSolutions(2)) == 1
}

func (puzzle *Puzzle) GetSolutions(max int) []*Puzzle {
	solutions := make([]*Puzzle, 0)
	solvers := make([]Solver, 0)
	solvers = append(solvers, puzzle.Solver())
	unique := map[string]bool{}

	for len(solvers) > 0 {
		solver := solvers[0]
		solvers = solvers[1:]
		solution, solved := solver.Solve()

		if !solved {
			min := solver.GetMinCandidateCount()
			minCell := solver.GetGroupWhere(func(group *CellGroups) bool {
				return group.cell.candidates.Count == min
			})
			if minCell != nil {
				for _, candidate := range minCell.cell.Candidates() {
					newSolver := solution.Solver()
					if newSolver.Set(minCell.cell.col, minCell.cell.row, candidate) {
						solvers = append(solvers, newSolver)
					}
				}
			}
		} else {
			id := solution.UniqueId()
			if !unique[id] {
				solutions = append(solutions, solution)
				if max > 0 && len(solutions) == max {
					break
				}
				unique[id] = true
			}
		}
	}

	return solutions
}
