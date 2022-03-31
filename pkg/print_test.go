package sudogo

import (
	"testing"
)

func TestPrint(t *testing.T) {
	original := Classic.Create([][]int{
		{0, 0, 0, 1, 0, 2, 0, 0, 0},
		{0, 6, 0, 0, 0, 8, 3, 0, 0},
		{5, 0, 0, 0, 0, 0, 0, 0, 9},
		{0, 0, 0, 4, 0, 7, 0, 0, 8},
		{6, 8, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 4, 0, 0, 5, 0, 1, 0},
		{0, 2, 0, 0, 0, 0, 5, 0, 0},
		{0, 0, 0, 0, 7, 0, 2, 0, 6},
		{0, 9, 0, 0, 0, 6, 7, 0, 0},
	})

	println(original.ToConsoleString())
}

func TestPrintCandidates(t *testing.T) {
	original := Classic.Create([][]int{
		{0, 0, 0, 1, 0, 2, 0, 0, 0},
		{0, 6, 0, 0, 0, 8, 3, 0, 0},
		{5, 0, 0, 0, 0, 0, 0, 0, 9},
		{0, 0, 0, 4, 0, 7, 0, 0, 8},
		{6, 8, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 4, 0, 0, 5, 0, 1, 0},
		{0, 2, 0, 0, 0, 0, 5, 0, 0},
		{0, 0, 0, 0, 7, 0, 2, 0, 6},
		{0, 9, 0, 0, 0, 6, 7, 0, 0},
	})

	println(original.ToConsoleCandidatesString())
}
