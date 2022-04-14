package rest

import (
	"strings"

	su "github.com/ClickerMonkey/sudogo/pkg"
)

type Index int

func (d Index) Validate(v Validator) {
	if kind, ok := v.Context["Kind"].(*su.Kind); ok {
		s := Index(kind.Size())
		if d < 0 {
			v.Add("Cannot be less than 0: %d", d)
		}
		if d >= s {
			v.Add("Cannot be greater than %d: %d", s, d)
		}
	}
}

type RelativeIndex int

func (d RelativeIndex) Validate(v Validator) {
	if kind, ok := v.Context["Kind"].(*su.Kind); ok {
		s := RelativeIndex(kind.Size())
		if d <= -s {
			v.Add("Cannot be less than %d: %d", -s, d)
		}
		if d >= s {
			v.Add("Cannot be greater than %d: %d", s, d)
		}
	}
}

type PuzzleDimension int

func (d PuzzleDimension) Validate(v Validator) {
	if d < 2 {
		v.Add("cannot be less than 2: %d", d)
	}
	if d > 64 {
		v.Add("cannot be greater than 64: %d", d)
	}
}

type Position struct {
	Col Index `json:"col"`
	Row Index `json:"row"`
}

func (p Position) toDomain() su.Position {
	return su.Position{Col: int(p.Col), Row: int(p.Row)}
}

type RelativePosition struct {
	Col RelativeIndex `json:"col"`
	Row RelativeIndex `json:"row"`
}

func (p RelativePosition) toDomain() su.Position {
	return su.Position{Col: int(p.Col), Row: int(p.Row)}
}

type ConstraintSumValue struct {
	Sum      int                 `json:"sum"`
	Cells    *[]Position         `json:"cells"`
	Relative *[]RelativePosition `json:"relative"`
}

func (p ConstraintSumValue) toDomain() su.Constraint {
	return &su.ConstraintSum{
		Sum:      su.SumConstant(p.Sum),
		Cells:    toDomainSlicePointer[su.Position](p.Cells),
		Relative: toDomainSlicePointer[su.Position](p.Relative),
	}
}

func (c ConstraintSumValue) Validate(v Validator) {
	if c.Sum <= 1 {
		v.Add("The sum %d is not a valid constraint sum.", c.Sum)
	}
	if c.Cells == nil && c.Relative == nil {
		v.Add("No cells specified.")
	}
}

type ConstraintSumCell struct {
	Sum         RelativePosition    `json:"sum"`
	SumRelative Trim[bool]          `json:"sumRelative"`
	Cells       *[]Position         `json:"cells"`
	Relative    *[]RelativePosition `json:"relative"`
}

func (p ConstraintSumCell) toDomain() su.Constraint {
	return &su.ConstraintSum{
		Sum:      su.SumCell(p.Sum.toDomain(), p.SumRelative.Value),
		Cells:    toDomainSlicePointer[su.Position](p.Cells),
		Relative: toDomainSlicePointer[su.Position](p.Relative),
	}
}

type ConstraintUnique struct {
	Cells    *[]Position         `json:"cells"`
	Relative *[]RelativePosition `json:"relative"`
	Same     Trim[bool]          `json:"same"`
}

func (p ConstraintUnique) toDomain() su.Constraint {
	return &su.ConstraintUnique{
		Cells:    toDomainSlicePointer[su.Position](p.Cells),
		Relative: toDomainSlicePointer[su.Position](p.Relative),
		Same:     p.Same.Value,
	}
}

type DirectionInt int

func (d DirectionInt) Validate(v Validator) {
	if d < -1 || d > 1 {
		v.Add("must be between -1 and 1: %d", d)
	}
}

type ConstraintOrder struct {
	Cells     *[]Position         `json:"cells"`
	Relative  *[]RelativePosition `json:"relative"`
	Direction DirectionInt        `json:"direction"`
}

func (p ConstraintOrder) toDomain() su.Constraint {
	return &su.ConstraintOrder{
		Cells:     toDomainSlicePointer[su.Position](p.Cells),
		Relative:  toDomainSlicePointer[su.Position](p.Relative),
		Direction: int(p.Direction),
	}
}

type ConstraintMagic struct {
	Center Position `json:"center"`
}

func (p ConstraintMagic) toDomain() su.Constraint {
	return &su.ConstraintMagic{
		Center: p.Center.toDomain(),
	}
}

type ConstraintScalePair struct {
	Scale  int      `json:"scale"`
	First  Position `json:"first"`
	Second Position `json:"second"`
}

func (p ConstraintScalePair) toDomain() su.Constraint {
	return &su.ConstraintScalePair{
		Scale:  p.Scale,
		First:  p.First.toDomain(),
		Second: p.Second.toDomain(),
	}
}

type ConstraintDifference struct {
	Min      int                 `json:"min"`
	Max      int                 `json:"max"`
	Cells    *[]Position         `json:"cells"`
	Relative *[]RelativePosition `json:"relative"`
	Exclude  *[]Position         `json:"exclude"`
}

func (p ConstraintDifference) toDomain() su.Constraint {
	return &su.ConstraintDifference{
		Min:      p.Min,
		Max:      p.Max,
		Cells:    toDomainSlicePointer[su.Position](p.Cells),
		Relative: toDomainSlicePointer[su.Position](p.Relative),
		Exclude:  toDomainSlicePointer[su.Position](p.Exclude),
	}
}

type ConstraintDivisible struct {
	By        int        `json:"by"`
	Remainder int        `json:"remainder"`
	Cells     []Position `json:"cells"`
}

func (p ConstraintDivisible) toDomain() su.Constraint {
	return &su.ConstraintDivisible{
		By:        p.By,
		Remainder: p.Remainder,
		Cells:     toDomainSlice[su.Position](p.Cells),
	}
}

type Constraints struct {
	SumValues   []ConstraintSumValue   `json:"sumValues"`
	SumCells    []ConstraintSumCell    `json:"sumCells"`
	Uniques     []ConstraintUnique     `json:"uniques"`
	Orders      []ConstraintOrder      `json:"orders"`
	Magics      []ConstraintMagic      `json:"magics"`
	ScalePairs  []ConstraintScalePair  `json:"scalePairs"`
	Differences []ConstraintDifference `json:"differences"`
	Divisibles  []ConstraintDivisible  `json:"divisibles"`
}

func (c Constraints) toDomain() []su.Constraint {
	d := make([]su.Constraint, 0)
	d = append(d, toDomainSlice[su.Constraint](c.SumValues)...)
	d = append(d, toDomainSlice[su.Constraint](c.SumCells)...)
	d = append(d, toDomainSlice[su.Constraint](c.Uniques)...)
	d = append(d, toDomainSlice[su.Constraint](c.Orders)...)
	d = append(d, toDomainSlice[su.Constraint](c.Magics)...)
	d = append(d, toDomainSlice[su.Constraint](c.ScalePairs)...)
	d = append(d, toDomainSlice[su.Constraint](c.Differences)...)
	d = append(d, toDomainSlice[su.Constraint](c.Divisibles)...)
	return d
}

type GenerateCount int

func (g GenerateCount) Validate(v Validator) {
	if g < 1 {
		v.Add("cannot be less than 1: %d", g)
	}
	if g > 1000 {
		v.Add("cannot be greater than 1000: %d", g)
	}
}

type GenerateSeed int64

func (g *GenerateSeed) Validate(v Validator) {
	seed, seedExists := v.Context["Seed"].(int64)
	if *g == 0 {
		if seedExists {
			*g = GenerateSeed(seed)
		} else {
			*g = GenerateSeed(su.RandomSeed())
		}
	}
	if !seedExists {
		v.Context["Seed"] = seed
	}
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

type GenerateKind struct {
	Count          Trim[GenerateCount]   `json:"count"`
	Seed           Trim[GenerateSeed]    `json:"seed"`
	Difficulty     string                `json:"difficulty"`
	MinCost        Trim[int]             `json:"minCost"`
	MaxCost        Trim[int]             `json:"maxCost"`
	MaxPlacements  Trim[int]             `json:"maxPlacements"`
	MaxSteps       Trim[int]             `json:"maxSteps"`
	MaxBatches     Trim[int]             `json:"maxBatches"`
	Symmetric      Trim[bool]            `json:"symmetric"`
	BoxWidth       Trim[PuzzleDimension] `json:"boxWidth"`
	BoxHeight      Trim[PuzzleDimension] `json:"boxHeight"`
	Constraints    Constraints           `json:"constraints"`
	Candidates     Trim[bool]            `json:"candidates"`
	State          Trim[bool]            `json:"state"`
	Solutions      Trim[bool]            `json:"solutions"`
	SolutionSteps  Trim[bool]            `json:"solutionSteps"`
	SolutionStates Trim[bool]            `json:"solutionStates"`
}

func (r *GenerateKind) Validate(v Validator) {
	if r.BoxWidth.Value == 0 {
		r.BoxWidth.Value = 3
	}
	if r.BoxHeight.Value == 0 {
		r.BoxHeight.Value = 3
	}
	if r.Count.Value == 0 {
		r.Count.Value = 1
	}

	v.Context["Kind"] = su.NewKind(int(r.BoxWidth.Value), int(r.BoxHeight.Value))
}

func (r GenerateKind) toDomain() (*su.Kind, su.ClearLimit) {
	boxWidth := su.Max(1, int(r.BoxWidth.Value))
	boxHeight := su.Max(1, int(r.BoxHeight.Value))
	kind := su.NewKind(boxWidth, boxHeight)
	kind.Constraints = r.Constraints.toDomain()
	limitScale := float32(kind.Area()) / 81.0

	clear := su.ClearLimit{}

	if d, ok := DifficultyMap[r.Difficulty]; ok {
		clear = d
	}
	if r.Symmetric.Value {
		clear.Symmetric = r.Symmetric.Value
	}

	applyValue := func(user int, out *int) {
		if user != 0 {
			*out = user
		} else if user < 0 {
			*out = 0
		} else {
			*out = int(float32(*out) * limitScale)
		}
	}

	applyValue(r.MaxBatches.Value, &clear.SolveLimit.MaxBatches)
	applyValue(r.MaxCost.Value, &clear.SolveLimit.MaxCost)
	applyValue(r.MinCost.Value, &clear.SolveLimit.MinCost)
	applyValue(r.MaxSteps.Value, &clear.SolveLimit.MaxLogs)
	applyValue(r.MaxPlacements.Value, &clear.SolveLimit.MaxPlacements)

	return kind, clear
}

type GenerateRequest struct {
	Seed  GenerateSeed   `json:"seed"`
	Kinds []GenerateKind `json:"kinds"`
}

func (r GenerateRequest) Validate(v Validator) {
	if len(r.Kinds) == 0 {
		v.AddField("kinds", "Kinds cannot be empty.")
	}
}

type Format string

var (
	FormatText Format = "text"
	FormatPDF  Format = "pdf"
	FormatJson Format = "json"
)

type FormatParam struct {
	Format Format `json:"format"`
}

func (f FormatParam) Validate(v Validator) {
	format := strings.ToLower(string(f.Format))
	switch format {
	case "text", "json", "pdf":
		f.Format = Format(format)
	default:
		v.AddField("format", "Invalid format: %s; Only txt, json, or pdf are supported.", format)
	}
}

type IDParam struct {
	ID string `json:"id"`
}

func (p IDParam) Validate(v Validator) {
	parsed := su.FromString(p.ID)
	if parsed == nil {
		parsed = su.FromEncoded(p.ID)
	}
	if parsed == nil {
		v.AddField("ID", "invalid puzzle identifier")
	}

	v.Context["Puzzle"] = parsed
}

type GeneratedPuzzle struct {
	Difficulty string            `json:"difficulty,omitempty"`
	Puzzle     PuzzleData        `json:"puzzle"`
	Duration   string            `json:"duration"`
	Solution   *PuzzleData       `json:"solution,omitempty"`
	Steps      []PuzzleSolveStep `json:"steps,omitempty"`

	Solver *su.Solver `json:"-"`
}

type PuzzleData struct {
	BoxWidth   int       `json:"boxWidth,omitempty"`
	BoxHeight  int       `json:"boxHeight,omitempty"`
	Values     [][]int   `json:"values"`
	Candidates [][][]int `json:"candidates,omitempty"`
	Encoded    string    `json:"encoded,omitempty"`
	State      string    `json:"state,omitempty"`

	Puzzle *su.Puzzle `json:"-"`
}

type PuzzleSolveStep struct {
	Technique         string      `json:"technique"`
	Index             int         `json:"index"`
	Batch             int         `json:"batch"`
	Cost              int         `json:"cost"`
	Placement         bool        `json:"placement"`
	Row               int         `json:"row"`
	Col               int         `json:"col"`
	Before            int         `json:"before"`
	BeforeCandidates  []int       `json:"beforeCandidates"`
	After             int         `json:"after"`
	AfterCandidates   []int       `json:"afterCandidates"`
	RunningCost       int         `json:"runningCost"`
	RunningPlacements int         `json:"runningPlacements"`
	State             *PuzzleData `json:"state,omitempty"`

	Log su.SolverLog `json:"-"`
}

type OptionsPDF struct {
	PuzzlesWide Trim[int] `json:"puzzlesWide"`
	PuzzlesHigh Trim[int] `json:"puzzlesHigh"`
}

func (o OptionsPDF) Validate(v Validator) {
	if o.PuzzlesWide.Value < 0 || o.PuzzlesWide.Value > 3 {
		v.AddField("puzzlesWide", "Puzzles wide must be between 1 and 3.")
	}
	if o.PuzzlesHigh.Value < 0 || o.PuzzlesHigh.Value > 4 {
		v.AddField("puzzlesHigh", "Puzzles high must be between 1 and 4.")
	}
}
