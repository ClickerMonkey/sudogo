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

func removeValue[T comparable](slice []T, value T) []T {
	for i, v := range slice {
		if v == value {
			return removeAtIndex(slice, i)
		}
	}
	return slice
}

func randomPointer[T any](random *rand.Rand, slice []*T) *T {
	n := len(slice)
	if n == 0 {
		return nil
	}
	i := random.Intn(n)
	return slice[i]
}

func randomElement[T any](random *rand.Rand, slice []T, notFound T) T {
	n := len(slice)
	if n == 0 {
		return notFound
	}
	i := random.Intn(n)
	return slice[i]
}

func pointerAt[T any](slice []*T, index int) *T {
	if index < 0 || index >= len(slice) {
		return nil
	}
	return slice[index]
}