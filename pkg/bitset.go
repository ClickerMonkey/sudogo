package sudogo

import (
	mathBits "math/bits"
)

type Bitset struct {
	Value uint64
	Count int
}

func (bits *Bitset) Fill(n int) {
	bits.Count = n
	bits.Value = (1 << n) - 1
}

func (bits *Bitset) Clear() {
	bits.Fill(0)
}

func (bits *Bitset) Has(i int) bool {
	return (bits.Value & (1 << i)) != 0
}

func (bits *Bitset) Set(i int, on bool) bool {
	change := bits.Has(i) != on
	if change {
		if on {
			bits.Count++
			bits.Value = bits.Value | (1 << i)
		} else {
			bits.Count--
			bits.Value = bits.Value & ^(1 << i)
		}
	}
	return change
}

func (bits *Bitset) ToSlice() []int {
	slice := make([]int, 0, bits.Count)
	remaining := bits.Value
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

func (bits *Bitset) First() int {
	return mathBits.TrailingZeros64(bits.Value)
}

func (bits *Bitset) UpdateCount() {
	bits.Count = mathBits.OnesCount64(bits.Value)
}

func (bits *Bitset) Remove(remove Bitset) int {
	original := bits.Count
	bits.Value = bits.Value & ^remove.Value
	bits.UpdateCount()
	return original - bits.Count
}

func (bits *Bitset) Or(or Bitset) {
	bits.Value = bits.Value | or.Value
	bits.UpdateCount()
}

func (bits *Bitset) And(and Bitset) {
	bits.Value = bits.Value & and.Value
	bits.UpdateCount()
}
