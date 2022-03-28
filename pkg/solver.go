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
	group.box = removeValue[*Cell](group.box, neighbor)
	group.row = removeValue[*Cell](group.row, neighbor)
	group.col = removeValue[*Cell](group.col, neighbor)
}

type Solver struct {
	puzzle   Puzzle
	steps    []SolveStep
	cells    []CellGroups
	unsolved []*CellGroups
}

// Solve steps are done sequentially repeatedly until all steps report back that nothing was changed.
type SolveStep func(solver *Solver, max int) (placements int, restart bool)

var StandardSolveSteps = []SolveStep{
	StepNakedSingle,
	StepHiddenSingle,
	StepRemovePointingCandidates,
	StepRemoveClaimingCandidates,
	StepRemoveNakedSubsetCandidates,
	StepRemoveHiddenSubsetCandidates,
}

var GenerateSolveSteps = []SolveStep{
	StepNakedSingle,
	StepHiddenSingle,
	StepRemovePointingCandidates,
	StepRemoveClaimingCandidates,
	StepRemoveNakedSubsetCandidates,
}

func NewSolver(starting Puzzle) Solver {
	puzzle := starting.Clone()
	steps := StandardSolveSteps
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

	return Solver{puzzle, steps, cells, unsolved}
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
	if group == nil || value <= 0 {
		return false
	}
	set := group.cell.SetValue(value)
	if set {
		for _, other := range group.all {
			other.RemoveCandidate(value)
		}

		solver.unsolved = removeValue[*CellGroups](solver.unsolved, group)

		for _, remaining := range solver.unsolved {
			remaining.Remove(group.cell)
		}
	}
	return set
}

type SetValueProvider func(solver *Solver) (*CellGroups, int)

func (solver *Solver) DoSets(max int, nextSet SetValueProvider) int {
	set := 0
	group, value := nextSet(solver)
	for group != nil {
		set++
		solver.SetGroup(group, value)
		if max > 0 && set == max {
			break
		}
		group, value = nextSet(solver)
	}
	return set
}

func (solver *Solver) Solved() bool {
	return len(solver.unsolved) == 0
}

func (solver *Solver) Solve() (solution *Puzzle, solved bool) {
	solver.Place(-1)
	return &solver.puzzle, solver.Solved()
}

func (solver *Solver) Place(count int) int {
	steps := solver.steps
	placed := 0
	placing := true
	for placing {
		placing = false
		for _, step := range steps {
			stepPlaced, stepRestart := step(solver, count-placed)

			placed += stepPlaced

			if count > 0 && placed >= count {
				placing = false
				break
			}
			if stepPlaced > 0 {
				placing = true
			}
			if stepRestart {
				placing = true
				break
			}
		}
	}

	return placed
}

// ==================================================
// Step: Naked Single
// 		http://hodoku.sourceforge.net/en/tech_singles.php
// ==================================================
var StepNakedSingle SolveStep = func(solver *Solver, max int) (int, bool) {
	return solver.DoSets(max, getNakedSingle), false
}

// A cell which has one possible candidate
func getNakedSingle(solver *Solver) (*CellGroups, int) {
	for _, group := range solver.unsolved {
		if group.cell.candidates.Count == 1 {
			return group, group.cell.FirstCandidate()
		}
	}
	return nil, 0
}

// ==================================================
// Step: Hidden Single
// 		http://hodoku.sourceforge.net/en/tech_singles.php
// ==================================================
var StepHiddenSingle SolveStep = func(solver *Solver, max int) (int, bool) {
	placed := solver.DoSets(max, getHiddenSingle)
	return placed, placed > 0
}

// A cell which has a candidate that is unique to the row, cell, or box
func getHiddenSingle(solver *Solver) (*CellGroups, int) {
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
	on := cell.candidates

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

// ==================================================
// Step: Remove Pointing Candidates
//		http://hodoku.sourceforge.net/en/tech_intersections.php
// ==================================================
var StepRemovePointingCandidates SolveStep = func(solver *Solver, max int) (int, bool) {
	return 0, doRemovePointingCandidates(solver, -1) > 0
}

// If in a box all candidates of a certain digit are confined to a row or column, that digit cannot appear outside of that box in that row or column.
func doRemovePointingCandidates(solver *Solver, max int) int {
	removed := 0

	for _, group := range solver.unsolved {
		cell := group.cell
		// all candidates in this box's row that are shared
		row := cell.candidates
		// all candidates in this box's column that are shared
		col := cell.candidates

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
			if other.row != cell.row {
				row.Remove(other.candidates)
			}
			if other.col != cell.col {
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

// ==================================================
// Step: Remove Pointing Candidates
//		http://hodoku.sourceforge.net/en/tech_intersections.php
// ==================================================
var StepRemoveClaimingCandidates SolveStep = func(solver *Solver, max int) (int, bool) {
	return 0, doRemoveClaimingCandidates(solver, -1) > 0
}

// If in a row or column a candidate only appears in a single box then that candidate can be removed from other cells in that box
func doRemoveClaimingCandidates(solver *Solver, max int) int {
	removed := 0

	for _, group := range solver.unsolved {
		cell := group.cell

		// all candidates in this row that are not shared outside of the box
		row := cell.candidates

		// remove candidates from row that exist in the cells row outside the box
		for _, other := range group.row {
			if other.box != cell.box {
				row.Remove(other.candidates)
			}
		}

		// what is remaining are the candidates unique to the row outside this box
		if row.Count > 0 {
			for _, other := range group.box {
				if other.row != cell.row {
					removed += other.candidates.Remove(row)
				}
			}
		}

		if max > 0 && removed >= max {
			break
		}

		// all candidates in this column that are not shared outside of the box
		col := cell.candidates

		// remove candidates from column that exist in the cells column outside the box
		for _, other := range group.col {
			if other.box != cell.box {
				col.Remove(other.candidates)
			}
		}

		// what is remaining are the candidates unique to the column outside this box
		if col.Count > 0 {
			for _, other := range group.box {
				if other.col != cell.col {
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

// ==================================================
// Step: Remove Naked Subset Candidates
//		http://hodoku.sourceforge.net/en/tech_naked.php
// ==================================================
func CreateStepRemoveNakedSubsetCandidates(subsets []int) SolveStep {
	return func(solver *Solver, max int) (int, bool) {
		return 0, doRemoveNakedSubsetCandidates(solver, -1, subsets) > 0
	}
}

var StepRemoveNakedSubsetCandidates SolveStep = CreateStepRemoveNakedSubsetCandidates([]int{2, 3})

// Find naked subsets and remove them as possible values for shared groups
func doRemoveNakedSubsetCandidates(solver *Solver, max int, subsets []int) int {
	removed := 0

	for _, subsetSize := range subsets {
		for _, group := range solver.unsolved {
			cell := group.cell

			if cell.candidates.Count != subsetSize {
				continue
			}
			removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, group.row)
			if max > 0 && removed >= max {
				break
			}
			removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, group.col)
			if max > 0 && removed >= max {
				break
			}
			removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, group.box)
			if max > 0 && removed >= max {
				break
			}
		}
	}

	return removed
}

// Remove naked subsets from group
func removeNakedSubsetCandidatesFromGroup(cellGroup *CellGroups, subsetSize int, group []*Cell) int {
	removed := 0
	matches := 1
	candidates := cellGroup.cell.candidates
	sameBox := true
	sameRow := true
	sameCol := true

	for _, other := range group {
		if other.candidates.Value == candidates.Value {
			matches++
			sameBox = sameBox && other.box == cellGroup.cell.box
			sameRow = sameRow && other.row == cellGroup.cell.row
			sameCol = sameCol && other.col == cellGroup.cell.col
		}
	}

	if matches == subsetSize {
		if sameBox {
			removed += removeCandidatesFromDifferent(cellGroup.box, candidates)
		}
		if sameRow {
			removed += removeCandidatesFromDifferent(cellGroup.row, candidates)
		}
		if sameCol {
			removed += removeCandidatesFromDifferent(cellGroup.col, candidates)
		}
	}
	return removed
}

func removeCandidatesFromDifferent(group []*Cell, candidates Candidates) int {
	removed := 0
	for _, other := range group {
		if other.candidates.Value != candidates.Value {
			removed += other.candidates.Remove(candidates)
		}
	}
	return removed
}

type candidateCells struct {
	candidate int
	cells     []*Cell
	size      int
}

func (cells *candidateCells) isSubset(other candidateCells) bool {
	matched := 0
	for i := 0; i < cells.size; i++ {
		if cells.cells[i] == other.cells[matched] {
			matched++
		}
	}
	return matched >= other.size
}

type candidateDistribution struct {
	candidates []candidateCells
}

func newDistribution(size int) candidateDistribution {
	candidates := make([]candidateCells, size)
	for i := 0; i < size; i++ {
		candidates[i].candidate = i + 1
		candidates[i].cells = make([]*Cell, size)
		candidates[i].size = 0
	}

	return candidateDistribution{candidates}
}

func (dist *candidateDistribution) set(cells []*Cell) {
	for i := range dist.candidates {
		dist.candidates[i].size = 0
	}
	for _, cell := range cells {
		for i := range dist.candidates {
			list := &dist.candidates[i]
			if cell.candidates.Has(list.candidate) {
				list.cells[list.size] = cell
				list.size++
			}
		}
	}
}

// ==================================================
// Step: Remove Hidden Subset Candidates
//		http://hodoku.sourceforge.net/en/tech_hidden.php
// ==================================================
func CreateStepRemoveHiddenSubsetCandidates(subsets []int) SolveStep {
	return func(solver *Solver, max int) (int, bool) {
		return 0, doRemoveHiddenSubsetCandidates(solver, -1, subsets) > 0
	}
}

var StepRemoveHiddenSubsetCandidates SolveStep = CreateStepRemoveHiddenSubsetCandidates([]int{2, 3})

// Find hidden subsets and remove them as possible values for shared groups
func doRemoveHiddenSubsetCandidates(solver *Solver, max int, subsets []int) int {
	dist := newDistribution(solver.puzzle.kind.Size())
	rowsTested := Bitset{}
	colsTested := Bitset{}
	boxsTested := Bitset{}
	removed := 0

	for _, group := range solver.unsolved {
		cell := group.cell

		// Only test a row/column/box of an unsolved cell once
		if !rowsTested.Has(cell.row) {
			rowsTested.Set(cell.row, true)
			fullRow := group.row[:]
			fullRow = append(fullRow, cell)
			dist.set(fullRow)
			removed += doRemoveHiddenSubset(&dist, max, subsets)
			if max > 0 && removed >= max {
				break
			}
		}
		if !colsTested.Has(cell.col) {
			colsTested.Set(cell.col, true)
			fullCol := group.col[:]
			fullCol = append(fullCol, cell)
			dist.set(fullCol)
			removed += doRemoveHiddenSubset(&dist, max, subsets)
			if max > 0 && removed >= max {
				break
			}
		}
		if !boxsTested.Has(cell.box) {
			boxsTested.Set(cell.box, true)
			fullBox := group.box[:]
			fullBox = append(fullBox, cell)
			dist.set(fullBox)
			removed += doRemoveHiddenSubset(&dist, max, subsets)
			if max > 0 && removed >= max {
				break
			}
		}
	}

	return removed
}

func doRemoveHiddenSubset(dist *candidateDistribution, max int, subsets []int) int {
	removed := 0
	for _, subsetSize := range subsets {
		n := len(dist.candidates)

		for listIndex := 0; listIndex < n; listIndex++ {
			list := dist.candidates[listIndex]

			if list.size == subsetSize {
				matchCandidates := Candidates{}
				matchCandidates.Set(list.candidate, true)

				for otherIndex := listIndex + 1; otherIndex < n; otherIndex++ {
					other := dist.candidates[otherIndex]

					if other.size > 0 && other.size <= subsetSize && list.isSubset(other) {
						matchCandidates.Set(other.candidate, true)
					}
				}
				if matchCandidates.Count >= subsetSize {
					for i := 0; i < list.size; i++ {
						removed += list.cells[i].candidates.And(matchCandidates)
					}
					if max > 0 && removed >= max {
						return removed
					}
					break
				}
			}
		}
	}

	return removed
}

// ==================================================
// Step: Remove Skyscraper Candidates
//		http://hodoku.sourceforge.net/en/tech_sdp.php
// ==================================================

// Find two rows that contain only two candidates for that digit. If two of those candidates are in the same column, one of the other two candidates must be true. All candidates that see both of those cells can therefore be eliminated.
