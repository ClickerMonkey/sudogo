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
		solution.PrintConsoleCandidates()
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

	if solution.Get(3, 2).Value != 6 {
		t.Errorf("The solver failed to use hidden single logic on r3c4.")
	}

	if !solved {
		solution.PrintConsoleCandidates()
		t.Fatal("The puzzle could no longer be solved.")
	}

	checkValid(solution, t)
}

func TestCandidateRemoveSteps(t *testing.T) {
	type CandidateTest struct {
		column int
		row    int
		before string
		after  string
		value  int
	}

	tests := []struct {
		puzzle Puzzle
		step   *SolveStep
		max    int
		tests  []CandidateTest
		solve  bool
	}{
		{
			puzzle: Classic.Create([][]int{
				{0, 2, 8, 0, 0, 7, 0, 0, 0},
				{0, 1, 6, 0, 8, 3, 0, 7, 0},
				{0, 0, 0, 0, 2, 0, 8, 5, 1},
				{1, 3, 7, 2, 9, 0, 0, 0, 0},
				{0, 0, 0, 7, 3, 0, 0, 0, 0},
				{0, 0, 0, 0, 4, 6, 3, 0, 7},
				{2, 9, 0, 0, 7, 0, 0, 0, 0},
				{0, 0, 0, 8, 6, 0, 1, 4, 0},
				{0, 0, 0, 3, 0, 0, 7, 0, 0},
			}),
			step: StepHiddenSingle,
			max:  1,
			tests: []CandidateTest{
				{
					column: 3,
					row:    2,
					before: "[4 6 9]",
					after:  "[]",
					value:  6,
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{9, 8, 4, 0, 0, 0, 0, 0, 0},
				{0, 0, 2, 5, 0, 0, 0, 4, 0},
				{0, 0, 1, 9, 0, 4, 0, 0, 2},
				{0, 0, 6, 0, 9, 7, 2, 3, 0},
				{0, 0, 3, 6, 0, 2, 0, 0, 0},
				{2, 0, 9, 0, 3, 5, 6, 1, 0},
				{1, 9, 5, 7, 6, 8, 4, 2, 3},
				{4, 2, 7, 3, 5, 1, 8, 9, 6},
				{6, 3, 8, 0, 0, 9, 7, 5, 1},
			}),
			step: StepRemovePointingCandidates,
			max:  1,
			tests: []CandidateTest{
				{
					column: 6,
					row:    2,
					before: "[3 5]",
					after:  "[3]",
				},
			},
			solve: true,
		},
		{
			puzzle: Classic.Create([][]int{
				{3, 1, 8, 0, 0, 5, 4, 0, 6},
				{0, 0, 0, 6, 0, 3, 8, 1, 0},
				{0, 0, 6, 0, 8, 0, 5, 0, 3},
				{8, 6, 4, 9, 5, 2, 1, 3, 7},
				{1, 2, 3, 4, 7, 6, 9, 5, 8},
				{7, 9, 5, 3, 1, 8, 2, 6, 4},
				{0, 3, 0, 5, 0, 0, 7, 8, 0},
				{0, 0, 0, 0, 0, 7, 3, 0, 5},
				{0, 0, 0, 0, 3, 9, 6, 4, 1},
			}),
			step: StepRemoveClaimingCandidates,
			max:  -1,
			tests: []CandidateTest{
				{
					column: 1,
					row:    2,
					before: "[4 7]",
					after:  "[4]",
				},
			},
			solve: true,
		},
		{
			puzzle: Classic.Create([][]int{
				{7, 0, 0, 8, 4, 9, 0, 3, 0},
				{9, 2, 8, 1, 3, 5, 0, 0, 6},
				{4, 0, 0, 2, 6, 7, 0, 8, 9},
				{6, 4, 2, 7, 8, 3, 9, 5, 1},
				{3, 9, 7, 4, 5, 1, 6, 2, 8},
				{8, 1, 5, 6, 9, 2, 3, 0, 0},
				{2, 0, 4, 5, 1, 6, 0, 9, 3},
				{1, 0, 0, 0, 0, 8, 0, 6, 0},
				{5, 0, 0, 0, 0, 4, 0, 1, 0},
			}),
			step: StepRemoveNakedSubsetCandidates2,
			max:  1,
			tests: []CandidateTest{
				{
					column: 1,
					row:    7,
					before: "[3 7]",
					after:  "[7]",
				},
			},
		},
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
		{
			puzzle: Classic.Create([][]int{
				{6, 9, 7, 0, 0, 0, 0, 0, 2},
				{0, 0, 1, 9, 7, 2, 0, 6, 3},
				{0, 0, 3, 0, 0, 6, 7, 9, 0},
				{9, 1, 2, 0, 0, 0, 6, 0, 7},
				{3, 7, 4, 2, 6, 0, 9, 5, 0},
				{8, 6, 5, 7, 0, 9, 0, 2, 4},
				{1, 4, 8, 6, 9, 3, 2, 7, 5},
				{7, 0, 9, 0, 2, 4, 0, 0, 6},
				{0, 0, 6, 8, 0, 7, 0, 0, 9},
			}),
			step: StepRemoveSkyscraperCandidates,
			max:  4,
			tests: []CandidateTest{
				{
					column: 3,
					row:    2,
					before: "[1 4 5]",
					after:  "[4 5]",
				},
				{
					column: 4,
					row:    2,
					before: "[1 4 5 8]",
					after:  "[4 5 8]",
				},
				{
					column: 6,
					row:    0,
					before: "[1 4 5 8]",
					after:  "[4 5 8]",
				},
				{
					column: 7,
					row:    0,
					before: "[1 4 8]",
					after:  "[4 8]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 1, 0, 2, 8, 7, 5, 9},
				{0, 8, 7, 9, 0, 5, 1, 3, 2},
				{9, 5, 2, 1, 7, 3, 4, 8, 6},
				{0, 2, 0, 7, 0, 0, 3, 4, 0},
				{0, 0, 0, 5, 0, 0, 2, 7, 0},
				{7, 1, 4, 8, 3, 2, 6, 9, 5},
				{0, 0, 0, 0, 9, 0, 8, 1, 7},
				{0, 7, 8, 0, 5, 1, 9, 6, 3},
				{1, 9, 0, 0, 8, 7, 5, 2, 4},
			}),
			step: StepRemoveSkyscraperCandidates,
			max:  1,
			tests: []CandidateTest{
				{
					column: 3,
					row:    0,
					before: "[4 6]",
					after:  "[6]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 8, 1, 0, 2, 0, 6, 0, 0},
				{0, 4, 2, 0, 6, 0, 0, 8, 9},
				{0, 5, 6, 8, 0, 0, 2, 4, 0},
				{6, 9, 3, 1, 4, 2, 7, 5, 8},
				{4, 2, 8, 3, 5, 7, 9, 1, 6},
				{1, 7, 5, 6, 8, 9, 3, 2, 4},
				{5, 1, 0, 0, 3, 6, 8, 9, 2},
				{2, 3, 0, 0, 0, 8, 4, 6, 0},
				{8, 6, 0, 2, 0, 0, 0, 0, 0},
			}),
			step: StepRemove2StringKiteCandidates,
			max:  1,
			tests: []CandidateTest{
				{
					column: 3,
					row:    1,
					before: "[5 7]",
					after:  "[7]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{3, 6, 1, 7, 0, 0, 2, 9, 5},
				{8, 4, 2, 3, 9, 5, 6, 7, 1},
				{0, 5, 0, 2, 6, 1, 4, 8, 3},
				{1, 0, 8, 5, 2, 6, 0, 3, 4},
				{6, 2, 5, 0, 0, 0, 0, 1, 8},
				{0, 3, 4, 1, 0, 0, 5, 2, 6},
				{4, 0, 0, 6, 1, 0, 8, 5, 2},
				{5, 8, 0, 0, 0, 2, 1, 6, 7},
				{2, 1, 6, 8, 5, 7, 3, 4, 9},
			}),
			step: StepRemove2StringKiteCandidates,
			max:  1,
			tests: []CandidateTest{
				{
					column: 5,
					row:    6,
					before: "[3 9]",
					after:  "[3]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{7, 2, 4, 9, 5, 6, 1, 3, 8},
				{1, 6, 8, 4, 2, 3, 5, 9, 7},
				{9, 3, 5, 7, 1, 8, 6, 2, 4},
				{5, 0, 0, 3, 0, 0, 8, 1, 0},
				{0, 4, 0, 0, 8, 1, 7, 5, 0},
				{0, 8, 1, 0, 7, 0, 2, 4, 0},
				{0, 1, 3, 0, 0, 0, 0, 7, 2},
				{0, 0, 0, 1, 0, 0, 0, 8, 5},
				{0, 5, 0, 0, 0, 7, 0, 6, 1},
			}),
			step: StepRemoveEmptyRectangleCandidates,
			max:  1,
			tests: []CandidateTest{
				{
					column: 5,
					row:    7,
					before: "[2 4 9]",
					after:  "[2 4]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{5, 9, 8, 6, 4, 3, 0, 0, 2},
				{0, 0, 3, 7, 5, 9, 6, 4, 8},
				{6, 7, 4, 1, 2, 8, 5, 9, 3},
				{4, 5, 7, 2, 0, 0, 8, 3, 0},
				{9, 0, 6, 3, 0, 7, 4, 2, 5},
				{0, 3, 2, 4, 0, 5, 0, 6, 0},
				{0, 0, 5, 9, 0, 4, 3, 8, 0},
				{3, 4, 1, 8, 7, 2, 9, 5, 6},
				{0, 0, 9, 5, 3, 0, 2, 0, 4},
			}),
			step: StepRemoveEmptyRectangleCandidates,
			max:  1,
			tests: []CandidateTest{
				{
					column: 8,
					row:    3,
					before: "[1 9]",
					after:  "[9]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{7, 0, 0, 0, 5, 4, 0, 1, 0},
				{0, 6, 3, 8, 7, 0, 4, 2, 5},
				{5, 0, 4, 0, 0, 0, 7, 0, 0},
				{2, 7, 0, 4, 0, 0, 0, 0, 1},
				{4, 0, 0, 9, 2, 0, 0, 0, 7},
				{0, 0, 0, 0, 0, 7, 5, 4, 2},
				{8, 5, 2, 0, 4, 3, 0, 7, 9},
				{3, 9, 0, 7, 8, 2, 0, 5, 4},
				{0, 4, 7, 5, 9, 0, 2, 8, 3},
			}),
			step: StepRemoveEmptyRectangleCandidates,
			max:  1,
			tests: []CandidateTest{
				{
					column: 3,
					row:    5,
					before: "[1 3 6]",
					after:  "[1 3]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 8, 1, 0, 2, 0, 6, 0, 0},
				{0, 4, 2, 0, 6, 0, 0, 8, 9},
				{0, 5, 6, 8, 0, 0, 2, 4, 0},
				{6, 9, 3, 1, 4, 2, 7, 5, 8},
				{4, 2, 8, 3, 5, 7, 9, 1, 6},
				{1, 7, 5, 6, 8, 9, 3, 2, 4},
				{5, 1, 0, 0, 3, 6, 8, 9, 2},
				{2, 3, 0, 0, 0, 8, 4, 6, 0},
				{8, 6, 0, 2, 0, 0, 0, 0, 0},
			}),
			step: StepRemoveEmptyRectangleCandidates,
			max:  3,
			tests: []CandidateTest{
				{
					column: 3,
					row:    1,
					before: "[5 7]",
					after:  "[7]",
				},
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{5, 8, 0, 1, 7, 9, 0, 0, 3},
				{0, 0, 0, 6, 0, 8, 9, 7, 5},
				{6, 9, 7, 3, 5, 0, 0, 0, 0},
				{9, 0, 0, 5, 3, 0, 7, 2, 8},
				{7, 0, 3, 8, 1, 0, 5, 0, 0},
				{8, 5, 0, 9, 0, 7, 1, 3, 0},
				{4, 6, 9, 2, 8, 1, 3, 5, 7},
				{0, 0, 8, 7, 6, 5, 0, 0, 0},
				{0, 7, 5, 4, 9, 3, 0, 0, 0},
			}),
			step: StepRemoveEmptyRectangleCandidates,
			max:  1,
			tests: []CandidateTest{
				{
					column: 4,
					row:    1,
					before: "[2 4]",
					after:  "[4]",
				},
				{
					column: 2,
					row:    5,
					before: "[2 4 6]",
					after:  "[4 6]",
				},
			},
			solve: true,
		},
	}

	for testIndex, test := range tests {
		solver := test.puzzle.Solver()
		puzzle := &solver.Puzzle

		for _, cellTest := range test.tests {
			testCell := puzzle.Get(cellTest.column, cellTest.row)
			actual := fmt.Sprint(testCell.Candidates())

			if actual != cellTest.before {
				puzzle.PrintConsoleCandidates()
				t.Fatalf("#%d (%s): Candidates for [%d,%d] are not %s they are %s", testIndex, test.step.Technique, cellTest.column, cellTest.row, cellTest.before, actual)
			}
		}

		removed, _ := test.step.Logic(&solver, SolverLimit{MaxBatches: test.max}, test.step)

		for _, cellTest := range test.tests {
			testCell := puzzle.Get(cellTest.column, cellTest.row)
			actual := fmt.Sprint(testCell.Candidates())

			if actual != cellTest.after {
				puzzle.PrintConsoleCandidates()
				t.Errorf("#%d (%s): Candidates for [%d,%d] are not %s they are %s. %d candidates removed.", testIndex, test.step.Technique, cellTest.column, cellTest.row, cellTest.after, actual, removed)
			}

			if cellTest.value != 0 && testCell.Value != cellTest.value {
				puzzle.PrintConsoleCandidates()
				t.Errorf("#%d (%s): Value for [%d,%d] is not %d it is %d.", testIndex, test.step.Technique, cellTest.column, cellTest.row, cellTest.value, testCell.Value)
			}
		}

		if test.solve {
			solution, solved := solver.Solve(SolverLimit{})

			if !solved {
				solution.PrintConsoleCandidates()
				t.Errorf("#%d (%s) solve failed.", testIndex, test.step.Technique)
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
		solutions := test.puzzle.GetSolutions(SolutionLimit{MaxSolutions: test.max})
		duration := time.Since(start)

		if len(solutions) != test.solutions {
			solver := test.puzzle.Solver()
			solver.LogState = true
			solver.LogEnabled = true
			solver.Solve(SolverLimit{})
			printSolveLogs(&solver, true)

			t.Errorf("An invalid number of solutions found %d expected %d in %s for %s.", len(solutions), test.solutions, duration, test.puzzle.String())
		}

		fmt.Printf("Solutions: %d in %s\n", len(solutions), duration)

		for _, solution := range solutions {
			checkValid(&solution.Puzzle, t)
		}
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
	solver.LogEnabled = true
	solver.Solve(SolverLimit{})

	fmt.Println("TestLogs")
	printSolveLogs(&solver, false)
}

func checkValid(puzzle *Puzzle, t *testing.T) {
	if !puzzle.IsValid() {
		puzzle.PrintConsoleCandidates()
		t.Fatal("The previous puzzle has invalid candidates")
	}
}

func printSolveLogs(solver *Solver, state bool) {
	for _, log := range solver.Logs {
		fmt.Println(log.String())
		if state && log.State != nil {
			log.State.PrintConsoleCandidates()
		}
	}
	last := solver.GetLastLog()
	fmt.Printf("Total cost = %d, placements = %d, batches = %d, logs = %d.\n", last.RunningCost, last.RunningPlacements, last.Batch, last.Index)
}
