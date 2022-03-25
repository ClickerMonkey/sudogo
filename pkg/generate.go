package sudogo

import "math/rand"

func (instance *PuzzleInstance) GetUnsolved() *Cell {
	if len(instance.unsolved) == 0 {
		return nil
	} else {
		return instance.unsolved[0]
	}
}

func (instance *PuzzleInstance) GetRandomUnsolved(random *rand.Rand) *Cell {
	n := len(instance.unsolved)
	if n == 0 {
		return nil
	}
	i := random.Intn(n)
	return instance.unsolved[i]
}

func (instance *PuzzleInstance) GetRandom(random *rand.Rand, match func(other *Cell) bool) *Cell {
	matches := int32(0)
	for i := range instance.cells {
		cell := &instance.cells[i]
		if match(cell) {
			matches++
		}
	}
	if matches == 0 {
		return nil
	}
	chosen := random.Int31n(matches)
	for i := range instance.cells {
		cell := &instance.cells[i]
		if match(cell) {
			chosen--
			if chosen < 0 {
				return cell
			}
		}
	}
	return nil
}

// The smallest number of candidates in a cell without a value. 0=solved, 1=naked
func (instance *PuzzleInstance) GetPressure() int {
	pressure := 0
	for i := range instance.unsolved {
		cell := instance.unsolved[i]
		if pressure == 0 || pressure > cell.candidates.Count {
			pressure = cell.candidates.Count
		}
	}
	return pressure
}

func (instance *PuzzleInstance) GetRandomPressured(random *rand.Rand) *Cell {
	minCount := instance.GetPressure()

	return instance.GetRandom(random, func(other *Cell) bool {
		return other.Empty() && other.candidates.Count == minCount
	})
}

func (instance *PuzzleInstance) Generate(random *rand.Rand) bool {
	for !instance.Solved() {
		instance.Solve()

		if instance.Solved() {
			break
		}

		rnd := instance.GetRandomUnsolved(random)

		if rnd == nil || rnd.candidates.Count == 0 {
			return false
		}

		instance.SetCell(rnd, *randomItem[int](random, rnd.Candidates()))
	}
	return true
}
