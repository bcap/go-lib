package collection

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	r := []int{5, 4, 3, 2, 1}

	m := SliceToMap(s, func(v int) (int, int) { return v, v })

	require.Equal(t, len(s), len(m))
	for i := 1; i <= len(s); i++ {
		require.Equal(t, i, m[i])
	}

	entries := MapEntries(m)
	require.Equal(t, len(s), len(entries))
	require.Equal(t, len(s), cap(entries))

	sort.Slice(entries, func(i, j int) bool { return entries[i].K < entries[j].K })
	for i := 0; i < len(s); i++ {
		require.Equal(t, i+1, entries[i].K)
		require.Equal(t, i+1, entries[i].V)
	}

	testInts := func(ints []int) {
		require.Equal(t, len(s), len(ints))
		require.Equal(t, len(s), cap(ints))
		sort.Ints(ints)
		require.Equal(t, s, ints)
	}

	testInts(entries.Keys())
	testInts(MapKeys(m))
	testInts(entries.Values())
	testInts(MapValues(m))

	require.Equal(t, s, SortedMap(m, func(a, b *MapEntry[int, int]) bool { return a.K < b.K }).Keys())
	require.Equal(t, s, SortedMap(m, func(a, b *MapEntry[int, int]) bool { return a.V < b.V }).Values())

	require.Equal(t, r, SortedMap(m, func(a, b *MapEntry[int, int]) bool { return a.K > b.K }).Keys())
	require.Equal(t, r, SortedMap(m, func(a, b *MapEntry[int, int]) bool { return a.V > b.V }).Values())

	require.Equal(t,
		map[string]int{"a": 1, "b": 2, "c": 3},
		JoinMaps(
			map[string]int{"a": 1},
			map[string]int{"b": 2},
			map[string]int{"c": 3},
		),
	)
	require.Equal(t,
		map[string]int{"a": 1, "b": 2, "c": 3},
		JoinMaps(
			map[string]int{"a": 10, "b": 2},
			map[string]int{"a": 1},
			map[string]int{"c": 3},
		),
	)
	require.Equal(t,
		map[string]int{"a": 10, "b": 2, "c": 3},
		JoinMaps(
			map[string]int{"a": 1},
			map[string]int{"a": 10, "b": 2},
			map[string]int{"c": 3},
		),
	)
}

func TestSet(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	r := []int{5, 4, 3, 2, 1}

	set := SlicesToSet(s)
	require.Equal(t, len(s), len(set))
	for i := 1; i <= len(s); i++ {
		v, ok := set[i]
		require.True(t, ok)
		require.Equal(t, struct{}{}, v)
	}
	require.Equal(t, set, SlicesToSet(r))
	require.Equal(t, set, SlicesToSet(s, r))
}
