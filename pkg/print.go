package sudogo

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func toString(write func(out io.Writer)) string {
	out := bytes.NewBufferString("")
	write(out)
	return out.String()
}

func (puzzle *Puzzle) Print() {
	print(puzzle.ToString())
}

func (puzzle *Puzzle) ToString() string {
	return toString(puzzle.Write)
}

func (puzzle *Puzzle) Write(out io.Writer) {
	boxsWide, boxsHigh, boxWidth, boxHeight, _ := puzzle.kind.GetDimensions()
	digitSize := puzzle.kind.DigitsSize()
	digitFormat := "%" + strconv.Itoa(digitSize) + "d"
	empty := strings.Repeat(" ", digitSize)
	thickH := strings.Repeat("\u2550", digitSize)
	thinH := strings.Repeat("\u2500", digitSize)
	thickV := "\u2551"
	thickX := "\u256C"
	thickTL := "\u2554"
	thickTR := "\u2557"
	thickBL := "\u255A"
	thickBR := "\u255D"
	thickT := "\u2566"
	thickL := "\u2560"
	thickR := "\u2563"
	thickB := "\u2569"
	thickTThin := "\u2564"
	thickLThin := "\u255F"
	thickRThin := "\u2562"
	thickBThin := "\u2567"
	thickVThin := "\u256B"
	thickHThin := "\u256A"
	thinV := "\u2502"
	thinX := "\u253C"
	newline := "\n"

	appendLine := func(left string, smallDivider string, bigDivider string, right string, getHorizontal func(column int) string) {
		io.WriteString(out, left)
		for box := 0; box < boxsWide; box++ {
			endOfBox := box == boxsWide-1

			for cell := 0; cell < boxWidth; cell++ {
				endOfCell := cell == boxWidth-1

				io.WriteString(out, getHorizontal(box*boxsWide+cell))

				if endOfBox && endOfCell {
					io.WriteString(out, right)
					io.WriteString(out, newline)
				} else if endOfCell {
					io.WriteString(out, bigDivider)
				} else {
					io.WriteString(out, smallDivider)
				}
			}
		}
	}

	appendTextLine := func(left string, smallDivider string, bigDivider string, right string, horizontal string) {
		appendLine(left, smallDivider, bigDivider, right, func(column int) string {
			return horizontal
		})
	}

	appendRow := func(row []int) {
		appendLine(thickV, thinV, thickV, thickV, func(column int) string {
			if row[column] == 0 {
				return empty
			} else {
				return fmt.Sprintf(digitFormat, row[column])
			}
		})
	}

	appendTextLine(thickTL, thickTThin, thickT, thickTR, thickH)

	for box := 0; box < boxsHigh; box++ {
		endOfBox := box == boxsHigh-1

		for cell := 0; cell < boxHeight; cell++ {
			endOfCell := cell == boxHeight-1

			appendRow(puzzle.GetRow(box*boxsHigh + cell))

			if endOfBox && endOfCell {
				appendTextLine(thickBL, thickBThin, thickB, thickBR, thickH)
			} else if endOfCell {
				appendTextLine(thickL, thickHThin, thickX, thickR, thickH)
			} else {
				appendTextLine(thickLThin, thinX, thickVThin, thickRThin, thinH)
			}
		}
	}
}

func (instance *Puzzle) PrintCells() {
	print(instance.ToCellsString())
}

func (cell *Cell) ToString() string {
	cellValue := "_"
	if cell.HasValue() {
		cellValue = strconv.Itoa(cell.value)
	}

	return fmt.Sprintf("{%d,%d} = %s %s", cell.col, cell.row, cellValue, fmt.Sprint(cell.Candidates()))
}

func (instance *Puzzle) ToCellsString() string {
	s := ""

	for i := range instance.cells {
		s += instance.cells[i].ToString() + "\n"
	}

	return s
}
func (puzzle *Puzzle) PrintCandidates() {
	print(puzzle.ToCandidatesString())
}

func (puzzle *Puzzle) ToCandidatesString() string {
	return toString(puzzle.WriteCandidates)
}

func (puzzle *Puzzle) WriteCandidates(out io.Writer) {
	boxsWide, boxsHigh, boxWidth, boxHeight, _ := puzzle.kind.GetDimensions()
	digitSize := puzzle.kind.DigitsSize()
	digitFormat := "%" + strconv.Itoa(digitSize) + "d"
	empty := strings.Repeat(" ", digitSize)
	solved := strings.Repeat("\u2591", digitSize)
	thickH := strings.Repeat("\u2550", digitSize*boxWidth)
	thinH := strings.Repeat("\u2500", digitSize*boxWidth)
	thickV := "\u2551"
	thickX := "\u256C"
	thickTL := "\u2554"
	thickTR := "\u2557"
	thickBL := "\u255A"
	thickBR := "\u255D"
	thickT := "\u2566"
	thickL := "\u2560"
	thickR := "\u2563"
	thickB := "\u2569"
	thickTThin := "\u2564"
	thickLThin := "\u255F"
	thickRThin := "\u2562"
	thickBThin := "\u2567"
	thickVThin := "\u256B"
	thickHThin := "\u256A"
	thinV := "\u2502"
	thinX := "\u253C"
	newline := "\n"

	appendLine := func(left string, smallDivider string, bigDivider string, right string, writeHorizontal func(column int)) {
		io.WriteString(out, left)
		for box := 0; box < boxsWide; box++ {
			endOfBox := box == boxsWide-1

			for cell := 0; cell < boxWidth; cell++ {
				endOfCell := cell == boxWidth-1

				writeHorizontal(box*boxsWide + cell)

				if endOfBox && endOfCell {
					io.WriteString(out, right)
					io.WriteString(out, newline)
				} else if endOfCell {
					io.WriteString(out, bigDivider)
				} else {
					io.WriteString(out, smallDivider)
				}
			}
		}
	}

	appendTextLine := func(left string, smallDivider string, bigDivider string, right string, horizontal string) {
		appendLine(left, smallDivider, bigDivider, right, func(column int) {
			io.WriteString(out, horizontal)
		})
	}

	centerX := boxWidth / 2
	centerY := (boxHeight - 1) / 2

	appendRow := func(row []*Cell) {
		for cellRow := 0; cellRow < boxHeight; cellRow++ {

			appendLine(thickV, thinV, thickV, thickV, func(column int) {
				cell := row[column]
				if cell.Empty() {
					for cellCol := 0; cellCol < boxWidth; cellCol++ {
						candidate := cellRow*boxHeight + cellCol + 1
						if cell.candidates.Has(candidate) {
							io.WriteString(out, fmt.Sprintf(digitFormat, candidate))
						} else {
							io.WriteString(out, empty)
						}
					}
				} else {
					for cellCol := 0; cellCol < boxWidth; cellCol++ {
						if cellCol == centerX && cellRow == centerY {
							io.WriteString(out, fmt.Sprintf(digitFormat, cell.value))
						} else {
							io.WriteString(out, solved)
						}
					}
				}
			})
		}
	}

	appendTextLine(thickTL, thickTThin, thickT, thickTR, thickH)

	for box := 0; box < boxsHigh; box++ {
		endOfBox := box == boxsHigh-1

		for cell := 0; cell < boxHeight; cell++ {
			endOfCell := cell == boxHeight-1

			appendRow(puzzle.GetRowCells(box*boxsHigh + cell))

			if endOfBox && endOfCell {
				appendTextLine(thickBL, thickBThin, thickB, thickBR, thickH)
			} else if endOfCell {
				appendTextLine(thickL, thickHThin, thickX, thickR, thickH)
			} else {
				appendTextLine(thickLThin, thinX, thickVThin, thickRThin, thinH)
			}
		}
	}
}
