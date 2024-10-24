package collection

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAggregate(t *testing.T) {
	type st struct {
		a int
		b string
	}

	s := []st{
		{1, "a"}, {2, "b"}, {3, "c"},
		{1, "a"}, {2, "b"}, {3, "c"}, // repeated line
		{4, "d"}, {5, "e"}, {6, "f"},
	}
	expected := map[st]int{
		{1, "a"}: 2,
		{2, "b"}: 4,
		{3, "c"}: 6,
		{4, "d"}: 4,
		{5, "e"}: 5,
		{6, "f"}: 6,
	}

	require.Equal(t,
		expected,
		Aggregate(s, func(acumm int, idx int, value st) int { return acumm + value.a }),
	)

	require.Equal(t,
		expected,
		AggregateI(s, func(acumm int, idx int) int { return acumm + s[idx].a }),
	)

	require.Equal(t,
		expected,
		AggregateV(s, func(acumm int, value st) int { return acumm + value.a }),
	)

	require.Equal(t,
		map[st]int{
			{1, "a"}: 2,
			{2, "b"}: 2,
			{3, "c"}: 2,
			{4, "d"}: 1,
			{5, "e"}: 1,
			{6, "f"}: 1,
		},
		Count(s),
	)
}

func TestGroupBy(t *testing.T) {
	type person struct {
		name string
		age  int
	}

	s := []person{
		{"Alice", 20},
		{"Bob", 30},
		{"Alice", 25},
		{"Eve", 35},
		{"Alice", 22},
		{"Bob", 40},
	}

	require.Equal(t,
		map[string][]person{
			"Alice": {{"Alice", 20}, {"Alice", 25}, {"Alice", 22}},
			"Bob":   {{"Bob", 30}, {"Bob", 40}},
			"Eve":   {{"Eve", 35}},
		},
		GroupBy(s, func(p person) string { return p.name }),
	)

	require.Equal(t,
		map[string][]person{
			"above 30": {{"Eve", 35}, {"Bob", 40}},
			"below 30": {{"Alice", 20}, {"Bob", 30}, {"Alice", 25}, {"Alice", 22}},
		},
		GroupBy(s, func(p person) string {
			if p.age > 30 {
				return "above 30"
			} else {
				return "below 30"
			}
		}),
	)
}
