package sudogo

import (
	mathBits "math/bits"
)

// A set of integers from 0 to 63 with O(1) add, remove, search, and union/substraction/intersection set operations.
type Bitset struct {
	Value uint64
	Count int
}

// Sets the set to the first N integers starting with 0.
func (bits *Bitset) Fill(n int) {
	bits.Count = n
	bits.Value = (1 << n) - 1
}

// Removes all integers from the set
func (bits *Bitset) Clear() {
	bits.Fill(0)
}

// Returns whether the integer exists in the set
func (bits *Bitset) Has(i int) bool {
	return (bits.Value & (1 << i)) != 0
}

// Adds or removes the integer from the set and returns true if the set has changed as a result of this call.
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

// Creates a slice of all integers in this set.
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

// The smallest integer in the set or 64 if the set is empty.
func (bits *Bitset) First() int {
	return mathBits.TrailingZeros64(bits.Value)
}

// Updates the number of integers in the set based on the Value.
func (bits *Bitset) UpdateCount() {
	bits.Count = mathBits.OnesCount64(bits.Value)
}

// Removes all integers in the given set from this set.
func (bits *Bitset) Remove(remove Bitset) int {
	original := bits.Count
	bits.Value = bits.Value & ^remove.Value
	bits.UpdateCount()
	return original - bits.Count
}

// Adds all integers in the given set to this set (union).
func (bits *Bitset) Or(or Bitset) {
	bits.Value = bits.Value | or.Value
	bits.UpdateCount()
}

// Removes any integers from this set that also don't exist in the given set (intersection).
func (bits *Bitset) And(and Bitset) {
	bits.Value = bits.Value & and.Value
	bits.UpdateCount()
}
