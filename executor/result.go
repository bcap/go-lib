package executor

import "fmt"

// Result represents the return of a function call, which is normally a value provided by the user or an error
// Inspired by Rust's Result type
type Result[T any] struct {
	Value T
	Err   error
}

func NewResult[T any](value T, err error) *Result[T] {
	return &Result[T]{
		Value: value,
		Err:   err,
	}
}

func (r *Result[T]) UnWrap() (T, error) {
	return r.Value, r.Err
}

func (r *Result[T]) Must() T {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.Value
}

func (r *Result[T]) IsOk() bool {
	return r.Err == nil
}

func (r *Result[T]) IsError() bool {
	return r.Err != nil
}

// Results is a slice of results
type Results[T any] []*Result[T]

func (rs Results[T]) Values() []T {
	values := make([]T, len(rs))
	for i, r := range rs {
		values[i] = r.Value
	}
	return values
}

func (rs Results[T]) Errors() []error {
	errors := make([]error, len(rs))
	for i, r := range rs {
		errors[i] = r.Err
	}
	return errors
}

func (rs Results[T]) ValuesOnly() []T {
	values := []T{}
	for _, r := range rs {
		if r.Err == nil {
			values = append(values, r.Value)
		}
	}
	return values
}

func (rs Results[T]) ErrorsOnly() []error {
	errors := []error{}
	for _, r := range rs {
		if r.Err != nil {
			errors = append(errors, r.Err)
		}
	}
	return errors
}

func (rs Results[T]) Stats() (int, int) {
	var errors int
	for _, r := range rs {
		if r.Err != nil {
			errors++
		}
	}
	return len(rs) - errors, errors
}

func (rs Results[T]) HasError() bool {
	for _, r := range rs {
		if r.Err != nil {
			return true
		}
	}
	return false
}

// ResultsMap is a map of key -> result
type ResultsMap[K comparable, T any] map[K]*Result[T]

func (rs ResultsMap[K, T]) Values() map[K]T {
	values := map[K]T{}
	for k, r := range rs {
		values[k] = r.Value
	}
	return values
}

func (rs ResultsMap[K, T]) Errors() map[K]error {
	errors := map[K]error{}
	for k, r := range rs {
		errors[k] = r.Err
	}
	return errors
}

func (rs ResultsMap[K, T]) ValuesOnly() map[K]T {
	values := map[K]T{}
	for k, r := range rs {
		if r.Err == nil {
			values[k] = r.Value
		}
	}
	return values
}

func (rs ResultsMap[K, T]) ErrorsOnly() map[K]error {
	errors := map[K]error{}
	for k, r := range rs {
		if r.Err != nil {
			errors[k] = r.Err
		}
	}
	return errors
}

func (rs ResultsMap[K, T]) Stats() (int, int) {
	var errors int
	for _, r := range rs {
		if r.Err != nil {
			errors++
		}
	}
	return len(rs) - errors, errors
}

func (rs ResultsMap[K, T]) HasError() bool {
	for _, r := range rs {
		if r.Err != nil {
			return true
		}
	}
	return false
}

// Error
type ResultsError struct {
	Errors []error
}

func (e ResultsError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	if len(e.Errors) == 1 {
		return fmt.Sprintf("error occurred: %v", e.Errors[0])
	}
	return fmt.Sprintf("multiple errors occurred (%d): [%v]", len(e.Errors), e.Errors)
}

func (rs Results[T]) Error() error {
	errors := rs.ErrorsOnly()
	if len(errors) == 0 {
		return nil
	}
	return &ResultsError{Errors: errors}
}

func (rs ResultsMap[K, T]) Error() error {
	errorsMap := rs.ErrorsOnly()
	if len(errorsMap) == 0 {
		return nil
	}
	errors := make([]error, 0, len(errorsMap))
	for _, err := range errorsMap {
		errors = append(errors, err)
	}
	return &ResultsError{Errors: errors}
}
