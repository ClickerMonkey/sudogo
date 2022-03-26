package sudogo

type Candidates struct {
	Bitset
}

func (cand *Candidates) Has(i int) bool {
	return cand.Bitset.Has(i - 1)
}

func (cand *Candidates) Set(i int, on bool) bool {
	return cand.Bitset.Set(i-1, on)
}

func (cand *Candidates) ToSlice() []int {
	slice := cand.Bitset.ToSlice()
	n := len(slice)
	for i := 0; i < n; i++ {
		slice[i]++
	}
	return slice
}

func (cand *Candidates) First() int {
	return cand.Bitset.First() + 1
}

func (cand *Candidates) Remove(remove Candidates) int {
	return cand.Bitset.Remove(remove.Bitset)
}

func (cand *Candidates) Or(or Candidates) {
	cand.Bitset.Or(or.Bitset)
}

func (cand *Candidates) And(and Candidates) {
	cand.Bitset.And(and.Bitset)
}
