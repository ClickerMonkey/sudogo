package rest

import (
	"strings"
	"time"

	su "github.com/ClickerMonkey/sudogo/pkg"
)

type PuzzleSolveLog struct {
	Step              string `json:"step"`
	Index             int    `json:"index"`
	Batch             int    `json:"batch"`
	Cost              int    `json:"cost"`
	Placement         bool   `json:"placement"`
	Row               int    `json:"row"`
	Col               int    `json:"col"`
	Before            int    `json:"before"`
	BeforeCandidates  []int  `json:"beforeCandidates"`
	After             int    `json:"after"`
	AfterCandidates   []int  `json:"afterCandidates"`
	RunningCost       int    `json:"runningCost"`
	RunningPlacements int    `json:"runningPlacements"`
}

func NewPuzzleSolveLog(log su.SolverLog) PuzzleSolveLog {
	return PuzzleSolveLog{
		Step:              log.Step.Technique,
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
	}
}

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

func DoGenerate(r JsonRequest[GenerateJsonRequest, None, None]) (any, int) {
	rsp := GenerateJsonResponse{}

	rsp.Seed = su.RandomSeed()
	if r.Body.Seed != 0 {
		rsp.Seed = r.Body.Seed
	}

	rand := su.RandomSeeded(rsp.Seed)

	for _, puzzle := range r.Body.Puzzles {
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
	Solution      [][]int          `json:"solution"`
	SolutionCount int              `json:"solutionCount"`
	Duration      string           `json:"duration"`
	Cost          int              `json:"cost"`
	Steps         []PuzzleSolveLog `json:"steps"`
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
	rsp.Steps = make([]PuzzleSolveLog, 0, len(first.Logs))

	for _, log := range first.Logs {
		rsp.Steps = append(rsp.Steps, NewPuzzleSolveLog(log))
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

/*
type TestBody struct {
	X int `json:"x"`
}

func (t TestBody) Validate(v Validator) {
	if t.X < 0 {
		v.Add("X cannot be less than 0: %d", t.X)
	}
}

type TestParams struct {
	Tag string `json:"tag"`
}

func (t TestParams) Validate(v Validator) {
	if len(t.Tag) > 10 {
		v.Add("Tag cannot be longer than 10 characters: %s", t.Tag)
	}
}

type TestQuery struct {
	OrderBy []struct {
		Field string `json:"field"`
		Desc  Trim[bool]   `json:"desc"`
	} `json:"orderBy"`
}

func doTest(r JsonRequest[TestBody, TestParams, TestQuery]) (any, int) {
	return []any{r.Body, r.Params, r.Query}, 200
}
*/
