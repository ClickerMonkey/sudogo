package sudogo

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

func removeAtIndex[T any](slice []T, index int) []T {
	last := len(slice) - 1
	if index >= 0 && index <= last {
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

func pointersWhere[T any](source []T, where func(item *T) bool) []*T {
	pointers := make([]*T, 0, len(source))
	for i := range source {
		item := &source[i]
		if where(item) {
			pointers = append(pointers, item)
		}
	}
	return pointers
}

func sliceIndex[T any](source []T, where func(item T) bool) int {
	for i, item := range source {
		if where(item) {
			return i
		}
	}
	return -1
}

func sliceClone[T any](source []T) []T {
	cloned := make([]T, len(source))
	copy(cloned, source)
	return cloned
}

func sliceLast[T any](source []T) *T {
	last := len(source) - 1
	if last == -1 {
		return nil
	}
	return &source[last]
}

func sliceRemoveLast[T any](source []T) []T {
	last := len(source) - 1
	if last == -1 {
		return source
	}
	return source[:last]
}

func Max[T constraints.Ordered](x T, y T) T {
	if x > y {
		return x
	}
	return y
}

func Min[T constraints.Ordered](x T, y T) T {
	if x < y {
		return x
	}
	return y
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func stringChunk(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)
	n := len(runes)

	if n == 0 {
		return []string{s}
	}

	for i := 0; i < n; i += chunkSize {
		nn := i + chunkSize
		if nn > n {
			nn = n
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}
