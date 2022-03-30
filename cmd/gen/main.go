package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"

	sudogo "github.com/ClickerMonkey/sudogo/pkg"
	gofpdf "github.com/jung-kurt/gofpdf"
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
	outputPdf := flag.String("pdf", "", "Output the puzzles to the given PDF file.")

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

	var pdf *gofpdf.Fpdf = nil

	if *outputPdf != "" {
		pdf = gofpdf.New("P", "mm", "A4", "")
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
			displaySolution := solution
			if !*solutions {
				displaySolution = nil
			}

			if pdf != nil {
				handlePDFOutput(puzzleIndex, puzzle, displaySolution, *candidates, *logSteps, pdf)
			} else {
				handleConsoleOutput(puzzleIndex, puzzle, displaySolution, *candidates, *logSteps)
			}
		} else {
			fmt.Printf("Puzzle #%d could not be generated.\n", puzzleIndex+1)
		}
	}

	if pdf != nil {
		pdf.OutputFileAndClose(*outputPdf)
	}
}

func handleConsoleOutput(puzzleIndex int, puzzle *sudogo.Puzzle, solution *sudogo.Puzzle, candidates bool, logSteps bool) {

	fmt.Printf("Puzzle #%d:\n", puzzleIndex+1)
	if candidates {
		puzzle.PrintCandidates()
	} else {
		puzzle.Print()
	}
	fmt.Println()

	if solution != nil {
		fmt.Printf("Solution for puzzle #%d:\n", puzzleIndex+1)
		solution.Print()
		fmt.Println()
	}

	if logSteps {
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
			fmt.Println("Solution steps could not be determined.")
		}
	}
}

func handlePDFOutput(puzzleIndex int, puzzle *sudogo.Puzzle, solution *sudogo.Puzzle, candidates bool, logSteps bool, pdf *gofpdf.Fpdf) {
	pdf.AddPage()

	size := float64(puzzle.Kind.Size())
	boxSize := puzzle.Kind.BoxSize
	pageW, pageH, _ := pdf.PageSize(0)
	marginL, marginT, marginR, marginB := pdf.GetMargins()
	innerW := pageW - marginL - marginR
	innerH := pageH - marginT - marginB
	cellW := innerW / size
	cellH := innerH / size
	cellSize := math.Min(cellW, cellH)
	boardSize := cellSize * size
	boardX := marginL + (innerW-boardSize)/2
	boardY := marginT + (innerH-boardSize)/2
	boxHeight := float64(boxSize.Height)
	boxWidth := float64(boxSize.Width)
	fontScale := cellSize / 21.1
	thickLineWidth := 1.0 * fontScale

	for y := 0.0; y < size; y++ {
		for x := 0.0; x < size; x++ {
			cell := puzzle.Get(int(y), int(x))
			cellValue := ""
			if cell.HasValue() {
				cellValue = strconv.Itoa(cell.Value)
			}
			pdf.SetFont("Arial", "B", 32*fontScale)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetXY(boardX+x*cellSize, boardY+y*cellSize)
			pdf.CellFormat(cellSize, cellSize, cellValue, "1", 0, "CM", false, 0, "")
			if candidates && cell.Empty() {
				cand := fmt.Sprintf("%v", cell.Candidates())
				cand = strings.Trim(cand, "[]")

				pdf.SetFont("Arial", "", 10*fontScale)
				pdf.SetTextColor(128, 128, 128)
				pdf.SetXY(boardX+x*cellSize, boardY+y*cellSize+thickLineWidth*2)
				pdf.CellFormat(cellSize, cellSize, cand, "0", 0, "CT", false, 0, "")
			}
		}
	}

	pdf.SetLineWidth(thickLineWidth)
	pdf.SetLineCapStyle("round")

	for y := 0.0; y <= size+0.0001; y += boxHeight {
		lineY := boardY + y*cellSize
		pdf.Line(boardX, lineY, boardX+boardSize, lineY)
	}

	for x := 0.0; x <= size+0.0001; x += boxWidth {
		lineX := boardX + x*cellSize
		pdf.Line(lineX, boardY, lineX, boardY+boardSize)
	}

	if solution != nil {
		solutionString := solution.UniqueId()
		pdf.SetFont("Arial", "B", 10*fontScale)
		pdf.SetTextColor(80, 80, 255)
		pdf.SetXY(boardX, boardY+boardSize+thickLineWidth*5)
		pdf.WriteAligned(boardSize, 12*fontScale, solutionString, "C")
	}
}
