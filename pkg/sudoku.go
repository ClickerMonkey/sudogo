package sudogo

import (
	"math"
	"math/rand"
)

type Size struct {
	Width  int
	Height int
}

func (size Size) Area() int {
	return size.Width * size.Height
}

type Bits struct {
	Value uint64
	Count int
}

func (bits *Bits) Fill(n int) {
	bits.Count = n
	bits.Value = (1 << n) - 1
}

func (bits *Bits) Clear() {
	bits.Fill(0)
}

func (bits *Bits) Clone() Bits {
	return Bits{bits.Value, bits.Count}
}

func (bits *Bits) On(i int) bool {
	return (bits.Value & (1 << i)) != 0
}

func (bits *Bits) Set(i int, on bool) bool {
	change := bits.On(i) != on
	if change {
		if on {
			bits.Count++
			bits.Value = bits.Value | (1 << i)
		} else {
			bits.Count--
			bits.Value = bits.Value & ^(1 << i)
		}
	}
	return change
}

func (bits *Bits) ToSlice() []int {
	slice := make([]int, 0, bits.Count)
	remaining := bits.Value
	current := 0
	for remaining > 0 {
		if (remaining & 1) == 1 {
			slice = append(slice, current)
		}
		remaining = remaining >> 1
		current++
	}
	return slice
}

func (bits *Bits) First() int {
	for i := 0; i < 64; i++ {
		if (bits.Value & (1 << i)) != 0 {
			return i
		}
	}
	return -1
}

func (bits *Bits) UpdateCount() {
	bits.Count = 0
	for i := 0; i < 64; i++ {
		if (bits.Value & (1 << i)) != 0 {
			bits.Count++
		}
	}
}

func (bits *Bits) Remove(remove uint64) int {
	original := bits.Count
	bits.Value = bits.Value & ^remove
	bits.UpdateCount()
	return original - bits.Count
}

func (bits *Bits) Or(or uint64) {
	bits.Value = bits.Value | or
	bits.UpdateCount()
}

func (bits *Bits) And(and uint64) {
	bits.Value = bits.Value & and
	bits.UpdateCount()
}

type Constraint func(puzzle *PuzzleInstance) bool

type Puzzle struct {
	Boxs        Size
	BoxSize     Size
	Constraints []Constraint
}

func (puzzle *Puzzle) Width() int {
	return puzzle.BoxSize.Width * puzzle.Boxs.Width
}

func (puzzle *Puzzle) Height() int {
	return puzzle.BoxSize.Height * puzzle.Boxs.Height
}

func (puzzle *Puzzle) Digits() int {
	return puzzle.BoxSize.Area()
}

func (puzzle *Puzzle) DigitsSize() int {
	return int(math.Floor(math.Log10(float64(puzzle.Digits())))) + 1
}

func (puzzle *Puzzle) DefaultCandidates() uint64 {
	return uint64((1 << (puzzle.Digits() + 1)) - 1)
}

func (puzzle *Puzzle) Create() PuzzleInstance {
	return NewInstance(puzzle)
}

var Classic = &Puzzle{
	Boxs:    Size{3, 3},
	BoxSize: Size{3, 3},
}

type Cell struct {
	value      int
	id         int
	row        int
	col        int
	box        int
	candidates Bits
}

func (cell *Cell) RemoveCandidate(value int) bool {
	return cell.candidates.Set(value, false)
}

func (cell *Cell) HasValue() bool {
	return cell.value != 0
}

func (cell *Cell) HasCandidate(value int) bool {
	return cell.candidates.On(value)
}

func (cell *Cell) SetValue(value int) bool {
	can := cell.candidates.On(value)
	if can {
		cell.value = value
		cell.candidates.Clear()
	}
	return can
}

func (cell *Cell) InGroup(other *Cell) bool {
	return cell.id != other.id && (cell.row == other.row || cell.col == other.col || cell.box == other.box)
}

func (cell *Cell) InBox(other *Cell) bool {
	return cell.id != other.id && cell.box == other.box
}

func (cell *Cell) InRow(other *Cell) bool {
	return cell.id != other.id && cell.row == other.row
}

func (cell *Cell) InColumn(other *Cell) bool {
	return cell.id != other.id && cell.col == other.col
}

func (cell *Cell) Candidates() []int {
	return cell.candidates.ToSlice()
}

func (cell *Cell) FirstCandidate() int {
	return cell.candidates.First()
}

type PuzzleInstance struct {
	puzzle *Puzzle
	cells  []Cell
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

	return PuzzleInstance{puzzle, cells}
}

func (instance *PuzzleInstance) Clone() PuzzleInstance {
	puzzle := instance.puzzle
	cells := make([]Cell, len(instance.cells))

	copy(cells, instance.cells)

	return PuzzleInstance{puzzle, cells}
}

func (instance *PuzzleInstance) Get(col int, row int) *Cell {
	return &instance.cells[row*instance.puzzle.Width()+col]
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

func (instance *PuzzleInstance) Set(col int, row int, value int) bool {
	return instance.SetCell(instance.Get(col, row), value)
}

func (instance *PuzzleInstance) SetCell(cell *Cell, value int) bool {
	set := cell.SetValue(value)

	if set {
		instance.ForInGroup(cell, func(other *Cell) bool {
			if !other.HasValue() {
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

func (instance *PuzzleInstance) GetNakedSingle() (*Cell, int) {
	for i := range instance.cells {
		cell := &instance.cells[i]
		if cell.candidates.Count == 1 {
			return cell, cell.FirstCandidate()
		}
	}
	return nil, 0
}

func (instance *PuzzleInstance) SolveNakedSingles(max int) int {
	return instance.DoSets(max, func() (*Cell, int) {
		return instance.GetNakedSingle()
	})
}

func (instance *PuzzleInstance) GetHiddenSingle() (*Cell, int) {
	// A cell which has a candidate that is unique to the row, cell, or box
	for i := range instance.cells {
		cell := &instance.cells[i]
		if !cell.HasValue() {
			box := instance.GetHiddenSingleBox(cell)
			if box != 0 {
				return cell, box
			}
			row := instance.GetHiddenSingleRow(cell)
			if row != 0 {
				return cell, row
			}
			col := instance.GetHiddenSingleColumn(cell)
			if col != 0 {
				return cell, col
			}
		}
	}
	return nil, 0
}

func (instance *PuzzleInstance) GetHiddenSingleBox(cell *Cell) int {
	on := cell.candidates
	instance.ForInBox(cell, func(other *Cell) bool {
		if !other.HasValue() {
			on.Remove(other.candidates.Value)
			if on.Count == 0 {
				return false
			}
		}
		return true
	})
	if on.Count == 1 {
		return on.First()
	}
	return 0
}

func (instance *PuzzleInstance) GetHiddenSingleRow(cell *Cell) int {
	on := cell.candidates
	instance.ForInRow(cell, func(other *Cell) bool {
		if !other.HasValue() {
			on.Remove(other.candidates.Value)
			if on.Count == 0 {
				return false
			}
		}
		return true
	})
	if on.Count == 1 {
		return on.First()
	}
	return 0
}

func (instance *PuzzleInstance) GetHiddenSingleColumn(cell *Cell) int {
	on := cell.candidates
	instance.ForInColumn(cell, func(other *Cell) bool {
		if !other.HasValue() {
			on.Remove(other.candidates.Value)
			if on.Count == 0 {
				return false
			}
		}
		return true
	})
	if on.Count == 1 {
		return on.First()
	}
	return 0
}

func (instance *PuzzleInstance) SolveHiddenSingles(max int) int {
	return instance.DoSets(max, func() (*Cell, int) {
		return instance.GetHiddenSingle()
	})
}

func (instance *PuzzleInstance) RemovePointingCandidates(max int) int {
	// If in a box all candidates of a certain digit are confined to a row or column, that digit cannot appear outside of that box in that row or column.
	removes := 0

	for i := range instance.cells {
		cell := &instance.cells[i]
		if cell.HasValue() {
			continue
		}

		// all candidates in this box's row that are shared
		row := cell.candidates.Clone()
		// all candidates in this box's column that are shared
		col := cell.candidates.Clone()

		// remove candidates that are not shared
		instance.ForInBox(cell, func(other *Cell) bool {
			if !other.HasValue() {
				if other.row == cell.row {
					row.And(other.candidates.Value)
				}
				if other.col == cell.col {
					col.And(other.candidates.Value)
				}
			}

			return row.Count > 0 || col.Count > 0
		})

		// remove candidates that exist outside the row or column
		instance.ForInBox(cell, func(other *Cell) bool {
			if !other.HasValue() {
				if other.row != cell.row && other.col != cell.col {
					row.Remove(other.candidates.Value)
					col.Remove(other.candidates.Value)
				}
			}

			return row.Count > 0 || col.Count > 0
		})

		// what is remaining are candidates confined to the cells row in the box
		if row.Count > 0 {
			instance.ForInRow(cell, func(other *Cell) bool {
				if !other.HasValue() && other.box != cell.box {
					removes += other.candidates.Remove(row.Value)
				}
				return true
			})
		}

		// what is remaining are candidates confined to the cells column in the box
		if col.Count > 0 {
			instance.ForInColumn(cell, func(other *Cell) bool {
				if !other.HasValue() && other.box != cell.box {
					removes += other.candidates.Remove(col.Value)
				}
				return true
			})
		}
	}

	return removes
}

func (instance *PuzzleInstance) RemoveClaimingCandidates(max int) int {
	// If in a row or column a candidate only appears in a single box then that candidate can be removed from other cells in that box
	removed := 0

	for i := range instance.cells {
		cell := &instance.cells[i]
		if cell.HasValue() {
			continue
		}

		// all candidates in this row that are not shared outside of the box
		row := cell.candidates.Clone()
		// all candidates in this column that are not shared outside of the box
		col := cell.candidates.Clone()

		// remove candidates from row that exist in the cells row outside the box
		instance.ForInRow(cell, func(other *Cell) bool {
			if !other.HasValue() && other.box != cell.box {
				row.Remove(other.candidates.Value)
			}

			return row.Count > 0
		})

		// remove candidates from column that exist in the cells column outside the box
		instance.ForInColumn(cell, func(other *Cell) bool {
			if !other.HasValue() && other.box != cell.box {
				col.Remove(other.candidates.Value)
			}

			return col.Count > 0
		})

		// what is remaining are the candidates unique to the row/column outside this box
		if row.Count > 0 || col.Count > 0 {
			instance.ForInBox(cell, func(other *Cell) bool {
				if !other.HasValue() {
					if row.Count > 0 && other.row != cell.row {
						removed += other.candidates.Remove(row.Value)
					}
					if col.Count > 0 && other.col != cell.col {
						removed += other.candidates.Remove(col.Value)
					}
				}
				return true
			})
		}
	}

	return removed
}

func (instance *PuzzleInstance) DoSets(max int, nextSet func() (*Cell, int)) int {
	set := 0
	cell, cellValue := nextSet()
	for cell != nil {
		set++
		instance.SetCell(cell, cellValue)
		if max > 0 && set == max {
			break
		}
		cell, cellValue = nextSet()
	}
	return set
}

func (instance *PuzzleInstance) Solved() bool {
	for i := range instance.cells {
		if !instance.cells[i].HasValue() {
			return false
		}
	}
	return true
}

func (instance *PuzzleInstance) Solve() int {
	placed := 0
	for {
		placed += instance.SolveNakedSingles(-1)

		if instance.SolveHiddenSingles(1) > 0 {
			placed++
			continue
		}

		pointing := instance.RemovePointingCandidates(-1)
		claiming := instance.RemoveClaimingCandidates(-1)

		if pointing > 0 || claiming > 0 {
			continue
		}

		break
	}

	return placed
}

func (instance *PuzzleInstance) GetUnsolved() *Cell {
	for i := range instance.cells {
		cell := &instance.cells[i]
		if !cell.HasValue() {
			return cell
		}
	}
	return nil
}

func (instance *PuzzleInstance) GetRandom(random *rand.Rand, match func(other *Cell) bool) *Cell {
	matches := int32(0)
	for i := range instance.cells {
		cell := &instance.cells[i]
		if match(cell) {
			matches++
		}
	}
	if matches == 0 {
		return nil
	}
	chosen := random.Int31n(matches)
	for i := range instance.cells {
		cell := &instance.cells[i]
		if match(cell) {
			chosen--
			if chosen < 0 {
				return cell
			}
		}
	}
	return nil
}

func (instance *PuzzleInstance) GetRandomUnsolved(random *rand.Rand) *Cell {
	return instance.GetRandom(random, func(other *Cell) bool {
		return !other.HasValue()
	})
}

func (instance *PuzzleInstance) GetRandomPressured(random *rand.Rand) *Cell {
	minCount := 0
	for i := range instance.cells {
		cell := &instance.cells[i]
		if !cell.HasValue() && (minCount == 0 || minCount > cell.candidates.Count) {
			minCount = cell.candidates.Count
		}
	}

	return instance.GetRandom(random, func(other *Cell) bool {
		return !other.HasValue() && other.candidates.Count == minCount
	})
}

func (instance *PuzzleInstance) Generate(random *rand.Rand) bool {
	for !instance.Solved() {
		instance.Solve()

		if instance.Solved() {
			break
		}

		rnd := instance.GetRandomUnsolved(random)

		if rnd == nil {
			return false
		}

		available := rnd.Candidates()
		max := int32(len(available))
		if max == 0 {
			return false
		}

		chosen := random.Int31n(max)

		instance.SetCell(rnd, available[chosen])
	}
	return true
}
