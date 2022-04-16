package rest

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	su "github.com/ClickerMonkey/sudogo/pkg"
)

type PuzzleParams struct {
	FormatParam
	IDParam
}

// =====================================================
// GET /puzzle/{format}/{id}
// =====================================================

type PuzzleFormatSingleQuery struct {
	Candidates Trim[bool] `json:"candidates"`
	State      Trim[bool] `json:"state"`
	Solution   Trim[bool] `json:"solution"`
}

type PuzzleFormatSingleJson struct {
	BoxWidth   int       `json:"boxWidth"`
	BoxHeight  int       `json:"boxHeight"`
	Puzzle     [][]int   `json:"puzzle"`
	Candidates [][][]int `json:"candidates,omitempty"`
}

func DoPuzzleFormatSingle(r JsonRequest[None, PuzzleParams, PuzzleFormatSingleQuery]) (any, int) {
	puzzle, puzzleExists := r.Validate["Puzzle"].(*su.Puzzle)

	if !puzzleExists {
		return nil, http.StatusNotFound
	}

	switch r.Params.Format {
	case FormatJson:
		rsp := PuzzleFormatSingleJson{}
		rsp.Puzzle = puzzle.GetAll()
		rsp.BoxWidth = puzzle.Kind.BoxSize.Width
		rsp.BoxHeight = puzzle.Kind.BoxSize.Height
		if r.Query.Candidates.Value {
			rsp.Candidates = puzzle.GetCandidates()
		}

		return rsp, http.StatusOK

	case FormatText:
		text := ""
		if r.Query.Candidates.Value {
			text = puzzle.ToConsoleCandidatesString()
		} else {
			text = puzzle.ToConsoleString()
		}

		return r.SendText(text, http.StatusOK)

	case FormatPDF:
		pdf := su.NewPDF()
		pdf.Add(puzzle, r.Query.Candidates.Value, r.Query.State.Value, r.Query.Solution.Value)

		return pdf.Send(r.Response, true)
	}

	return nil, http.StatusOK
}

// =====================================================
// GET /solve/{format}/{id}
// =====================================================

type SolveFormatSingleJson struct {
	BoxWidth  int               `json:"boxWidth"`
	BoxHeight int               `json:"boxHeight"`
	Puzzle    [][]int           `json:"puzzle"`
	Solution  [][]int           `json:"solution"`
	Unique    bool              `json:"unique"`
	Duration  string            `json:"duration"`
	Cost      int               `json:"cost"`
	Steps     []PuzzleSolveStep `json:"steps,omitempty"`
}

type SolveFormatSingleQuery struct {
	Steps      Trim[bool] `json:"steps"`
	States     Trim[bool] `json:"stepStates"`
	Candidates Trim[bool] `json:"stepCandidates"`
}

func DoSolveFormatSingle(r JsonRequest[None, PuzzleParams, SolveFormatSingleQuery]) (any, int) {
	puzzle, puzzleExists := r.Validate["Puzzle"].(*su.Puzzle)

	if !puzzleExists {
		return nil, http.StatusNotFound
	}

	steps := r.Query.Steps.Value
	states := r.Query.States.Value

	start := time.Now()

	solutions := puzzle.GetSolutions(su.SolutionsLimit{
		MaxSolutions: 2,
		LogEnabled:   steps,
		LogState:     states,
	})

	duration := time.Since(start)

	if len(solutions) == 0 {
		return nil, http.StatusNotFound
	}

	first := solutions[0]

	switch r.Params.Format {
	case FormatJson:
		rsp := SolveFormatSingleJson{}
		rsp.Puzzle = puzzle.GetAll()
		rsp.BoxWidth = puzzle.Kind.BoxSize.Width
		rsp.BoxHeight = puzzle.Kind.BoxSize.Height
		rsp.Solution = first.Puzzle.GetAll()
		rsp.Unique = len(solutions) == 1
		rsp.Duration = duration.String()
		rsp.Cost = first.GetLastLog().RunningCost
		if steps {
			rsp.Steps = toPuzzleSteps(first, false)
		}

		return rsp, http.StatusOK

	case FormatText:
		sb := strings.Builder{}

		sb.WriteString(first.Puzzle.ToConsoleString())
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("Duration: %v\n", duration))
		sb.WriteString(fmt.Sprintf("Unique: %v\n", len(solutions) == 1))
		sb.WriteString(fmt.Sprintf("Cost: %d\n", first.GetLastLog().RunningCost))
		sb.WriteString("\n")

		for _, log := range first.Logs {
			sb.WriteString(log.String())
			sb.WriteString("\n")

			if states && log.State != nil {
				if r.Query.Candidates.Value {
					sb.WriteString(log.State.ToConsoleCandidatesString())
				} else {
					sb.WriteString(log.State.ToConsoleString())
				}
				sb.WriteString("\n")
				sb.WriteString("\n")
			}
		}

		return r.SendText(sb.String(), http.StatusOK)

	case FormatPDF:
		pdf := su.NewPDF()
		pdf.Add(&first.Puzzle, false, false, false)

		return pdf.Send(r.Response, true)
	}

	return nil, http.StatusOK
}

// =====================================================
// GET /generate/{format}
// =====================================================

type GenerateFormatSingleJson struct {
	Seed    int64             `json:"seed"`
	Puzzles []GeneratedPuzzle `json:"puzzles"`
}

type GenerateFormatSingleQuery struct {
	GenerateKind
	PDF OptionsPDF `json:"pdf"`
}

func DoGenerateFormatSingle(r JsonRequest[None, FormatParam, GenerateFormatSingleQuery]) (any, int) {
	kind := r.Query.GenerateKind
	generated := doKindGeneration(kind, nil, nil)

	switch r.Params.Format {
	case FormatJson:
		rsp := GenerateFormatSingleJson{}
		rsp.Seed = int64(kind.Seed.Value)
		for _, g := range generated {
			rsp.Puzzles = append(rsp.Puzzles, g)
		}

		return rsp, http.StatusOK

	case FormatText:
		sb := strings.Builder{}

		writeGeneratedKind(&sb, kind, generated)

		return r.SendText(sb.String(), http.StatusOK)

	case FormatPDF:
		pdf := su.NewPDF()
		pdf.PuzzlesWide = r.Query.PDF.PuzzlesWide.Value
		pdf.PuzzlesHigh = r.Query.PDF.PuzzlesHigh.Value
		for _, g := range generated {
			pdf.Add(g.Puzzle.Puzzle, kind.Candidates.Value, kind.State.Value, kind.Solutions.Value)
		}

		return pdf.Send(r.Response, true)
	}

	return nil, -1
}

// =====================================================
// POST /generate/{format}
// =====================================================

type GenerateFormatManyJson struct {
	Seed    int64             `json:"seed"`
	Kind    GenerateKind      `json:"kind"`
	Puzzles []GeneratedPuzzle `json:"puzzles"`
}

type GenerateFormatManyQuery struct {
	PDF OptionsPDF `json:"pdf"`
}

func DoGenerateFormatMany(r JsonRequest[GenerateRequest, FormatParam, GenerateFormatManyQuery]) (any, int) {
	rsp := doKindsGeneration(r.Body.Kinds, r.Body.Seed)

	switch r.Params.Format {
	case FormatJson:
		return rsp, http.StatusOK

	case FormatText:
		sb := strings.Builder{}

		for _, run := range rsp {
			writeGeneratedKind(&sb, run.Kind, run.Puzzles)
		}

		return r.SendText(sb.String(), http.StatusOK)

	case FormatPDF:
		pdf := su.NewPDF()
		pdf.PuzzlesWide = r.Query.PDF.PuzzlesWide.Value
		pdf.PuzzlesHigh = r.Query.PDF.PuzzlesHigh.Value
		for _, run := range rsp {
			for _, g := range run.Puzzles {
				pdf.Add(g.Puzzle.Puzzle, run.Kind.Candidates.Value, run.Kind.State.Value, run.Kind.Solutions.Value)
			}
		}

		return pdf.Send(r.Response, true)
	}

	return nil, -1
}

// =====================================================
// POST /solution/{format}
// =====================================================

type SolveFormatComplexJson struct {
	BoxWidth  int               `json:"boxWidth"`
	BoxHeight int               `json:"boxHeight"`
	Puzzle    [][]int           `json:"puzzle"`
	Solution  [][]int           `json:"solution"`
	Unique    bool              `json:"unique"`
	Duration  string            `json:"duration"`
	Cost      int               `json:"cost"`
	Steps     []PuzzleSolveStep `json:"steps,omitempty"`
}

type SolveFormatComplexQuery struct {
	Steps      Trim[bool] `json:"steps"`
	States     Trim[bool] `json:"stepStates"`
	Candidates Trim[bool] `json:"stepCandidates"`
}

func DoSolveFormatComplex(r JsonRequest[SolveKind, FormatParam, SolveFormatSingleQuery]) (any, int) {
	puzzle, solveLimit := r.Body.toDomain()

	steps := r.Query.Steps.Value
	states := r.Query.States.Value

	start := time.Now()

	solutions := puzzle.GetSolutions(su.SolutionsLimit{
		SolveLimit:   solveLimit,
		MaxSolutions: 2,
		LogEnabled:   steps,
		LogState:     states,
	})

	duration := time.Since(start)

	if len(solutions) == 0 {
		return nil, http.StatusNotFound
	}

	first := solutions[0]

	switch r.Params.Format {
	case FormatJson:
		rsp := SolveFormatComplexJson{}
		rsp.Puzzle = puzzle.GetAll()
		rsp.BoxWidth = puzzle.Kind.BoxSize.Width
		rsp.BoxHeight = puzzle.Kind.BoxSize.Height
		rsp.Solution = first.Puzzle.GetAll()
		rsp.Unique = len(solutions) == 1
		rsp.Duration = duration.String()
		rsp.Cost = first.GetLastLog().RunningCost
		if steps {
			rsp.Steps = toPuzzleSteps(first, false)
		}

		return rsp, http.StatusOK

	case FormatText:
		sb := strings.Builder{}

		sb.WriteString(first.Puzzle.ToConsoleString())
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("Duration: %v\n", duration))
		sb.WriteString(fmt.Sprintf("Unique: %v\n", len(solutions) == 1))
		sb.WriteString(fmt.Sprintf("Cost: %d\n", first.GetLastLog().RunningCost))
		sb.WriteString("\n")

		for _, log := range first.Logs {
			sb.WriteString(log.String())
			sb.WriteString("\n")
			if states && log.State != nil {
				if r.Query.Candidates.Value {
					sb.WriteString(log.State.ToConsoleCandidatesString())
				} else {
					sb.WriteString(log.State.ToConsoleString())
				}
				sb.WriteString("\n")
				sb.WriteString("\n")
			}
		}

		return r.SendText(sb.String(), http.StatusOK)

	case FormatPDF:
		pdf := su.NewPDF()
		pdf.Add(&first.Puzzle, false, false, false)

		return pdf.Send(r.Response, true)
	}

	return nil, http.StatusOK
}

// =====================================================
// POST /solutions/{format}
// =====================================================

type SolutionsFormatComplexJson struct {
	BoxWidth  int                              `json:"boxWidth"`
	BoxHeight int                              `json:"boxHeight"`
	Puzzle    [][]int                          `json:"puzzle"`
	Duration  string                           `json:"duration"`
	Solutions []SolutionsFormatComplexSolution `json:"solutions"`
}

type SolutionsFormatComplexSolution struct {
	Solution [][]int           `json:"solution"`
	Cost     int               `json:"cost"`
	Steps    []PuzzleSolveStep `json:"steps,omitempty"`
}

type SolutionsFormatComplexQuery struct {
	Steps  Trim[bool] `json:"steps"`
	States Trim[bool] `json:"stepStates"`
	Limit  Trim[int]  `json:"limit"`
}

func (q *SolutionsFormatComplexQuery) Validate(v Validator) {
	initAndValidate(&q.Limit.Value, 256, 1, 1024, v.Field("limit"))
}

func DoSolutionsFormatComplex(r JsonRequest[SolveKind, FormatParam, SolutionsFormatComplexQuery]) (any, int) {
	puzzle, solveLimit := r.Body.toDomain()

	steps := r.Query.Steps.Value
	states := r.Query.States.Value
	limit := r.Query.Limit.Value

	start := time.Now()

	solutions := puzzle.GetSolutions(su.SolutionsLimit{
		SolveLimit:   solveLimit,
		MaxSolutions: limit,
		LogEnabled:   steps,
		LogState:     states,
	})

	duration := time.Since(start)

	if len(solutions) == 0 {
		return nil, http.StatusNotFound
	}

	switch r.Params.Format {
	case FormatJson:
		rsp := SolutionsFormatComplexJson{}
		rsp.Puzzle = puzzle.GetAll()
		rsp.BoxWidth = puzzle.Kind.BoxSize.Width
		rsp.BoxHeight = puzzle.Kind.BoxSize.Height
		rsp.Duration = duration.String()
		rsp.Solutions = make([]SolutionsFormatComplexSolution, len(solutions))

		for i, solution := range solutions {
			rsp.Solutions[i].Cost = solution.GetLastLog().Cost
			rsp.Solutions[i].Solution = solution.Puzzle.GetAll()
			if steps {
				rsp.Solutions[i].Steps = toPuzzleSteps(solution, false)
			}
		}

		return rsp, http.StatusOK

	case FormatText:
		sb := strings.Builder{}

		sb.WriteString(fmt.Sprintf("Solutions: %v\n", len(solutions)))
		sb.WriteString(fmt.Sprintf("Duration: %v\n\n", duration))

		for i, solution := range solutions {
			sb.WriteString(fmt.Sprintf("\nSolution #%d\n", i+1))
			sb.WriteString(solution.Puzzle.ToConsoleString())
			sb.WriteString("\n")
			sb.WriteString(fmt.Sprintf("Cost: %d\n\n", solution.GetLastLog().RunningCost))

			for _, log := range solution.Logs {
				sb.WriteString(log.String())
				sb.WriteString("\n")
				if states && log.State != nil {
					sb.WriteString(log.State.String())
					sb.WriteString("\n")
					sb.WriteString("\n")
				}
			}
		}

		return r.SendText(sb.String(), http.StatusOK)

	case FormatPDF:
		pdf := su.NewPDF()
		for _, solution := range solutions {
			pdf.Add(&solution.Puzzle, false, false, false)
		}

		return pdf.Send(r.Response, true)
	}

	return nil, http.StatusOK
}

// Helper functions

func doKindsGeneration(gens []GenerateKind, seed GenerateSeed) []GenerateFormatManyJson {
	kinds := make([]GenerateFormatManyJson, len(gens))
	previousSeed := int64(seed)
	previousRandom := su.RandomSeeded(previousSeed)

	for index, g := range gens {
		kinds[index].Seed = int64(g.Seed.Value)
		kinds[index].Kind = g
		kinds[index].Puzzles = doKindGeneration(g, &previousSeed, previousRandom)
	}

	return kinds
}

func doKindGeneration(gen GenerateKind, previousSeed *int64, previousRandom *rand.Rand) []GeneratedPuzzle {
	rand := previousRandom
	randSeed := int64(gen.Seed.Value)

	if rand == nil || previousSeed == nil || *previousSeed != randSeed {
		rand = su.RandomSeeded(randSeed)

		if previousSeed != nil {
			*previousSeed = randSeed
		}
		if previousRandom != nil {
			*previousRandom = *rand
		}
	}

	puzzles := []GeneratedPuzzle{}
	kind, clear := gen.toDomain()
	candidates := gen.Candidates.Value

	clear.Fast = true
	clear.MaxStates = gen.TryClears.Value

	count := int(gen.Count.Value)

	type GenerationAttempt struct {
		Generated *su.Puzzle
		Cleared   *su.Puzzle
		Solver    *su.Solver
	}

	for i := 0; i < count; i++ {

		start := time.Now()
		result := make(chan *GenerationAttempt)

		var wg sync.WaitGroup
		wg.Add(gen.TryCount.Value)

		for k := 0; k < gen.TryCount.Value; k++ {
			go func() {
				defer wg.Done()

				g := su.NewRandomGenerator(kind, rand)
				g.Solver().LogEnabled = gen.SolutionSteps.Value
				g.Solver().LogState = gen.SolutionStates.Value

				generated, _ := g.Attempts(gen.TryAttempts.Value)

				if generated != nil {
					cleared, _ := g.ClearCells(generated, clear)

					if cleared != nil {
						solver := g.Solver()

						result <- &GenerationAttempt{
							Generated: generated,
							Cleared:   cleared,
							Solver:    solver,
						}
					}
				}
			}()
		}

		go func() {
			wg.Wait()
			close(result)
		}()

		found, success := <-result

		if success {
			duration := time.Since(start)

			pg := GeneratedPuzzle{}
			pg.Difficulty = gen.Difficulty
			pg.Duration = duration.String()
			pg.Puzzle = toPuzzleData(found.Cleared, candidates)
			pg.Solver = found.Solver

			if gen.Solutions.Value {
				solution := toPuzzleData(found.Generated, false)
				pg.Solution = &solution
			}
			if found.Solver.LogEnabled || found.Solver.LogState {
				pg.Steps = toPuzzleSteps(found.Solver, candidates)
			}

			puzzles = append(puzzles, pg)
		}
	}

	return puzzles
}

func toPuzzleData(puzzle *su.Puzzle, candidates bool) PuzzleData {
	pd := PuzzleData{}
	pd.Puzzle = puzzle
	pd.BoxWidth = puzzle.Kind.BoxSize.Width
	pd.BoxHeight = puzzle.Kind.BoxSize.Height
	pd.Values = puzzle.GetAll()
	pd.Encoded = puzzle.EncodedString()
	pd.State = puzzle.String()

	if candidates {
		pd.Candidates = puzzle.GetCandidates()
	}

	return pd
}

func toPuzzleSteps(solver *su.Solver, candidates bool) []PuzzleSolveStep {
	steps := make([]PuzzleSolveStep, 0, len(solver.Logs))

	for _, log := range solver.Logs {
		steps = append(steps, toPuzzleStep(log, solver.LogState, candidates))
	}

	return steps
}

func toPuzzleStep(log su.SolverLog, state bool, candidates bool) PuzzleSolveStep {
	step := PuzzleSolveStep{
		Technique:         log.Step.Technique,
		Index:             log.Index,
		Batch:             log.Batch,
		Cost:              log.Cost,
		Placement:         log.Placement,
		Row:               log.Before.Row,
		Col:               log.Before.Col,
		Before:            log.Before.Value,
		BeforeCandidates:  log.Before.Candidates(),
		After:             log.After.Value,
		AfterCandidates:   log.After.Candidates(),
		RunningCost:       log.RunningCost,
		RunningPlacements: log.RunningPlacements,
		Log:               log,
	}

	if state {
		s := toPuzzleData(log.State, candidates)
		step.State = &s
	}

	return step
}

func writeGeneratedKind(sb *strings.Builder, kind GenerateKind, generated []GeneratedPuzzle) {
	sb.WriteString(fmt.Sprintf("Seed: %d\n", kind.Seed.Value))

	for index, g := range generated {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("Puzzle #%d\n", index+1))

		if kind.State.Value {
			sb.WriteString(g.Puzzle.Puzzle.ToStateString(true, "."))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")

		if kind.Candidates.Value {
			sb.WriteString(g.Puzzle.Puzzle.ToConsoleCandidatesString())
		} else {
			sb.WriteString(g.Puzzle.Puzzle.ToConsoleString())
		}
		sb.WriteString("\n\n")

		if g.Difficulty != "" {
			sb.WriteString(fmt.Sprintf("Difficulty: %v\n", g.Difficulty))
		}
		sb.WriteString(fmt.Sprintf("Duration: %v\n\n", g.Duration))

		if kind.Solutions.Value {
			sb.WriteString("Solution:\n")

			sb.WriteString(g.Solution.Puzzle.ToConsoleString())
			sb.WriteString("\n\n")
		}

		if kind.SolutionSteps.Value || kind.SolutionStates.Value {
			sb.WriteString("Steps:\n")

			for _, log := range g.Steps {
				if kind.SolutionSteps.Value {
					sb.WriteString(log.Log.String())
					sb.WriteString("\n")
				}
				if kind.SolutionStates.Value {
					sb.WriteString(log.Log.State.ToConsoleString())
					sb.WriteString("\n")
				}
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
	}
}
