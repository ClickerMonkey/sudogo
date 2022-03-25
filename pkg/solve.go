package sudogo

func (instance *PuzzleInstance) GetNakedSingle() (*Cell, int) {
	for i := range instance.cells {
		cell := &instance.cells[i]
		if cell.candidates.Count == 1 {
			return cell, cell.FirstCandidate()
		}
	}
	return nil, 0
}

func (instance *PuzzleInstance) SolveNakedSingles(max int) int {
	return instance.DoSets(max, func() (*Cell, int) {
		return instance.GetNakedSingle()
	})
}

func (instance *PuzzleInstance) GetHiddenSingle() (*Cell, int) {
	// A cell which has a candidate that is unique to the row, cell, or box
	for i := range instance.cells {
		cell := &instance.cells[i]
		if cell.Empty() {
			box := instance.GetHiddenSingleBox(cell)
			if box != 0 {
				return cell, box
			}
			row := instance.GetHiddenSingleRow(cell)
			if row != 0 {
				return cell, row
			}
			col := instance.GetHiddenSingleColumn(cell)
			if col != 0 {
				return cell, col
			}
		}
	}
	return nil, 0
}

func (instance *PuzzleInstance) GetHiddenSingleBox(cell *Cell) int {
	on := cell.candidates
	instance.ForInBox(cell, func(other *Cell) bool {
		if other.Empty() {
			on.Remove(other.candidates.Value)
			if on.Count == 0 {
				return false
			}
		}
		return true
	})
	if on.Count == 1 {
		return on.First()
	}
	return 0
}

func (instance *PuzzleInstance) GetHiddenSingleRow(cell *Cell) int {
	on := cell.candidates
	instance.ForInRow(cell, func(other *Cell) bool {
		if other.Empty() {
			on.Remove(other.candidates.Value)
			if on.Count == 0 {
				return false
			}
		}
		return true
	})
	if on.Count == 1 {
		return on.First()
	}
	return 0
}

func (instance *PuzzleInstance) GetHiddenSingleColumn(cell *Cell) int {
	on := cell.candidates
	instance.ForInColumn(cell, func(other *Cell) bool {
		if other.Empty() {
			on.Remove(other.candidates.Value)
			if on.Count == 0 {
				return false
			}
		}
		return true
	})
	if on.Count == 1 {
		return on.First()
	}
	return 0
}

func (instance *PuzzleInstance) SolveHiddenSingles(max int) int {
	return instance.DoSets(max, func() (*Cell, int) {
		return instance.GetHiddenSingle()
	})
}

func (instance *PuzzleInstance) RemovePointingCandidates(max int) int {
	// If in a box all candidates of a certain digit are confined to a row or column, that digit cannot appear outside of that box in that row or column.
	removed := 0

	for i := range instance.cells {
		cell := &instance.cells[i]
		if cell.HasValue() {
			continue
		}

		// all candidates in this box's row that are shared
		row := cell.candidates.Clone()
		// all candidates in this box's column that are shared
		col := cell.candidates.Clone()

		// remove candidates that are not shared
		instance.ForInBox(cell, func(other *Cell) bool {
			if other.Empty() {
				if other.row == cell.row {
					row.And(other.candidates.Value)
				}
				if other.col == cell.col {
					col.And(other.candidates.Value)
				}
			}

			return row.Count > 0 || col.Count > 0
		})

		// remove candidates that exist outside the row or column
		instance.ForInBox(cell, func(other *Cell) bool {
			if other.Empty() {
				if other.row != cell.row && other.col != cell.col {
					row.Remove(other.candidates.Value)
					col.Remove(other.candidates.Value)
				}
			}

			return row.Count > 0 || col.Count > 0
		})

		// what is remaining are candidates confined to the cells row in the box
		if row.Count > 0 {
			instance.ForInRow(cell, func(other *Cell) bool {
				if other.Empty() && other.box != cell.box {
					removed += other.candidates.Remove(row.Value)
				}
				return true
			})
		}

		if max > 0 && removed >= max {
			break
		}

		// what is remaining are candidates confined to the cells column in the box
		if col.Count > 0 {
			instance.ForInColumn(cell, func(other *Cell) bool {
				if other.Empty() && other.box != cell.box {
					removed += other.candidates.Remove(col.Value)
				}
				return true
			})
		}

		if max > 0 && removed >= max {
			break
		}
	}

	return removed
}

func (instance *PuzzleInstance) RemoveClaimingCandidates(max int) int {
	// If in a row or column a candidate only appears in a single box then that candidate can be removed from other cells in that box
	removed := 0

	for i := range instance.cells {
		cell := &instance.cells[i]
		if cell.HasValue() {
			continue
		}

		// all candidates in this row that are not shared outside of the box
		row := cell.candidates.Clone()
		// all candidates in this column that are not shared outside of the box
		col := cell.candidates.Clone()

		// remove candidates from row that exist in the cells row outside the box
		instance.ForInRow(cell, func(other *Cell) bool {
			if other.Empty() && other.box != cell.box {
				row.Remove(other.candidates.Value)
			}

			return row.Count > 0
		})

		// remove candidates from column that exist in the cells column outside the box
		instance.ForInColumn(cell, func(other *Cell) bool {
			if other.Empty() && other.box != cell.box {
				col.Remove(other.candidates.Value)
			}

			return col.Count > 0
		})

		// what is remaining are the candidates unique to the row/column outside this box
		if row.Count > 0 || col.Count > 0 {
			instance.ForInBox(cell, func(other *Cell) bool {
				if other.Empty() {
					if row.Count > 0 && other.row != cell.row {
						removed += other.candidates.Remove(row.Value)
					}
					if col.Count > 0 && other.col != cell.col {
						removed += other.candidates.Remove(col.Value)
					}
				}
				return true
			})
		}

		if max > 0 && removed >= max {
			break
		}
	}

	return removed
}

func (instance *PuzzleInstance) RemoveHiddenSubsetCandidates(max int) int {
	removed := 0
	// subsets := [2]int{2, 3}

	// for _, subsetSize := range subsets {
	// 	for i := range instance.cells {
	// 		cell := &instance.cells[i]
	// 		if cell.HasValue() {
	// 			continue
	// 		}

	// 	}
	// }

	return removed
}

func (instance *PuzzleInstance) RemoveNakedSubsetCandidates(max int) int {
	removed := 0
	subsets := [2]int{2, 3}

	for _, subsetSize := range subsets {
		for i := range instance.cells {
			cell := &instance.cells[i]
			if cell.HasValue() || cell.candidates.Count != subsetSize {
				continue
			}

			candidates := cell.candidates.Value
			rowMatches := 1

			instance.ForInRow(cell, func(other *Cell) bool {
				if other.Empty() && other.candidates.Value == candidates {
					rowMatches++
				}
				return rowMatches <= subsetSize
			})

			if rowMatches == subsetSize {
				instance.ForInRow(cell, func(other *Cell) bool {
					if other.Empty() && other.candidates.Value != candidates {
						removed += other.candidates.Remove(candidates)
					}
					return true
				})

				if max > 0 && removed >= max {
					break
				}
			}

			colMatches := 1

			instance.ForInColumn(cell, func(other *Cell) bool {
				if other.Empty() && other.candidates.Value == candidates {
					colMatches++
				}
				return colMatches <= subsetSize
			})

			if colMatches == subsetSize {
				instance.ForInColumn(cell, func(other *Cell) bool {
					if other.Empty() && other.candidates.Value != candidates {
						removed += other.candidates.Remove(candidates)
					}
					return true
				})

				if max > 0 && removed >= max {
					break
				}
			}

			boxMatches := 1

			instance.ForInBox(cell, func(other *Cell) bool {
				if other.Empty() && other.candidates.Value == candidates {
					boxMatches++
				}
				return boxMatches <= subsetSize
			})

			if boxMatches == subsetSize {
				instance.ForInBox(cell, func(other *Cell) bool {
					if other.Empty() && other.candidates.Value != candidates {
						removed += other.candidates.Remove(candidates)
					}
					return true
				})

				if max > 0 && removed >= max {
					break
				}
			}
		}
	}

	return removed
}

func (instance *PuzzleInstance) DoSets(max int, nextSet func() (*Cell, int)) int {
	set := 0
	cell, cellValue := nextSet()
	for cell != nil {
		set++
		instance.SetCell(cell, cellValue)
		if max > 0 && set == max {
			break
		}
		cell, cellValue = nextSet()
	}
	return set
}

func (instance *PuzzleInstance) Solved() bool {
	for i := range instance.cells {
		if !instance.cells[i].HasValue() {
			return false
		}
	}
	return true
}

func (instance *PuzzleInstance) Solve() int {
	placed := 0
	for {
		placed += instance.SolveNakedSingles(-1)

		if instance.SolveHiddenSingles(1) > 0 {
			placed++
			continue
		}

		if instance.RemovePointingCandidates(1) > 0 {
			continue
		}

		if instance.RemoveClaimingCandidates(1) > 0 {
			continue
		}

		if instance.RemoveHiddenSubsetCandidates(1) > 0 {
			continue
		}

		if instance.RemoveNakedSubsetCandidates(1) > 0 {
			continue
		}

		break
	}

	return placed
}
