package rest

import (
	su "github.com/ClickerMonkey/sudogo/pkg"
)

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

func DoGenerate(body GenerateJsonRequest, params None, query None) (any, int) {
	rsp := GenerateJsonResponse{}

	rsp.Seed = su.RandomSeed()
	if body.Seed != 0 {
		rsp.Seed = body.Seed
	}

	rand := su.RandomSeeded(rsp.Seed)

	for _, puzzle := range body.Puzzles {
		kind, _, clear := puzzle.toDomain()

		for i := 0; i < int(puzzle.Count); i++ {
			gen := su.NewRandomGenerator(kind, rand)
			gen.Solver().LogEnabled = bool(puzzle.SolutionLogs)
			gen.Solver().LogState = bool(puzzle.SolutionStates)

			if generated, _ := gen.Generate(); generated != nil {
				if cleared, _ := gen.ClearCells(generated, clear); cleared != nil {

					final := PuzzleGenerated{
						Values:          cleared.GetAll(),
						SolutionEncoded: generated.EncodedString(),
					}

					if puzzle.Solutions {
						solution := generated.GetAll()
						final.Solution = &solution
					}

					if puzzle.SolutionLogs {
						solverLogs := gen.Solver().Logs
						logs := make([]string, 0, len(solverLogs))
						for _, log := range solverLogs {
							logs = append(logs, log.String())
						}
						final.Logs = &logs
					}

					if puzzle.Candidates {
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

type TestBody struct {
	X int `json:"x"`
}
type TestParams struct {
	Tag string `json:"tag"`
}
type TestQuery struct {
	OrderBy []struct {
		Field string `json:"field"`
		Desc  Bool   `json:"desc"`
	} `json:"orderBy"`
}

func doTest(body TestBody, param TestParams, query TestQuery) (any, int) {
	return []any{body, param, query}, 200
}
