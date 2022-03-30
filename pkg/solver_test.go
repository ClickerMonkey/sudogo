package sudogo

import (
	"fmt"
	"testing"
	"time"
)

func TestSolveSimple(t *testing.T) {
	original := Classic.Create([][]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	})

	s := original.Solver()
	solution, solved := s.Solve(SolverLimit{})

	if !solved {
		solution.PrintCandidates()
		t.Errorf("Failed to solve TestSolveSimple")
	}

	checkValid(solution, t)
}

func TestSolveHiddenSingle(t *testing.T) {
	original := Classic.Create([][]int{
		{0, 2, 8, 0, 0, 7, 0, 0, 0},
		{0, 1, 6, 0, 8, 3, 0, 7, 0},
		{0, 0, 0, 0, 2, 0, 8, 5, 1},
		{1, 3, 7, 2, 9, 0, 0, 0, 0},
		{0, 0, 0, 7, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 4, 6, 3, 0, 7},
		{2, 9, 0, 0, 7, 0, 0, 0, 0},
		{0, 0, 0, 8, 6, 0, 1, 4, 0},
		{0, 0, 0, 3, 0, 0, 7, 0, 0},
	})

	s := original.Solver()
	solution, solved := s.Solve(SolverLimit{})

	if solution.Get(3, 2).value != 6 {
		t.Errorf("The solver failed to use hidden single logic on r3c4.")
	}

	if !solved {
		solution.PrintCandidates()
		t.Fatal("The puzzle could no longer be solved.")
	}

	checkValid(solution, t)
}

func TestPointing(t *testing.T) {
	original := Classic.Create([][]int{
		{9, 8, 4, 0, 0, 0, 0, 0, 0},
		{0, 0, 2, 5, 0, 0, 0, 4, 0},
		{0, 0, 1, 9, 0, 4, 0, 0, 2},
		{0, 0, 6, 0, 9, 7, 2, 3, 0},
		{0, 0, 3, 6, 0, 2, 0, 0, 0},
		{2, 0, 9, 0, 3, 5, 6, 1, 0},
		{1, 9, 5, 7, 6, 8, 4, 2, 3},
		{4, 2, 7, 3, 5, 1, 8, 9, 6},
		{6, 3, 8, 0, 0, 9, 7, 5, 1},
	})

	s := original.Solver()
	p := &s.puzzle

	if fmt.Sprint(p.Get(6, 2).Candidates()) != "[3 5]" {
		t.Errorf("Invalid candidates for r3c7 in Pointing")
	}

	StepRemovePointingCandidates.logic(&s, SolverLimit{}, StepRemovePointingCandidates)

	if fmt.Sprint(p.Get(6, 2).Candidates()) != "[3]" {
		t.Errorf("Invalid candidates for r3c7 in Pointing after step")
	}

	solution, solved := s.Solve(SolverLimit{})

	if solution.Get(6, 2).value != 3 {
		t.Errorf("Test Pointing solve failed")
	}

	if !solved {
		solution.PrintCandidates()
		t.Fatal("Puzzle not solved")
	}

	checkValid(solution, t)
}

func TestClaiming(t *testing.T) {
	original := Classic.Create([][]int{
		{3, 1, 8, 0, 0, 5, 4, 0, 6},
		{0, 0, 0, 6, 0, 3, 8, 1, 0},
		{0, 0, 6, 0, 8, 0, 5, 0, 3},
		{8, 6, 4, 9, 5, 2, 1, 3, 7},
		{1, 2, 3, 4, 7, 6, 9, 5, 8},
		{7, 9, 5, 3, 1, 8, 2, 6, 4},
		{0, 3, 0, 5, 0, 0, 7, 8, 0},
		{0, 0, 0, 0, 0, 7, 3, 0, 5},
		{0, 0, 0, 0, 3, 9, 6, 4, 1},
	})

	s := original.Solver()
	p := &s.puzzle

	if fmt.Sprint(p.Get(1, 2).Candidates()) != "[4 7]" {
		t.Errorf("Invalid candidates for r3c2 in Claiming")
	}

	StepRemoveClaimingCandidates.logic(&s, SolverLimit{}, StepRemoveClaimingCandidates)

	if fmt.Sprint(p.Get(1, 2).Candidates()) != "[4]" {
		t.Errorf("Invalid candidates for r3c2 in Claiming after step")
	}

	solution, solved := s.Solve(SolverLimit{})

	if solution.Get(1, 2).value != 4 {
		t.Errorf("Test Claiming solve failed")
	}

	if !solved {
		solution.PrintCandidates()
		t.Fatal("Puzzle not solved")
	}

	checkValid(solution, t)
}

func TestNakedPair(t *testing.T) {
	original := Classic.Create([][]int{
		{7, 0, 0, 8, 4, 9, 0, 3, 0},
		{9, 2, 8, 1, 3, 5, 0, 0, 6},
		{4, 0, 0, 2, 6, 7, 0, 8, 9},
		{6, 4, 2, 7, 8, 3, 9, 5, 1},
		{3, 9, 7, 4, 5, 1, 6, 2, 8},
		{8, 1, 5, 6, 9, 2, 3, 0, 0},
		{2, 0, 4, 5, 1, 6, 0, 9, 3},
		{1, 0, 0, 0, 0, 8, 0, 6, 0},
		{5, 0, 0, 0, 0, 4, 0, 1, 0},
	})

	s := original.Solver()
	p := &s.puzzle

	r8c2 := p.Get(1, 7)

	c0 := fmt.Sprint(r8c2.Candidates())
	if c0 != "[3 7]" {
		t.Fatalf("Candidates for r8c2 are not [3 7] they are %s", c0)
	}

	StepRemoveNakedSubsetCandidates2.logic(&s, SolverLimit{}, StepRemoveNakedSubsetCandidates2)

	c1 := fmt.Sprint(r8c2.Candidates())
	if c1 != "[7]" {
		t.Fatalf("Candidates for r8c2 are not [7] they are %s", c1)
	}

	checkValid(p, t)
}

func TestHiddenPair(t *testing.T) {
	type CandidateTest struct {
		column int
		row    int
		before string
		after  string
	}

	tests := []struct {
		puzzle Puzzle
		step   *SolveStep
		max    int
		tests  []CandidateTest
	}{
		{
			puzzle: Classic.Create([][]int{
				{0, 4, 9, 1, 3, 2, 0, 0, 0},
				{0, 8, 1, 4, 7, 9, 0, 0, 0},
				{3, 2, 7, 6, 8, 5, 9, 1, 4},
				{0, 9, 6, 0, 5, 1, 8, 0, 0},
				{0, 7, 5, 0, 2, 8, 0, 0, 0},
				{0, 3, 8, 0, 4, 6, 0, 0, 5},
				{8, 5, 3, 2, 6, 7, 0, 0, 0},
				{7, 1, 2, 8, 9, 4, 5, 6, 3},
				{9, 6, 4, 5, 1, 3, 0, 0, 0},
			}),
			step: StepRemoveHiddenSubsetCandidates2,
			max:  19,
			tests: []CandidateTest{
				{
					column: 8,
					row:    4,
					before: "[1 6 9]",
					after:  "[1 9]",
				},
				{
					column: 8,
					row:    6,
					before: "[1 9]",
					after:  "[1 9]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 0, 0, 6, 0, 0, 0, 0},
				{0, 0, 0, 0, 4, 2, 7, 3, 6},
				{0, 0, 6, 7, 3, 0, 0, 4, 0},
				{0, 9, 4, 0, 0, 0, 0, 6, 8},
				{0, 0, 0, 0, 9, 6, 4, 0, 7},
				{6, 0, 7, 0, 5, 0, 9, 2, 3},
				{1, 0, 0, 0, 0, 0, 0, 8, 5},
				{0, 6, 0, 0, 8, 0, 2, 7, 1},
				{0, 0, 5, 0, 1, 0, 0, 9, 4},
			}),
			step: StepRemoveHiddenSubsetCandidates2,
			max:  1,
			tests: []CandidateTest{
				{
					column: 0,
					row:    0,
					before: "[2 3 4 5 7 8 9]",
					after:  "[4 7]",
				},
				{
					column: 1,
					row:    0,
					before: "[1 2 3 4 5 7 8]",
					after:  "[4 7]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{5, 0, 0, 6, 2, 0, 0, 3, 7},
				{0, 0, 4, 8, 9, 0, 0, 0, 0},
				{0, 0, 0, 0, 5, 0, 0, 0, 0},
				{9, 3, 0, 0, 0, 0, 0, 0, 0},
				{0, 2, 0, 0, 0, 0, 6, 0, 5},
				{7, 0, 0, 0, 0, 0, 0, 0, 3},
				{0, 0, 0, 0, 0, 9, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 7, 0, 0},
				{6, 8, 0, 5, 7, 0, 0, 0, 2},
			}),
			step: StepRemoveHiddenSubsetCandidates3,
			max:  23,
			tests: []CandidateTest{
				{
					column: 5,
					row:    3,
					before: "[1 2 4 5 6 7 8]",
					after:  "[2 5 6]",
				},
				{
					column: 5,
					row:    5,
					before: "[1 2 4 5 6 8]",
					after:  "[2 5 6]",
				},
				{
					column: 5,
					row:    7,
					before: "[1 2 3 4 6 8]",
					after:  "[2 6]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{2, 8, 0, 0, 0, 0, 4, 7, 3},
				{5, 3, 4, 8, 2, 7, 1, 9, 6},
				{0, 7, 1, 0, 3, 4, 0, 8, 0},
				{3, 0, 0, 5, 0, 0, 0, 4, 0},
				{0, 0, 0, 3, 4, 0, 0, 6, 0},
				{4, 6, 0, 7, 9, 0, 3, 1, 0},
				{0, 9, 0, 2, 0, 3, 6, 5, 4},
				{0, 0, 3, 0, 0, 9, 8, 2, 1},
				{0, 0, 0, 0, 8, 0, 9, 3, 7},
			}),
			step: StepRemoveHiddenSubsetCandidates3,
			max:  3,
			tests: []CandidateTest{
				{
					column: 1,
					row:    7,
					before: "[4 5]",
					after:  "[4 5]",
				},
				{
					column: 1,
					row:    8,
					before: "[1 2 4 5]",
					after:  "[2 4 5]",
				},
				{
					column: 2,
					row:    8,
					before: "[2 5 6]",
					after:  "[2 5]",
				},
			},
		},
	}

	for _, test := range tests {
		solver := test.puzzle.Solver()
		puzzle := &solver.puzzle

		for _, cellTest := range test.tests {
			testCell := puzzle.Get(cellTest.column, cellTest.row)
			actual := fmt.Sprint(testCell.Candidates())

			if actual != cellTest.before {
				puzzle.PrintCandidates()
				t.Fatalf("Candidates for [%d,%d] are not %s they are %s", cellTest.column, cellTest.row, cellTest.before, actual)
			}
		}

		removed, _ := test.step.logic(&solver, SolverLimit{maxBatches: test.max}, test.step)

		for _, cellTest := range test.tests {
			testCell := puzzle.Get(cellTest.column, cellTest.row)
			actual := fmt.Sprint(testCell.Candidates())

			if actual != cellTest.after {
				puzzle.PrintCandidates()
				t.Fatalf("Candidates for [%d,%d] are not %s they are %s. %d candidates removed.", cellTest.column, cellTest.row, cellTest.after, actual, removed)
			}
		}

		checkValid(puzzle, t)
	}
}

func TestSolveHard(t *testing.T) {
	// https://www.thonky.com/sudoku/solution-count
	tests := []struct {
		puzzle    Puzzle
		max       int
		solutions int
	}{
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 0, 1, 0, 2, 0, 0, 0},
				{0, 6, 0, 0, 0, 8, 3, 0, 0},
				{5, 0, 0, 0, 0, 0, 0, 0, 9},
				{0, 0, 0, 4, 0, 7, 0, 0, 8},
				{6, 8, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 4, 0, 0, 0, 0, 1, 0},
				{0, 2, 0, 0, 0, 0, 5, 0, 0},
				{0, 0, 0, 0, 7, 1, 2, 0, 6},
				{0, 9, 0, 0, 0, 6, 7, 0, 0},
			}),
			max:       0,
			solutions: 13,
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 0, 1, 0, 2, 0, 0, 0},
				{0, 6, 0, 0, 0, 8, 3, 0, 0},
				{5, 0, 0, 0, 0, 0, 0, 0, 9},
				{0, 0, 0, 4, 0, 7, 0, 0, 8},
				{6, 8, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 4, 0, 0, 0, 0, 1, 0},
				{0, 2, 0, 0, 0, 0, 5, 0, 0},
				{0, 0, 0, 0, 7, 1, 2, 0, 6},
				{0, 9, 0, 0, 0, 6, 7, 8, 0},
			}),
			max:       0,
			solutions: 1,
		},
	}

	for _, test := range tests {
		start := time.Now()
		solutions := test.puzzle.GetSolutions(SolutionLimit{maxSolutions: test.max})
		duration := time.Since(start)

		if len(solutions) != test.solutions {
			t.Errorf("An invalid number of solutions found %d expected %d in %s for %s.", len(solutions), test.solutions, duration, test.puzzle.UniqueId())
		}

		fmt.Printf("Solutions: %d in %s\n", len(solutions), duration)

		for _, solution := range solutions {
			checkValid(&solution.puzzle, t)
		}
	}
}

func checkValid(puzzle *Puzzle, t *testing.T) {
	if !puzzle.IsValid() {
		puzzle.PrintCandidates()
		t.Fatal("The previous puzzle has invalid candidates")
	}
}

func TestLogs(t *testing.T) {
	original := Classic.Create([][]int{
		{9, 8, 4, 0, 0, 0, 0, 0, 0},
		{0, 0, 2, 5, 0, 0, 0, 4, 0},
		{0, 0, 1, 9, 0, 4, 0, 0, 2},
		{0, 0, 6, 0, 9, 7, 2, 3, 0},
		{0, 0, 3, 6, 0, 2, 0, 0, 0},
		{2, 0, 9, 0, 3, 5, 6, 1, 0},
		{1, 9, 5, 7, 6, 8, 4, 2, 3},
		{4, 2, 7, 3, 5, 1, 8, 9, 6},
		{6, 3, 8, 0, 0, 9, 7, 5, 1},
	})

	solver := original.Solver()
	solver.logEnabled = true
	solver.Solve(SolverLimit{})

	fmt.Println("TestLogs")
	printSolveLogs(&solver)
}

func printSolveLogs(solver *Solver) {
	for _, log := range solver.logs {
		if log.placement {
			fmt.Printf("%s placed %d at r%dc%d (cost=%d)\n", log.step.technique, log.after.value, log.after.row+1, log.after.col+1, log.cost)
		} else {
			fmt.Printf("%s candidates set %v from %v at r%dc%d (cost=%d)\n", log.step.technique, log.before.candidates.ToSlice(), log.after.candidates.ToSlice(), log.before.row+1, log.before.col+1, log.cost)
		}
	}
	last := solver.GetLastLog()
	fmt.Printf("Total cost = %d, placements = %d, batches = %d, logs = %d.\n", last.runningCost, last.runningPlacements, last.batch, last.index)
}
