package collection

import "sort"

func ReverseSlice[V any](slice []V) []V {
	reversed := make([]V, len(slice))
	copy(reversed, slice)
	ReverseSliceInplace(reversed)
	return reversed
}

func ReverseSliceInplace[V any](slice []V) {
	i := 0
	j := len(slice) - 1
	for i < j {
		slice[i], slice[j] = slice[j], slice[i]
		i++
		j--
	}
}

func SortedSlice[V any](slice []V, less func(a, b V) bool) []V {
	result := make([]V, len(slice))
	copy(result, slice)
	SortSlice(result, less)
	return result
}

func SortSlice[V any](slice []V, less func(a, b V) bool) {
	sort.Slice(slice, func(i, j int) bool { return less(slice[i], slice[j]) })
}

func SliceToMap[K comparable, V any](slice []V, keyFn func(V) K) map[K]V {
	m := map[K]V{}
	for _, v := range slice {
		m[keyFn(v)] = v
	}
	return m
}
