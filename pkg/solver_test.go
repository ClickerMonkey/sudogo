package sudogo

import (
	"fmt"
	"testing"
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
	solution, solved := s.Solve()

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
	solution, solved := s.Solve()

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

	StepRemovePointingCandidates(&s, -1)

	if fmt.Sprint(p.Get(6, 2).Candidates()) != "[3]" {
		t.Errorf("Invalid candidates for r3c7 in Pointing after step")
	}

	solution, solved := s.Solve()

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

	StepRemoveClaimingCandidates(&s, -1)

	if fmt.Sprint(p.Get(1, 2).Candidates()) != "[4]" {
		t.Errorf("Invalid candidates for r3c2 in Claiming after step")
	}

	solution, solved := s.Solve()

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

	StepRemoveNakedSubsetCandidates(&s, 1)

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
		puzzle  Puzzle
		subsets []int
		max     int
		tests   []CandidateTest
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
			subsets: []int{2},
			max:     19,
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
			subsets: []int{2},
			max:     1,
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
			subsets: []int{3},
			max:     23,
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

		removed := doRemoveHiddenSubsetCandidates(&solver, test.max, test.subsets)

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
	original := Classic.Create([][]int{
		{0, 0, 0, 1, 0, 2, 0, 0, 0},
		{0, 6, 0, 0, 0, 8, 3, 0, 0},
		{5, 0, 0, 0, 0, 0, 0, 0, 9},
		{0, 0, 0, 4, 0, 7, 0, 0, 8},
		{6, 8, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 4, 0, 0, 0, 0, 1, 0},
		{0, 2, 0, 0, 0, 0, 5, 0, 0},
		{0, 0, 0, 0, 7, 1, 2, 0, 6},
		{0, 9, 0, 0, 0, 6, 7, 0, 0},
	})

	s := original.Solver()
	solution, solved := s.Solve()

	if !solved {
		solution.PrintCandidates()
		t.Fatal("Could not SolveHard puzzle")
	}

	checkValid(solution, t)
}

func checkValid(puzzle *Puzzle, t *testing.T) {
	if !puzzle.IsValid() {
		puzzle.PrintCandidates()
		t.Fatal("The previous puzzle has invalid candidates")
	}
}
