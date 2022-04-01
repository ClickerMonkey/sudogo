package sudogo

// A cell holds a value or possible values in a Sudoku puzzle.
type Cell struct {
	// The Value in the cell or 0 if a no Value exists yet.
	Value int
	// The unique Id of the cell. This is also the index of this cell in the puzzle's cells slice.
	Id int
	// The zero-based Row this cell is in.
	Row int
	// The zero-based column this cell is in.
	Col int
	// The zero-based Box this cell is in which starts in the top left of the puzzle and goes right and then restarts on the left side for the next row.
	Box int
	// The possible values in the cell if there is no value.
	candidates Candidates
}

// Returns whether this cell has a value in it.
func (cell *Cell) HasValue() bool {
	return cell.Value != 0
}

// Returns this cell is empty (does not have a value in it).
func (cell *Cell) Empty() bool {
	return cell.Value == 0
}

// Returns if this cell is valid, meaning the value and candidates it has match. If this returns false then there is a logical error in the software.
func (cell *Cell) Valid() bool {
	return (cell.Value != 0) == (cell.candidates.Value == 0)
}

// Returns whether this cell and the given cell are in the same group (box, column, or row).
func (cell *Cell) InGroup(other *Cell) bool {
	return cell.Id != other.Id && (cell.Row == other.Row || cell.Col == other.Col || cell.Box == other.Box)
}

// Returns whether this cell and the given cell are in the same box.
func (cell *Cell) InBox(other *Cell) bool {
	return cell.Id != other.Id && cell.Box == other.Box
}

// Returns whether this cell and the given cell are in the same row.
func (cell *Cell) InRow(other *Cell) bool {
	return cell.Id != other.Id && cell.Row == other.Row
}

// Returns whether this cell and the given cell are in the same column.
func (cell *Cell) InColumn(other *Cell) bool {
	return cell.Id != other.Id && cell.Col == other.Col
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
		cell.Value = value
		cell.candidates.Clear()
	}
	return can
}

// Returns the candidates that exists in this cell as a slice of ints.
func (cell *Cell) Candidates() []int {
	return cell.candidates.ToSlice()
}

// Returns the smallest candidate available in this cell. If none exist then 64 is returned.
func (cell *Cell) FirstCandidate() int {
	return cell.candidates.First()
}

// Returns the largest candidate available in this cell. If none exist then 0 is returned.
func (cell *Cell) LastCandidate() int {
	return cell.candidates.Last()
}

// The minimum value this cell can possibly be, taking into account it might
// already have a value and therefore no candidates.
func (cell *Cell) MinValue() int {
	if cell.HasValue() {
		return cell.Value
	} else {
		return cell.FirstCandidate()
	}
}

// The maximum value this cell can possibly be, taking into account it might
// already have a value and therefore no candidates.
func (cell *Cell) MaxValue() int {
	if cell.HasValue() {
		return cell.Value
	} else {
		return cell.LastCandidate()
	}
}
