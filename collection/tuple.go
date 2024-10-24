package collection

type Tuple2[T1, T2 any] struct {
	V1 T1
	V2 T2
}

type Tuple3[T1, T2, T3 any] struct {
	V1 T1
	V2 T2
	V3 T3
}

type Tuple4[T1, T2, T3, T4 any] struct {
	V1 T1
	V2 T2
	V3 T3
	V4 T4
}

type Tuple5[T1, T2, T3, T4, T5 any] struct {
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
}

func Zip2[T1, T2 any](slice1 []T1, slice2 []T2) []Tuple2[T1, T2] {
	minLength := len(slice1)
	if len(slice2) < len(slice1) {
		minLength = len(slice2)
	}
	result := []Tuple2[T1, T2]{}
	for i := 0; i < minLength; i++ {
		entry := Tuple2[T1, T2]{
			V1: slice1[i],
			V2: slice2[i],
		}
		result = append(result, entry)
	}
	return result
}

func Zip3[T1, T2, T3 any](slice1 []T1, slice2 []T2, slice3 []T3) []Tuple3[T1, T2, T3] {
	minLength := len(slice1)
	if len(slice2) < minLength {
		minLength = len(slice2)
	}
	if len(slice3) < minLength {
		minLength = len(slice3)
	}
	result := []Tuple3[T1, T2, T3]{}
	for i := 0; i < minLength; i++ {
		entry := Tuple3[T1, T2, T3]{
			V1: slice1[i],
			V2: slice2[i],
			V3: slice3[i],
		}
		result = append(result, entry)
	}
	return result
}

func Zip4[T1, T2, T3, T4 any](slice1 []T1, slice2 []T2, slice3 []T3, slice4 []T4) []Tuple4[T1, T2, T3, T4] {
	minLength := len(slice1)
	if len(slice2) < minLength {
		minLength = len(slice2)
	}
	if len(slice3) < minLength {
		minLength = len(slice3)
	}
	if len(slice4) < minLength {
		minLength = len(slice4)
	}
	result := []Tuple4[T1, T2, T3, T4]{}
	for i := 0; i < minLength; i++ {
		entry := Tuple4[T1, T2, T3, T4]{
			V1: slice1[i],
			V2: slice2[i],
			V3: slice3[i],
			V4: slice4[i],
		}
		result = append(result, entry)
	}
	return result
}

func Zip5[T1, T2, T3, T4, T5 any](slice1 []T1, slice2 []T2, slice3 []T3, slice4 []T4, slice5 []T5) []Tuple5[T1, T2, T3, T4, T5] {
	minLength := len(slice1)
	if len(slice2) < minLength {
		minLength = len(slice2)
	}
	if len(slice3) < minLength {
		minLength = len(slice3)
	}
	if len(slice4) < minLength {
		minLength = len(slice4)
	}
	if len(slice5) < minLength {
		minLength = len(slice5)
	}
	result := []Tuple5[T1, T2, T3, T4, T5]{}
	for i := 0; i < minLength; i++ {
		entry := Tuple5[T1, T2, T3, T4, T5]{
			V1: slice1[i],
			V2: slice2[i],
			V3: slice3[i],
			V4: slice4[i],
			V5: slice5[i],
		}
		result = append(result, entry)
	}
	return result
}

func Tuple2ToMap[T1 comparable, T2 any](tuples []Tuple2[T1, T2]) map[T1]T2 {
	result := map[T1]T2{}
	for _, tuple := range tuples {
		result[tuple.V1] = tuple.V2
	}
	return result
}

func Tuple3ToMap[T1 comparable, T2, T3 any](tuples []Tuple3[T1, T2, T3]) map[T1]Tuple2[T2, T3] {
	result := map[T1]Tuple2[T2, T3]{}
	for _, tuple := range tuples {
		result[tuple.V1] = Tuple2[T2, T3]{V1: tuple.V2, V2: tuple.V3}
	}
	return result
}

func Tuple3ToMap2[T1, T2 comparable, T3 any](tuples []Tuple3[T1, T2, T3]) map[Tuple2[T1, T2]]T3 {
	result := map[Tuple2[T1, T2]]T3{}
	for _, tuple := range tuples {
		result[Tuple2[T1, T2]{V1: tuple.V1, V2: tuple.V2}] = tuple.V3
	}
	return result
}

func Tuple4ToMap[T1 comparable, T2, T3, T4 any](tuples []Tuple4[T1, T2, T3, T4]) map[T1]Tuple3[T2, T3, T4] {
	result := map[T1]Tuple3[T2, T3, T4]{}
	for _, tuple := range tuples {
		result[tuple.V1] = Tuple3[T2, T3, T4]{V1: tuple.V2, V2: tuple.V3, V3: tuple.V4}
	}
	return result
}

func Tuple4ToMap2[T1, T2 comparable, T3, T4 any](tuples []Tuple4[T1, T2, T3, T4]) map[Tuple2[T1, T2]]Tuple2[T3, T4] {
	result := map[Tuple2[T1, T2]]Tuple2[T3, T4]{}
	for _, tuple := range tuples {
		result[Tuple2[T1, T2]{V1: tuple.V1, V2: tuple.V2}] = Tuple2[T3, T4]{V1: tuple.V3, V2: tuple.V4}
	}
	return result
}

func Tuple4ToMap3[T1, T2, T3 comparable, T4 any](tuples []Tuple4[T1, T2, T3, T4]) map[Tuple3[T1, T2, T3]]T4 {
	result := map[Tuple3[T1, T2, T3]]T4{}
	for _, tuple := range tuples {
		result[Tuple3[T1, T2, T3]{V1: tuple.V1, V2: tuple.V2, V3: tuple.V3}] = tuple.V4
	}
	return result
}

func Tuple5ToMap[T1 comparable, T2, T3, T4, T5 any](tuples []Tuple5[T1, T2, T3, T4, T5]) map[T1]Tuple4[T2, T3, T4, T5] {
	result := map[T1]Tuple4[T2, T3, T4, T5]{}
	for _, tuple := range tuples {
		result[tuple.V1] = Tuple4[T2, T3, T4, T5]{V1: tuple.V2, V2: tuple.V3, V3: tuple.V4, V4: tuple.V5}
	}
	return result
}

func Tuple5ToMap2[T1, T2 comparable, T3, T4, T5 any](tuples []Tuple5[T1, T2, T3, T4, T5]) map[Tuple2[T1, T2]]Tuple3[T3, T4, T5] {
	result := map[Tuple2[T1, T2]]Tuple3[T3, T4, T5]{}
	for _, tuple := range tuples {
		result[Tuple2[T1, T2]{V1: tuple.V1, V2: tuple.V2}] = Tuple3[T3, T4, T5]{V1: tuple.V3, V2: tuple.V4, V3: tuple.V5}
	}
	return result
}

func Tuple5ToMap3[T1, T2, T3 comparable, T4, T5 any](tuples []Tuple5[T1, T2, T3, T4, T5]) map[Tuple3[T1, T2, T3]]Tuple2[T4, T5] {
	result := map[Tuple3[T1, T2, T3]]Tuple2[T4, T5]{}
	for _, tuple := range tuples {
		result[Tuple3[T1, T2, T3]{V1: tuple.V1, V2: tuple.V2, V3: tuple.V3}] = Tuple2[T4, T5]{V1: tuple.V4, V2: tuple.V5}
	}
	return result
}

func Tuple5ToMap4[T1, T2, T3, T4 comparable, T5 any](tuples []Tuple5[T1, T2, T3, T4, T5]) map[Tuple4[T1, T2, T3, T4]]T5 {
	result := map[Tuple4[T1, T2, T3, T4]]T5{}
	for _, tuple := range tuples {
		result[Tuple4[T1, T2, T3, T4]{V1: tuple.V1, V2: tuple.V2, V3: tuple.V3, V4: tuple.V4}] = tuple.V5
	}
	return result
}
