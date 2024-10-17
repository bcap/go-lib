package namedvalue

type NV[T any] struct {
	Name  string
	Value T
}

func New[T any](name string, value T) NV[T] {
	return NV[T]{Name: name, Value: value}
}
