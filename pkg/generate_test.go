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
		attempts  int
	}{
		{
			name:      "Easy",
			clear:     30,
			symmetric: true,
			attempts:  128,
		},
		{
			name:      "Medium",
			clear:     40,
			symmetric: true,
			attempts:  128,
		},
		{
			name:      "Hard",
			clear:     50,
			symmetric: false,
			attempts:  128,
		},
	}

	for _, test := range tests {
		start := time.Now()
		puzzle, _ := gen.Generate()
		cleared, actualAttempts := gen.ClearCells(test.clear, test.symmetric, test.attempts)
		duration := time.Since(start)

		if !cleared {
			t.Errorf("Failed to generate %s puzzle after %d attempts in %s", test.name, actualAttempts, duration)
		} else {
			puzzle.Print()
			fmt.Printf("Generated %s in %s after %d attempts\n", test.name, duration, actualAttempts)
		}
	}
}
