package collection

func Aggregate[T comparable, V any](slice []T, aggrFn func(acumm V, idx int, value T) V) map[T]V {
	result := map[T]V{}
	for idx, value := range slice {
		result[value] = aggrFn(result[value], idx, value)
	}
	return result
}

func AggregateI[T comparable, V any](slice []T, aggrFn func(acumm V, idx int) V) map[T]V {
	result := map[T]V{}
	for idx, value := range slice {
		result[value] = aggrFn(result[value], idx)
	}
	return result
}

func AggregateV[T comparable, V any](slice []T, aggrFn func(acumm V, value T) V) map[T]V {
	result := map[T]V{}
	for _, value := range slice {
		result[value] = aggrFn(result[value], value)
	}
	return result
}

func Count[T comparable](slice []T) map[T]int {
	result := map[T]int{}
	for _, value := range slice {
		result[value]++
	}
	return result
}

func CountUnique[T comparable](slice []T) int {
	set := map[T]struct{}{}
	for _, value := range slice {
		set[value] = struct{}{}
	}
	return len(set)
}

func GroupBy[T, V comparable](slice []T, selectorFn func(value T) V) map[V][]T {
	result := map[V][]T{}
	for _, value := range slice {
		key := selectorFn(value)
		result[key] = append(result[key], value)
	}
	return result
}
