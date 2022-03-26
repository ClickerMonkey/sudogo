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
	return &puzzle.cells[row*puzzle.kind.Width()+col]
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
	width := puzzle.kind.Width()
	height := puzzle.kind.Height()
	sets := 0

	puzzle.Clear()

	for y, row := range values {
		if y >= height {
			break
		}
		for x, value := range row {
			if x >= width {
				break
			}
			if value > 0 && puzzle.Set(x, y, value) {
				sets++
			}
		}
	}

	return sets
}
