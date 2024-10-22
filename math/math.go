package math

import "cmp"

func Max[T cmp.Ordered](v ...T) T {
	return choose(v, func(a, b T) bool { return a > b })
}

func Min[T cmp.Ordered](v ...T) T {
	return choose(v, func(a, b T) bool { return a < b })
}

func choose[T cmp.Ordered](values []T, pickLeft func(T, T) bool) T {
	if len(values) == 0 {
		var zeroVal T
		return zeroVal
	}
	if len(values) == 1 {
		return values[0]
	}
	current := values[0]
	for _, v := range values[1:] {
		if pickLeft(v, current) {
			current = v
		}
	}
	return current
}
