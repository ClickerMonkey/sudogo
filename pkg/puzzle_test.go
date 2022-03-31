package sudogo

import (
	"fmt"
	"testing"
)

func TestStateString(t *testing.T) {
	tests := []struct {
		puzzle      Puzzle
		includeKind bool
		emptyValue  string
		expected    string
	}{
		{
			puzzle: Kind2x2.Create([][]int{
				{1, 2, 3, 4},
				{0, 3, 2, 1},
				{0, 0, 1, 2},
				{0, 0, 0, 3},
			}),
			includeKind: false,
			emptyValue:  ".",
			expected:    "1234.321..12...3",
		},
		{
			puzzle: Kind2x2.Create([][]int{
				{1, 2, 3, 4},
				{0, 3, 2, 1},
				{0, 0, 1, 2},
				{0, 0, 0, 3},
			}),
			includeKind: true,
			emptyValue:  " ",
			expected:    "2x2,1234 321  12   3",
		},
		{
			puzzle: Kind3x2.Create([][]int{
				{0, 0, 3, 0, 1, 0},
				{5, 6, 0, 3, 2, 0},
				{0, 5, 4, 2, 0, 3},
				{2, 0, 6, 4, 5, 0},
				{0, 1, 2, 0, 4, 5},
				{0, 4, 0, 1, 0, 0},
			}),
			includeKind: false,
			emptyValue:  ".",
			expected:    "..3.1.56.32..542.32.645..12.45.4.1..",
		},
		{
			puzzle: Kind3x2.Create([][]int{
				{0, 0, 3, 0, 1, 0},
				{5, 6, 0, 3, 2, 0},
				{0, 5, 4, 2, 0, 3},
				{2, 0, 6, 4, 5, 0},
				{0, 1, 2, 0, 4, 5},
				{0, 4, 0, 1, 0, 0},
			}),
			includeKind: true,
			emptyValue:  "_",
			expected:    "3x2,__3_1_56_32__542_32_645__12_45_4_1__",
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
				{0, 9, 0, 0, 0, 6, 7, 0, 0},
			}),
			includeKind: false,
			emptyValue:  ".",
			expected:    "...1.2....6...83..5.......9...4.7..868.........4....1..2....5......712.6.9...67..",
		},
	}

	for testIndex, test := range tests {
		if !test.puzzle.IsValid() {
			t.Errorf("Puzzle #%d for TestStateString is not valid.", testIndex)
			continue
		}

		actual := test.puzzle.ToStateString(test.includeKind, test.emptyValue)

		if actual != test.expected {
			t.Errorf("Puzzle state string expected %s was %s.", test.expected, actual)
			continue
		}

		parsed := FromString(actual)

		if parsed == nil {
			t.Errorf("Puzzle failed to parse %s.", actual)
		} else {
			if !comparePuzzles(t, parsed, &test.puzzle) {
				fmt.Printf("Expected:\n%s\n", test.puzzle.ToConsoleString())
				fmt.Printf("Actual:\n%s\n", parsed.ToConsoleString())
			}
		}
	}
}

func TestEncodedString(t *testing.T) {
	tests := []struct {
		puzzle   Puzzle
		expected string
	}{
		{
			puzzle: Kind2x2.Create([][]int{
				{1, 2, 3, 4},
				{0, 3, 2, 1},
				{0, 0, 1, 2},
				{0, 0, 0, 3},
			}),
			expected: "LCRL8EhC",
		},
		{
			puzzle: Kind3x2.Create([][]int{
				{0, 0, 3, 0, 1, 0},
				{5, 6, 0, 3, 2, 0},
				{0, 5, 4, 2, 0, 3},
				{2, 0, 6, 4, 5, 0},
				{0, 1, 2, 0, 4, 5},
				{0, 4, 0, 1, 0, 0},
			}),
			expected: "AS4LrNo62oiyLSRNAEM=",
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
				{0, 9, 0, 0, 0, 6, 7, 0, 0},
			}),
			expected: "A4YH884hh1mk+w1YDPY1/WqoJc6fhK6UGsKkov4kDYKwYw==",
		},
	}

	for testIndex, test := range tests {
		if !test.puzzle.IsValid() {
			t.Errorf("Puzzle #%d for TestEncodedString is not valid.", testIndex)
			continue
		}

		actual := test.puzzle.EncodedString()

		if actual != test.expected {
			t.Errorf("Puzzle encoded string expected %s was %s.", test.expected, actual)
			continue
		}

		parsed := FromEncoded(actual)

		if parsed == nil {
			t.Errorf("Puzzle failed to parse %s.", actual)
		} else {
			if !comparePuzzles(t, parsed, &test.puzzle) {
				fmt.Printf("Expected:\n%s\n", test.puzzle.ToConsoleString())
				fmt.Printf("Actual:\n%s\n", parsed.ToConsoleString())
			}
		}
	}
}

func comparePuzzles(t *testing.T, actual *Puzzle, expected *Puzzle) bool {
	if actual.Kind.BoxSize.Width != expected.Kind.BoxSize.Width {
		t.Errorf("Expected puzzle box width is %d but was %d", expected.Kind.BoxSize.Width, actual.Kind.BoxSize.Width)
		return false
	}
	if actual.Kind.BoxSize.Height != expected.Kind.BoxSize.Height {
		t.Errorf("Expected puzzle box height is %d but was %d", expected.Kind.BoxSize.Height, actual.Kind.BoxSize.Height)
		return false
	}
	valid := true
	for i := range actual.Cells {
		a := &actual.Cells[i]
		e := &expected.Cells[i]

		if a.Value != e.Value {
			t.Errorf("Expected cell [%d,%d] is %d but was %d.", a.Col, a.Row, e.Value, a.Value)
			valid = false
		}
		if a.candidates.Value != e.candidates.Value {
			t.Errorf("Expected candidates [%d,%d] is %v but was %v.", a.Col, a.Row, e.Candidates(), a.Candidates())
			valid = false
		}
	}
	return valid
}
