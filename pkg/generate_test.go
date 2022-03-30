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
		clear     ClearLimits
		symmetric bool
		maxStates int
	}{
		{
			name: "Easy",
			clear: ClearLimits{
				SolverLimit: SolverLimit{
					MaxPlacements: 30,
				},
				Symmetric: true,
			},
		},
		{
			name: "Medium",
			clear: ClearLimits{
				SolverLimit: SolverLimit{
					MaxPlacements: 40,
				},
				Symmetric: true,
			},
		},
		{
			name: "Hard",
			clear: ClearLimits{
				SolverLimit: SolverLimit{
					MaxPlacements: 50,
				},
				Symmetric: true,
			},
		},
		{
			name:  "DifficultyBeginner",
			clear: DifficultyBeginner,
		},
		{
			name:  "DifficultyEasy",
			clear: DifficultyEasy,
		},
		{
			name:  "DifficultyMedium",
			clear: DifficultyMedium,
		},
		// {
		// 	name:  "DifficultyTricky",
		// 	clear: DifficultyTricky,
		// },
		// {
		// 	name:  "DifficultyFiendish",
		// 	clear: DifficultyFiendish,
		// },
		// {
		// 	name:  "DifficultyDiabolical",
		// 	clear: DifficultyDiabolical,
		// },
	}

	for _, test := range tests {
		limits := test.clear.Extend(ClearLimits{
			MaxStates: 64,
		})
		start := time.Now()

		var attempts int = 32
		var cleared *Puzzle = nil
		var states int = 0

		for cleared == nil && attempts > 0 {
			puzzle, _ := gen.Generate()
			cleared, states = gen.ClearCells(puzzle, limits)
			attempts--
		}

		duration := time.Since(start)

		if cleared == nil || !cleared.HasUniqueSolution() {
			t.Errorf("Failed to generate unique %s puzzle after %d states in %s", test.name, states, duration)
		} else {
			cleared.Print()
			fmt.Printf("Generated %s (%d empty cells) in %s after %d states\n", test.name, test.clear.MaxPlacements, duration, states)
		}
	}
}
