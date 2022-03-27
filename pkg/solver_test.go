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

	println("TestSolveSimple")
	original.Print()

	s := NewSolver(original)
	p := s.puzzle
	s.Solve()
	p.Print()

	if !s.Solved() {
		t.Errorf("The solver failed to solve the naked single puzzle.")
	}
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

	println("TestSolveHiddenSingle")
	s := original.Solver()
	p := s.puzzle
	s.Solve()
	p.Print()

	if p.Get(3, 2).value != 6 {
		t.Errorf("The solver failed to use hidden single logic on r3c4.")
	}
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
	p := s.puzzle

	if fmt.Sprint(p.Get(6, 2).Candidates()) != "[3 5]" {
		t.Errorf("Invalid candidates for r3c7 in Pointing")
	}

	StepRemovePointingCandidates(&s, -1)

	if fmt.Sprint(p.Get(6, 2).Candidates()) != "[3]" {
		t.Errorf("Invalid candidates for r3c7 in Pointing after step")
	}

	s.Solve()

	if p.Get(6, 2).value != 3 {
		t.Errorf("Test Pointing solve failed")
	}

	if !p.IsSolved() {
		p.PrintCandidates()
		t.Fatal("Puzzle not solved")
	}
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
	p := s.puzzle

	if fmt.Sprint(p.Get(1, 2).Candidates()) != "[4 7]" {
		t.Errorf("Invalid candidates for r3c2 in Claiming")
	}

	StepRemoveClaimingCandidates(&s, -1)

	if fmt.Sprint(p.Get(1, 2).Candidates()) != "[4]" {
		t.Errorf("Invalid candidates for r3c2 in Claiming after step")
	}

	s.Solve()

	if p.Get(1, 2).value != 4 {
		t.Errorf("Test Claiming solve failed")
	}

	if !p.IsSolved() {
		p.PrintCandidates()
		t.Fatal("Puzzle not solved")
	}
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
	p := s.puzzle

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
}
