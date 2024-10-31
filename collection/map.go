package collection

import (
	"sort"
)

func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	result := map[K]V{}
	for k, v := range m {
		result[k] = v
	}
	return result
}

func MapMap[K1, K2 comparable, V1, V2 any](m map[K1]V1, mapFn func(K1, V1) (K2, V2)) map[K2]V2 {
	result := map[K2]V2{}
	for k1, v1 := range m {
		k2, v2 := mapFn(k1, v1)
		result[k2] = v2
	}
	return result
}

func MapMapK[K1, K2 comparable, V any](m map[K1]V, mapFn func(K1) K2) map[K2]V {
	result := map[K2]V{}
	for k1, v := range m {
		result[mapFn(k1)] = v
	}
	return result
}

func MapMapV[K comparable, V1, V2 any](m map[K]V1, mapFn func(V1) V2) map[K]V2 {
	result := map[K]V2{}
	for k, v := range m {
		result[k] = mapFn(v)
	}
	return result
}

func ReduceMap[K comparable, T1, T2 any](m map[K]T1, reduceFn func(accum T2, key K, value T1) T2) T2 {
	var accum T2
	for key := range m {
		accum = reduceFn(accum, key, m[key])
	}
	return accum
}

func ReduceMapK[K comparable, T1, T2 any](m map[K]T1, reduceFn func(accum T2, key K) T2) T2 {
	var accum T2
	for key := range m {
		accum = reduceFn(accum, key)
	}
	return accum
}

func ReduceMapV[K comparable, T1, T2 any](m map[K]T1, reduceFn func(accum T2, value T1) T2) T2 {
	var accum T2
	for key := range m {
		accum = reduceFn(accum, m[key])
	}
	return accum
}

func FilterMap[K comparable, V any](m map[K]V, filterFn func(K, V) bool) map[K]V {
	result := map[K]V{}
	for k, v := range m {
		if filterFn(k, v) {
			result[k] = v
		}
	}
	return result
}

func FilterMapK[K comparable, V any](m map[K]V, filterFn func(K) bool) map[K]V {
	result := map[K]V{}
	for k := range m {
		if filterFn(k) {
			result[k] = m[k]
		}
	}
	return result
}

func FilterMapV[K comparable, V any](m map[K]V, filterFn func(V) bool) map[K]V {
	result := map[K]V{}
	for k, v := range m {
		if filterFn(v) {
			result[k] = v
		}
	}
	return result
}

type MapEntry[K comparable, V any] struct {
	K K
	V V
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	slice := make([]K, 0, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	return slice
}

func MapValues[K comparable, V any](m map[K]V) []V {
	slice := make([]V, 0, len(m))
	for _, v := range m {
		slice = append(slice, v)
	}
	return slice
}

func MapEntries[K comparable, V any](m map[K]V) MapEntriesS[K, V] {
	slice := make([]MapEntry[K, V], 0, len(m))
	for k, v := range m {
		slice = append(slice, MapEntry[K, V]{K: k, V: v})
	}
	return slice
}

type MapEntriesS[K comparable, V any] []MapEntry[K, V]

func (s MapEntriesS[K, V]) Keys() []K {
	slice := make([]K, 0, len(s))
	for _, e := range s {
		slice = append(slice, e.K)
	}
	return slice
}

func (s MapEntriesS[K, V]) Values() []V {
	slice := make([]V, 0, len(s))
	for _, e := range s {
		slice = append(slice, e.V)
	}
	return slice
}

func SortedMap[K comparable, V any](m map[K]V, less func(a *MapEntry[K, V], b *MapEntry[K, V]) bool) MapEntriesS[K, V] {
	entries := MapEntries(m)
	sort.Slice(entries, func(i, j int) bool { return less(&entries[i], &entries[j]) })
	return entries
}

func JoinMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	var totalLength int
	for _, m := range maps {
		totalLength += len(m)
	}
	result := make(map[K]V, totalLength)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
