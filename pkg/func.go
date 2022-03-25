package sudogo

import "math/rand"

func removeAtIndex[T any](slice []T, index int) []T {
	last := len(slice) - 1
	if (index >= 0 && index <= last) {
		slice[index] = slice[last]
		slice = slice[:last]
	}
	return slice
}

func bitsOn(value uint64) int {
	on := 0
	for i := 0; i < 64; i++ {
		if (value & (1 << i)) != 0 {
			on++
		}
	}
	return on
}

func randomItem[T any](random *rand.Rand, slice []T) *T {
	i := random.Intn(len(slice))
	return &slice[i]
}