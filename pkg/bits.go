package sudogo

type Bits struct {
	Value uint64
	Count int
}

func (bits *Bits) Fill(n int) {
	bits.Count = n
	bits.Value = (1 << n) - 1
}

func (bits *Bits) Clear() {
	bits.Fill(0)
}

func (bits *Bits) Clone() Bits {
	return Bits{bits.Value, bits.Count}
}

func (bits *Bits) On(i int) bool {
	return (bits.Value & (1 << i)) != 0
}

func (bits *Bits) Set(i int, on bool) bool {
	change := bits.On(i) != on
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

func (bits *Bits) ToSlice() []int {
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

func (bits *Bits) First() int {
	for i := 0; i < 64; i++ {
		if (bits.Value & (1 << i)) != 0 {
			return i
		}
	}
	return -1
}

func (bits *Bits) UpdateCount() {
	bits.Count = bitsOn(bits.Value)
}

func (bits *Bits) Remove(remove uint64) int {
	original := bits.Count
	bits.Value = bits.Value & ^remove
	bits.UpdateCount()
	return original - bits.Count
}

func (bits *Bits) Or(or uint64) {
	bits.Value = bits.Value | or
	bits.UpdateCount()
}

func (bits *Bits) And(and uint64) {
	bits.Value = bits.Value & and
	bits.UpdateCount()
}
