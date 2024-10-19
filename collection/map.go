package collection

import (
	"sort"
)

type MapEntry[K comparable, V any] struct {
	K K
	V V
}

func MapEntries[K comparable, V any](m map[K]V) MapEntriesS[K, V] {
	slice := make([]MapEntry[K, V], 0, len(m))
	for k, v := range m {
		slice = append(slice, MapEntry[K, V]{K: k, V: v})
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

func MapKeys[K comparable, V any](m map[K]V) []K {
	slice := make([]K, 0, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	return slice
}

func SortedMap[K comparable, V any](m map[K]V, less func(a *MapEntry[K, V], b *MapEntry[K, V]) bool) MapEntriesS[K, V] {
	entries := MapEntries(m)
	sort.Slice(entries, func(i, j int) bool { return less(&entries[i], &entries[j]) })
	return entries
}

func SortedMapKeys[K comparable, V any](m map[K]V, less func(a *MapEntry[K, V], b *MapEntry[K, V]) bool) []K {
	return SortedMap(m, less).Keys()
}

func SortedMapValues[K comparable, V any](m map[K]V, less func(a *MapEntry[K, V], b *MapEntry[K, V]) bool) []V {
	return SortedMap(m, less).Values()
}
