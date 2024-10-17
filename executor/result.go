package executor

import "fmt"

type Result[T any] struct {
	Value T
	Err   error
}

func (r *Result[T]) UnWrap() (T, error) {
	return r.Value, r.Err
}

func (r *Result[T]) IsError() bool {
	return r.Err != nil
}

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
