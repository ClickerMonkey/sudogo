package sudogo

// A cell holds a value or possible values in a Sudoku puzzle.
type Cell struct {
	// The value in the cell or 0 if a no value exists yet.
	value int
	// The unique id of the cell. This is also the index of this cell in the puzzle's cells slice.
	id int
	// The zero-based row this cell is in.
	row int
	// The zero-based column this cell is in.
	col int
	// The zero-based box this cell is in which starts in the top left of the puzzle and goes right and then restarts on the left side for the next row.
	box int
	// The possible values in the cell if there is no value.
	candidates Candidates
}

// Returns whether this cell has a value in it.
func (cell *Cell) HasValue() bool {
	return cell.value != 0
}

// Returns this cell is empty (does not have a value in it).
func (cell *Cell) Empty() bool {
	return cell.value == 0
}

// Returns if this cell is valid, meaning the value and candidates it has match. If this returns false then there is a logical error in the software.
func (cell *Cell) Valid() bool {
	return (cell.value != 0) == (cell.candidates.Value == 0)
}

// Returns whether this cell and the given cell are in the same group (box, column, or row).
func (cell *Cell) InGroup(other *Cell) bool {
	return cell.id != other.id && (cell.row == other.row || cell.col == other.col || cell.box == other.box)
}

// Returns whether this cell and the given cell are in the same box.
func (cell *Cell) InBox(other *Cell) bool {
	return cell.id != other.id && cell.box == other.box
}

// Returns whether this cell and the given cell are in the same row.
func (cell *Cell) InRow(other *Cell) bool {
	return cell.id != other.id && cell.row == other.row
}

// Returns whether this cell and the given cell are in the same column.
func (cell *Cell) InColumn(other *Cell) bool {
	return cell.id != other.id && cell.col == other.col
}

// Removes the given candidate from this cell and returns whether it existed in the first place.
func (cell *Cell) RemoveCandidate(value int) bool {
	return cell.candidates.Set(value, false)
}

// Determines if this cell has the given candidate.
func (cell *Cell) HasCandidate(value int) bool {
	return cell.candidates.Has(value)
}

// Attempts to set the value of this cell. If this cell already has a value or doesn't have the given value as a candidate
// then false will be returned. Otherwise if the value is applied then true is returned.
func (cell *Cell) SetValue(value int) bool {
	can := cell.candidates.Has(value)
	if can {
		cell.value = value
		cell.candidates.Clear()
	}
	return can
}

// Returns the candidates that exists in this cell as a slice of ints.
func (cell *Cell) Candidates() []int {
	return cell.candidates.ToSlice()
}

// Returns the small candidate available in this cell. If none exist then 64 is returned.
func (cell *Cell) FirstCandidate() int {
	return cell.candidates.First()
}
