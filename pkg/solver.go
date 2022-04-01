package sudogo

import "fmt"

type CellGroups struct {
	Cell        *Cell
	All         []*Cell
	Box         []*Cell
	Row         []*Cell
	Col         []*Cell
	Constraints []Constraint
}

func (group *CellGroups) Remove(neighbor *Cell) {
	group.All = removeValue(group.All, neighbor)
	group.Box = removeValue(group.Box, neighbor)
	group.Row = removeValue(group.Row, neighbor)
	group.Col = removeValue(group.Col, neighbor)
}

type Solver struct {
	Puzzle   Puzzle
	Steps    []*SolveStep
	Cells    []CellGroups
	Unsolved []*CellGroups

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

type SolverLimit struct {
	MinCost       int
	MaxCost       int
	MaxPlacements int
	MaxLogs       int
	MaxBatches    int
}

type SolveStepLogic func(solver *Solver, limits SolverLimit, step *SolveStep) (placements int, restart bool)

type SolveStep struct {
	Technique      string
	FirstCost      int
	SubsequentCost int
	Logic          SolveStepLogic
}

var StandardSolveSteps = []*SolveStep{
	StepNakedSingle,
	StepHiddenSingle,
	StepRemovePointingCandidates,
	StepRemoveClaimingCandidates,
	StepConstraints,
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
	StepConstraints,
	StepRemoveNakedSubsetCandidates2,
	StepRemoveNakedSubsetCandidates3,
}

func NewSolver(starting Puzzle) Solver {
	puzzle := starting.Clone()
	cells := make([]CellGroups, puzzle.Kind.Area())
	unsolved := make([]*CellGroups, 0, puzzle.Kind.Area())
	groupCapacity := puzzle.Kind.Digits()
	allCapacity := groupCapacity * 3
	constraints := puzzle.Kind.Constraints
	constraintCapacity := len(constraints)

	for i := range puzzle.Cells {
		cell := &puzzle.Cells[i]
		group := &cells[i]

		group.Cell = cell
		group.All = make([]*Cell, 0, allCapacity)
		group.Box = make([]*Cell, 0, groupCapacity)
		group.Row = make([]*Cell, 0, groupCapacity)
		group.Col = make([]*Cell, 0, groupCapacity)
		group.Constraints = make([]Constraint, 0, constraintCapacity)
		if cell.Empty() {
			unsolved = append(unsolved, group)
			for k := range constraints {
				constrain := constraints[k]
				if constrain.Affects(cell) {
					group.Constraints = append(group.Constraints, constrain)
				}
			}
		}
		for k := range puzzle.Cells {
			other := &puzzle.Cells[k]
			if i != k && other.Empty() {
				if cell.InGroup(other) {
					group.All = append(group.All, other)
				}
				if cell.InBox(other) {
					group.Box = append(group.Box, other)
				}
				if cell.InRow(other) {
					group.Row = append(group.Row, other)
				}
				if cell.InColumn(other) {
					group.Col = append(group.Col, other)
				}
			}
		}
	}

	return Solver{
		Puzzle:        puzzle,
		Steps:         StandardSolveSteps,
		Cells:         cells,
		Unsolved:      unsolved,
		LogEnabled:    false,
		LogState:      false,
		LogTechniques: map[string]int{},
		Logs:          []SolverLog{},
		logTemplate:   SolverLog{},
	}
}

func (solver *Solver) Set(col int, row int, value int) bool {
	return solver.SetCell(solver.Puzzle.Get(col, row), value)
}

func (solver *Solver) SetCell(cell *Cell, value int) bool {
	if cell == nil {
		return false
	}
	return solver.SetGroup(&solver.Cells[cell.Id], value)
}

func (solver *Solver) SetGroup(group *CellGroups, value int) bool {
	if group == nil || value <= 0 {
		return false
	}
	set := group.Cell.SetValue(value)
	if set {
		for _, other := range group.All {
			other.RemoveCandidate(value)
		}

		solver.Unsolved = removeValue(solver.Unsolved, group)

		for _, remaining := range solver.Unsolved {
			remaining.Remove(group.Cell)
		}
	}
	return set
}

func (solver *Solver) GetMinCandidateCount() int {
	min := 0
	for _, group := range solver.Unsolved {
		if min == 0 || min > group.Cell.candidates.Count {
			min = group.Cell.candidates.Count
		}
	}
	return min
}

func (solver *Solver) GetGroupWhere(where func(group *CellGroups) bool) *CellGroups {
	for _, group := range solver.Unsolved {
		if where(group) {
			return group
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

func (solver *Solver) canContinue(limits SolverLimit, cost int) bool {
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

func (solver *Solver) canContinueStep(limits SolverLimit, step *SolveStep) bool {
	return solver.canContinue(limits, solver.getCost(step))
}

func (solver *Solver) getCost(step *SolveStep) int {
	techniqueCount := solver.LogTechniques[step.Technique]
	cost := step.FirstCost
	if techniqueCount > 0 {
		cost = step.SubsequentCost
	}
	return cost
}

func (solver *Solver) logStep(step *SolveStep) {
	cost := solver.getCost(step)
	solver.LogTechniques[step.Technique]++
	solver.logTemplate.Batch++
	solver.logTemplate.Step = step
	solver.logTemplate.Cost = cost
	solver.logTemplate.RunningCost += cost
}

func (solver *Solver) logBefore(before *Cell) {
	if solver.LogEnabled {
		log := solver.logTemplate
		log.Before = *before
		solver.Logs = append(solver.Logs, log)
	}
	solver.logTemplate.Index++
}

func (solver *Solver) logAfter(after *Cell) {
	last := solver.GetLastLog()
	last.After = *after
	if solver.LogState {
		state := solver.Puzzle.Clone()
		last.State = &state
	}
}

func (solver *Solver) logPlacement(after *Cell) {
	last := solver.GetLastLog()
	last.After = *after
	last.Placement = true
	last.RunningPlacements = last.RunningPlacements + 1
	solver.logTemplate.RunningPlacements = last.RunningPlacements
}

func (solver *Solver) Solved() bool {
	return len(solver.Unsolved) == 0
}

func (solver *Solver) Solve(limits SolverLimit) (*Puzzle, bool) {
	steps := solver.Steps
	placing := true
	for placing {
		placing = false
		for _, step := range steps {
			stepPlaced, stepRestart := step.Logic(solver, limits, step)

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
	return &solver.Puzzle, solver.Solved()
}

// ==================================================
// Step: Naked Single
// 		http://hodoku.sourceforge.net/en/tech_singles.php
// ==================================================
var StepNakedSingle = &SolveStep{
	Technique:      "Naked Single",
	FirstCost:      100,
	SubsequentCost: 100,
	Logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
		placements := 0
		for solver.canContinueStep(limits, step) {
			group, groupValue := getNakedSingle(solver)
			if group != nil {
				solver.logStep(step)
				solver.logBefore(group.Cell)
				solver.SetGroup(group, groupValue)
				solver.logPlacement(group.Cell)
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
	for _, group := range solver.Unsolved {
		if group.Cell.candidates.Count == 1 {
			return group, group.Cell.FirstCandidate()
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
	Logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
		placements := 0
		for solver.canContinueStep(limits, step) {
			group, groupValue := getHiddenSingle(solver)
			if group != nil {
				solver.logStep(step)
				solver.logBefore(group.Cell)
				solver.SetGroup(group, groupValue)
				solver.logPlacement(group.Cell)
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
	for _, group := range solver.Unsolved {
		box := getHiddenSingleFromGroup(group.Cell, group.Box)
		if box != 0 {
			return group, box
		}
		row := getHiddenSingleFromGroup(group.Cell, group.Row)
		if row != 0 {
			return group, row
		}
		col := getHiddenSingleFromGroup(group.Cell, group.Col)
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
// Step: Constraints
// ==================================================
var StepConstraints = &SolveStep{
	Technique:      "Constraints",
	FirstCost:      0,
	SubsequentCost: 0,
	Logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
		removed := 0

		for _, group := range solver.Unsolved {
			if len(group.Constraints) == 0 {
				continue
			}

			cell := group.Cell
			candidates := cell.candidates

			for _, constraint := range group.Constraints {
				constraint.RemoveCandidates(cell, &solver.Puzzle, &candidates)
			}

			if candidates.Value != cell.candidates.Value {
				solver.logStep(step)
				solver.logBefore(group.Cell)
				removed += cell.candidates.Count - candidates.Count
				cell.candidates = candidates
				solver.logPlacement(group.Cell)

				if !solver.canContinueStep(limits, step) {
					break
				}
			}
		}

		return 0, removed > 0
	},
}

// ==================================================
// Step: Remove Pointing Candidates
//		http://hodoku.sourceforge.net/en/tech_intersections.php
// ==================================================
var StepRemovePointingCandidates = &SolveStep{
	Technique:      "Pointing Candidates",
	FirstCost:      350,
	SubsequentCost: 200,
	Logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
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

	for _, group := range solver.Unsolved {
		cell := group.Cell
		// all candidates in this box's row that are shared
		row := cell.candidates
		// all candidates in this box's column that are shared
		col := cell.candidates

		// remove candidates that are not shared
		for _, other := range group.Box {
			if other.Row == cell.Row {
				row.And(other.candidates)
			}
			if other.Col == cell.Col {
				col.And(other.candidates)
			}
		}

		// remove candidates that exist outside the row or column
		for _, other := range group.Box {
			if other.Row != cell.Row {
				row.Remove(other.candidates)
			}
			if other.Col != cell.Col {
				col.Remove(other.candidates)
			}
		}

		// what is remaining are candidates confined to the cells row in the box
		if row.Count > 0 {
			hasOverlap := false
			for _, other := range group.Row {
				if other.Box != cell.Box && other.candidates.Overlaps(row) {
					hasOverlap = true
					break
				}
			}
			if hasOverlap {
				solver.logStep(step)
				for _, other := range group.Row {
					if other.Box != cell.Box && other.candidates.Overlaps(row) {
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
			for _, other := range group.Col {
				if other.Box != cell.Box && other.candidates.Overlaps(row) {
					hasOverlap = true
					break
				}
			}
			if hasOverlap {
				solver.logStep(step)
				for _, other := range group.Col {
					if other.Box != cell.Box && other.candidates.Overlaps(row) {
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
	Technique:      "Claiming Candidates",
	FirstCost:      350,
	SubsequentCost: 200,
	Logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
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

	for _, group := range solver.Unsolved {
		cell := group.Cell

		// all candidates in this row that are not shared outside of the box
		row := cell.candidates

		// remove candidates from row that exist in the cells row outside the box
		for _, other := range group.Row {
			if other.Box != cell.Box {
				row.Remove(other.candidates)
			}
		}

		// what is remaining are the candidates unique to the row outside this box
		if row.Count > 0 {
			solver.logStep(step)
			for _, other := range group.Box {
				if other.Row != cell.Row && other.candidates.Overlaps(row) {
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
		for _, other := range group.Col {
			if other.Box != cell.Box {
				col.Remove(other.candidates)
			}
		}

		// what is remaining are the candidates unique to the column outside this box
		if col.Count > 0 {
			solver.logStep(step)
			for _, other := range group.Box {
				if other.Col != cell.Col && other.candidates.Overlaps(col) {
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
		Technique:      technique,
		FirstCost:      firstCost,
		SubsequentCost: subsequentCost,
		Logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
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

	for _, group := range solver.Unsolved {
		cell := group.Cell

		if cell.candidates.Count != subsetSize {
			continue
		}
		removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, solver, limits, step, group.Row)
		if !solver.canContinueStep(limits, step) {
			break
		}
		removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, solver, limits, step, group.Col)
		if !solver.canContinueStep(limits, step) {
			break
		}
		removed += removeNakedSubsetCandidatesFromGroup(group, subsetSize, solver, limits, step, group.Box)
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
	candidates := cellGroup.Cell.candidates
	sameBox := true
	sameRow := true
	sameCol := true

	for _, other := range group {
		if other.candidates.Value == candidates.Value {
			matches++
			sameBox = sameBox && other.Box == cellGroup.Cell.Box
			sameRow = sameRow && other.Row == cellGroup.Cell.Row
			sameCol = sameCol && other.Col == cellGroup.Cell.Col
		}
	}

	if matches == subsetSize {
		if sameBox {
			removed += removeCandidatesFromDifferent(cellGroup.Box, candidates, solver, limits, step)
		}
		if sameRow && solver.canContinueStep(limits, step) {
			removed += removeCandidatesFromDifferent(cellGroup.Row, candidates, solver, limits, step)
		}
		if sameCol && solver.canContinueStep(limits, step) {
			removed += removeCandidatesFromDifferent(cellGroup.Col, candidates, solver, limits, step)
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
		Technique:      technique,
		FirstCost:      firstCost,
		SubsequentCost: subsequentCost,
		Logic: func(solver *Solver, limits SolverLimit, step *SolveStep) (int, bool) {
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
	dist := newDistribution(solver.Puzzle.Kind.Size())
	rowsTested := Bitset{}
	colsTested := Bitset{}
	boxsTested := Bitset{}
	removed := 0

	for _, group := range solver.Unsolved {
		cell := group.Cell

		// Only test a row/column/box of an unsolved cell once
		if !rowsTested.Has(cell.Row) {
			rowsTested.Set(cell.Row, true)
			fullRow := group.Row[:]
			fullRow = append(fullRow, cell)
			dist.set(fullRow)
			removed += doRemoveHiddenSubset(&dist, subsetSize, solver, limits, step)
			if !solver.canContinueStep(limits, step) {
				break
			}
		}
		if !colsTested.Has(cell.Col) {
			colsTested.Set(cell.Col, true)
			fullCol := group.Col[:]
			fullCol = append(fullCol, cell)
			dist.set(fullCol)
			removed += doRemoveHiddenSubset(&dist, subsetSize, solver, limits, step)
			if !solver.canContinueStep(limits, step) {
				break
			}
		}
		if !boxsTested.Has(cell.Box) {
			boxsTested.Set(cell.Box, true)
			fullBox := group.Box[:]
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
