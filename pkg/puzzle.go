package sudogo

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Puzzle struct {
	Kind  *Kind
	Cells []Cell
}

func New(kind *Kind) Puzzle {
	boxsWide, _, boxWidth, boxHeight, size := kind.GetDimensions()
	cellCount := size * size
	cells := make([]Cell, cellCount)

	for i := 0; i < cellCount; i++ {
		cell := &cells[i]
		cell.Id = i
		cell.Value = 0
		cell.Row = i / size
		cell.Col = i % size
		cell.Box = ((cell.Row / boxHeight) * boxsWide) + (cell.Col / boxWidth)
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

func (puzzle *Puzzle) Contains(col int, row int) bool {
	size := puzzle.Kind.Size()
	return col >= 0 && col < size && row >= 0 && row < size
}

func (puzzle *Puzzle) IsCandidate(value int) bool {
	return value >= puzzle.MinCandidate() && value <= puzzle.MaxCandidate()
}

func (puzzle *Puzzle) MinCandidate() int {
	return 1
}
func (puzzle *Puzzle) MaxCandidate() int {
	return puzzle.Kind.Size()
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

func (puzzle *Puzzle) String() string {
	return puzzle.ToStateString(false, ".")
}

func (puzzle *Puzzle) ToStateString(includeKind bool, emptyValue string) string {
	sb := strings.Builder{}
	if includeKind {
		sb.WriteString(strconv.Itoa(puzzle.Kind.BoxSize.Width))
		sb.WriteString("x")
		sb.WriteString(strconv.Itoa(puzzle.Kind.BoxSize.Height))
		sb.WriteString(",")
	}
	digitSize := puzzle.Kind.DigitsSize()
	digitFormat := fmt.Sprintf("%%0%dd", digitSize)
	emptyCell := strings.Repeat(emptyValue, digitSize)
	for k := range puzzle.Cells {
		other := &puzzle.Cells[k]
		if other.Empty() {
			sb.WriteString(emptyCell)
		} else {
			sb.WriteString(fmt.Sprintf(digitFormat, other.Value))
		}
	}
	return sb.String()
}

func FromString(input string) *Puzzle {
	hasKind := strings.Index(input, ",") != -1
	cells := input
	boxWidth := 3
	boxHeight := 3
	if hasKind {
		parts := strings.Split(input, ",")
		size := strings.Split(parts[0], "x")
		var err error
		if boxWidth, err = strconv.Atoi(size[0]); err != nil {
			return nil
		}
		if boxHeight, err = strconv.Atoi(size[1]); err != nil {
			return nil
		}
		cells = parts[1]
	}
	n := len(cells)
	if n != boxWidth*boxWidth*boxHeight*boxHeight {
		size := math.Round(math.Sqrt(float64(n)))
		boxSize := int(math.Ceil(math.Sqrt(size)))
		boxWidth = boxSize
		boxHeight = int(size) / boxSize
	}
	kind := NewKind(boxWidth, boxHeight)
	puzzle := New(kind)
	digits := stringChunk(cells, kind.DigitsSize())
	for i := range puzzle.Cells {
		if i < len(digits) {
			if cellValue, err := strconv.Atoi(digits[i]); err == nil && cellValue > 0 {
				puzzle.SetCell(&puzzle.Cells[i], cellValue)
			}
		}
	}
	return &puzzle
}

func (puzzle *Puzzle) EncodedString() string {
	i := big.NewInt(0)
	scale := big.NewInt(int64(puzzle.Kind.Digits() + 1))
	boxSizeScale := big.NewInt(32)
	for k := range puzzle.Cells {
		cell := &puzzle.Cells[k]
		i.Mul(i, scale)
		i.Add(i, big.NewInt(int64(cell.Value)))
	}
	i.Mul(i, boxSizeScale)
	i.Add(i, big.NewInt(int64(puzzle.Kind.BoxSize.Height)))
	i.Mul(i, boxSizeScale)
	i.Add(i, big.NewInt(int64(puzzle.Kind.BoxSize.Width)))
	return base64.StdEncoding.EncodeToString(i.Bytes())
}

func FromEncoded(input string) *Puzzle {
	bytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil
	}
	i := big.NewInt(0).SetBytes(bytes)
	boxSizeScale := big.NewInt(32)
	boxWidth := int(big.NewInt(0).Mod(i, boxSizeScale).Int64())
	i.Div(i, boxSizeScale)
	boxHeight := int(big.NewInt(0).Mod(i, boxSizeScale).Int64())
	i.Div(i, boxSizeScale)
	kind := NewKind(boxWidth, boxHeight)
	scale := big.NewInt(int64(kind.Digits() + 1))
	puzzle := New(kind)
	cellCount := len(puzzle.Cells) - 1
	for k := range puzzle.Cells {
		cellValue := int(big.NewInt(0).Mod(i, scale).Int64())
		if cellValue > 0 {
			puzzle.SetCell(&puzzle.Cells[cellCount-k], cellValue)
		}
		i.Div(i, scale)
	}
	return &puzzle
}

func (puzzle *Puzzle) UniqueId() string {
	sb := strings.Builder{}
	for k := range puzzle.Cells {
		other := &puzzle.Cells[k]
		sb.WriteString(strconv.Itoa(other.Value))
	}
	return sb.String()
}

func (puzzle *Puzzle) HasUniqueSolution() bool {
	return len(puzzle.GetSolutions(SolutionsLimit{MaxSolutions: 2})) == 1
}

type SolutionsLimit struct {
	SolveLimit
	MaxSolutions int
	LogEnabled   bool
	LogState     bool
}

func (puzzle *Puzzle) GetSolutions(limit SolutionsLimit) []*Solver {
	solutions := make([]*Solver, 0, Max(0, limit.MaxSolutions))
	unique := map[string]bool{}

	solvers := NewQueue[Solver]()
	solvers.Offer(puzzle.Solver())

	for !solvers.Empty() {
		solver := solvers.Poll()
		solution, solved := solver.Solve(limit.SolveLimit)

		if !solved {
			min := solver.GetMinCandidateCount()
			minCell := solver.GetGroupWhere(func(group *CellGroups) bool {
				return group.Cell.candidates.Count == min
			})
			if minCell != nil {
				for _, candidate := range minCell.Cell.Candidates() {
					newSolver := solution.Solver()
					newSolver.LogEnabled = limit.LogEnabled
					newSolver.LogState = limit.LogState
					if newSolver.Set(minCell.Cell.Col, minCell.Cell.Row, candidate) {
						solvers.Offer(newSolver)
					}
				}
			}
		} else {
			id := solution.UniqueId()
			if !unique[id] {
				solutions = append(solutions, solver)
				if limit.MaxSolutions > 0 && len(solutions) == limit.MaxSolutions {
					break
				}
				unique[id] = true
			}
		}
	}

	return solutions
}
