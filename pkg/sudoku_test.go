package sudogo

import (
	"fmt"
	"testing"
	"time"
)

func TestBits(t *testing.T) {
	b := Bits{}

	if b.Count != 0 {
		t.Error("default count is not zero")
	}
	if b.Value != 0 {
		t.Error("default value is not zero")
	}

	b.Fill(4)

	if b.Value != 15 {
		t.Error("fill has wrong value")
	}
	if b.Count != 4 {
		t.Error("fill has wrong count")
	}

	if b.Set(0, true) {
		t.Error("bit 0 is already set and I was allowed to set it")
	}
	if b.Set(1, true) {
		t.Error("bit 1 is already set and I was allowed to set it")
	}
	if b.Set(2, true) {
		t.Error("bit 2 is already set and I was allowed to set it")
	}
	if b.Set(3, true) {
		t.Error("bit 3 is already set and I was allowed to set it")
	}
	if b.Count != 4 {
		t.Error("set of already set bits affected count")
	}

	if !b.Set(0, false) {
		t.Error("bit 0 could not be set")
	}
	if b.Count != 3 {
		t.Error("count wrong after set", 0)
	}
	if fmt.Sprint(b.ToSlice()) != "[1 2 3]" {
		t.Error("setting bit 0 resulted in the wrong slice")
	}

	if !b.Set(1, false) {
		t.Error("bit 1 could not be set")
	}
	if b.Count != 2 {
		t.Error("count wrong after set 1")
	}
	if fmt.Sprint(b.ToSlice()) != "[2 3]" {
		t.Error("setting bit 1 resulted in the wrong slice")
	}

	if !b.Set(2, false) {
		t.Error("bit 2 could not be set")
	}
	if b.Count != 1 {
		t.Error("count wrong after set 2")
	}
	if fmt.Sprint(b.ToSlice()) != "[3]" {
		t.Error("setting bit 2 resulted in the wrong slice")
	}

	if !b.Set(3, false) {
		t.Error("bit 3 could not be set")
	}
	if b.Count != 0 {
		t.Error("count wrong after set 3")
	}
	if fmt.Sprint(b.ToSlice()) != "[]" {
		t.Error("setting bit 3 resulted in the wrong slice")
	}
}

func TestSolveSimple(t *testing.T) {
	p := Classic.Create([][]int{
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
	p.Print()
	p.Solve()
	p.Print()

	if !p.Solved() {
		t.Errorf("The solver failed to solve the naked single puzzle.")
	}
}

func TestSolveHiddenSingle(t *testing.T) {
	p := Classic.Create([][]int{
		{0, 2, 8, 0, 0, 7, 0, 0, 0},
		{0, 1, 6, 0, 8, 3, 0, 7, 0},
		{0, 0, 0, 6, 2, 0, 8, 5, 1},
		{1, 3, 7, 2, 9, 0, 0, 0, 0},
		{0, 0, 0, 7, 3, 0, 0, 0, 0},
		{0, 0, 0, 0, 4, 6, 3, 0, 7},
		{2, 9, 0, 0, 7, 0, 0, 0, 0},
		{0, 0, 0, 8, 6, 0, 1, 4, 0},
		{0, 0, 0, 3, 0, 0, 7, 0, 0},
	})

	println("TestSolveHiddenSingle")
	p.Solve()
	p.Print()

	if p.Get(3, 2).value != 6 {
		t.Errorf("The solver failed to use hidden single logic on r3c4.")
	}
}

func TestNakedPair(t *testing.T) {
	p := Classic.Create([][]int{
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

	r8c2 := p.Get(1, 7)

	c0 := fmt.Sprint(r8c2.Candidates())
	if c0 != "[3 7]" {
		t.Fatalf("Candidates for r8c2 are not [3 7] they are %s", c0)
	}

	p.RemoveNakedSubsetCandidates(1)

	c1 := fmt.Sprint(r8c2.Candidates())
	if c1 != "[7]" {
		t.Fatalf("Candidates for r8c2 are not [7] they are %s", c1)
	}
}

func TestGenerate(t *testing.T) {
	start := time.Now()
	p := Classic.Generate()
	duration := time.Since(start)
	fmt.Printf("TestGenerate in %s\n", duration)
	p.Print()
}
