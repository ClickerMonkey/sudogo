# sudogo
Advanced Sudoku solving &amp; generating with Go

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
  puzzle.Print()
  // puzzle.ToString()

  // Print to stdout showing candidates in each cell
  puzzle.PrintCandidates()
  // puzzle.ToCandidatesString()

  // Solve
  solver := puzzle.Solver()
  solution, solved := solver.Solve()
  if solved {
    println("Solved!")
  }

  // Generate
  gen := su.Classic.Generator()
  new, attempts := gen.Generate()

  // Puzzles of different sizes
  s4x4 := su.Kind2x2.Generator().Generate()
  s6x6 := su.Kind3x2.Generator().Generate()
  s12x12 := su.Kind4x3.Generator().Generate()
  s16x16 := su.Kind4x4.Generator().Generate()
}
```

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
- [ ] Logging solve operations (type, score, text, cell before/after)
- [ ] Sum Constraint (cells add up to a number)
  - Cells in path add up to another cell
  - Path adds up to constant
  - Knights move cells add up to number
- [ ] Value Constraints (cells must contain certain values in any order)
  - Diagonals contain 1-9
- [ ] Match Constraints (cells in relative positions can't have same value)
  - Knights move cells can't be same number
  - Kings move cells can't be same number
- [ ] Solve step: Skyscraper (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [ ] Solve step: 2-String Kite/Dual 2-String Kite (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [ ] Solve step: Turbot Fish (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [ ] Solve step: Empty Rectangle (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [ ] Solve step: Dual Empty Rectangle (http://hodoku.sourceforge.net/en/tech_sdp.php)
- [ ] Solve step: X-Wing (http://hodoku.sourceforge.net/en/tech_fishb.php)
- [ ] Solve step: Swordfish (http://hodoku.sourceforge.net/en/tech_fishb.php)
- [ ] Solve step: Jellyfish (http://hodoku.sourceforge.net/en/tech_fishb.php)
- [ ] Solve step: Finned/Sashimi X-Wing/Swordfish/Jellyfish (http://hodoku.sourceforge.net/en/tech_fishfs.php)
- [ ] Solve step: Franken Fish (http://hodoku.sourceforge.net/en/tech_fishc.php)
- [ ] Solve step: Mutant Fish (http://hodoku.sourceforge.net/en/tech_fishc.php)
- [ ] Solve step: Siamese Fish (http://hodoku.sourceforge.net/en/tech_fishc.php)
- [ ] Solve step: Constraints