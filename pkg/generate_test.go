package sudogo

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	g := Classic.Generator()

	start := time.Now()
	p, attempts := g.Generate()

	if p == nil {
		t.Fatalf("Failed to generate a puzzle")
	} else {
		p.Print()
	}

	duration := time.Since(start)
	fmt.Printf("TestGenerate in %s after %d attempts.\n", duration, attempts)
}

func TestGenerateClear(t *testing.T) {
	gen := Classic.Generator()

	tests := []struct {
		name      string
		clear     int
		symmetric bool
		maxStates int
	}{
		{
			name:      "Easy",
			clear:     30,
			symmetric: true,
			maxStates: 1 << 10,
		},
		{
			name:      "Medium",
			clear:     40,
			symmetric: true,
			maxStates: 1 << 10,
		},
		{
			name:      "Hard",
			clear:     50,
			symmetric: false,
			maxStates: 1 << 10,
		},
	}

	for _, test := range tests {
		start := time.Now()
		puzzle, _ := gen.Generate()
		cleared, states := gen.ClearCells(puzzle, test.clear, test.symmetric, test.maxStates)
		duration := time.Since(start)

		if cleared == nil || !cleared.HasUniqueSolution() {
			t.Errorf("Failed to generate unique %s puzzle after %d states in %s", test.name, states, duration)
		} else {
			cleared.Print()
			fmt.Printf("Generated %s (%d empty cells) in %s after %d states\n", test.name, test.clear, duration, states)
		}
	}
}
