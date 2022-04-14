package sudogo

import "fmt"

type Solver struct {
	Puzzle   Puzzle
	Steps    []*SolveStep
	Unsolved []*Cell
	Boxs     [][]*Cell
	Rows     [][]*Cell
	Cols     [][]*Cell

	LogEnabled    bool
	LogState      bool
	LogTechniques map[string]int
	Logs          []SolverLog
	logTemplate   SolverLog
}

type SolverLog struct {
	Step              *SolveStep
	Index             int
	Batch             int
	Cost              int
	Placement         bool
	Before            Cell
	After             Cell
	State             *Puzzle
	RunningCost       int
	RunningPlacements int
}

func (log SolverLog) String() string {
	if log.Placement {
		return fmt.Sprintf("[%d]: %s placed %d at r%dc%d (cost=%d)", log.Index, log.Step.Technique, log.After.Value, log.After.Row+1, log.After.Col+1, log.Cost)
	} else {
		return fmt.Sprintf("[%d]: %s candidates %v to %v at r%dc%d (cost=%d)", log.Index, log.Step.Technique, log.Before.candidates.ToSlice(), log.After.candidates.ToSlice(), log.Before.Row+1, log.Before.Col+1, log.Cost)
	}
}

type SolveLimit struct {
	MinCost       int
	MaxCost       int
	MaxPlacements int
	MaxLogs       int
	MaxBatches    int
}

type SolveStepLogic func(solver *Solver, limits SolveLimit, step *SolveStep) (placements int, restart bool)

type SolveStep struct {
	Technique      string
	FirstCost      int
	SubsequentCost int
	Logic          SolveStepLogic
}

var StandardSolveSteps = []*SolveStep{
	StepNakedSingle,
	StepHiddenSingle,
	StepPointingCandidates,
	StepClaimingCandidates,
	StepConstraints,
	StepSkyscraper,
	Step2StringKite,
	StepNakedSubsets2,
	StepHiddenSubsets2,
	StepEmptyRectangle,
	StepXWing,
	StepNakedSubsets3,
	StepHiddenSubsets3,
	StepSwordfish,
	StepNakedSubsets4,
	StepHiddenSubsets4,
	StepJellyfish,
}

var GenerateSolveSteps = []*SolveStep{
	StepNakedSingle,
	StepHiddenSingle,
	StepPointingCandidates,
	StepClaimingCandidates,
	StepConstraints,
	StepNakedSubsets2,
	StepNakedSubsets3,
}

func NewSolver(starting Puzzle) Solver {
	puzzle := starting.Clone()
	unsolved := make([]*Cell, 0, puzzle.Kind.Area())
	groupCapacity := puzzle.Kind.Digits()
	rows := make([][]*Cell, groupCapacity)
	cols := make([][]*Cell, groupCapacity)
	boxs := make([][]*Cell, groupCapacity)

	for i := 0; i < groupCapacity; i++ {
		rows[i] = make([]*Cell, 0, groupCapacity)
		cols[i] = make([]*Cell, 0, groupCapacity)
		boxs[i] = make([]*Cell, 0, groupCapacity)
	}

	for i := range puzzle.Cells {
		cell := &puzzle.Cells[i]

		if cell.Empty() {
			unsolved = append(unsolved, cell)
			rows[cell.Row] = append(rows[cell.Row], cell)
			cols[cell.Col] = append(cols[cell.Col], cell)
			boxs[cell.Box] = append(boxs[cell.Box], cell)
		}
	}

	return Solver{
		Puzzle:        puzzle,
		Steps:         StandardSolveSteps,
		Unsolved:      unsolved,
		Rows:          rows,
		Cols:          cols,
		Boxs:          boxs,
		LogEnabled:    false,
		LogState:      false,
		LogTechniques: map[string]int{},
		Logs:          []SolverLog{},
		logTemplate:   SolverLog{},
	}
}

func (solver *Solver) Row(row int) []*Cell {
	return solver.Rows[row]
}

func (solver *Solver) Col(col int) []*Cell {
	return solver.Cols[col]
}

func (solver *Solver) Box(box int) []*Cell {
	return solver.Boxs[box]
}

func (solver *Solver) Group(groupIndex Group, cell *Cell) []*Cell {
	if groupIndex == GroupCol {
		return solver.Cols[cell.Col]
	} else if groupIndex == GroupRow {
		return solver.Rows[cell.Row]
	} else {
		return solver.Boxs[cell.Box]
	}
}

func (solver *Solver) Set(col int, row int, value int) bool {
	return solver.SetCell(solver.Puzzle.Get(col, row), value)
}

func (solver *Solver) SetCell(cell *Cell, value int) bool {
	if cell == nil || value <= 0 {
		return false
	}
	set := cell.SetValue(value)
	if set {
		rows := solver.Rows[cell.Row]
		cols := solver.Cols[cell.Col]
		boxs := solver.Boxs[cell.Box]

		for _, other := range rows {
			other.RemoveCandidate(value)
		}
		for _, other := range cols {
			other.RemoveCandidate(value)
		}
		for _, other := range boxs {
			other.RemoveCandidate(value)
		}

		solver.Unsolved = removeValue(solver.Unsolved, cell)
		solver.Rows[cell.Row] = removeValue(rows, cell)
		solver.Cols[cell.Col] = removeValue(cols, cell)
		solver.Boxs[cell.Box] = removeValue(boxs, cell)
	}

	return set
}

func (solver *Solver) GetMinCandidateCount() int {
	min := 0
	for _, cell := range solver.Unsolved {
		if min == 0 || min > cell.candidates.Count {
			min = cell.candidates.Count
		}
	}
	return min
}

func (solver *Solver) GetCellWhere(where func(cell *Cell) bool) *Cell {
	for _, cell := range solver.Unsolved {
		if where(cell) {
			return cell
		}
	}
	return nil
}

func (solver *Solver) GetLastLog() *SolverLog {
	n := len(solver.Logs) - 1
	if n == -1 {
		return &solver.logTemplate
	}
	return &solver.Logs[n]
}

func (solver *Solver) CanContinue(limits SolveLimit, cost int) bool {
	lastLog := solver.GetLastLog()
	if limits.MaxLogs > 0 && lastLog.Index >= limits.MaxLogs {
		return false
	}
	if limits.MaxBatches > 0 && lastLog.Batch > limits.MaxBatches {
		return false
	}
	if limits.MaxCost > 0 && lastLog.RunningCost+cost > limits.MaxCost {
		return false
	}
	if limits.MinCost > 0 && lastLog.RunningCost >= limits.MinCost {
		return false
	}
	if limits.MaxPlacements > 0 && lastLog.RunningPlacements >= limits.MaxPlacements {
		return false
	}
	return true
}

func (solver *Solver) CanContinueStep(limits SolveLimit, step *SolveStep) bool {
	return solver.CanContinue(limits, solver.GetCost(step))
}

func (solver *Solver) GetCost(step *SolveStep) int {
	techniqueCount := solver.LogTechniques[step.Technique]
	cost := step.FirstCost
	if techniqueCount > 0 {
		cost = step.SubsequentCost
	}
	return cost
}

func (solver *Solver) LogStep(step *SolveStep) {
	cost := solver.GetCost(step)
	solver.LogTechniques[step.Technique]++
	solver.logTemplate.Batch++
	solver.logTemplate.Step = step
	solver.logTemplate.Cost = cost
	solver.logTemplate.RunningCost += cost
}

func (solver *Solver) LogBefore(before *Cell) {
	if solver.LogEnabled {
		log := solver.logTemplate
		log.Before = *before
		solver.Logs = append(solver.Logs, log)
	}
	solver.logTemplate.Index++
}

func (solver *Solver) LogAfter(after *Cell) {
	last := solver.GetLastLog()
	last.After = *after
	if solver.LogState {
		state := solver.Puzzle.Clone()
		last.State = &state
	}
}

func (solver *Solver) LogPlacement(after *Cell) {
	last := solver.GetLastLog()
	last.After = *after
	last.Placement = true
	last.RunningPlacements = last.RunningPlacements + 1
	solver.logTemplate.RunningPlacements = last.RunningPlacements
	if solver.LogState {
		state := solver.Puzzle.Clone()
		last.State = &state
	}
}

func (solver *Solver) Solved() bool {
	return len(solver.Unsolved) == 0
}

func (solver *Solver) Solve(limits SolveLimit) (*Puzzle, bool) {
	steps := solver.Steps
	placing := true
	for placing {
		placing = false
		for _, step := range steps {
			stepPlaced, stepRestart := step.Logic(solver, limits, step)

			if !solver.CanContinue(limits, 0) {
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
	return &solver.Puzzle, solver.Solved()
}

// ==================================================
// Step: Brute Force
// ==================================================
var StepBruteForce = &SolveStep{
	Technique:      "Brute Force",
	FirstCost:      0,
	SubsequentCost: 0,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		return 0, false
	},
}

// ==================================================
// Step: Naked Single
// 		http://hodoku.sourceforge.net/en/tech_singles.php
// ==================================================
var StepNakedSingle = &SolveStep{
	Technique:      "Naked Single",
	FirstCost:      100,
	SubsequentCost: 100,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		placements := 0
		for solver.CanContinueStep(limits, step) {
			cell, cellValue := getNakedSingle(solver)
			if cell != nil {
				solver.LogStep(step)
				solver.LogBefore(cell)
				solver.SetCell(cell, cellValue)
				solver.LogPlacement(cell)
				placements++
			} else {
				break
			}
		}

		return placements, false
	},
}

// A cell which has one possible candidate
func getNakedSingle(solver *Solver) (*Cell, int) {
	for _, cell := range solver.Unsolved {
		if cell.candidates.Count == 1 {
			return cell, cell.FirstCandidate()
		}
	}
	return nil, 0
}

// ==================================================
// Step: Hidden Single
// 		http://hodoku.sourceforge.net/en/tech_singles.php
// ==================================================
var StepHiddenSingle = &SolveStep{
	Technique:      "Hidden Single",
	FirstCost:      100,
	SubsequentCost: 100,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		placements := 0
		for solver.CanContinueStep(limits, step) {
			cell, cellValue := getHiddenSingle(solver)
			if cell != nil {
				solver.LogStep(step)
				solver.LogBefore(cell)
				solver.SetCell(cell, cellValue)
				solver.LogPlacement(cell)
				placements++
			} else {
				break
			}
		}

		return placements, placements > 0
	},
}

// A cell which has a candidate that is unique to the row, cell, or box
func getHiddenSingle(solver *Solver) (*Cell, int) {
	for _, cell := range solver.Unsolved {
		box := getHiddenSingleFromGroup(cell, solver.Box(cell.Box))
		if box != 0 {
			return cell, box
		}
		row := getHiddenSingleFromGroup(cell, solver.Row(cell.Row))
		if row != 0 {
			return cell, row
		}
		col := getHiddenSingleFromGroup(cell, solver.Col(cell.Col))
		if col != 0 {
			return cell, col
		}
	}
	return nil, 0
}

// Get the candidate hidden single found in the given group, or 0 if none found.
func getHiddenSingleFromGroup(cell *Cell, group []*Cell) int {
	on := cell.candidates

	for _, other := range group {
		if other.Id == cell.Id {
			continue
		}
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
// Step: Constraints
// ==================================================
var StepConstraints = &SolveStep{
	Technique:      "Constraints",
	FirstCost:      0,
	SubsequentCost: 0,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		removed := 0

		for _, cell := range solver.Unsolved {
			if len(cell.Constraints) == 0 {
				continue
			}

			candidates := cell.candidates

			for _, constraint := range cell.Constraints {
				constraint.RemoveCandidates(cell, &solver.Puzzle, &candidates)
			}

			if candidates.Value != cell.candidates.Value {
				solver.LogStep(step)
				solver.LogBefore(cell)
				removed += cell.candidates.Count - candidates.Count
				cell.candidates = candidates
				solver.LogPlacement(cell)

				if !solver.CanContinueStep(limits, step) {
					break
				}
			}
		}

		return 0, removed > 0
	},
}

// ==================================================
// Step: Pointing Candidates
//		http://hodoku.sourceforge.net/en/tech_intersections.php
// ==================================================
var StepPointingCandidates = &SolveStep{
	Technique:      "Pointing Candidates",
	FirstCost:      350,
	SubsequentCost: 200,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		removed := false
		if solver.CanContinueStep(limits, step) {
			removed = doRemovePointingCandidates(solver, limits, step) > 0
		}
		return 0, removed
	},
}

// If in a box all candidates of a certain digit are confined to a row or column, that digit cannot appear outside of that box in that row or column.
func doRemovePointingCandidates(solver *Solver, limits SolveLimit, step *SolveStep) int {
	removed := 0

	for _, group := range solver.Unsolved {
		removed += doRemovePointingCandidatesGroup(solver, limits, step, group, GroupCol)
		if !solver.CanContinueStep(limits, step) {
			break
		}
		removed += doRemovePointingCandidatesGroup(solver, limits, step, group, GroupRow)
		if !solver.CanContinueStep(limits, step) {
			break
		}
	}

	return removed
}

func doRemovePointingCandidatesGroup(solver *Solver, limits SolveLimit, step *SolveStep, cell *Cell, groupIndex Group) int {
	// all candidates in this box's group that are shared
	cand := cell.candidates
	removed := 0

	// remove candidates that are not shared
	for _, other := range solver.Box(cell.Box) {
		if other.Id == cell.Id {
			continue
		}
		if other.GetGroup(groupIndex) == cell.GetGroup(groupIndex) {
			cand.And(other.candidates)
		}
	}

	// remove candidates that exist outside the row or column
	for _, other := range solver.Box(cell.Box) {
		if other.Id == cell.Id {
			continue
		}
		if other.GetGroup(groupIndex) != cell.GetGroup(groupIndex) {
			cand.Remove(other.candidates)
		}
	}

	// what is remaining are candidates confined to the cells column in the box
	if cand.Count > 0 {
		hasOverlap := false
		for _, other := range solver.Group(groupIndex, cell) {
			if other.Id == cell.Id {
				continue
			}
			if other.Box != cell.Box && other.candidates.Overlaps(cand) {
				hasOverlap = true
				break
			}
		}
		if hasOverlap {
			solver.LogStep(step)
			for _, other := range solver.Group(groupIndex, cell) {
				if other.Id == cell.Id {
					continue
				}
				if other.Box != cell.Box && other.candidates.Overlaps(cand) {
					solver.LogBefore(other)
					removed += other.candidates.Remove(cand)
					solver.LogAfter(other)
				}
			}
		}
	}

	return removed
}

// ==================================================
// Step: Claiming Candidates
//		http://hodoku.sourceforge.net/en/tech_intersections.php
// ==================================================
var StepClaimingCandidates = &SolveStep{
	Technique:      "Claiming Candidates",
	FirstCost:      350,
	SubsequentCost: 200,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		removed := false
		if solver.CanContinueStep(limits, step) {
			removed = doRemoveClaimingCandidates(solver, limits, step) > 0
		}
		return 0, removed
	},
}

// If in a row or column a candidate only appears in a single box then that candidate can be removed from other cells in that box
func doRemoveClaimingCandidates(solver *Solver, limits SolveLimit, step *SolveStep) int {
	removed := 0
	removed += doRemoveClaimingCandidatesGroups(solver, limits, step, GroupCol)
	if solver.CanContinueStep(limits, step) {
		removed += doRemoveClaimingCandidatesGroups(solver, limits, step, GroupRow)
	}
	return removed
}

func doRemoveClaimingCandidatesGroups(solver *Solver, limits SolveLimit, step *SolveStep, groupIndex Group) int {
	removed := 0

	for _, cell := range solver.Unsolved {

		// all candidates in this cand that are not shared outside of the box
		cand := cell.candidates

		// remove candidates from row that exist in the cells row outside the box
		for _, other := range solver.Group(groupIndex, cell) {
			if other.Id == cell.Id {
				continue
			}
			if other.Box != cell.Box {
				cand.Remove(other.candidates)
			}
		}

		// what is remaining are the candidates unique to the row outside this box
		if cand.Count > 0 {
			solver.LogStep(step)
			for _, other := range solver.Box(cell.Box) {
				if other.Id == cell.Id {
					continue
				}
				if other.GetGroup(groupIndex) != cell.GetGroup(groupIndex) && other.candidates.Overlaps(cand) {
					solver.LogBefore(other)
					removed += other.candidates.Remove(cand)
					solver.LogAfter(other)
				}
			}
			if !solver.CanContinueStep(limits, step) {
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
func CreateStepNakedSubsets(subsetSize int, technique string, firstCost int, subsequentCost int) *SolveStep {
	return &SolveStep{
		Technique:      technique,
		FirstCost:      firstCost,
		SubsequentCost: subsequentCost,
		Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
			removed := false
			if solver.CanContinueStep(limits, step) {
				removed = doRemoveNakedSubsetCandidates(solver, subsetSize, limits, step) > 0
			}
			return 0, removed
		},
	}
}

var StepNakedSubsets2 = CreateStepNakedSubsets(2, "Naked Pair", 750, 500)
var StepNakedSubsets3 = CreateStepNakedSubsets(3, "Naked Triplet", 2000, 1400)
var StepNakedSubsets4 = CreateStepNakedSubsets(4, "Naked Quad", 5000, 4000)

// Find naked subsets and remove them as possible values for shared groups
func doRemoveNakedSubsetCandidates(solver *Solver, subsetSize int, limits SolveLimit, step *SolveStep) int {
	removed := 0

	for _, cell := range solver.Unsolved {

		if cell.candidates.Count != subsetSize {
			continue
		}
		removed += removeNakedSubsetCandidatesFromGroup(cell, subsetSize, solver, limits, step, solver.Row(cell.Row))
		if !solver.CanContinueStep(limits, step) {
			break
		}
		removed += removeNakedSubsetCandidatesFromGroup(cell, subsetSize, solver, limits, step, solver.Col(cell.Col))
		if !solver.CanContinueStep(limits, step) {
			break
		}
		removed += removeNakedSubsetCandidatesFromGroup(cell, subsetSize, solver, limits, step, solver.Box(cell.Box))
		if !solver.CanContinueStep(limits, step) {
			break
		}
	}

	return removed
}

// Remove naked subsets from group
func removeNakedSubsetCandidatesFromGroup(cell *Cell, subsetSize int, solver *Solver, limits SolveLimit, step *SolveStep, group []*Cell) int {
	removed := 0
	matches := 1
	candidates := cell.candidates
	sameBox := true
	sameRow := true
	sameCol := true

	for _, other := range group {
		if other.Id == cell.Id {
			continue
		}
		if other.candidates.Value == candidates.Value {
			matches++
			sameBox = sameBox && other.Box == cell.Box
			sameRow = sameRow && other.Row == cell.Row
			sameCol = sameCol && other.Col == cell.Col
		}
	}

	if matches == subsetSize {
		if sameBox {
			removed += removeCandidatesFromDifferent(cell, solver.Box(cell.Box), candidates, solver, limits, step)
		}
		if sameRow && solver.CanContinueStep(limits, step) {
			removed += removeCandidatesFromDifferent(cell, solver.Row(cell.Row), candidates, solver, limits, step)
		}
		if sameCol && solver.CanContinueStep(limits, step) {
			removed += removeCandidatesFromDifferent(cell, solver.Col(cell.Col), candidates, solver, limits, step)
		}
	}
	return removed
}

func removeCandidatesFromDifferent(cell *Cell, group []*Cell, candidates Candidates, solver *Solver, limits SolveLimit, step *SolveStep) int {
	removed := 0
	hasOverlap := false
	for _, other := range group {
		if other.Id == cell.Id {
			continue
		}
		if other.candidates.Value != candidates.Value && other.candidates.Overlaps(candidates) {
			hasOverlap = true
			break
		}
	}
	if hasOverlap {
		solver.LogStep(step)
		for _, other := range group {
			if other.Id == cell.Id {
				continue
			}
			if other.candidates.Value != candidates.Value && other.candidates.Overlaps(candidates) {
				solver.LogBefore(other)
				removed += other.candidates.Remove(candidates)
				solver.LogAfter(other)
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

func (dist *candidateDistribution) clear() {
	for i := range dist.candidates {
		dist.candidates[i].size = 0
	}
}

func (dist *candidateDistribution) reset(cells []*Cell) {
	dist.clear()
	dist.addCells(cells)
}

func (dist *candidateDistribution) addCells(cells []*Cell) {
	for _, cell := range cells {
		dist.addCell(cell)
	}
}

func (dist *candidateDistribution) addCell(cell *Cell) {
	for i := range dist.candidates {
		list := &dist.candidates[i]
		if cell.candidates.Has(list.candidate) {
			list.cells[list.size] = cell
			list.size++
		}
	}
}

// ==================================================
// Step: Remove Hidden Subset Candidates
//		http://hodoku.sourceforge.net/en/tech_hidden.php
// ==================================================
func CreateStepHiddenSubsets(subsetSize int, technique string, firstCost int, subsequentCost int) *SolveStep {
	return &SolveStep{
		Technique:      technique,
		FirstCost:      firstCost,
		SubsequentCost: subsequentCost,
		Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
			removed := false
			if solver.CanContinueStep(limits, step) {
				removed = doRemoveHiddenSubsetCandidates(solver, subsetSize, limits, step) > 0
			}
			return 0, removed
		},
	}
}

var StepHiddenSubsets2 = CreateStepHiddenSubsets(2, "Hidden Pair", 1500, 1200)
var StepHiddenSubsets3 = CreateStepHiddenSubsets(3, "Hidden Triplet", 2400, 1600)
var StepHiddenSubsets4 = CreateStepHiddenSubsets(4, "Hidden Quad", 7000, 5000)

// Find hidden subsets and remove them as possible values for shared groups
func doRemoveHiddenSubsetCandidates(solver *Solver, subsetSize int, limits SolveLimit, step *SolveStep) int {
	dist := newDistribution(solver.Puzzle.Kind.Size())
	tested := [3]Bitset{}
	removed := 0

	for _, cell := range solver.Unsolved {

		for g := GroupCol; g <= GroupBox; g++ {
			// Only test a row/column/box of an unsolved cell once
			cellGroup := cell.GetGroup(g)
			if !tested[g].Has(cellGroup) {
				tested[g].Set(cellGroup, true)

				dist.reset(solver.Group(g, cell))

				removed += doRemoveHiddenSubset(&dist, subsetSize, solver, limits, step)

				if !solver.CanContinueStep(limits, step) {
					return removed
				}
			}
		}
	}

	return removed
}

func doRemoveHiddenSubset(dist *candidateDistribution, subsetSize int, solver *Solver, limits SolveLimit, step *SolveStep) int {
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
					solver.LogStep(step)
					for i := 0; i < list.size; i++ {
						other := list.cells[i]
						if other.candidates.Differences(matchCandidates) {
							solver.LogBefore(other)
							removed += other.candidates.And(matchCandidates)
							solver.LogAfter(other)
						}
					}
					if !solver.CanContinueStep(limits, step) {
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

var StepSkyscraper = &SolveStep{
	Technique:      "Skyscraper",
	FirstCost:      2800,
	SubsequentCost: 1600,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		return 0, doSkyscraper(solver, limits, step) > 0
	},
}

// Find two rows that contain only two candidates for that digit.
// If two of those candidates are in the same column, one of the other two candidates must be true.
// All candidates that see both of those cells can therefore be eliminated.
func doSkyscraper(solver *Solver, limits SolveLimit, step *SolveStep) int {
	removed := 0
	removed += doSkyscraperRemoveGroup(solver, limits, step, GroupCol)
	if solver.CanContinueStep(limits, step) {
		removed += doSkyscraperRemoveGroup(solver, limits, step, GroupRow)
	}
	return removed
}

func doSkyscraperRemoveGroup(solver *Solver, limits SolveLimit, step *SolveStep, groupIndex Group) int {
	size := solver.Puzzle.Kind.Size()
	removed := 0

	groups := getGroupCandidateDistributions(solver, groupIndex)

	oppositeGroup := 1 - groupIndex
	groupsLast := len(groups) - 1
	for a := 0; a < groupsLast; a++ {
		groupA := groups[a].candidates
		for b := a + 1; b <= groupsLast; b++ {
			groupB := groups[b].candidates
			for candidate := 0; candidate < size; candidate++ {
				if groupA[candidate].size == 2 && groupB[candidate].size == 2 {
					a0 := groupA[candidate].cells[0]
					a1 := groupA[candidate].cells[1]
					b0 := groupB[candidate].cells[0]
					b1 := groupB[candidate].cells[1]

					if a0.GetGroup(oppositeGroup) == b0.GetGroup(oppositeGroup) {
						removed += removeCandidateInGroups(solver, limits, step, candidate+1, a1, b1)
					} else if a0.GetGroup(oppositeGroup) == b1.GetGroup(oppositeGroup) {
						removed += removeCandidateInGroups(solver, limits, step, candidate+1, a1, b0)
					} else if a1.GetGroup(oppositeGroup) == b0.GetGroup(oppositeGroup) {
						removed += removeCandidateInGroups(solver, limits, step, candidate+1, a0, b1)
					} else if a1.GetGroup(oppositeGroup) == b1.GetGroup(oppositeGroup) {
						removed += removeCandidateInGroups(solver, limits, step, candidate+1, a0, b0)
					}
					if !solver.CanContinueStep(limits, step) {
						return removed
					}
				}
			}
		}
	}

	return removed
}

func removeCandidateInGroups(solver *Solver, limits SolveLimit, step *SolveStep, candidate int, a *Cell, b *Cell) int {
	removed := 0

	for _, cell := range solver.Unsolved {
		if cell.HasCandidate(candidate) && cell.InGroup(a) && cell.InGroup(b) {
			solver.LogStep(step)
			solver.LogBefore(cell)
			cell.RemoveCandidate(candidate)
			removed++
			solver.LogAfter(cell)

			if !solver.CanContinueStep(limits, step) {
				return removed
			}
		}
	}

	return removed
}

func getGroupCandidateDistributions(solver *Solver, groupIndex Group) []*candidateDistribution {
	size := solver.Puzzle.Kind.Size()
	groupsTested := Bitset{}
	groups := []*candidateDistribution{}

	for _, cell := range solver.Unsolved {
		cellGroup := cell.GetGroup(groupIndex)

		if !groupsTested.Has(cellGroup) {
			groupsTested.Set(cellGroup, true)
			dist := newDistribution(size)
			dist.reset(solver.Group(groupIndex, cell))
			groups = append(groups, &dist)
		}
	}

	return groups
}

// ==================================================
// Step: 2-String Kite Candidates
//		http://hodoku.sourceforge.net/en/tech_sdp.php
// ==================================================

var Step2StringKite = &SolveStep{
	Technique:      "2-String Kite",
	FirstCost:      2800,
	SubsequentCost: 1600,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		return 0, do2StringKite(solver, limits, step) > 0
	},
}

// Concentrate again on one digit.
// Find a row and a column that have only two candidates left (the "strings").
// One candidate from the row and one candidate from the column have to be in the same block.
// The candidate that sees the two other cells can be eliminated.
func do2StringKite(solver *Solver, limits SolveLimit, step *SolveStep) int {
	size := solver.Puzzle.Kind.Size()
	removed := 0
	rows := getGroupCandidateDistributions(solver, GroupRow)
	cols := getGroupCandidateDistributions(solver, GroupCol)

	for candidate := 0; candidate < size; candidate++ {
		for _, row := range rows {
			rowCands := row.candidates[candidate]
			for _, col := range cols {
				colCands := col.candidates[candidate]

				if rowCands.size == 2 && colCands.size == 2 {
					r0 := rowCands.cells[0]
					r1 := rowCands.cells[1]
					c0 := colCands.cells[0]
					c1 := colCands.cells[1]

					if r0.InGroup(c0) {
						removed += removeCandidateInGroups(solver, limits, step, candidate+1, r1, c1)
					} else if r0.InGroup(c1) {
						removed += removeCandidateInGroups(solver, limits, step, candidate+1, r1, c0)
					} else if r1.InGroup(c0) {
						removed += removeCandidateInGroups(solver, limits, step, candidate+1, r0, c1)
					} else if r1.InGroup(c1) {
						removed += removeCandidateInGroups(solver, limits, step, candidate+1, r0, c0)
					}
					if !solver.CanContinueStep(limits, step) {
						return removed
					}
				}
			}
		}
	}

	return removed
}

// ==================================================
// Step: Empty Rectangle
//		http://hodoku.sourceforge.net/en/tech_sdp.php#er
// ==================================================

var StepEmptyRectangle = &SolveStep{
	Technique:      "Empty Rectangle",
	FirstCost:      2800,
	SubsequentCost: 1600,
	Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
		return 0, doEmptyRectangle(solver, limits, step) > 0
	},
}

// Concentrate again on one digit.
// Find a row and a column that have only two candidates left (the "strings").
// One candidate from the row and one candidate from the column have to be in the same block.
// The candidate that sees the two other cells can be eliminated.
func doEmptyRectangle(solver *Solver, limits SolveLimit, step *SolveStep) int {
	removed := 0

	boxTested := Bitset{}

	for _, cell := range solver.Unsolved {

		if !boxTested.Has(cell.Box) {
			boxTested.Set(cell.Box, true)

			getEmptyRectangles(solver.Box(cell.Box), func(candidate, col, row int) bool {
				can := false
				can = findPerpendicularPair(solver, candidate, row, GroupRow, cell.Box, func(groupFound, otherGroup *Cell) bool {
					dual := countCandidateInGroup(solver, candidate, otherGroup.Row, GroupRow) == 2
					check := removeCandidate(solver, limits, step, col, otherGroup.Row, candidate)
					if check {
						removed++
						if !solver.CanContinueStep(limits, step) {
							return false
						}
						if dual {
							check = removeCandidate(solver, limits, step, groupFound.Col, groupFound.Row, candidate)
							if check {
								removed++
								if !solver.CanContinueStep(limits, step) {
									return false
								}
							}
						}
					}
					return true
				})
				if !can {
					return false
				}
				can = findPerpendicularPair(solver, candidate, col, GroupCol, cell.Box, func(groupFound, otherGroup *Cell) bool {
					dual := countCandidateInGroup(solver, candidate, otherGroup.Col, GroupCol) == 2
					check := removeCandidate(solver, limits, step, otherGroup.Col, row, candidate)
					if check {
						removed++
						if !solver.CanContinueStep(limits, step) {
							return false
						}
						if dual {
							check = removeCandidate(solver, limits, step, groupFound.Col, groupFound.Row, candidate)
							if check {
								removed++
								if !solver.CanContinueStep(limits, step) {
									return false
								}
							}
						}
					}
					return true
				})
				return can
			})
		}
	}

	return removed
}

func getEmptyRectangles(box []*Cell, onEmptyRectangle func(candidate int, col int, row int) bool) {
	remaining := Candidates{}
	minRow := box[0].Row
	maxRow := box[0].Row
	minCol := box[0].Col
	maxCol := box[0].Col

	for _, cell := range box {
		remaining.Or(cell.candidates)
		minRow = Min(minRow, cell.Row)
		maxRow = Max(maxRow, cell.Row)
		minCol = Min(minCol, cell.Col)
		maxCol = Max(maxCol, cell.Col)
	}

	for remaining.Count > 0 {
		candidate := remaining.First()
		remaining.Set(candidate, false)

		for row := minRow; row <= maxRow; row++ {
			for col := minCol; col <= maxCol; col++ {
				rowCount := 0
				colCount := 0

				for _, cell := range box {
					if cell.HasCandidate(candidate) {
						inRow := cell.Row == row
						inCol := cell.Col == col

						if inRow && inCol {
							// don't count
						} else if inRow {
							rowCount++
						} else if inCol {
							colCount++
						} else {
							rowCount = 0
							colCount = 0
							break
						}
					}
				}

				if rowCount > 0 && colCount > 0 {
					if !onEmptyRectangle(candidate, col, row) {
						return
					}
				}
			}
		}
	}
}

func findPerpendicularPair(solver *Solver, candidate int, groupSearch int, groupType Group, notBox int, onPair func(groupFound *Cell, otherGroup *Cell) bool) bool {
	for _, cell := range solver.Unsolved {

		if cell.Box != notBox && cell.GetGroup(groupType) == groupSearch && cell.HasCandidate(candidate) {
			var match *Cell
			matches := 0
			for _, other := range solver.Group(1-groupType, cell) {
				if other.Id == cell.Id {
					continue
				}
				if other.HasCandidate(candidate) {
					match = other
					matches++
				}
			}
			if matches == 1 && match.Box != notBox && match.Box != cell.Box { // the algorithm doesn't require a different box, but it doesn't work in same box
				if !onPair(cell, match) {
					return false
				}
			}
		}
	}

	return true
}

func removeCandidate(solver *Solver, limit SolveLimit, step *SolveStep, col int, row int, candidate int) bool {
	inter := solver.Puzzle.Get(col, row)
	canRemove := inter.Empty() && inter.HasCandidate(candidate)
	if canRemove {
		solver.LogStep(step)
		solver.LogBefore(inter)
		inter.RemoveCandidate(candidate)
		solver.LogAfter(inter)
	}
	return canRemove
}

func countCandidateInGroup(solver *Solver, candidate int, groupSearch int, groupType Group) int {
	for _, cell := range solver.Unsolved {

		if cell.GetGroup(groupType) == groupSearch {
			searchGroup := solver.Group(groupType, cell)
			matches := 0

			if cell.HasCandidate(candidate) {
				matches++
			}

			for _, other := range searchGroup {
				if other.Id == cell.Id {
					continue
				}
				if other.HasCandidate(candidate) {
					matches++
				}
			}

			return matches
		}
	}

	return 0
}

// ==================================================
// Step: Remove Naked Subset Candidates
//		http://hodoku.sourceforge.net/en/tech_naked.php
// ==================================================
func CreateStepBasicFish(setSize int, technique string, firstCost int, subsequentCost int) *SolveStep {
	return &SolveStep{
		Technique:      technique,
		FirstCost:      firstCost,
		SubsequentCost: subsequentCost,
		Logic: func(solver *Solver, limits SolveLimit, step *SolveStep) (int, bool) {
			removed := false
			if solver.CanContinueStep(limits, step) {
				removed = doBasicFish(solver, limits, step, setSize) > 0
			}
			return 0, removed
		},
	}
}

var StepXWing = CreateStepBasicFish(2, "X-Wing", 2800, 1600)
var StepSwordfish = CreateStepBasicFish(3, "Swordfish", 8000, 6000)
var StepJellyfish = CreateStepBasicFish(4, "Jellyfish", 10000, 8000)

func doBasicFish(solver *Solver, limits SolveLimit, step *SolveStep, setSize int) int {
	removed := 0
	removed += doBasicFishGroups(solver, limits, step, setSize, GroupRow)
	if solver.CanContinueStep(limits, step) {
		removed += doBasicFishGroups(solver, limits, step, setSize, GroupCol)
	}
	return removed
}

func doBasicFishGroups(solver *Solver, limits SolveLimit, step *SolveStep, setSize int, groupType Group) int {
	size := solver.Puzzle.Kind.Size()
	removed := 0
	groups := getGroupCandidateDistributions(solver, groupType)

	for candidateIndex := 0; candidateIndex < size; candidateIndex++ {
		candidate := candidateIndex + 1
		matched := make([]candidateCells, 0)
		for _, group := range groups {
			groupCandidates := group.candidates[candidateIndex]
			if groupCandidates.size > 0 && groupCandidates.size <= setSize {
				matched = append(matched, groupCandidates)
			}
		}
		if len(matched) >= setSize {
			combs := combinations[candidateCells]{}
			for i := 0; i < setSize; i++ {
				combs.add(matched)
			}

			set := combs.start()
			for combs.next(set) {
				rowsHit := Bitset{}
				colsHit := Bitset{}
				for _, candCells := range set {
					for k := 0; k < candCells.size; k++ {
						c := candCells.cells[k]
						rowsHit.Set(c.Row, true)
						colsHit.Set(c.Col, true)
					}
				}
				if rowsHit.Count == setSize && colsHit.Count == setSize {
					for cellIndex := range solver.Puzzle.Cells {
						cell := &solver.Puzzle.Cells[cellIndex]
						inColumn := colsHit.Has(cell.Col)
						inRow := rowsHit.Has(cell.Row)
						if cell.HasCandidate(candidate) && ((inRow && !inColumn) || (!inRow && inColumn)) {
							solver.LogStep(step)
							solver.LogBefore(cell)
							cell.RemoveCandidate(candidate)
							removed++
							solver.LogAfter(cell)

							if !solver.CanContinueStep(limits, step) {
								return removed
							}
						}
					}
					break
				}
			}
		}
	}

	return removed
}
