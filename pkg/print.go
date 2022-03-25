package sudogo

import (
	"fmt"
	"strconv"
	"strings"
)

func (instance *Puzzle) Print() {
	print(instance.ToString())
}

func (instance *Puzzle) ToString() string {
	cellSpace := instance.kind.DigitsSize()
	cellEmpty := strings.Repeat(" ", cellSpace)
	cellTop := strings.Repeat("-", cellSpace)
	cellsWide := instance.kind.Width()
	cellsHigh := instance.kind.Height()
	boxWidth := instance.kind.BoxSize.Width
	boxHeight := instance.kind.BoxSize.Height
	s := ""

	appendRow := func() {
		for x := 0; x < cellsWide; x++ {
			if x%boxWidth == 0 {
				s += "+"
			}
			s += cellTop
		}
		s += "+\n"
	}

	for y := 0; y < cellsHigh; y++ {
		if y%boxHeight == 0 {
			appendRow()
		}

		for x := 0; x < cellsWide; x++ {
			cell := instance.Get(x, y)
			if x%boxWidth == 0 {
				s += "|"
			}
			if cell.value != 0 {
				s += strconv.Itoa(cell.value)
			} else {
				s += cellEmpty
			}
		}
		s += "|\n"
	}

	appendRow()

	return s
}

func (instance *Puzzle) PrintCells() {
	print(instance.ToCellsString())
}

func (cell *Cell) ToString() string {
	cellValue := "_"
	if cell.value != 0 {
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
