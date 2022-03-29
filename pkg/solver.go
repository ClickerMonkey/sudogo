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
	puzzle        Puzzle
	steps         []*SolveStep
	cells         []CellGroups
	unsolved      []*CellGroups
	logEnabled    bool
	logTemplate   SolverLog
	logs          []SolverLog
	logTechniques map[string]int
}

type SolverLog struct {
	step              *SolveStep
	index             int
	batch             int
	cost              int
	placement         bool
	before            Cell
	after             Cell
	runningCost       int
	runningPlacements int
}

type SolverLimit struct {
	minCost       int
	maxCost       int
	maxPlacements int
	maxLogs       int
	maxBatches    int
}

type SolveStepLogic func(solver *Solver, limits SolverLimit, step *SolveStep) (placements int, restart bool)

type SolveStep struct {
	technique      string
	firstCost      int
	subsequentCost int
	logic          SolveStepLogic
}

var StandardSolveSteps = []*SolveStep{
	StepNakedSingle,
	StepHiddenSingle,
	StepRemovePointingCandidates,
	StepRemoveClaimingCandidates,
	StepRemoveNakedSubsetCandidates2,
	StepRemoveHiddenSubsetCandidates2,
	StepRemoveNakedSubsetCandidates3,
	StepRemoveHiddenSubsetCandidates3,
	StepRemoveNakedSubsetCandidates4,
	StepRemoveHiddenSubsetCandidates4,
}

var GenerateSolveSteps = []*SolveStep{
	StepNakedSingle,
	StepHiddenSingle,
	StepRemovePointingCandidates,
	StepRemoveClaimingCandidates,
	StepRemoveNakedSubsetCandidates2,
	StepRemoveNakedSubsetCandidates3,
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

	logEnabled := false
	logTemplate := SolverLog{}
	logs := []SolverLog{}
	logTechniques := map[string]int{}

	return Solver{puzzle, steps, cells, unsolved, logEnabled, logTemplate, logs, logTechniques}
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

func (solver *Solver) GetMinCandidateCount() int {
	min := 0
	for _, group := range solver.unsolved {
		if min == 0 || min > group.cell.candidates.Count {
			min = group.cell.candidates.Count
		}
	}
	return min
}

func (solver *Solver) GetGroupWhere(where func(group *CellGroups) bool) *CellGroups {
	for _, group := range solver.unsolved {
		if where(group) {
			return group
		}
	}
	return nil
}

func (solver *Solver) GetLogs() []SolverLog {
	return solver.logs
}

func (solver *Solver) GetLastLog() *SolverLog {
	n := len(solver.logs) - 1
	if n == -1 {
		return &solver.logTemplate
	}
	return &solver.logs[n]
}

func (solver *Solver) canContinue(limits SolverLimit, cost int) bool {
	lastLog := solver.GetLastLog()
	if limits.maxLogs > 0 && lastLog.index >= limits.maxLogs {
		return false
	}
	if limits.maxBatches > 0 && lastLog.batch > limits.maxBatches {
		return false
	}
	if limits.maxCost > 0 && lastLog.runningCost+cost > limits.maxCost {
		return false
	}
	if limits.minCost > 0 && lastLog.runningCost >= limits.minCost {
		return false
	}
	if limits.maxPlacements > 0 && lastLog.runningPlacements >= limits.maxPlacements {
		return false
	}
	return true
}

func (solver *Solver) canContinueStep(limits SolverLimit, step *SolveStep) bool {
	return solver.canContinue(limits, solver.getCost(step))
}

func (solver *Solver) getCost(step *SolveStep) int {
	techniqueCount := solver.logTechniques[step.technique]
	cost := step.firstCost
	if techniqueCount > 0 {
		cost = step.subsequentCost
	}
	return cost
}

func (solver *Solver) logStep(step *SolveStep) {
	cost := solver.getCost(step)
	solver.logTechniques[step.technique]++
	solver.logTemplate.batch++
	solver.logTemplate.step = step
	solver.logTemplate.cost = cost
	solver.logTemplate.runningCost += cost
}

func (solver *Solver) logBefore(before *Cell) {
	if solver.logEnabled {
		log := solver.logTemplate
		log.before = *before
		solver.logs = append(solver.logs, log)
	}
	solver.logTemplate.index++
}

func (solver *Solver) logAfter(after *Cell) {
	last := solver.GetLastLog()
	last.after = *after
}

func (solver *Solver) logPlacement(after *Cell) {
	last := solver.GetLastLog()
	last.after = *after
	last.placement = true
	last.runningPlacements = last.runningPlacements + 1
	solver.logTemplate.runningPlacements = last.runningPlacements
}

func (solver *Solver) Solved() bool {
	return len(solver.unsolved) == 0
}

func (solver *Solver) Solve() (solution *Puzzle, solved bool) {
	solver.Place(SolverLimit{})
	return &solver.puzzle, solver.Solved()
}

func (solver *Solver) Place(limits SolverLimit) int {
	steps := solver.steps
	placed := 0
	placing := true
	for placing {
		placing = false
		for _, step := range steps {
			stepPlaced, stepRestart := step.logic(solver, limits, step)
			placed += stepPlaced

			if !solver.canContinue(limits, 0) {
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
var StepNakedSingle = &SolveStep{
	technique:      "Naked Single",
	firstCost:      100,
	subsequentCost: 100,
	logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
		placements := 0
		for solver.canContinueStep(limits, step) {
			group, groupValue := getNakedSingle(solver)
			if group != nil {
				solver.logStep(step)
				solver.logBefore(group.cell)
				solver.SetGroup(group, groupValue)
				solver.logPlacement(group.cell)
				placements++
			} else {
				break
			}
		}

		return placements, false
	},
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
var StepHiddenSingle = &SolveStep{
	technique:      "Hidden Single",
	firstCost:      100,
	subsequentCost: 100,
	logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
		placements := 0
		for solver.canContinueStep(limits, step) {
			group, groupValue := getHiddenSingle(solver)
			if group != nil {
				solver.logStep(step)
				solver.logBefore(group.cell)
				solver.SetGroup(group, groupValue)
				solver.logPlacement(group.cell)
				placements++
			} else {
				break
			}
		}

		return placements, placements > 0
	},
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
var StepRemovePointingCandidates = &SolveStep{
	technique:      "Pointing Candidates",
	firstCost:      350,
	subsequentCost: 200,
	logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
		removed := false
		if solver.canContinueStep(limits, step) {
			removed = doRemovePointingCandidates(solver, limits, step) > 0
		}
		return 0, removed
	},
}

// If in a box all candidates of a certain digit are confined to a row or column, that digit cannot appear outside of that box in that row or column.
func doRemovePointingCandidates(solver *Solver, limits SolverLimit, step *SolveStep) int {
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
			hasOverlap := false
			for _, other := range group.row {
				if other.box != cell.box && other.candidates.Overlaps(row) {
					hasOverlap = true
					break
				}
			}
			if hasOverlap {
				solver.logStep(step)
				for _, other := range group.row {
					if other.box != cell.box && other.candidates.Overlaps(row) {
						solver.logBefore(other)
						removed += other.candidates.Remove(row)
						solver.logAfter(other)
					}
				}
				if !solver.canContinueStep(limits, step) {
					break
				}
			}
		}

		// what is remaining are candidates confined to the cells column in the box
		if col.Count > 0 {
			hasOverlap := false
			for _, other := range group.col {
				if other.box != cell.box && other.candidates.Overlaps(row) {
					hasOverlap = true
					break
				}
			}
			if hasOverlap {
				solver.logStep(step)
				for _, other := range group.col {
					if other.box != cell.box && other.candidates.Overlaps(row) {
						solver.logBefore(other)
						removed += other.candidates.Remove(col)
						solver.logAfter(other)
					}
				}
				if !solver.canContinueStep(limits, step) {
					break
				}
			}
		}
	}

	return removed
}

// ==================================================
// Step: Remove Pointing Candidates
//		http://hodoku.sourceforge.net/en/tech_intersections.php
// ==================================================

var StepRemoveClaimingCandidates = &SolveStep{
	technique:      "Claiming Candidates",
	firstCost:      350,
	subsequentCost: 200,
	logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
		removed := false
		if solver.canContinueStep(limits, step) {
			removed = doRemoveClaimingCandidates(solver, limits, step) > 0
		}
		return 0, removed
	},
}

// If in a row or column a candidate only appears in a single box then that candidate can be removed from other cells in that box
func doRemoveClaimingCandidates(solver *Solver, limits SolverLimit, step *SolveStep) int {
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
			solver.logStep(step)
			for _, other := range group.box {
				if other.row != cell.row && other.candidates.Overlaps(row) {
					solver.logBefore(other)
					removed += other.candidates.Remove(row)
					solver.logAfter(other)
				}
			}
			if !solver.canContinueStep(limits, step) {
				break
			}
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
			solver.logStep(step)
			for _, other := range group.box {
				if other.col != cell.col && other.candidates.Overlaps(row) {
					solver.logBefore(other)
					removed += other.candidates.Remove(col)
					solver.logAfter(other)
				}
			}
			if !solver.canContinueStep(limits, step) {
				break
			}
		}
	}

	return removed
}

// ==================================================
// Step: Remove Naked Subset Candidates
//		http://hodoku.sourceforge.net/en/tech_naked.php
// ==================================================
func CreateStepRemoveNakedSubsetCandidates(subsetSize int, technique string, firstCost int, subsequentCost int) *SolveStep {
	return &SolveStep{
		technique:      technique,
		firstCost:      firstCost,
		subsequentCost: subsequentCost,
		logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
			removed := false
			if solver.canContinueStep(limits, step) {
				removed = doRemoveNakedSubsetCandidates(solver, subsetSize, limits, step) > 0
			}
			return 0, removed
		},
	}
}

var StepRemoveNakedSubsetCandidates2 = CreateStepRemoveNakedSubsetCandidates(2, "Naked Pair", 750, 500)
var StepRemoveNakedSubsetCandidates3 = CreateStepRemoveNakedSubsetCandidates(3, "Naked Triplet", 2000, 1400)
var StepRemoveNakedSubsetCandidates4 = CreateStepRemoveNakedSubsetCandidates(4, "Naked Quadruplet", 5000, 4000)

// Find naked subsets and remove them as possible values for shared groups
func doRemoveNakedSubsetCandidates(solver *Solver, subsetSize int, limits SolverLimit, step *SolveStep) int {
	removed := 0

	for _, group := range solver.unsolved {
		cell := group.cell

		if cell.candidates.Count != subsetSize {
			continue
		}
		removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, solver, limits, step, group.row)
		if !solver.canContinueStep(limits, step) {
			break
		}
		removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, solver, limits, step, group.col)
		if !solver.canContinueStep(limits, step) {
			break
		}
		removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, solver, limits, step, group.box)
		if !solver.canContinueStep(limits, step) {
			break
		}
	}

	return removed
}

// Remove naked subsets from group
func removeNakedSubsetCandidatesFromGroup(cellGroup *CellGroups, subsetSize int, solver *Solver, limits SolverLimit, step *SolveStep, group []*Cell) int {
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
			removed += removeCandidatesFromDifferent(cellGroup.box, candidates, solver, limits, step)
		}
		if sameRow && solver.canContinueStep(limits, step) {
			removed += removeCandidatesFromDifferent(cellGroup.row, candidates, solver, limits, step)
		}
		if sameCol && solver.canContinueStep(limits, step) {
			removed += removeCandidatesFromDifferent(cellGroup.col, candidates, solver, limits, step)
		}
	}
	return removed
}

func removeCandidatesFromDifferent(group []*Cell, candidates Candidates, solver *Solver, limits SolverLimit, step *SolveStep) int {
	removed := 0
	hasOverlap := false
	for _, other := range group {
		if other.candidates.Value != candidates.Value && other.candidates.Overlaps(candidates) {
			hasOverlap = true
			break
		}
	}
	if hasOverlap {
		solver.logStep(step)
		for _, other := range group {
			if other.candidates.Value != candidates.Value && other.candidates.Overlaps(candidates) {
				solver.logBefore(other)
				removed += other.candidates.Remove(candidates)
				solver.logAfter(other)
			}
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
func CreateStepRemoveHiddenSubsetCandidates(subsetSize int, technique string, firstCost int, subsequentCost int) *SolveStep {
	return &SolveStep{
		technique:      technique,
		firstCost:      firstCost,
		subsequentCost: subsequentCost,
		logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
			removed := false
			if solver.canContinueStep(limits, step) {
				removed = doRemoveHiddenSubsetCandidates(solver, subsetSize, limits, step) > 0
			}
			return 0, removed
		},
	}
}

var StepRemoveHiddenSubsetCandidates2 = CreateStepRemoveHiddenSubsetCandidates(2, "Hidden Pair", 1500, 1200)
var StepRemoveHiddenSubsetCandidates3 = CreateStepRemoveHiddenSubsetCandidates(3, "Hidden Triplet", 2400, 1600)
var StepRemoveHiddenSubsetCandidates4 = CreateStepRemoveHiddenSubsetCandidates(4, "Hidden Quadruplet", 7000, 5000)

// Find hidden subsets and remove them as possible values for shared groups
func doRemoveHiddenSubsetCandidates(solver *Solver, subsetSize int, limits SolverLimit, step *SolveStep) int {
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
			removed += doRemoveHiddenSubset(&dist, subsetSize, solver, limits, step)
			if !solver.canContinueStep(limits, step) {
				break
			}
		}
		if !colsTested.Has(cell.col) {
			colsTested.Set(cell.col, true)
			fullCol := group.col[:]
			fullCol = append(fullCol, cell)
			dist.set(fullCol)
			removed += doRemoveHiddenSubset(&dist, subsetSize, solver, limits, step)
			if !solver.canContinueStep(limits, step) {
				break
			}
		}
		if !boxsTested.Has(cell.box) {
			boxsTested.Set(cell.box, true)
			fullBox := group.box[:]
			fullBox = append(fullBox, cell)
			dist.set(fullBox)
			removed += doRemoveHiddenSubset(&dist, subsetSize, solver, limits, step)
			if !solver.canContinueStep(limits, step) {
				break
			}
		}
	}

	return removed
}

func doRemoveHiddenSubset(dist *candidateDistribution, subsetSize int, solver *Solver, limits SolverLimit, step *SolveStep) int {
	removed := 0
	n := len(dist.candidates)

	for listIndex := 0; listIndex < n; listIndex++ {
		list := dist.candidates[listIndex]

		if list.size == subsetSize {
			matchCandidates := Candidates{}
			matchCandidates.Set(list.candidate, true)

			for otherIndex := 0; otherIndex < n; otherIndex++ {
				other := dist.candidates[otherIndex]

				if other.size > 0 && other.size <= subsetSize && list.isSubset(other) {
					matchCandidates.Set(other.candidate, true)
				}
			}
			if matchCandidates.Count >= subsetSize {
				hasOverlap := false
				for i := 0; i < list.size; i++ {
					other := list.cells[i]
					if other.candidates.Differences(matchCandidates) {
						hasOverlap = true
						break
					}
				}
				if hasOverlap {
					solver.logStep(step)
					for i := 0; i < list.size; i++ {
						other := list.cells[i]
						if other.candidates.Differences(matchCandidates) {
							solver.logBefore(other)
							removed += other.candidates.And(matchCandidates)
							solver.logAfter(other)
						}
					}
					if !solver.canContinueStep(limits, step) {
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
