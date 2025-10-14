package collection

import (
	"sort"
)

func CopySlice[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}

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

func SliceToMap[K comparable, V any](slice []V, mapFn func(V) (K, V)) map[K]V {
	m := map[K]V{}
	for _, v := range slice {
		k, v := mapFn(v)
		m[k] = v
	}
	return m
}

func SlicesToSet[V comparable](slices ...[]V) map[V]struct{} {
	set := map[V]struct{}{}
	for _, slice := range slices {
		for _, v := range slice {
			set[v] = struct{}{}
		}
	}
	return set
}

func JoinSlices[T any](slices ...[]T) []T {
	var totalLen int
	for _, slice := range slices {
		totalLen += len(slice)
	}
	merged := make([]T, totalLen)
	idx := 0
	for _, slice := range slices {
		copy(merged[idx:], slice)
		idx += len(slice)
	}
	return merged
}

func MapSlice[T1, T2 any](slice []T1, mapFn func(idx int, v T1) T2) []T2 {
	result := make([]T2, len(slice))
	for idx := range slice {
		result[idx] = mapFn(idx, slice[idx])
	}
	return result
}

func MapSliceI[T1, T2 any](slice []T1, mapFn func(idx int) T2) []T2 {
	result := make([]T2, len(slice))
	for idx := range slice {
		result[idx] = mapFn(idx)
	}
	return result
}

func MapSliceV[T1, T2 any](slice []T1, mapFn func(v T1) T2) []T2 {
	result := make([]T2, len(slice))
	for idx := range slice {
		result[idx] = mapFn(slice[idx])
	}
	return result
}

func MapSliceInplace[T1 any](slice []T1, mapFn func(idx int, v T1) T1) {
	for idx := range slice {
		slice[idx] = mapFn(idx, slice[idx])
	}
}

func MapSliceInplaceI[T1 any](slice []T1, mapFn func(idx int) T1) {
	for idx := range slice {
		slice[idx] = mapFn(idx)
	}
}

func MapSliceInplaceV[T1 any](slice []T1, mapFn func(v T1) T1) {
	for idx := range slice {
		slice[idx] = mapFn(slice[idx])
	}
}

func ReduceSlice[T1, T2 any](slice []T1, reduceFn func(accum T2, idx int, value T1) T2) T2 {
	var accum T2
	for idx := range slice {
		accum = reduceFn(accum, idx, slice[idx])
	}
	return accum
}

func ReduceSliceI[T1, T2 any](slice []T1, reduceFn func(accum T2, idx int) T2) T2 {
	var accum T2
	for idx := range slice {
		accum = reduceFn(accum, idx)
	}
	return accum
}

func ReduceSliceV[T1, T2 any](slice []T1, reduceFn func(accum T2, v T1) T2) T2 {
	var accum T2
	for idx := range slice {
		accum = reduceFn(accum, slice[idx])
	}
	return accum
}

func FilterSlice[T any](slice []T, filterFn func(idx int, v T) bool) []T {
	result := make([]T, 0, len(slice))
	for idx, v := range slice {
		if filterFn(idx, v) {
			result = append(result, v)
		}
	}
	return ClipSlice(result)
}

func FilterSliceI[T any](slice []T, filterFn func(idx int) bool) []T {
	result := make([]T, 0, len(slice))
	for idx, v := range slice {
		if filterFn(idx) {
			result = append(result, v)
		}
	}
	return ClipSlice(result)
}

func FilterSliceV[T any](slice []T, filterFn func(v T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if filterFn(v) {
			result = append(result, v)
		}
	}
	return ClipSlice(result)
}

func FlattenSlice2[T any](slice [][]T) []T {
	var result []T = make([]T, 0, len(slice))
	for _, item := range slice {
		result = append(result, item...)
	}
	return ClipSlice(result)
}

func FlattenSlice3[T any](slice [][][]T) []T {
	var result []T = make([]T, 0, len(slice))
	for _, item := range slice {
		for _, item := range item {
			result = append(result, item...)
		}
	}
	return ClipSlice(result)
}

func ClipSlice[T any](slice []T) []T {
	return slice[:len(slice):len(slice)]
}
