package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	su "github.com/ClickerMonkey/sudogo/pkg"
)

type Validator struct {
	Path        []string
	Validations *[]Validation
}

func (v Validator) NextInt(next int) Validator {
	return v.Next(strconv.Itoa(next))
}

func (v Validator) Next(next string) Validator {
	return Validator{
		Path:        append(v.Path, next),
		Validations: v.Validations,
	}
}

func (v *Validator) Add(format string, args ...any) {
	*v.Validations = append(*v.Validations, Validation{
		Path:    strings.Join(v.Path, "."),
		Message: fmt.Sprintf(format, args...),
	})
}

func (v *Validator) AddNext(next string, format string, args ...any) {
	*v.Validations = append(*v.Validations, Validation{
		Path:    strings.Join(append(v.Path, next), "."),
		Message: fmt.Sprintf(format, args...),
	})
}

type Validation struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type ValidateContext struct {
	Kind su.Kind
}

type CanValidate interface {
	Validate(ctx *ValidateContext, v Validator)
}

type hasDomain[T any] interface {
	toDomain() T
}

func domainPointer[T any, D hasDomain[T]](domain *D) *T {
	if domain == nil {
		return nil
	}
	value := (*domain).toDomain()
	return &value
}

func domainSlice[T any, D hasDomain[T]](domains []D) []T {
	slice := make([]T, len(domains))
	for i := range domains {
		slice[i] = domains[i].toDomain()
	}
	return slice
}

func domainSlicePointer[T any, D hasDomain[T]](domains *[]D) *[]T {
	if domains == nil {
		return nil
	}
	actual := *domains
	slice := make([]T, len(actual))
	for i := range actual {
		slice[i] = actual[i].toDomain()
	}
	return &slice
}

type Index int

func (d Index) Validate(ctx *ValidateContext, v Validator) {
	s := Index(ctx.Kind.Size())
	if d < 0 {
		v.Add("Cannot be less than 0: %d", d)
	}
	if d >= s {
		v.Add("Cannot be greater than %d: %d", s, d)
	}
}

type RelativeIndex int

func (d RelativeIndex) Validate(ctx *ValidateContext, v Validator) {
	s := RelativeIndex(ctx.Kind.Size())
	if d <= -s {
		v.Add("Cannot be less than %d: %d", -s, d)
	}
	if d >= s {
		v.Add("Cannot be greater than %d: %d", s, d)
	}
}

type PuzzleDimension int

func (d PuzzleDimension) Validate(ctx *ValidateContext, v Validator) {
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
		Cells:    domainSlicePointer[su.Position](p.Cells),
		Relative: domainSlicePointer[su.Position](p.Relative),
	}
}

func (c ConstraintSumValue) Validate(ctx *ValidateContext, v Validator) {
	if c.Sum <= 1 {
		v.Add("The sum %d is not a valid constraint sum.", c.Sum)
	}
	if c.Cells == nil && c.Relative == nil {
		v.Add("No cells specified.")
	}
}

type ConstraintSumCell struct {
	Sum         RelativePosition    `json:"sum"`
	SumRelative bool                `json:"sumRelative"`
	Cells       *[]Position         `json:"cells"`
	Relative    *[]RelativePosition `json:"relative"`
}

func (p ConstraintSumCell) toDomain() su.Constraint {
	return &su.ConstraintSum{
		Sum:      su.SumCell(p.Sum.toDomain(), p.SumRelative),
		Cells:    domainSlicePointer[su.Position](p.Cells),
		Relative: domainSlicePointer[su.Position](p.Relative),
	}
}

type ConstraintUnique struct {
	Cells    *[]Position         `json:"cells"`
	Relative *[]RelativePosition `json:"relative"`
	Same     bool                `json:"same"`
}

func (p ConstraintUnique) toDomain() su.Constraint {
	return &su.ConstraintUnique{
		Cells:    domainSlicePointer[su.Position](p.Cells),
		Relative: domainSlicePointer[su.Position](p.Relative),
		Same:     p.Same,
	}
}

type DirectionInt int

func (d DirectionInt) Validate(ctx *ValidateContext, v Validator) {
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
		Cells:     domainSlicePointer[su.Position](p.Cells),
		Relative:  domainSlicePointer[su.Position](p.Relative),
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
		Cells:    domainSlicePointer[su.Position](p.Cells),
		Relative: domainSlicePointer[su.Position](p.Relative),
		Exclude:  domainSlicePointer[su.Position](p.Exclude),
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
		Cells:     domainSlice[su.Position](p.Cells),
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
	d = append(d, domainSlice[su.Constraint](c.SumValues)...)
	d = append(d, domainSlice[su.Constraint](c.SumCells)...)
	d = append(d, domainSlice[su.Constraint](c.Uniques)...)
	d = append(d, domainSlice[su.Constraint](c.Orders)...)
	d = append(d, domainSlice[su.Constraint](c.Magics)...)
	d = append(d, domainSlice[su.Constraint](c.ScalePairs)...)
	d = append(d, domainSlice[su.Constraint](c.Differences)...)
	d = append(d, domainSlice[su.Constraint](c.Divisibles)...)
	return d
}

type GenerateCount int

func (g GenerateCount) Validate(ctx *ValidateContext, v Validator) {
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
	Symmetric     bool            `json:"symmetric"`
	BoxWidth      PuzzleDimension `json:"boxWidth"`
	BoxHeight     PuzzleDimension `json:"boxHeight"`
	Constraints   Constraints     `json:"constraints"`
	Candidates    bool            `json:"candidates"`
	Solutions     bool            `json:"solutions"`
}

func (r GenerateBase) Validate(ctx *ValidateContext, v Validator) {
	ctx.Kind = *su.NewKind(int(r.BoxWidth), int(r.BoxHeight))
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
	clear.Symmetric = r.Symmetric

	return kind, solve, clear
}

type GenerateJson struct {
	GenerateBase
	SolutionLogs   bool `json:"solutionLogs"`
	SolutionStates bool `json:"solutionStates"`
}

type GenerateJsonRequest struct {
	Puzzles []GenerateJson `json:"types"`
}

func handleGenerateJson(w http.ResponseWriter, r *http.Request) {
	out, err := ParseBody[GenerateJsonRequest](r)
	if err != nil {
		fmt.Fprintf(w, "Error %v", err)
		return
	}

	ctx := &ValidateContext{}
	errs := GetValidations(out, ctx)
	if len(errs) > 0 {
		json.NewEncoder(w).Encode(errs)
		return
	}

	type PuzzleGenerated struct {
		Values          [][]int    `json:"values"`
		Solution        *[][]int   `json:"solution"`
		SolutionEncoded string     `json:"solutionEncoded"`
		Logs            *[]string  `json:"logs"`
		Candidates      *[][][]int `json:"candidates"`
	}

	puzzles := []PuzzleGenerated{}

	for _, puzzle := range out.Puzzles {
		kind, _, clear := puzzle.toDomain()
		gen := kind.Generator()
		gen.Solver().LogEnabled = puzzle.SolutionLogs
		gen.Solver().LogState = puzzle.SolutionStates

		for i := 0; i < int(puzzle.Count); i++ {
			gen.Reset()

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

					puzzles = append(puzzles, final)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(puzzles)
}

func Register() {
	// mux := http.NewServeMux()
	http.HandleFunc("/generate", handleGenerateJson)

	// http.Handle("/api", mux)
}

func GetValidations(value any, ctx *ValidateContext) []Validation {
	validator := Validator{
		Validations: &[]Validation{},
	}
	Validate(reflect.ValueOf(value), ctx, validator)
	return *validator.Validations
}

func Validate(value reflect.Value, ctx *ValidateContext, validator Validator) {
	if v, ok := value.Interface().(CanValidate); ok {
		v.Validate(ctx, validator)
	}

	switch value.Kind() {
	case reflect.Pointer:
		if !value.IsNil() {
			Validate(value.Elem(), ctx, validator)
		}
	case reflect.Interface:
		Validate(value.Elem(), ctx, validator)
	case reflect.Struct:
		valueType := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := valueType.Field(i)
			next := field.Name
			nextJson := field.Tag.Get("json")
			if nextJson != "" {
				next = nextJson
			}
			if field.Anonymous {
				Validate(value.Field(i), ctx, validator)
			} else {
				Validate(value.Field(i), ctx, validator.Next(next))
			}
		}
	case reflect.Array:
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			Validate(value.Index(i), ctx, validator.NextInt(i))
		}
	}
}

func ParseBody[T any](r *http.Request) (*T, error) {
	var parsed T

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parsed)

	return &parsed, err
}
