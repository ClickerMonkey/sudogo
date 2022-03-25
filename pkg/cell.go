package sudogo

type Cell struct {
	value      int
	id         int
	row        int
	col        int
	box        int
	candidates Candidates
}

func (cell *Cell) HasValue() bool {
	return cell.value != 0
}

func (cell *Cell) Empty() bool {
	return cell.value == 0
}

func (cell *Cell) InGroup(other *Cell) bool {
	return cell.id != other.id && (cell.row == other.row || cell.col == other.col || cell.box == other.box)
}

func (cell *Cell) InBox(other *Cell) bool {
	return cell.id != other.id && cell.box == other.box
}

func (cell *Cell) InRow(other *Cell) bool {
	return cell.id != other.id && cell.row == other.row
}

func (cell *Cell) InColumn(other *Cell) bool {
	return cell.id != other.id && cell.col == other.col
}

func (cell *Cell) RemoveCandidate(value int) bool {
	return cell.candidates.Set(value, false)
}

func (cell *Cell) HasCandidate(value int) bool {
	return cell.candidates.Has(value)
}

func (cell *Cell) SetValue(value int) bool {
	can := cell.candidates.Has(value)
	if can {
		cell.value = value
		cell.candidates.Clear()
	}
	return can
}

func (cell *Cell) Candidates() []int {
	return cell.candidates.ToSlice()
}

func (cell *Cell) FirstCandidate() int {
	return cell.candidates.First()
}
