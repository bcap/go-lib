package collection

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	r := []int{5, 4, 3, 2, 1}

	assert.Equal(t, r, ReverseSlice(s))
	assert.Equal(t, s, ReverseSlice(r))
	assert.Equal(t, len(s), cap(ReverseSlice(s)))

	assert.Equal(t, s, SortedSlice(s, func(a, b int) bool { return a < b }))
	assert.Equal(t, r, SortedSlice(s, func(a, b int) bool { return a > b }))
	assert.Equal(t, s, SortedSlice(r, func(a, b int) bool { return a < b }))
	assert.Equal(t, len(s), cap(SortedSlice(r, func(a, b int) bool { return a < b })))

	assert.Equal(t, s, JoinSlices(s))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 5, 4, 3, 2, 1}, JoinSlices(s, r))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1}, JoinSlices(s, s, r))
}

func TestMap(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	r := []int{5, 4, 3, 2, 1}

	m := SliceToMap(s, func(v int) (int, int) { return v, v })

	assert.Equal(t, len(s), len(m))
	for i := 1; i <= len(s); i++ {
		assert.Equal(t, i, m[i])
	}

	entries := MapEntries(m)
	assert.Equal(t, len(s), len(entries))
	assert.Equal(t, len(s), cap(entries))

	sort.Slice(entries, func(i, j int) bool { return entries[i].K < entries[j].K })
	for i := 0; i < len(s); i++ {
		assert.Equal(t, i+1, entries[i].K)
		assert.Equal(t, i+1, entries[i].V)
	}

	testInts := func(ints []int) {
		assert.Equal(t, len(s), len(ints))
		assert.Equal(t, len(s), cap(ints))
		sort.Ints(ints)
		assert.Equal(t, s, ints)
	}

	testInts(entries.Keys())
	testInts(MapKeys(m))
	testInts(entries.Values())
	testInts(MapValues(m))

	assert.Equal(t, s, SortedMapKeys(m, func(a, b *MapEntry[int, int]) bool { return a.K < b.K }))
	assert.Equal(t, s, SortedMapValues(m, func(a, b *MapEntry[int, int]) bool { return a.V < b.V }))

	assert.Equal(t, r, SortedMapKeys(m, func(a, b *MapEntry[int, int]) bool { return a.K > b.K }))
	assert.Equal(t, r, SortedMapValues(m, func(a, b *MapEntry[int, int]) bool { return a.V > b.V }))

	assert.Equal(t,
		map[string]int{"a": 1, "b": 2, "c": 3},
		JoinMaps(
			map[string]int{"a": 1},
			map[string]int{"b": 2},
			map[string]int{"c": 3},
		),
	)
	assert.Equal(t,
		map[string]int{"a": 1, "b": 2, "c": 3},
		JoinMaps(
			map[string]int{"a": 10, "b": 2},
			map[string]int{"a": 1},
			map[string]int{"c": 3},
		),
	)
	assert.Equal(t,
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
	assert.Equal(t, len(s), len(set))
	for i := 1; i <= len(s); i++ {
		v, ok := set[i]
		assert.True(t, ok)
		assert.Equal(t, struct{}{}, v)
	}
	assert.Equal(t, set, SlicesToSet(r))
	assert.Equal(t, set, SlicesToSet(s, r))
}
