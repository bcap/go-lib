package safe

func Deref[T any](v *T) T {
	if v == nil {
		var zeroVal T
		return zeroVal
	}
	return *v
}

func DerefD[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}
