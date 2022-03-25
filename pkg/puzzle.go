package sudogo

type Puzzle struct {
	kind  *Kind
	cells []Cell
}

func New(kind *Kind) Puzzle {
	boxsWide := kind.Boxs.Width
	boxsHigh := kind.Boxs.Height
	boxWidth := kind.BoxSize.Width
	boxHeight := kind.BoxSize.Height
	cellsWide := boxsWide * boxWidth
	cellsHigh := boxsHigh * boxHeight
	cellCount := cellsWide * cellsHigh
	candidates := kind.Digits()
	cells := make([]Cell, cellCount)

	for i := 0; i < cellCount; i++ {
		cell := &cells[i]
		cell.id = i
		cell.value = 0
		cell.row = i / cellsWide
		cell.col = i % cellsWide
		cell.box = ((cell.row / boxsHigh) * boxHeight) + (cell.col / boxsWide)
		cell.candidates.Fill(candidates)
		cell.candidates.Set(0, false)
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

func (puzzle *Puzzle) Solver() Solver {
	return NewSolver(*puzzle)
}

func (puzzle *Puzzle) Get(col int, row int) *Cell {
	return &puzzle.cells[row*puzzle.kind.Width()+col]
}

func (puzzle *Puzzle) Set(col int, row int, value int) bool {
	return puzzle.SetCell(puzzle.Get(col, row), value)
}

func (puzzle *Puzzle) SetCell(cell *Cell, value int) bool {
	set := cell.SetValue(value)

	if set {
		puzzle.ForInGroup(cell, func(other *Cell) bool {
			if other.Empty() {
				other.RemoveCandidate(value)
			}

			return true
		})
	}

	return set
}

func (puzzle *Puzzle) SetAll(values [][]int) int {
	width := puzzle.kind.Width()
	height := puzzle.kind.Height()
	sets := 0

	for y, row := range values {
		if y >= height {
			break
		}
		for x, value := range row {
			if x >= width {
				break
			}
			if puzzle.Set(x, y, value) {
				sets++
			}
		}
	}

	return sets
}

func (puzzle *Puzzle) ForInGroup(cell *Cell, call func(other *Cell) bool) bool {
	for i := range puzzle.cells {
		c := &puzzle.cells[i]
		if cell.InGroup(c) {
			if !call(c) {
				return false
			}
		}
	}
	return true
}

func (puzzle *Puzzle) ForInBox(cell *Cell, call func(other *Cell) bool) bool {
	for i := range puzzle.cells {
		c := &puzzle.cells[i]
		if cell.InBox(c) {
			if !call(c) {
				return false
			}
		}
	}
	return true
}

func (puzzle *Puzzle) ForInRow(cell *Cell, call func(other *Cell) bool) bool {
	width := puzzle.kind.Width()
	start := cell.row * width

	for i := 0; i < width; i++ {
		if cell.col != i && !call(&puzzle.cells[i+start]) {
			return false
		}
	}
	return true
}

func (puzzle *Puzzle) ForInColumn(cell *Cell, call func(other *Cell) bool) bool {
	width := puzzle.kind.Width()
	height := puzzle.kind.Height()

	for i := 0; i < height; i++ {
		if cell.row != i && !call(&puzzle.cells[i*width+cell.col]) {
			return false
		}
	}
	return true
}
