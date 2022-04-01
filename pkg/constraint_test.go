package sudogo

import (
	"fmt"
	"testing"
	"time"
)

func TestConstraintSum(t *testing.T) {
	tests := []struct {
		puzzle          Puzzle
		constraint      ConstraintSum
		cellsExpected   []string
		affects         []Position
		affectsExpected []bool
	}{
		{
			puzzle: Classic.Empty(),
			constraint: ConstraintSum{
				Sum: SumConstant(8),
				Cells: &[]Position{
					{0, 0},
					{1, 0},
					{2, 0},
				},
			},
			cellsExpected: []string{
				"[1 2 3 4 5]",
				"[1 2 3 4 5]",
				"[1 2 3 4 5]",
			},
			affects: []Position{
				{0, 0}, {1, 0}, {2, 0}, {3, 0},
				{0, 1}, {1, 1}, {2, 1}, {3, 1},
			},
			affectsExpected: []bool{
				true, true, true, false,
				false, false, false, false,
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			}),
			constraint: ConstraintSum{
				Sum: SumConstant(15),
				Cells: &[]Position{
					{0, 0},
					{1, 0},
					{2, 0},
				},
			},
			cellsExpected: []string{
				"[1 2 3 4 5 6 7 8 9]",
				"[1 2 3 4 5 6 7 8 9]",
				"[1 2 3 4 5 6 7 8 9]",
			},
			affects:         []Position{},
			affectsExpected: []bool{},
		},
		{
			puzzle: Classic.Create([][]int{
				{7, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			}),
			constraint: ConstraintSum{
				Sum: SumConstant(15),
				Cells: &[]Position{
					{0, 0},
					{1, 0},
					{2, 0},
				},
			},
			cellsExpected: []string{
				"[]",
				"[2 3 4 5 6]",
				"[2 3 4 5 6]",
			},
			affects:         []Position{},
			affectsExpected: []bool{},
		},
		{
			puzzle: Classic.Create([][]int{
				{7, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 5, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 4, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 6, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			}),
			constraint: ConstraintSum{
				Sum: SumConstant(15),
				Cells: &[]Position{
					{0, 0},
					{1, 0},
					{2, 0},
				},
			},
			cellsExpected: []string{
				"[]",
				"[6]", // 1 2 3 6
				"[2]", // 1 2 3
			},
			affects:         []Position{},
			affectsExpected: []bool{},
		},
	}

	for testIndex, test := range tests {
		for i, pos := range test.affects {
			cell := test.puzzle.Get(pos.Col, pos.Row)
			expected := test.affectsExpected[i]
			actual := test.constraint.Affects(cell)

			if expected != actual {
				t.Errorf("TestSumConstraint Affects failed at index %d for test %d", i, testIndex)
			}
		}

		for cellIndex, pos := range *test.constraint.Cells {
			cell := test.puzzle.Get(pos.Col, pos.Row)
			cand := cell.candidates

			if cell.Empty() {
				test.constraint.RemoveCandidates(cell, &test.puzzle, &cand)
			}

			expected := test.cellsExpected[cellIndex]
			actual := fmt.Sprintf("%v", cand.ToSlice())

			if expected != actual {
				t.Errorf("TestConstraintSum cells failed at index %d for test %d, actual: %s, expected: %s", cellIndex, testIndex, actual, expected)
			}
		}
	}
}

func TestConstraintUnique(t *testing.T) {
	tests := []struct {
		puzzle        Puzzle
		constraint    ConstraintUnique
		cellsExpected []string
	}{
		{
			puzzle: Classic.Create([][]int{
				{7, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 5, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 4, 0, 0, 0, 0, 8, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 6, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 2, 0, 0},
				{9, 0, 0, 0, 0, 0, 1, 0, 0},
			}),
			constraint: ConstraintUnique{
				Cells: &[]Position{
					{0, 0},
					{1, 1},
					{2, 2},
					{3, 3},
					{4, 4},
					{5, 5},
					{6, 6},
					{7, 7},
					{8, 8},
				},
			},
			cellsExpected: []string{
				"[]",
				"[]",
				"[]",
				"[1 2 3 6 8 9]",
				"[1 2 3 8 9]",
				"[1 2 3 6 8 9]",
				"[3 6 8 9]",
				"[3 6 9]",
				"[3 6 8]",
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{7, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 5, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 4, 0, 0, 0, 0, 8, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 6, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 2, 0, 0},
				{9, 0, 0, 0, 0, 0, 1, 0, 0},
			}),
			constraint: ConstraintUnique{
				Cells: &[]Position{
					{0, 0},
					{3, 3},
					{6, 6},
				},
				Same: true,
			},
			cellsExpected: []string{
				"[]",
				"[7]",
				"[7]",
			},
		},
	}

	for testIndex, test := range tests {
		for cellIndex, pos := range *test.constraint.Cells {
			cell := test.puzzle.Get(pos.Col, pos.Row)
			cand := cell.candidates

			if cell.Empty() {
				test.constraint.RemoveCandidates(cell, &test.puzzle, &cand)
			}

			expected := test.cellsExpected[cellIndex]
			actual := fmt.Sprintf("%v", cand.ToSlice())

			if expected != actual {
				t.Errorf("TestConstraintUnique cells failed at index %d for test %d, actual: %s, expected: %s", cellIndex, testIndex, actual, expected)
			}
		}
	}
}

func TestConstraintOrder(t *testing.T) {
	tests := []struct {
		puzzle        Puzzle
		constraint    ConstraintOrder
		cellsExpected []string
	}{
		{
			puzzle: Classic.Create([][]int{
				{7, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 5, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 4, 0, 0, 0, 0, 8, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 6, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 2, 0, 0},
				{9, 0, 0, 0, 0, 0, 1, 0, 0},
			}),
			constraint: ConstraintOrder{
				Cells: &[]Position{
					{0, 0},
					{1, 1},
					{2, 2},
					{3, 3},
					{4, 4},
				},
			},
			cellsExpected: []string{
				"[]",
				"[]",
				"[]",
				"[2 3]",
				"[1 2]",
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 2, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 8, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 6, 0, 7, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 2, 0, 0},
				{9, 0, 0, 0, 0, 0, 1, 0, 0},
			}),
			constraint: ConstraintOrder{
				Direction: 1,
				Cells: &[]Position{
					{0, 0},
					{1, 1},
					{2, 2},
					{3, 3},
					{4, 4},
				},
			},
			cellsExpected: []string{
				"[1]",
				"[]",
				"[3 4 5]",
				"[4 5 6]",
				"[]",
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 8, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 6, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 2, 0, 0},
				{9, 0, 0, 0, 0, 0, 1, 0, 0},
			}),
			constraint: ConstraintOrder{
				Direction: 1,
				Cells: &[]Position{
					{0, 0},
					{1, 1},
					{2, 2},
					{3, 3},
					{4, 4},
				},
			},
			cellsExpected: []string{
				"[1 2 3 4 5]",
				"[2 3 4 5 6]",
				"[3 4 5 7]",
				"[4 5 6 7 8]",
				"[5 7 8 9]",
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 8, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 6, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 2, 0, 0},
				{9, 0, 0, 0, 0, 0, 1, 0, 0},
			}),
			constraint: ConstraintOrder{
				Cells: &[]Position{
					{0, 0},
					{1, 1},
					{2, 2},
					{3, 3},
					{4, 4},
				},
			},
			cellsExpected: []string{
				"[1 2 3 4 5 6 7 8]",
				"[1 2 3 4 5 6 7 8 9]",
				"[1 2 3 4 5 7 9]",
				"[1 2 3 4 5 6 7 8 9]",
				"[1 2 3 4 5 7 8 9]",
			},
		},
		{
			puzzle: Classic.Create([][]int{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 4, 0, 0, 0, 0, 8, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 6, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 2, 0, 0},
				{9, 0, 0, 0, 0, 0, 1, 0, 0},
			}),
			constraint: ConstraintOrder{
				Cells: &[]Position{
					{0, 0},
					{1, 1},
					{2, 2},
					{3, 3},
					{4, 4},
				},
			},
			cellsExpected: []string{
				"[1 2 3 5 6 7 8]",
				"[1 2 3 5 6 7 8 9]",
				"[]",
				"[1 2 3 5 6 7 8 9]",
				"[1 2 3 5 7 8 9]",
			},
		},
	}

	for testIndex, test := range tests {
		for cellIndex, pos := range *test.constraint.Cells {
			cell := test.puzzle.Get(pos.Col, pos.Row)
			cand := cell.candidates

			if cell.Empty() {
				test.constraint.RemoveCandidates(cell, &test.puzzle, &cand)
			}

			expected := test.cellsExpected[cellIndex]
			actual := fmt.Sprintf("%v", cand.ToSlice())

			if expected != actual {
				t.Errorf("TestConstraintUnique cells failed at index %d for test %d, actual: %s, expected: %s", cellIndex, testIndex, actual, expected)
			}
		}
	}
}

func Test4GivenDigits(t *testing.T) {
	// https://www.youtube.com/watch?v=hAyZ9K2EBF0
	kind := NewKind(3, 3)
	kind.Constraints = []Constraint{
		// Diagonal from top left to botttom right must have 1-9
		&ConstraintUnique{
			Cells: &[]Position{
				{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8},
			},
		},
		// Diagonal from top right to botttom left must have 1-9
		&ConstraintUnique{
			Cells: &[]Position{
				{8, 0}, {7, 1}, {6, 2}, {5, 3}, {4, 4}, {3, 5}, {2, 6}, {1, 7}, {0, 8},
			},
		},
		// Knights move apart cannot contain same digit
		&ConstraintUnique{
			Relative: &[]Position{
				{2, 1}, {2, -1},
				{-2, 1}, {-2, -1},
				{1, 2}, {1, -2},
				{-1, 2}, {-1, -2},
			},
		},
		// Magic Square
		&ConstraintMagic{
			Center: Position{4, 4},
		},
	}

	puzzle := kind.Create([][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{3, 8, 4, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 2},
	})

	// solver := puzzle.Solver()
	// solver.Solve(SolverLimit{})
	// solver.puzzle.PrintConsoleCandidates()

	start := time.Now()
	solutions := puzzle.GetSolutions(SolutionLimit{})
	duration := time.Since(start)

	if len(solutions) != 1 {
		t.Fatalf("Incorrect number of solutions found for Test4GivenDigits: %d", len(solutions))
	} else {
		t.Logf("Solution found for Test4GivenDigits in %s", duration)

		solutions[0].Puzzle.PrintConsole()
	}
}
