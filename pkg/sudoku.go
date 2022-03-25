package sudogo

type PuzzleInstance struct {
	puzzle   *Puzzle
	cells    []Cell
	unsolved []*Cell
}

func NewInstance(puzzle *Puzzle) PuzzleInstance {
	boxsWide := puzzle.Boxs.Width
	boxsHigh := puzzle.Boxs.Height
	boxWidth := puzzle.BoxSize.Width
	boxHeight := puzzle.BoxSize.Height
	cellsWide := boxsWide * boxWidth
	cellsHigh := boxsHigh * boxHeight
	cellCount := cellsWide * cellsHigh
	candidates := puzzle.Digits() + 1
	cells := make([]Cell, cellCount)
	unsolved := make([]*Cell, cellCount)

	for i := 0; i < cellCount; i++ {
		cell := &cells[i]
		cell.id = i
		cell.value = 0
		cell.row = i / cellsWide
		cell.col = i % cellsWide
		cell.box = ((cell.row / boxsHigh) * boxHeight) + (cell.col / boxsWide)
		cell.candidates.Fill(candidates)
		cell.candidates.Set(0, false)
		unsolved[i] = cell
	}

	return PuzzleInstance{puzzle, cells, unsolved}
}

func (instance *PuzzleInstance) Clone() PuzzleInstance {
	puzzle := instance.puzzle
	n := len(instance.cells)
	cells := make([]Cell, n)
	unsolved := make([]*Cell, n)

	copy(cells, instance.cells)

	for i := 0; i < n; i++ {
		unsolved[i] = &cells[i]
	}

	return PuzzleInstance{puzzle, cells, unsolved}
}

func (instance *PuzzleInstance) Get(col int, row int) *Cell {
	return &instance.cells[row*instance.puzzle.Width()+col]
}

func (instance *PuzzleInstance) Set(col int, row int, value int) bool {
	return instance.SetCell(instance.Get(col, row), value)
}

func (instance *PuzzleInstance) SetCell(cell *Cell, value int) bool {
	set := cell.SetValue(value)

	if set {
		if cell.HasValue() {
			for i := range instance.unsolved {
				if instance.unsolved[i].id == cell.id {
					instance.unsolved = removeAtIndex(instance.unsolved, i)
					break
				}
			}
		}

		instance.ForInGroup(cell, func(other *Cell) bool {
			if other.Empty() {
				other.RemoveCandidate(value)
			}

			return true
		})
	}

	return set
}

func (instance *PuzzleInstance) SetAll(values [][]int) int {
	width := instance.puzzle.Width()
	height := instance.puzzle.Height()
	sets := 0

	for y, row := range values {
		if y >= height {
			break
		}
		for x, value := range row {
			if x >= width {
				break
			}
			if instance.Set(x, y, value) {
				sets++
			}
		}
	}

	return sets
}

func (instance *PuzzleInstance) ForInGroup(cell *Cell, call func(other *Cell) bool) bool {
	for i := range instance.cells {
		c := &instance.cells[i]
		if cell.InGroup(c) {
			if !call(c) {
				return false
			}
		}
	}
	return true
}

func (instance *PuzzleInstance) ForInBox(cell *Cell, call func(other *Cell) bool) bool {
	for i := range instance.cells {
		c := &instance.cells[i]
		if cell.InBox(c) {
			if !call(c) {
				return false
			}
		}
	}
	return true
}

func (instance *PuzzleInstance) ForInRow(cell *Cell, call func(other *Cell) bool) bool {
	width := instance.puzzle.Width()
	start := cell.row * width

	for i := 0; i < width; i++ {
		if cell.col != i && !call(&instance.cells[i+start]) {
			return false
		}
	}
	return true
}

func (instance *PuzzleInstance) ForInColumn(cell *Cell, call func(other *Cell) bool) bool {
	width := instance.puzzle.Width()
	height := instance.puzzle.Height()

	for i := 0; i < height; i++ {
		if cell.row != i && !call(&instance.cells[i*width+cell.col]) {
			return false
		}
	}
	return true
}
