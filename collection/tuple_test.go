package collection

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTuple(t *testing.T) {
	type st struct {
		v int
	}

	//
	// Test empty zip in its many forms
	//

	require.Equal(t,
		[]Tuple2[int, string]{},
		Zip2([]int{}, []string{}),
	)

	require.Equal(t,
		[]Tuple3[int, string, st]{},
		Zip3([]int{}, []string{}, []st{}),
	)

	require.Equal(t,
		[]Tuple4[int, string, st, any]{},
		Zip4([]int{}, []string{}, []st{}, []any{}),
	)

	require.Equal(t,
		[]Tuple5[int, string, st, any, []string]{},
		Zip5([]int{}, []string{}, []st{}, []any{}, [][]string{}),
	)

	//
	// Test zipping in its many forms
	//

	require.Equal(t,
		[]Tuple2[int, string]{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		},
		Zip2(
			[]int{1, 2, 3},
			[]string{"a", "b", "c"},
		),
	)

	require.Equal(t,
		[]Tuple3[int, string, st]{
			{1, "a", st{10}},
			{2, "b", st{20}},
			{3, "c", st{30}},
		},
		Zip3(
			[]int{1, 2, 3},
			[]string{"a", "b", "c"},
			[]st{{10}, {20}, {30}},
		),
	)

	require.Equal(t,
		[]Tuple4[int, string, st, any]{
			{1, "a", st{10}, nil},
			{2, "b", st{20}, []string{}},
			{3, "c", st{30}, true},
		},
		Zip4(
			[]int{1, 2, 3},
			[]string{"a", "b", "c"},
			[]st{{10}, {20}, {30}},
			[]any{nil, []string{}, true},
		),
	)

	require.Equal(t,
		[]Tuple5[int, string, st, any, []string]{
			{1, "a", st{10}, nil, nil},
			{2, "b", st{20}, []string{}, []string{"b1", "b2"}},
			{3, "c", st{30}, true, []string{"c1", "c2"}},
		},
		Zip5(
			[]int{1, 2, 3},
			[]string{"a", "b", "c"},
			[]st{{10}, {20}, {30}},
			[]any{nil, []string{}, true},
			[][]string{nil, {"b1", "b2"}, {"c1", "c2"}},
		),
	)

	//
	// Ignore extra elements when zipping
	//
	require.Equal(t,
		[]Tuple2[int, string]{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		},
		Zip2(
			[]int{1, 2, 3},
			[]string{"a", "b", "c", "d", "e", "f"},
		),
	)

	require.Equal(t,
		[]Tuple2[int, string]{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		},
		Zip2(
			[]int{1, 2, 3, 4, 5, 6},
			[]string{"a", "b", "c"},
		),
	)

	require.Equal(t,
		[]Tuple3[int, string, st]{
			{1, "a", st{10}},
			{2, "b", st{20}},
			{3, "c", st{30}},
		},
		Zip3(
			[]int{1, 2, 3, 4, 5, 6},
			[]string{"a", "b", "c", "d", "e", "f"},
			[]st{{10}, {20}, {30}},
		),
	)

	require.Equal(t,
		[]Tuple3[int, string, st]{
			{1, "a", st{10}},
			{2, "b", st{20}},
			{3, "c", st{30}},
		},
		Zip3(
			[]int{1, 2, 3},
			[]string{"a", "b", "c", "d", "e", "f"},
			[]st{{10}, {20}, {30}, {40}, {50}, {60}, {70}},
		),
	)

	require.Equal(t,
		[]Tuple3[int, string, st]{
			{1, "a", st{10}},
			{2, "b", st{20}},
			{3, "c", st{30}},
		},
		Zip3(
			[]int{1, 2, 3, 4, 5, 6},
			[]string{"a", "b", "c"},
			[]st{{10}, {20}, {30}, {40}, {50}, {60}, {70}},
		),
	)

	//
	// Test tuple to map in its many forms
	//

	tuple2s := []Tuple2[int, string]{
		{1, "a"},
		{2, "b"},
		{3, "c"},
	}
	require.Equal(t,
		map[int]string{
			1: "a",
			2: "b",
			3: "c",
		},
		Tuple2ToMap(tuple2s),
	)

	tuple3s := []Tuple3[int, string, float64]{
		{1, "a", 1.1},
		{2, "b", 2.2},
		{3, "c", 3.3},
	}
	require.Equal(t,
		map[int]Tuple2[string, float64]{
			1: {"a", 1.1},
			2: {"b", 2.2},
			3: {"c", 3.3},
		},
		Tuple3ToMap(tuple3s),
	)
	require.Equal(t,
		map[Tuple2[int, string]]float64{
			{1, "a"}: 1.1,
			{2, "b"}: 2.2,
			{3, "c"}: 3.3,
		},
		Tuple3ToMap2(tuple3s),
	)

	tuple4s := []Tuple4[int, string, float64, any]{
		{1, "a", 1.1, nil},
		{2, "b", 2.2, []string{}},
		{3, "c", 3.3, true},
	}
	require.Equal(t,
		map[int]Tuple3[string, float64, any]{
			1: {"a", 1.1, nil},
			2: {"b", 2.2, []string{}},
			3: {"c", 3.3, true},
		},
		Tuple4ToMap(tuple4s),
	)
	require.Equal(t,
		map[Tuple2[int, string]]Tuple2[float64, any]{
			{1, "a"}: {1.1, nil},
			{2, "b"}: {2.2, []string{}},
			{3, "c"}: {3.3, true},
		},
		Tuple4ToMap2(tuple4s),
	)
	require.Equal(t,
		map[Tuple3[int, string, float64]]any{
			{1, "a", 1.1}: nil,
			{2, "b", 2.2}: []string{},
			{3, "c", 3.3}: true,
		},
		Tuple4ToMap3(tuple4s),
	)

	tuple5s := []Tuple5[int, string, float64, any, []string]{
		{1, "a", 1.1, nil, nil},
		{2, "b", 2.2, 5, []string{"b1", "b2"}},
		{3, "c", 3.3, true, []string{"c1", "c2"}},
	}
	require.Equal(t,
		map[int]Tuple4[string, float64, any, []string]{
			1: {"a", 1.1, nil, nil},
			2: {"b", 2.2, 5, []string{"b1", "b2"}},
			3: {"c", 3.3, true, []string{"c1", "c2"}},
		},
		Tuple5ToMap(tuple5s),
	)
	require.Equal(t,
		map[Tuple2[int, string]]Tuple3[float64, any, []string]{
			{1, "a"}: {1.1, nil, nil},
			{2, "b"}: {2.2, 5, []string{"b1", "b2"}},
			{3, "c"}: {3.3, true, []string{"c1", "c2"}},
		},
		Tuple5ToMap2(tuple5s),
	)
	require.Equal(t,
		map[Tuple3[int, string, float64]]Tuple2[any, []string]{
			{1, "a", 1.1}: {nil, nil},
			{2, "b", 2.2}: {5, []string{"b1", "b2"}},
			{3, "c", 3.3}: {true, []string{"c1", "c2"}},
		},
		Tuple5ToMap3(tuple5s),
	)
	require.Equal(t,
		map[Tuple4[int, string, float64, any]][]string{
			{1, "a", 1.1, nil}:  nil,
			{2, "b", 2.2, 5}:    {"b1", "b2"},
			{3, "c", 3.3, true}: {"c1", "c2"},
		},
		Tuple5ToMap4(tuple5s),
	)
}
