package rest

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	su "github.com/ClickerMonkey/sudogo/pkg"
)

type PuzzleParams struct {
	FormatParam
	IDParam
}

/*

type PuzzleGenerated struct {
	Values          [][]int    `json:"values"`
	Solution        *[][]int   `json:"solution"`
	SolutionEncoded string     `json:"solutionEncoded"`
	Logs            *[]string  `json:"logs"`
	Candidates      *[][][]int `json:"candidates"`
}

type GenerateJsonResponse struct {
	Seed    int64             `json:"seed"`
	Puzzles []PuzzleGenerated `json:"puzzles"`
}

func DoGenerate(r JsonRequest[GenerateRequest, None, None]) (any, int) {
	rsp := GenerateJsonResponse{}

	rsp.Seed = su.RandomSeed()
	if r.Body.Seed != 0 {
		rsp.Seed = int64(r.Body.Seed)
	}

	rand := su.RandomSeeded(rsp.Seed)

	for _, puzzle := range r.Body.Kinds {
		kind, _, clear := puzzle.toDomain()

		for i := 0; i < int(puzzle.Count); i++ {
			gen := su.NewRandomGenerator(kind, rand)
			gen.Solver().LogEnabled = puzzle.SolutionLogs.Value
			gen.Solver().LogState = puzzle.SolutionStates.Value

			if generated, _ := gen.Generate(); generated != nil {
				if cleared, _ := gen.ClearCells(generated, clear); cleared != nil {

					final := PuzzleGenerated{
						Values:          cleared.GetAll(),
						SolutionEncoded: generated.EncodedString(),
					}

					if puzzle.Solutions.Value {
						solution := generated.GetAll()
						final.Solution = &solution
					}

					if puzzle.SolutionLogs.Value {
						solverLogs := gen.Solver().Logs
						logs := make([]string, 0, len(solverLogs))
						for _, log := range solverLogs {
							logs = append(logs, log.String())
						}
						final.Logs = &logs
					}

					if puzzle.Candidates.Value {
						size := kind.Size()
						cand := make([][][]int, size)
						for r := 0; r < size; r++ {
							cand[r] = make([][]int, size)
							for c := 0; c < size; c++ {
								cand[r][c] = cleared.Get(c, r).Candidates()
							}
						}
						final.Candidates = &cand
					}

					rsp.Puzzles = append(rsp.Puzzles, final)
				}
			}
		}
	}

	return rsp, 200
}

type PuzzleGet struct {
	ID string `json:"id"`
}

func DoPuzzleGet(r JsonRequest[None, PuzzleGet, None]) (any, int) {
	p := su.FromString(r.Params.ID)
	if p != nil {
		return p.GetAll(), 200
	}
	return nil, 400
}

type PuzzleGenerateSimpleQuery struct {
	Difficulty string `json:"d"`
}

type PuzzleGenerateSimpleResponse struct {
	Difficulty string  `json:"difficulty"`
	State      string  `json:"state"`
	Encoded    string  `json:"encoded"`
	Duration   string  `json:"duration"`
	Puzzle     [][]int `json:"puzzle"`
	Solution   [][]int `json:"solution"`
}

var DifficultyMap = map[string]su.ClearLimit{
	"":           su.DifficultyMedium,
	"beginner":   su.DifficultyBeginner,
	"easy":       su.DifficultyEasy,
	"medium":     su.DifficultyMedium,
	"hard":       su.DifficultyHard,
	"tricky":     su.DifficultyTricky,
	"diabolical": su.DifficultyDiabolical,
	"fiendish":   su.DifficultyFiendish,
}

func DoPuzzleGenerateSimple(r JsonRequest[None, None, PuzzleGenerateSimpleQuery]) (any, int) {
	key := strings.ToLower(r.Query.Difficulty)
	limits, exists := DifficultyMap[key]

	if !exists {
		return nil, 404
	}

	start := time.Now()

	generator := su.Classic.Generator()
	generated, _ := generator.Generate()

	if generated == nil {
		return nil, 408
	}

	fastLimits := limits
	fastLimits.Fast = true

	cleared, _ := generator.ClearCells(generated, fastLimits)

	duration := time.Since(start)

	if cleared == nil {
		return nil, 408
	}

	rsp := PuzzleGenerateSimpleResponse{}
	rsp.Difficulty = key
	rsp.Puzzle = cleared.GetAll()
	rsp.Solution = generated.GetAll()
	rsp.State = cleared.ToStateString(true, ".")
	rsp.Encoded = cleared.EncodedString()
	rsp.Duration = duration.String()

	return rsp, 200
}

type PuzzleSolveSimpleParams struct {
	ID string `json:"id"`
}

type PuzzleSolveSimpleResponse struct {
	Solution      [][]int           `json:"solution"`
	SolutionCount int               `json:"solutionCount"`
	Duration      string            `json:"duration"`
	Cost          int               `json:"cost"`
	Steps         []PuzzleSolveStep `json:"steps"`
}

func DoPuzzleSolveSimple(r JsonRequest[None, PuzzleSolveSimpleParams, None]) (any, int) {
	puzzle := su.FromString(r.Params.ID)
	if puzzle == nil {
		return nil, 400
	}

	rsp := PuzzleSolveSimpleResponse{}

	start := time.Now()

	solutions := puzzle.GetSolutions(su.SolutionsLimit{
		MaxSolutions: 20,
		LogEnabled:   true,
	})

	duration := time.Since(start)

	if len(solutions) == 0 {
		return rsp, 200
	}

	first := solutions[0]

	rsp.Solution = first.Puzzle.GetAll()
	rsp.SolutionCount = len(solutions)
	rsp.Duration = duration.String()
	rsp.Cost = first.GetLastLog().RunningCost
	rsp.Steps = make([]PuzzleSolveStep, 0, len(first.Logs))

	for _, log := range first.Logs {
		rsp.Steps = append(rsp.Steps, NewPuzzleSolveStep(log))
	}

	return rsp, 200
}

type PuzzlePDFSimpleQuery struct {
	Difficulty string     `json:"d"`
	Candidates Trim[bool] `json:"c"`
	Count      Trim[int]  `json:"n"`
	Wide       Trim[int]  `json:"w"`
	High       Trim[int]  `json:"h"`
}

func DoPuzzlePDFSimple(r JsonRequest[None, None, PuzzlePDFSimpleQuery]) (any, int) {
	key := strings.ToLower(r.Query.Difficulty)
	limits, exists := DifficultyMap[key]

	if !exists {
		return nil, 404
	}

	pdf := su.NewPDF()
	pdf.Candidates = bool(r.Query.Candidates.Value)
	pdf.PuzzlesHigh = su.Max(1, r.Query.High.Value)
	pdf.PuzzlesWide = su.Max(1, r.Query.Wide.Value)

	count := su.Max(1, r.Query.Count.Value)

	for len(pdf.Puzzles) < count {
		generator := su.Classic.Generator()
		generated, _ := generator.Generate()

		if generated == nil {
			continue
		}

		fastLimits := limits
		fastLimits.Fast = true

		cleared, _ := generator.ClearCells(generated, fastLimits)

		if cleared == nil {
			continue
		}

		pdf.Add(cleared)
	}

	pdf.Send(r.Response, false)

	return nil, -1
}
*/

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
	Steps  Trim[bool] `json:"steps"`
	States Trim[bool] `json:"stepStates"`
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
		text := first.Puzzle.ToConsoleString() + "\n"
		text += fmt.Sprintf("Duration: %v\n", duration)
		text += fmt.Sprintf("Unique: %v\n", len(solutions) == 1)
		text += fmt.Sprintf("Cost: %d\n", first.GetLastLog().RunningCost)
		text += "\n"

		for _, log := range first.Logs {
			text += log.String() + "\n"
		}

		return r.SendText(text, http.StatusOK)

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

		sb.WriteString(fmt.Sprintf("Seed: %d\n", kind.Seed.Value))

		for index, g := range generated {
			sb.WriteString("\n")
			sb.WriteString(fmt.Sprintf("Puzzle #%d\n\n", index+1))

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

		return r.SendText(sb.String(), http.StatusOK)

	case FormatPDF:
		pdf := su.NewPDF()
		pdf.PuzzlesWide = su.Max(1, r.Query.PDF.PuzzlesWide.Value)
		pdf.PuzzlesHigh = su.Max(1, r.Query.PDF.PuzzlesHigh.Value)
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

// =====================================================
// POST /solution/{format}
// =====================================================

// =====================================================
// POST /solutions/{format}
// =====================================================

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
	clear.MaxStates = 256

	count := su.Max(1, int(gen.Count.Value))

	generateTries := 48
	generateAttempts := 256

	for i := 0; i < count; i++ {
		g := su.NewRandomGenerator(kind, rand)
		g.Solver().LogEnabled = gen.SolutionSteps.Value
		g.Solver().LogState = gen.SolutionStates.Value

		start := time.Now()

		for k := 0; k < generateTries; k++ {
			generated, _ := g.Attempts(generateAttempts)

			if generated != nil {
				cleared, _ := g.ClearCells(generated, clear)

				duration := time.Since(start)

				if cleared != nil {
					solver := g.Solver()

					pg := GeneratedPuzzle{}
					pg.Difficulty = gen.Difficulty
					pg.Duration = duration.String()
					pg.Puzzle = toPuzzleData(cleared, candidates)
					pg.Solver = solver

					if gen.Solutions.Value {
						solution := toPuzzleData(generated, false)
						pg.Solution = &solution
					}
					if solver.LogEnabled || solver.LogState {
						pg.Steps = toPuzzleSteps(g.Solver(), candidates)
					}

					puzzles = append(puzzles, pg)

					break
				}
			}
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
