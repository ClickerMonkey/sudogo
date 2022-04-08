package rest

import (
	"strconv"
	"strings"

	su "github.com/ClickerMonkey/sudogo/pkg"
)

type Index int

func (d Index) Validate(v Validator) {
	if kind, ok := (*v.Context)["Kind"].(*su.Kind); ok {
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
	if kind, ok := (*v.Context)["Kind"].(*su.Kind); ok {
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

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	asString := strings.Trim(string(data), `"`)
	parsed, err := strconv.ParseBool(asString)
	if err == nil {
		*b = Bool(parsed)
	}
	return err
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
	SumRelative Bool                `json:"sumRelative"`
	Cells       *[]Position         `json:"cells"`
	Relative    *[]RelativePosition `json:"relative"`
}

func (p ConstraintSumCell) toDomain() su.Constraint {
	return &su.ConstraintSum{
		Sum:      su.SumCell(p.Sum.toDomain(), bool(p.SumRelative)),
		Cells:    toDomainSlicePointer[su.Position](p.Cells),
		Relative: toDomainSlicePointer[su.Position](p.Relative),
	}
}

type ConstraintUnique struct {
	Cells    *[]Position         `json:"cells"`
	Relative *[]RelativePosition `json:"relative"`
	Same     Bool                `json:"same"`
}

func (p ConstraintUnique) toDomain() su.Constraint {
	return &su.ConstraintUnique{
		Cells:    toDomainSlicePointer[su.Position](p.Cells),
		Relative: toDomainSlicePointer[su.Position](p.Relative),
		Same:     bool(p.Same),
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

type GenerateBase struct {
	Count         GenerateCount   `json:"count"`
	MinCost       int             `json:"minCost"`
	MaxCost       int             `json:"maxCost"`
	MaxPlacements int             `json:"maxPlacements"`
	MaxLogs       int             `json:"maxLogs"`
	MaxBatches    int             `json:"maxBatches"`
	Symmetric     Bool            `json:"symmetric"`
	BoxWidth      PuzzleDimension `json:"boxWidth"`
	BoxHeight     PuzzleDimension `json:"boxHeight"`
	Constraints   Constraints     `json:"constraints"`
	Candidates    Bool            `json:"candidates"`
	Solutions     Bool            `json:"solutions"`
}

func (r GenerateBase) Validate(v Validator) {
	(*v.Context)["Kind"] = su.NewKind(int(r.BoxWidth), int(r.BoxHeight))
}

func (r GenerateBase) toDomain() (*su.Kind, su.SolveLimit, su.ClearLimit) {
	kind := su.NewKind(int(r.BoxWidth), int(r.BoxHeight))
	kind.Constraints = r.Constraints.toDomain()

	solve := su.SolveLimit{}
	solve.MaxBatches = r.MaxBatches
	solve.MaxCost = r.MaxCost
	solve.MinCost = r.MinCost
	solve.MaxLogs = r.MaxLogs
	solve.MaxPlacements = r.MaxPlacements

	clear := su.ClearLimit{}
	clear.SolveLimit = solve
	clear.Symmetric = bool(r.Symmetric)

	return kind, solve, clear
}

type GenerateJson struct {
	GenerateBase
	SolutionLogs   Bool `json:"solutionLogs"`
	SolutionStates Bool `json:"solutionStates"`
}

type GenerateJsonRequest struct {
	Seed    int64          `json:"seed"`
	Puzzles []GenerateJson `json:"types"`
}
