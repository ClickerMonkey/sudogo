package sudogo

// A constraint is an added rule for solving the puzzle. This enables more complex puzzles with fewer givens.
// The following basic constraints are supported
// - A collection of cells sum up to a value
//		- Diagonal lines: https://www.youtube.com/watch?v=Vc-FYo_nur4
//		- Groups/cages: https://www.youtube.com/watch?v=2v6Lf3Q5AEo&t=885s
//		- X/V: https://www.youtube.com/watch?v=9ATC_uBF8ow
//		- Squares: https://www.youtube.com/watch?v=u6Le6f9d0KU&t=602s
//		- Knights move: https://www.youtube.com/watch?v=3FMNh-_FNlk
//		- Magic square: https://www.youtube.com/watch?v=La7Yg_rav24
// - A collection of cells sum up to the value of another cell
//		- Path: https://www.youtube.com/watch?v=Vc-FYo_nur4
// - A collection of cells go from increasing to decreasing order (known direction or not)
//		- Thermo(meter): https://www.youtube.com/watch?v=KTth49YrQVU
//		- Circled ends: https://www.youtube.com/watch?v=Tpk3ga2T9Ps&t=159s
// - A collection of cells even & odd digits sum to same value
// - A cell cannot have the same value as a collection of cells
//		- Kings move
//		- Knights move: https://www.youtube.com/watch?v=hAyZ9K2EBF0
// 		- No repeats in group/age: https://www.youtube.com/watch?v=u6Le6f9d0KU&t=602s, https://www.youtube.com/watch?v=hAyZ9K2EBF0
// - Or constraint (multiple constraints can dictate which candidates are available)
// 		- Sum is square: https://www.youtube.com/watch?v=u6Le6f9d0KU&t=602s
// - And constraint
//		- Circled ends & no duplicates: https://www.youtube.com/watch?v=Tpk3ga2T9Ps&t=159s
type Constraint interface {
	Affects(cell *Cell) bool
	RemoveCandidates(cell *Cell, puzzle *Puzzle, remove *Candidates)
}

// A cell position on a puzzle.
type Position struct {
	Col int
	Row int
}

// ==================================================
// Constraint: Sum Value
// ==================================================

// A function which returns a value other cells should sum to given
// the cell being constrained and the puzzle. If there is no sum
// to constrain then 0 should be returned.
type ConstraintSumProvider func(cell *Cell, puzzle *Puzzle) int

// A constraint on a set of cells that states that set or relative cells
// should sum to a value.
type ConstraintSum struct {
	// The function which returns a value the cells need to sum to.
	Sum ConstraintSumProvider
	// The list of cells that are constrained to some sum. If this is nil then
	// all cells will be a part of the constraint.
	Cells    *[]Position
	Relative *[]Position
}

// A sum provider which returns a constant value
func SumConstant(value int) ConstraintSumProvider {
	return func(cell *Cell, puzzle *Puzzle) int {
		return value
	}
}

// A sum provider which returns the value (or largest candidate) of a cell
// at the given position
func SumCell(pos Position, relative bool) ConstraintSumProvider {
	return func(cell *Cell, puzzle *Puzzle) int {
		col := pos.Col
		row := pos.Row
		if relative {
			col += cell.Col
			row += cell.Row
		}
		value := 0
		if puzzle.Contains(col, row) {
			other := puzzle.Get(col, row)
			value = other.MaxValue()
		}
		return value
	}
}

// A sum provider which returns the value (or largest candidate) of a cell
// at the given position
func SumCells(positions []Position, relative bool) ConstraintSumProvider {
	return func(cell *Cell, puzzle *Puzzle) int {
		value := 0
		for _, pos := range positions {
			col := pos.Col
			row := pos.Row
			if relative {
				col += cell.Col
				row += cell.Row
			}
			if puzzle.Contains(col, row) {
				other := puzzle.Get(col, row)
				value += other.MaxValue()
			}
		}
		return value
	}
}

func (c *ConstraintSum) Affects(cell *Cell) bool {
	return containsCell(cell, c.Cells, nil)
}

func (c *ConstraintSum) RemoveCandidates(cell *Cell, puzzle *Puzzle, remove *Candidates) {
	target := c.Sum(cell, puzzle)
	if target == 0 {
		return
	}
	cells := getCells(puzzle, cell, c.Cells, c.Relative)

	sum := 0
	combos := combinations[int]{}

	for _, other := range cells {
		if other.HasValue() {
			sum += other.Value
		} else if other.Id != cell.Id {
			combos.add(other.Candidates())
		}
	}

	sumEmpty := target - sum

	if combos.empty() {
		chosen := sumEmpty
		if remove.Has(chosen) {
			remove.Clear()
			remove.Set(chosen, true)
		}
		return
	}

	candidates := remove.ToSlice()
	values := combos.start()
	for _, candidate := range candidates {
		comboSum := sumEmpty - candidate
		foundSum := false

		for combos.next(values) {
			if intsUnique(values) && intsSum(values) == comboSum {
				foundSum = true
				break
			}
		}

		combos.reset()

		if !foundSum {
			remove.Set(candidate, false)
		}
	}
}

// ==================================================
// Constraint: Uniqueness
// ==================================================

// A constraint on a set of cells where they can't have the same value OR
// they all need to have the same value.
type ConstraintUnique struct {
	Cells    *[]Position
	Relative *[]Position
	Same     bool
}

func (c *ConstraintUnique) Affects(cell *Cell) bool {
	return containsCell(cell, c.Cells, nil)
}

func (c *ConstraintUnique) RemoveCandidates(cell *Cell, puzzle *Puzzle, remove *Candidates) {
	traverseCells(puzzle, cell, c.Cells, c.Relative, func(other *Cell, index int) {
		if other.HasValue() {
			if c.Same {
				remove.Clear()
				remove.Set(other.Value, true)
			} else {
				remove.Set(other.Value, false)
			}
		}
	})
}

// ==================================================
// Constraint: Uniqueness
// ==================================================
type ConstraintOrder struct {
	Cells     *[]Position
	Relative  *[]Position
	Direction int
}

func (c *ConstraintOrder) Affects(cell *Cell) bool {
	return containsCell(cell, c.Cells, nil)
}

func (c *ConstraintOrder) RemoveCandidates(cell *Cell, puzzle *Puzzle, remove *Candidates) {
	cells := getCells(puzzle, cell, c.Cells, c.Relative)
	i := -1
	for k, other := range cells {
		if other.Id == cell.Id {
			i = k
			break
		}
	}
	if i == -1 {
		return
	}

	var firstValue *Cell
	dir := c.Direction
	if dir == 0 {
		for _, other := range cells {
			if other.HasValue() {
				if firstValue == nil {
					firstValue = other
				} else if firstValue.Value < other.Value {
					dir = 1
					break
				} else {
					dir = -1
					break
				}
			}
		}
	}

	if dir != 0 {
		puzzleMin := 1
		puzzleMax := puzzle.Kind.Digits()

		min := puzzleMin
		max := puzzleMax

		for k, other := range cells {
			if k == i {
				continue
			}
			d := (i - k) * dir
			if d > 0 {
				otherMin := other.MinValue() + d
				min = Max(min, otherMin)
			}
			if d < 0 {
				otherMax := other.MaxValue() + d
				max = Min(max, otherMax)
			}
		}

		for i := puzzleMin; i < min; i++ {
			remove.Set(i, false)
		}
		for i := max + 1; i <= puzzleMax; i++ {
			remove.Set(i, false)
		}
	} else if firstValue != nil {
		remove.Set(firstValue.Value, false)
	}
}

// ==================================================
// Constraint: Magic Square
// ==================================================
type ConstraintMagic struct {
	Center Position
}

func (c *ConstraintMagic) Affects(cell *Cell) bool {
	dx := cell.Col - c.Center.Col
	dy := cell.Row - c.Center.Row
	return dx <= 1 && dx >= -1 && dy <= 1 && dy >= -1
}

func (c *ConstraintMagic) RemoveCandidates(cell *Cell, puzzle *Puzzle, remove *Candidates) {
	dx := cell.Col - c.Center.Col
	dy := cell.Row - c.Center.Row
	if dx == 0 && dy == 0 {
		remove.Clear()
		remove.Set(5, true)
	} else if dx*dy == 0 {
		remove.Set(2, false)
		remove.Set(4, false)
		remove.Set(6, false)
		remove.Set(8, false)
	} else {
		remove.Set(1, false)
		remove.Set(3, false)
		remove.Set(7, false)
		remove.Set(9, false)
	}
	if remove.Count > 1 {
		verSum := 0
		verCount := 0
		for ry := 0; ry < 3; ry++ {
			other := puzzle.Get(cell.Col, cell.Row-(dy+1)+ry)
			if other.HasValue() {
				verCount++
				verSum += other.Value
			}
		}

		if verCount == 2 {
			remove.Clear()
			remove.Set(15-verSum, true)
			return
		}

		horSum := 0
		horCount := 0
		for rx := 0; rx < 3; rx++ {
			other := puzzle.Get(cell.Col-(dx+1)+rx, cell.Row)
			if other.HasValue() {
				horCount++
				horSum += other.Value
			}
		}

		if horCount == 2 {
			remove.Clear()
			remove.Set(15-horSum, true)
			return
		}
	}
}

// ==================================================
// Constraint: Scale
// ==================================================
type ConstraintScalePair struct {
	Scale  int
	First  Position
	Second Position
}

func (c *ConstraintScalePair) Affects(cell *Cell) bool {
	return isSame(cell, c.First) || isSame(cell, c.Second)
}

func (c *ConstraintScalePair) RemoveCandidates(cell *Cell, puzzle *Puzzle, remove *Candidates) {
	var other *Cell
	if isSame(cell, c.First) {
		other = getAbsoluteCell(puzzle, c.Second)
	} else if isSame(cell, c.Second) {
		other = getAbsoluteCell(puzzle, c.First)
	}
	if other == nil {
		return
	}
	possible := other.candidates
	if other.HasValue() {
		possible.Clear()
		possible.Set(other.Value, true)
	}

	available := *remove
	for available.Count > 0 {
		candidate := available.First()
		available.Set(candidate, false)

		up := candidate * c.Scale
		upValid := puzzle.IsCandidate(up) && possible.Has(up)

		down := candidate / c.Scale
		downValid := down*c.Scale == candidate && puzzle.IsCandidate(down) && possible.Has(down)

		if !upValid && !downValid {
			remove.Set(candidate, false)
		}
	}
}

func ConstraintScalePairs(scale int, pairs [][2]Position) []ConstraintScalePair {
	constraints := make([]ConstraintScalePair, len(pairs))
	for pairIndex, pair := range pairs {
		constraints[pairIndex].Scale = scale
		constraints[pairIndex].First = pair[0]
		constraints[pairIndex].Second = pair[1]
	}
	return constraints
}

// ==================================================
// Constraint: Difference
// ==================================================
type ConstraintDifference struct {
	// The minimum difference that should exist between the given cells.
	// If the cells are in the same groups then the minimum is already technically 1
	// since the same value can't exist in the same group. A value of 2 means all
	// cells will need to be atleast 2 apart.
	Min int
	// The maximum difference that should exist between the given cells.
	// For example if the Max is 4 and one of the cells is 2 then the other cells
	// are constrained to 1, 2, 3, 4, 5, and 6.
	Max int
	// The cells which are affected by this constaint. If nil all cells in the puzzle
	// are affected (minus what's given in Exclude).
	Cells *[]Position
	// The cells which are looked at by the constraint. If nil the cells involved in
	// the logic are the Cells given.
	Relative *[]Position
	// The cells to exclude from being constrained when Cells is nil
	// (meaning all cells are constrained int he puzzle).
	Exclude *[]Position
}

func (c *ConstraintDifference) Affects(cell *Cell) bool {
	return containsCell(cell, c.Cells, c.Exclude)
}

func (c *ConstraintDifference) RemoveCandidates(cell *Cell, puzzle *Puzzle, remove *Candidates) {
	surrounding := Candidates{}
	surroundingValues := Candidates{}

	traverseCells(puzzle, cell, c.Cells, c.Relative, func(other *Cell, index int) {
		if other.HasValue() {
			surroundingValues.Set(other.Value, true)
		} else {
			surrounding.Or(other.candidates)
		}
	})

	candidateMin := puzzle.MinCandidate()
	candidateMax := puzzle.MaxCandidate()

	for surroundingValues.Count > 0 {
		candidate := surroundingValues.First()
		surroundingValues.Set(candidate, false)

		doMinMaxDifference(candidate, c.Min, c.Max, candidateMin, candidateMax, remove, false)
	}

	if surrounding.Count > 0 {
		common := Candidates{}
		common.Fill(candidateMax)

		for surrounding.Count > 0 {
			candidate := surrounding.First()
			surrounding.Set(candidate, false)

			unique := Candidates{}
			doMinMaxDifference(candidate, c.Min, c.Max, candidateMin, candidateMax, &unique, true)
			common.And(unique)
		}

		if common.Count > 0 {
			remove.Remove(common)
		}
	}
}

func doMinMaxDifference(candidate int, min int, max int, candidateMin int, candidateMax int, out *Candidates, set bool) {
	if min > 1 {
		minMin := Max(candidate-min+1, candidateMin)
		minMax := Min(candidate+min-1, candidateMax)
		for c := minMin; c <= minMax; c++ {
			out.Set(c, set)
		}
	}

	if max > 0 {
		maxMin := candidate - max
		maxMax := candidate + max

		for c := candidateMin; c < maxMin; c++ {
			out.Set(c, set)
		}
		for c := maxMax + 1; c <= candidateMax; c++ {
			out.Set(c, set)
		}
	}
}

// ==================================================
// Constraint: Divisible
// ==================================================
type ConstraintDivisible struct {
	By        int
	Remainder int
	Cells     []Position
}

func (c *ConstraintDivisible) Affects(cell *Cell) bool {
	return containsCell(cell, &c.Cells, nil)
}

func (c *ConstraintDivisible) RemoveCandidates(cell *Cell, puzzle *Puzzle, remove *Candidates) {
	cand := remove
	for cand.Count > 0 {
		candidate := cand.First()
		cand.Set(candidate, false)

		if candidate%c.By != c.Remainder {
			remove.Set(candidate, false)
		}
	}
}

func ConstraintEven(cells []Position) ConstraintDivisible {
	return ConstraintDivisible{
		By:        2,
		Remainder: 0,
		Cells:     cells,
	}
}

func ConstraintOdd(cells []Position) ConstraintDivisible {
	return ConstraintDivisible{
		By:        2,
		Remainder: 1,
		Cells:     cells,
	}
}

// Functions

func traverseCells(puzzle *Puzzle, cell *Cell, absolute *[]Position, relative *[]Position, traverse func(other *Cell, index int)) {
	if relative != nil {
		for i := range *relative {
			pos := (*relative)[i]
			cell := getRelativeCell(puzzle, pos, cell)
			if cell != nil {
				traverse(cell, i)
			}
		}

	} else if absolute != nil {
		for i := range *absolute {
			pos := (*absolute)[i]
			cell := getAbsoluteCell(puzzle, pos)
			traverse(cell, i)
		}
	}
}

func isSame(cell *Cell, pos Position) bool {
	return cell.Col == pos.Col && cell.Row == pos.Row
}

func getAbsoluteCell(puzzle *Puzzle, pos Position) *Cell {
	return puzzle.Get(pos.Col, pos.Row)
}

func getRelativeCell(puzzle *Puzzle, pos Position, relative *Cell) *Cell {
	col := relative.Col + pos.Col
	row := relative.Row + pos.Row
	if !puzzle.Contains(col, row) {
		return nil
	}
	return puzzle.Get(col, row)
}

func getCells(puzzle *Puzzle, cell *Cell, absolute *[]Position, relative *[]Position) []*Cell {
	n := 0
	if relative != nil {
		n = len(*relative)
	} else if absolute != nil {
		n = len(*absolute)
	}

	cells := make([]*Cell, n)
	traverseCells(puzzle, cell, absolute, relative, func(other *Cell, index int) {
		cells[index] = other
	})
	return cells
}

func containsCell(cell *Cell, cells *[]Position, exclude *[]Position) bool {
	if cells == nil {
		if exclude == nil {
			return true
		}

		return !cellExists(cell, *exclude)
	}
	return cellExists(cell, *cells)
}

func cellExists(cell *Cell, cells []Position) bool {
	for _, p := range cells {
		if isSame(cell, p) {
			return true
		}
	}
	return false
}

func intsSum(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
}

func intsUnique(values []int) bool {
	cand := Candidates{}
	for _, v := range values {
		if cand.Has(v) {
			return false
		}
		cand.Set(v, true)
	}
	return true
}

type combinations[T any] struct {
	groups  [][]T
	current []int
}

func (c *combinations[T]) add(item []T) {
	c.groups = append(c.groups, item)
	c.current = append(c.current, 0)
}

func (c *combinations[T]) empty() bool {
	return len(c.groups) == 0
}

func (c *combinations[T]) reset() {
	for i := range c.current {
		c.current[i] = 0
	}
}

func (c *combinations[T]) done() bool {
	last := len(c.current) - 1
	return c.current[last] == len(c.groups[last])
}

func (c *combinations[T]) increment(k int) bool {
	c.current[k]++
	if c.current[k] == len(c.groups[k]) {
		if k == len(c.current)-1 {
			return false
		} else {
			c.current[k] = 0
			return c.increment(k + 1)
		}
	}
	return true
}

func (c *combinations[T]) start() []T {
	return make([]T, len(c.groups))
}

func (c *combinations[T]) next(out []T) bool {
	if c.done() {
		return false
	}
	for i, g := range c.groups {
		out[i] = g[c.current[i]]
	}
	c.increment(0)
	return true
}
