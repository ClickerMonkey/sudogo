package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	sudogo "github.com/ClickerMonkey/sudogo/pkg"
)

func main() {
	var (
		types = map[string]*sudogo.ClearLimits{
			"beginner": &sudogo.DifficultyBeginner,
			"easy":     &sudogo.DifficultyEasy,
			"medium":   &sudogo.DifficultyMedium,
			"hard": {
				SolverLimit: sudogo.SolverLimit{
					MaxPlacements: 58,
				},
			},
			"custom": {},
		}
		chosenType   string              = "medium"
		chosenLimits *sudogo.ClearLimits = types[chosenType]
	)

	flag.Func("type", "One of beginner, easy, medium, hard, or custom.", func(value string) error {
		chosenType = strings.ToLower(value)
		chosenLimits = types[chosenType]
		if chosenLimits == nil {
			return errors.New("must be one of beginner, easy, medium, hard, or custom")
		}
		return nil
	})

	boxWidth := flag.Int("boxWidth", 3, "The width of a box.")
	boxHeight := flag.Int("boxHeight", 3, "The height of a box.")
	symmetric := flag.Bool("symmetric", true, "If the easier puzzles should be symmetric.")
	count := flag.Int("count", 1, "The number of puzzles to generate.")
	clearDepth := flag.Int("clearDepth", 64, "How many times we should try to clear out a generated puzzle to match the requested criteria before generating a new puzzle.")
	tries := flag.Int("tries", 128, "How many times we should generate a new puzzle to try to clear.")
	minCost := flag.Int("minCost", -1, "Override the minCost value for generation.")
	maxCost := flag.Int("maxCost", -1, "Override the maxCost value for generation.")
	maxPlacements := flag.Int("maxPlacements", -1, "Override the maxPlacements value for generation.")
	maxLogs := flag.Int("maxLogs", -1, "Override the maxLogs value for generation.")
	maxBatches := flag.Int("maxBatches", -1, "Override the maxBatches value for generation.")
	candidates := flag.Bool("candidates", false, "If the candidate puzzles should be printed.")
	solutions := flag.Bool("solutions", false, "If the solutions should be printed as well.")
	logSteps := flag.Bool("steps", false, "If the soluton steps should be logged.")

	flag.Parse()

	kind := sudogo.Kind{
		BoxSize: sudogo.Size{
			Width:  *boxWidth,
			Height: *boxHeight,
		},
	}

	chosenLimits.MaxStates = *clearDepth

	if *symmetric == false || chosenType == "custom" {
		chosenLimits.Symmetric = *symmetric
	}
	if *minCost != -1 {
		chosenLimits.MinCost = *minCost
	}
	if *maxCost != -1 {
		chosenLimits.MaxCost = *maxCost
	}
	if *maxPlacements != -1 {
		chosenLimits.MaxPlacements = *maxPlacements
	}
	if *maxLogs != -1 {
		chosenLimits.MaxLogs = *maxLogs
	}
	if *maxBatches != -1 {
		chosenLimits.MaxBatches = *maxBatches
	}

	gen := kind.Generator()

	for puzzleIndex := 0; puzzleIndex < *count; puzzleIndex++ {
		var puzzle, solution *sudogo.Puzzle

		for tryIndex := 0; tryIndex < *tries; tryIndex++ {
			generated, _ := gen.Generate()
			if generated == nil {
				continue
			}
			cleared, _ := gen.ClearCells(generated, *chosenLimits)
			if cleared != nil {
				puzzle = cleared
				solution = generated
				break
			}
		}

		if puzzle != nil {
			fmt.Printf("Puzzle #%d:\n", puzzleIndex+1)
			if *candidates {
				puzzle.PrintCandidates()
			} else {
				puzzle.Print()
			}
			fmt.Println()
		} else {
			fmt.Printf("Puzzle #%d could not be generated.\n", puzzleIndex+1)
		}
		if *solutions && solution != nil {
			fmt.Printf("Solution for puzzle #%d:\n", puzzleIndex+1)
			solution.Print()
			fmt.Println()
		}
		if puzzle != nil && *logSteps {
			solver := puzzle.Solver()
			solver.LogEnabled = true
			_, solved := solver.Solve(sudogo.SolverLimit{})
			if solved {
				fmt.Printf("Steps for puzzle #%d:\n", puzzleIndex+1)
				for _, log := range solver.GetLogs() {
					fmt.Println(log.String())
				}
				fmt.Println()
			} else {
				fmt.Printf("Solution steps could not be determined.\n")
			}
		}

	}

}
