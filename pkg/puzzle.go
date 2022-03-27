package sudogo

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

func (puzzle *Puzzle) IsSolved() bool {
	size := puzzle.kind.Size()
	rows := make([]Candidates, size)
	cols := make([]Candidates, size)
	boxs := make([]Candidates, size)

	complete := Candidates{}
	complete.Fill(size)

	for i := range puzzle.cells {
		cell := &puzzle.cells[i]
		if cell.value == 0 {
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
