package sudogo

type CellGroups struct {
	cell *Cell
	all  []*Cell
	box  []*Cell
	row  []*Cell
	col  []*Cell
}

func (group *CellGroups) Remove(neighbor *Cell) {
	group.all = removeValue[*Cell](group.all, neighbor)
	group.box = removeValue[*Cell](group.all, neighbor)
	group.row = removeValue[*Cell](group.row, neighbor)
	group.col = removeValue[*Cell](group.col, neighbor)
}

type Solver struct {
	puzzle   Puzzle
	cells    []CellGroups
	unsolved []*CellGroups
}

func NewSolver(starting Puzzle) Solver {
	puzzle := starting.Clone()
	cells := make([]CellGroups, puzzle.kind.Area())
	unsolved := make([]*CellGroups, 0, puzzle.kind.Area())
	groupCapacity := puzzle.kind.Digits()
	allCapacity := groupCapacity * 3

	for i := range puzzle.cells {
		cell := &puzzle.cells[i]
		group := &cells[i]

		group.cell = cell
		group.all = make([]*Cell, 0, allCapacity)
		group.box = make([]*Cell, 0, groupCapacity)
		group.row = make([]*Cell, 0, groupCapacity)
		group.col = make([]*Cell, 0, groupCapacity)
		if cell.Empty() {
			unsolved = append(unsolved, group)
		}
		for k := range puzzle.cells {
			other := &puzzle.cells[k]
			if i != k && other.Empty() {
				if cell.InGroup(other) {
					group.all = append(group.all, other)
				}
				if cell.InBox(other) {
					group.box = append(group.box, other)
				}
				if cell.InRow(other) {
					group.row = append(group.row, other)
				}
				if cell.InColumn(other) {
					group.col = append(group.col, other)
				}
			}
		}
	}

	return Solver{puzzle, cells, unsolved}
}

func (solver *Solver) Set(col int, row int, value int) bool {
	return solver.SetCell(solver.puzzle.Get(col, row), value)
}

func (solver *Solver) SetCell(cell *Cell, value int) bool {
	if cell == nil {
		return false
	}
	return solver.SetGroup(&solver.cells[cell.id], value)
}

func (solver *Solver) SetGroup(group *CellGroups, value int) bool {
	if group == nil {
		return false
	}
	set := solver.puzzle.SetCell(group.cell, value)
	if set {
		solver.unsolved = removeValue[*CellGroups](solver.unsolved, group)
		for _, remaining := range solver.unsolved {
			remaining.Remove(group.cell)
		}
	}
	return set
}

// A cell which has one possible candidate
func (solver *Solver) GetNakedSingle() (*CellGroups, int) {
	for _, group := range solver.unsolved {
		if group.cell.candidates.Count == 1 {
			return group, group.cell.FirstCandidate()
		}
	}
	return nil, 0
}

// Solve N naked singles if that many are available. If max is less than 1 then all singles will be solved.
func (solver *Solver) SolveNakedSingles(max int) int {
	return solver.DoSets(max, func() (*CellGroups, int) {
		return solver.GetNakedSingle()
	})
}

// A cell which has a candidate that is unique to the row, cell, or box
func (solver *Solver) GetHiddenSingle() (*CellGroups, int) {
	for _, group := range solver.unsolved {
		box := getHiddenSingleFromGroup(group.cell, group.box)
		if box != 0 {
			return group, box
		}
		row := getHiddenSingleFromGroup(group.cell, group.row)
		if row != 0 {
			return group, row
		}
		col := getHiddenSingleFromGroup(group.cell, group.col)
		if col != 0 {
			return group, col
		}
	}
	return nil, 0
}

// Get the candidate hidden single found in the given group, or 0 if none found.
func getHiddenSingleFromGroup(cell *Cell, group []*Cell) int {
	on := cell.candidates.Clone()

	for _, other := range group {
		on.Remove(other.candidates)
		if on.Count == 0 {
			break
		}
	}
	if on.Count == 1 {
		return on.First()
	}
	return 0
}

// Solve N hidden singles if that many are available. If max is less than 1 then all singles will be solved.
func (solver *Solver) SolveHiddenSingles(max int) int {
	return solver.DoSets(max, func() (*CellGroups, int) {
		return solver.GetHiddenSingle()
	})
}

// If in a box all candidates of a certain digit are confined to a row or column, that digit cannot appear outside of that box in that row or column.
func (solver *Solver) RemovePointingCandidates(max int) int {
	removed := 0

	for _, group := range solver.unsolved {
		cell := group.cell
		// all candidates in this box's row that are shared
		row := cell.candidates.Clone()
		// all candidates in this box's column that are shared
		col := cell.candidates.Clone()

		// remove candidates that are not shared
		for _, other := range group.box {
			if other.row == cell.row {
				row.And(other.candidates)
			}
			if other.col == cell.col {
				col.And(other.candidates)
			}
		}

		// remove candidates that exist outside the row or column
		for _, other := range group.box {
			if other.row != cell.row && other.col != cell.col {
				row.Remove(other.candidates)
				col.Remove(other.candidates)
			}
		}

		// what is remaining are candidates confined to the cells row in the box
		if row.Count > 0 {
			for _, other := range group.row {
				if other.box != cell.box {
					removed += other.candidates.Remove(row)
				}
			}
		}

		if max > 0 && removed >= max {
			break
		}

		// what is remaining are candidates confined to the cells column in the box
		if col.Count > 0 {
			for _, other := range group.col {
				if other.box != cell.box {
					removed += other.candidates.Remove(col)
				}
			}
		}

		if max > 0 && removed >= max {
			break
		}
	}

	return removed
}

// If in a row or column a candidate only appears in a single box then that candidate can be removed from other cells in that box
func (solver *Solver) RemoveClaimingCandidates(max int) int {
	removed := 0

	for _, group := range solver.unsolved {
		cell := group.cell
		// all candidates in this row that are not shared outside of the box
		row := cell.candidates.Clone()
		// all candidates in this column that are not shared outside of the box
		col := cell.candidates.Clone()

		// remove candidates from row that exist in the cells row outside the box
		for _, other := range group.row {
			if other.box != cell.box {
				row.Remove(other.candidates)
			}
		}

		// remove candidates from column that exist in the cells column outside the box
		for _, other := range group.col {
			if other.box != cell.box {
				col.Remove(other.candidates)
			}
		}

		// what is remaining are the candidates unique to the row/column outside this box
		if row.Count > 0 || col.Count > 0 {
			for _, other := range group.box {
				if row.Count > 0 && other.row != cell.row {
					removed += other.candidates.Remove(row)
				}
				if col.Count > 0 && other.col != cell.col {
					removed += other.candidates.Remove(col)
				}
			}
		}

		if max > 0 && removed >= max {
			break
		}
	}

	return removed
}

func (solver *Solver) RemoveHiddenSubsetCandidates(max int) int {
	removed := 0
	// subsets := [2]int{2, 3}

	// for _, subsetSize := range subsets {
	// 	for i := range instance.cells {
	// 		cell := &instance.cells[i]
	// 		if cell.HasValue() {
	// 			continue
	// 		}

	// 	}
	// }

	return removed
}

func (solver *Solver) RemoveNakedSubsetCandidates(max int) int {
	removed := 0
	subsets := [2]int{2, 3}

	for _, subsetSize := range subsets {
		for _, group := range solver.unsolved {
			cell := group.cell

			if cell.candidates.Count != subsetSize {
				continue
			}

			removed += removeNakedSubsetCandidatesFromGroup(cell.candidates, subsetSize, group.row)
			if max > 0 && removed >= max {
				break
			}

			removed += removeNakedSubsetCandidatesFromGroup(cell.candidates, subsetSize, group.col)
			if max > 0 && removed >= max {
				break
			}

			removed += removeNakedSubsetCandidatesFromGroup(cell.candidates, subsetSize, group.box)
			if max > 0 && removed >= max {
				break
			}

		}
	}

	return removed
}

func removeNakedSubsetCandidatesFromGroup(candidates Candidates, subsetSize int, group []*Cell) int {
	removed := 0
	matches := 1

	for _, other := range group {
		if other.candidates.Value == candidates.Value {
			matches++
		}
	}

	if matches == subsetSize {
		for _, other := range group {
			if other.candidates.Value != candidates.Value {
				removed += other.candidates.Remove(candidates)
			}
		}
	}
	return removed
}

func (solver *Solver) DoSets(max int, nextSet func() (*CellGroups, int)) int {
	set := 0
	group, value := nextSet()
	for group != nil {
		set++
		solver.SetGroup(group, value)
		if max > 0 && set == max {
			break
		}
		group, value = nextSet()
	}
	return set
}

func (solver *Solver) Solved() bool {
	return len(solver.unsolved) == 0
}

func (solver *Solver) Solve() int {
	placed := 0
	for {
		placed += solver.SolveNakedSingles(-1)

		if solver.SolveHiddenSingles(1) > 0 {
			placed++
			continue
		}

		if solver.RemovePointingCandidates(1) > 0 {
			continue
		}

		if solver.RemoveClaimingCandidates(1) > 0 {
			continue
		}

		if solver.RemoveHiddenSubsetCandidates(1) > 0 {
			continue
		}

		if solver.RemoveNakedSubsetCandidates(1) > 0 {
			continue
		}

		break
	}

	return placed
}
