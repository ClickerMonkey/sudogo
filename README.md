# sudogo
Advanced Sudoku solving &amp; generating with Go

### Features
- Handles sudoku puzzles up to 64x64.
- Generates any number of puzzles with configurable difficulty to the console or PDF with solutions, candidates, and solution steps optionally included.
- Lists the steps it took to solve a puzzle and the techniques used.
- Finds all solutions for invalid puzzles.
- Is extendable and very configurable.

## Importing

```go
import (
  su "github.com/ClickerMonkey/sudogo"
)

func main() {
  // Create from an existing one
  puzzle := su.Classic.Create([][]int{
    {3, 1, 8, 0, 0, 5, 4, 0, 6},
    {0, 0, 0, 6, 0, 3, 8, 1, 0},
    {0, 0, 6, 0, 8, 0, 5, 0, 3},
    {8, 6, 4, 9, 5, 2, 1, 3, 7},
    {1, 2, 3, 4, 7, 6, 9, 5, 8},
    {7, 9, 5, 3, 1, 8, 2, 6, 4},
    {0, 3, 0, 5, 0, 0, 7, 8, 0},
    {0, 0, 0, 0, 0, 7, 3, 0, 5},
    {0, 0, 0, 0, 3, 9, 6, 4, 1},
  })

  // Print to stdout
  puzzle.PrintConsole()
  // puzzle.ToString()

  // Print to stdout showing candidates in each cell
  puzzle.PrintConsoleCandidates()
  // puzzle.ToCandidatesString()

  // Solve
  solver := puzzle.Solver()
  solution, solved := solver.Solve(su.SolverLimit{})
  if solved {
    println("Solved!")
  }

  // Get all solutions
  solutions := puzzle.GetSolutions(su.SolutionLimit{})

  // Generate
  gen := su.Classic.Generator()
  new, attempts := gen.Generate()

  // Clear out a generated puzzle so we can solve it
  cleared, _ := gen.ClearCells(new, su.DifficultyMedium)
  cleared.Print() // solve this!

  // Puzzles of different sizes
  s4x4 := su.Kind2x2.Generator().Generate()
  s6x6 := su.Kind3x2.Generator().Generate()
  s12x12 := su.Kind4x3.Generator().Generate()
  s16x16 := su.Kind4x4.Generator().Generate()
}
```

## API

Notable functions for the main types

### Puzzle
- Solver() Solver
- Get(col, row) \*Cell
- Set(col, row, value) bool
- SetAll(values) int
- IsSolved() bool
- IsValid() bool
- UniqueId() string
- HasUniqueSolution() bool
- GetSolutions(limits) []\*Solver
- Print() / ToString() / Write(out)
- PrintCandidates() / ToCandidatesString() / WriteCandidates(out)

### Solver
- Set(col, row, value) bool
- Solved() bool
- Solve(limits) (\*Puzzle, bool)

### Generator
- Attempt() \*Puzzle
- Attempts(tries) (\*Puzzle,int)
- Generate() (\*Puzzle,int)
- ClearCells(puzzle,limits) (\*Puzzle,int)
 
### Kind
- Create(values) Puzzle
- Generator()

## Output

Puzzle.Print()
```
╔═╤═╤═╦═╤═╤═╦═╤═╤═╗
║5│3│ ║ │7│ ║ │ │ ║
╟─┼─┼─╫─┼─┼─╫─┼─┼─╢
║6│ │ ║1│9│5║ │ │ ║
╟─┼─┼─╫─┼─┼─╫─┼─┼─╢
║ │9│8║ │ │ ║ │6│ ║
╠═╪═╪═╬═╪═╪═╬═╪═╪═╣
║8│ │ ║ │6│ ║ │ │3║
╟─┼─┼─╫─┼─┼─╫─┼─┼─╢
║4│ │ ║8│ │3║ │ │1║
╟─┼─┼─╫─┼─┼─╫─┼─┼─╢
║7│ │ ║ │2│ ║ │ │6║
╠═╪═╪═╬═╪═╪═╬═╪═╪═╣
║ │6│ ║ │ │ ║2│8│ ║
╟─┼─┼─╫─┼─┼─╫─┼─┼─╢
║ │ │ ║4│1│9║ │ │5║
╟─┼─┼─╫─┼─┼─╫─┼─┼─╢
║ │ │ ║ │8│ ║ │7│9║
╚═╧═╧═╩═╧═╧═╩═╧═╧═╝
```

Puzzle.PrintCandidates()
```
╔═══╤═══╤═══╦═══╤═══╤═══╦═══╤═══╤═══╗
║  3│  3│  3║░░░│  3│░░░║   │   │   ║
║4  │4  │   ║░1░│456│░2░║4 6│456│45 ║
║789│7  │789║░░░│  9│░░░║ 8 │78 │7  ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║12 │░░░│12 ║   │   │░░░║░░░│ 2 │12 ║
║4  │░6░│   ║ 5 │45 │░8░║░3░│45 │45 ║
║7 9│░░░│7 9║7 9│  9│░░░║░░░│7  │7  ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║░░░│1 3│123║  3│  3│  3║1  │ 2 │░░░║
║░5░│4  │   ║  6│4 6│4  ║4 6│4 6│░9░║
║░░░│7  │78 ║7  │   │   ║ 8 │78 │░░░║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║123│1 3│123║░░░│123│░░░║   │ 23│░░░║
║   │ 5 │ 5 ║░4░│  6│░7░║  6│ 56│░8░║
║  9│   │  9║░░░│  9│░░░║  9│  9│░░░║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║░░░│░░░│123║ 23│123│1 3║   │ 23│ 23║
║░6░│░8░│ 5 ║   │   │   ║4  │45 │45 ║
║░░░│░░░│7 9║  9│  9│  9║  9│7 9│7  ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║ 23│  3│░░░║ 23│ 23│░░░║   │░░░│ 23║
║   │   │░4░║  6│  6│░5░║  6│░1░│   ║
║7 9│7  │░░░║ 89│ 89│░░░║  9│░░░│7  ║
╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣
║1 3│░░░│1 3║  3│1 3│1 3║░░░│  3│1 3║
║4  │░2░│  6║   │4  │4  ║░5░│4  │4  ║
║78 │░░░│78 ║ 89│ 89│  9║░░░│ 89│   ║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║1 3│1 3│1 3║  3│░░░│1 3║░░░│  3│░░░║
║4  │45 │ 5 ║ 5 │░7░│4  ║░2░│4  │░6░║
║ 8 │   │ 8 ║ 89│░░░│  9║░░░│ 89│░░░║
╟───┼───┼───╫───┼───┼───╫───┼───┼───╢
║1 3│░░░│1 3║ 23│123│░░░║░░░│  3│1 3║
║4  │░9░│ 5 ║ 5 │45 │░6░║░7░│4  │4  ║
║ 8 │░░░│ 8 ║ 8 │ 8 │░░░║░░░│ 8 │   ║
╚═══╧═══╧═══╩═══╧═══╧═══╩═══╧═══╧═══╝
```

### TODO

- [X] More efficient cell clearing after generation
- [x] Logging solve operations (type, score, text, cell before/after)
- [x] Sum Constraint (cells add up to a number)
  - Cells in path add up to another cell
  - Path adds up to constant
  - Knights move cells add up to number
- [x] Value Constraints (cells must contain certain values in any order)
  - Diagonals contain 1-9
- [x] Match Constraints (cells in relative positions can't have same value)
  - Knights move cells can't be same number
  - Kings move cells can't be same number
- [x] Solve step: Skyscraper (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [x] Solve step: 2-String Kite/Dual 2-String Kite (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [x] Solve step: Empty Rectangle (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [x] Solve step: Dual Empty Rectangle (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [ ] Solve step: X-Wing (http://hodoku.sourceforge.net/en/tech_fishb.php)
- [ ] Solve step: Swordfish (http://hodoku.sourceforge.net/en/tech_fishb.php)
- [ ] Solve step: Jellyfish (http://hodoku.sourceforge.net/en/tech_fishb.php)
- [ ] Solve step: Finned/Sashimi X-Wing/Swordfish/Jellyfish (http://hodoku.sourceforge.net/en/tech_fishfs.php)
- [ ] Solve step: Franken Fish (http://hodoku.sourceforge.net/en/tech_fishc.php)
- [ ] Solve step: Mutant Fish (http://hodoku.sourceforge.net/en/tech_fishc.php)
- [ ] Solve step: Siamese Fish (http://hodoku.sourceforge.net/en/tech_fishc.php)
- [x] Solve step: Constraints

### Resources
- http://hodoku.sourceforge.net/en/techniques.php
- https://www.sudokuoftheday.com/difficulty
