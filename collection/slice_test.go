package collection

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	r := []int{5, 4, 3, 2, 1}

	require.Equal(t, r, ReverseSlice(s))
	require.Equal(t, s, ReverseSlice(r))
	require.Equal(t, len(s), cap(ReverseSlice(s)))

	require.Equal(t, s, SortedSlice(s, func(a, b int) bool { return a < b }))
	require.Equal(t, r, SortedSlice(s, func(a, b int) bool { return a > b }))
	require.Equal(t, s, SortedSlice(r, func(a, b int) bool { return a < b }))
	require.Equal(t, len(s), cap(SortedSlice(r, func(a, b int) bool { return a < b })))

	require.Equal(t, s, JoinSlices(s))
	require.Equal(t, []int{1, 2, 3, 4, 5, 5, 4, 3, 2, 1}, JoinSlices(s, r))
	require.Equal(t, []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1}, JoinSlices(s, s, r))

	copy1 := []int{1, 2, 3, 4, 5}
	copy2 := CopySlice(copy1)
	require.Equal(t, s, copy1)
	require.Equal(t, s, copy2)
	require.NotEqual(t, unsafe.Pointer(&copy1), unsafe.Pointer(&copy2))
	copy2[0] = 999
	require.Equal(t, s, copy1)
	require.NotEqual(t, s, copy2)

	copy1 = CopySlice(s)
	MapSliceInplace(copy1, func(idx int, v int) int { return v * 10 })
	require.Equal(t, []int{10, 20, 30, 40, 50}, copy1)

	copy1 = CopySlice(s)
	MapSliceInplaceI(copy1, func(idx int) int { return s[idx] * 10 })
	require.Equal(t, []int{10, 20, 30, 40, 50}, copy1)

	copy1 = CopySlice(s)
	MapSliceInplaceV(copy1, func(v int) int { return v * 10 })
	require.Equal(t, []int{10, 20, 30, 40, 50}, copy1)

	require.Equal(t, []int{10, 20, 30, 40, 50}, MapSlice(s, func(idx int, v int) int { return v * 10 }))
	require.Equal(t, []int{10, 20, 30, 40, 50}, MapSliceI(s, func(idx int) int { return s[idx] * 10 }))
	require.Equal(t, []int{10, 20, 30, 40, 50}, MapSliceV(s, func(v int) int { return v * 10 }))

	require.Equal(t, []any{10, 20, 30, 40, 50}, MapSlice(s, func(idx int, v int) any { return v * 10 }))
	require.Equal(t, []any{10, 20, 30, 40, 50}, MapSliceI(s, func(idx int) any { return s[idx] * 10 }))
	require.Equal(t, []any{10, 20, 30, 40, 50}, MapSliceV(s, func(v int) any { return v * 10 }))

	require.Equal(t, []int{1, 2, 3}, FilterSlice(s, func(idx int, v int) bool { return v < 4 }))
	require.Equal(t, []int{1, 2, 3}, FilterSliceI(s, func(idx int) bool { return s[idx] < 4 }))
	require.Equal(t, []int{1, 2, 3}, FilterSliceV(s, func(v int) bool { return v < 4 }))

	require.Equal(t, 15, ReduceSlice(s, func(acc int, idx int, v int) int { return acc + v }))
	require.Equal(t, 15, ReduceSliceI(s, func(acc int, idx int) int { return acc + s[idx] }))
	require.Equal(t, 15, ReduceSliceV(s, func(acc int, v int) int { return acc + v }))

	require.Equal(t, 15.0, ReduceSlice(s, func(acc float64, idx int, v int) float64 { return acc + float64(v) }))
	require.Equal(t, 15.0, ReduceSliceI(s, func(acc float64, idx int) float64 { return acc + float64(s[idx]) }))
	require.Equal(t, 15.0, ReduceSliceV(s, func(acc float64, v int) float64 { return acc + float64(v) }))
}

func TestSliceFlatten(t *testing.T) {
	s := [][]int{
		{},
		{1, 2, 3},
		{4, 5},
		{},
		{6},
	}
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, FlattenSlice2(s))

	s2 := [][][]int{
		{
			{1, 2},
			{3},
		},
		{},
		{
			{},
			{},
		},
		{
			{4, 5},
			{},
		},
		{
			{},
			{6},
		},
	}
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, FlattenSlice3(s2))

	require.Equal(t, [][]int{{1, 2}, {3}, {}, {}, {4, 5}, {}, {}, {6}}, FlattenSlice2(s2))
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, FlattenSlice2(FlattenSlice2(s2)))
}
