package sudogo

import "math/bits"

type Candidates struct {
	Value uint64
	Count int
}

func (cand *Candidates) Fill(candidates int) {
	cand.Count = candidates
	cand.Value = (1 << (candidates + 1)) - 2
}

func (cand *Candidates) Clear() {
	cand.Fill(0)
}

func (cand *Candidates) Clone() Candidates {
	return Candidates{cand.Value, cand.Count}
}

func (cand *Candidates) Has(i int) bool {
	return (cand.Value & (1 << i)) != 0
}

func (cand *Candidates) Set(i int, on bool) bool {
	change := cand.Has(i) != on
	if change {
		if on {
			cand.Count++
			cand.Value = cand.Value | (1 << i)
		} else {
			cand.Count--
			cand.Value = cand.Value & ^(1 << i)
		}
	}
	return change
}

func (cand *Candidates) ToSlice() []int {
	slice := make([]int, 0, cand.Count)
	remaining := cand.Value
	current := 0
	for remaining > 0 {
		if (remaining & 1) == 1 {
			slice = append(slice, current)
		}
		remaining = remaining >> 1
		current++
	}
	return slice
}

func (cand *Candidates) First() int {
	return bits.TrailingZeros64(cand.Value) // 70043210 1
}

func (cand *Candidates) UpdateCount() {
	cand.Count = bits.OnesCount64(cand.Value)
}

func (cand *Candidates) Remove(remove Candidates) int {
	original := cand.Count
	cand.Value = cand.Value & ^remove.Value
	cand.UpdateCount()
	return original - cand.Count
}

func (cand *Candidates) Or(or Candidates) {
	cand.Value = cand.Value | or.Value
	cand.UpdateCount()
}

func (cand *Candidates) And(and Candidates) {
	cand.Value = cand.Value & and.Value
	cand.UpdateCount()
}
