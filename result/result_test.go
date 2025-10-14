package result

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResult(t *testing.T) {
	r := NewResult(42, nil)
	require.Equal(t, 42, r.Value)
	require.Equal(t, 42, r.Must())
	require.Nil(t, r.Err)
	require.True(t, r.IsOk())
	require.False(t, r.IsError())

	r = NewResult(0, errors.New("an error"))
	require.Equal(t, 0, r.Value)
	require.Panics(t, func() { r.Must() })
	require.Equal(t, errors.New("an error"), r.Err)
	require.False(t, r.IsOk())
	require.True(t, r.IsError())
}

func TestResults(t *testing.T) {
	results := Results[int]{
		NewResult(1, nil),
		NewResult(2, nil),
		NewResult(0, errors.New("an error")),
		NewResult(4, nil),
	}

	values, err := results.Unwrap()
	require.Error(t, err)
	require.Equal(t, []int{1, 2, 0, 4}, values)

	require.Panics(t, func() { results.Must() })

	require.Equal(t, []int{1, 2, 0, 4}, results.Values())
	require.Equal(t, []int{1, 2, 4}, results.ValuesOnly())
	require.Equal(t, []error{nil, nil, errors.New("an error"), nil}, results.Errors())
	require.Equal(t, []error{errors.New("an error")}, results.ErrorsOnly())
	require.True(t, results.HasError())

	ok, nok := results.Stats()
	require.Equal(t, 3, ok)
	require.Equal(t, 1, nok)
}
