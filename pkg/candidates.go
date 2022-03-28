package sudogo

// A set of possible values for a cell. This is similar to a Bitset but is 1 based.
type Candidates struct {
	Bitset
}

// Returns whether the candidate exists in the set.
func (cand *Candidates) Has(i int) bool {
	return cand.Bitset.Has(i - 1)
}

// Adds or removes the candidate from the set and returns true if the set has changed as a result of this call.
func (cand *Candidates) Set(i int, on bool) bool {
	return cand.Bitset.Set(i-1, on)
}

// Creates a slice of all candidates in this set.
func (cand *Candidates) ToSlice() []int {
	slice := cand.Bitset.ToSlice()
	n := len(slice)
	for i := 0; i < n; i++ {
		slice[i]++
	}
	return slice
}

// The smallest candidate in the set or 64 if the set is empty
func (cand *Candidates) First() int {
	return cand.Bitset.First() + 1
}

// The largest candidate in the set or 0 if the set is empty
func (cand *Candidates) Last() int {
	return cand.Bitset.Last() + 1
}

// Removes all candidates in the given set from this set.
func (cand *Candidates) Remove(remove Candidates) int {
	return cand.Bitset.Remove(remove.Bitset)
}

// Adds all candidates in the given set to this set (union).
func (cand *Candidates) Or(or Candidates) int {
	return cand.Bitset.Or(or.Bitset)
}

// Removes any candidates from this set that also don't exist in the given set (intersection).
func (cand *Candidates) And(and Candidates) int {
	return cand.Bitset.And(and.Bitset)
}
