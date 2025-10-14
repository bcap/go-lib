package executor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bcap/go-lib/result"
	"github.com/stretchr/testify/require"
)

func TestExecutor(t *testing.T) {
	var testRunning atomic.Bool
	testRunning.Store(true)
	defer testRunning.Store(false)

	parallelism := 5
	length := 1000

	e := New[int](parallelism)
	e.testSyncCheckpoint = make(chan struct{})

	var currentParallelism atomic.Int32
	var maxParallelismAchieved atomic.Int32

	var wg sync.WaitGroup
	wg.Add(1)

	futures := make([]*Future[int], length)
	for i := range futures {
		i := i
		futures[i] = e.Submit(context.Background(), func() (int, error) {
			current := currentParallelism.Add(1)
			if current > maxParallelismAchieved.Load() {
				maxParallelismAchieved.Store(current)
			}
			defer currentParallelism.Add(-1)
			wg.Wait()
			time.Sleep(5 * time.Millisecond)
			return i, nil
		})

		// Constantly try to read results from the futures during the whole test duration
		// This is to check if we trigger race conditions when running `go test -race`
		go func() {
			for testRunning.Load() {
				futures[i].GetNoBlock()
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}

	// Validate that all futures are still in the initial state
	for _, f := range futures {
		require.Equal(t, AwaitingExecution, f.State())
		require.False(t, f.IsDone())
		r, ok := f.GetNoBlock()
		require.Nil(t, r)
		require.False(t, ok)
	}

	// allow goroutines that execute the futures to run
	close(e.testSyncCheckpoint)

	// wait until all futures executions are launched
	start := time.Now()
	for e.launched.Load() < int64(length) {
		if time.Since(start) > 5*time.Second {
			t.Fatal("timeout while waiting for futures to be launched")
		}
		time.Sleep(1 * time.Millisecond)
	}

	require.Equal(t, length, int(e.launched.Load()))
	require.Equal(t, 0, int(e.done.Load()))

	// Allow futures to move forward with their executions
	wg.Done()

	// wait until all futures are done
	start = time.Now()
	for e.done.Load() < int64(length) {
		if time.Since(start) > 5*time.Second {
			t.Fatal("timeout while waiting for futures to be done")
		}
		time.Sleep(1 * time.Millisecond)
	}

	// Validate futures
	for i, f := range futures {
		require.True(t, f.IsDone())
		r := f.Get(context.Background())
		require.Equal(t, ResultStored, f.State())
		require.True(t, f.IsDone())
		require.NoError(t, r.Err)
		require.Equal(t, i, r.Value)

		r2, ok := f.GetNoBlock()
		require.True(t, ok)
		require.Equal(t, r, r2)
		r3 := f.Get(context.Background())
		require.Equal(t, r, r3)
	}

	require.Equal(t, int32(parallelism), maxParallelismAchieved.Load())
	require.Equal(t, int32(0), currentParallelism.Load())

	require.Equal(t, int64(length), e.launched.Load())
	require.Equal(t, int64(length), e.done.Load())
	require.Equal(t, int64(0), e.inFlight.Load())
}

func TestExecutorCollect(t *testing.T) {
	parallelism := 5
	length := 1000

	inputs := make([]int, length)
	for i := range inputs {
		inputs[i] = i
	}

	outputs := Collect(context.Background(), parallelism, inputs, func(i int) (int, error) {
		time.Sleep(5 * time.Millisecond)
		var err error
		// generate an error for every 10th element
		if i%10 == 0 {
			err = fmt.Errorf("error %d", i)
		}
		return i, err
	})

	require.Equal(t, length, len(outputs))
	require.Equal(t, length, len(outputs.Values()))
	require.Equal(t, length, len(outputs.Errors()))
	require.Equal(t, length-length/10, len(outputs.ValuesOnly()))
	require.Equal(t, length/10, len(outputs.ErrorsOnly()))

	ok, nok := outputs.Stats()
	require.Equal(t, length-length/10, ok)
	require.Equal(t, length/10, nok)

	for i, r := range outputs {
		if i%10 == 0 {
			require.Error(t, r.Err)
		} else {
			require.NoError(t, r.Err)
			require.Equal(t, i, r.Value)
		}
	}

	valuesOnly := []int{}
	for i, v := range outputs.Values() {
		if i%10 != 0 {
			valuesOnly = append(valuesOnly, v)
		}
		require.Equal(t, i, v)
	}

	errorsOnly := []error{}
	for i, v := range outputs.Errors() {
		if i%10 == 0 {
			require.Error(t, v)
			errorsOnly = append(errorsOnly, v)
		} else {
			require.NoError(t, v)
		}
	}

	require.Equal(t, valuesOnly, outputs.ValuesOnly())
	require.Equal(t, errorsOnly, outputs.ErrorsOnly())
	require.True(t, outputs.HasError())
}

func TestExecutorCollectMap(t *testing.T) {
	parallelism := 5
	length := 1000

	inputs := make([]int, length)
	for i := range inputs {
		inputs[i] = i
	}

	outputs := CollectMap(context.Background(), parallelism, inputs, func(k int) (int, error) {
		time.Sleep(5 * time.Millisecond)
		var err error
		// generate an error for every 10th element
		if k%10 == 0 {
			err = fmt.Errorf("error %d", k)
		}
		return k, err
	})

	require.Equal(t, length, len(outputs))
	require.Equal(t, length, len(outputs.Values()))
	require.Equal(t, length, len(outputs.Errors()))
	require.Equal(t, length-length/10, len(outputs.ValuesOnly()))
	require.Equal(t, length/10, len(outputs.ErrorsOnly()))

	ok, nok := outputs.Stats()
	require.Equal(t, length-length/10, ok)
	require.Equal(t, length/10, nok)

	for i, r := range outputs {
		if i%10 == 0 {
			require.Error(t, r.Err)
		} else {
			require.NoError(t, r.Err)
			require.Equal(t, i, r.Value)
		}
	}

	valuesOnly := map[int]int{}
	for k, v := range outputs.Values() {
		if k%10 != 0 {
			valuesOnly[k] = v
		}
		require.Equal(t, k, v)
	}

	errorsOnly := map[int]error{}
	for k, v := range outputs.Errors() {
		if k%10 == 0 {
			require.Error(t, v)
			errorsOnly[k] = v
		} else {
			require.NoError(t, v)
		}
	}

	require.Equal(t, valuesOnly, outputs.ValuesOnly())
	require.Equal(t, errorsOnly, outputs.ErrorsOnly())
	require.True(t, outputs.HasError())
}

func TestExecutorCollectMapReplace(t *testing.T) {
	parallelism := 5
	length := 1000

	m := map[int]*result.Result[int]{}
	for i := 0; i < length; i++ {
		m[i] = nil
	}

	CollectMapReplace(context.Background(), parallelism, m, func(k int) (int, error) {
		time.Sleep(5 * time.Millisecond)
		var err error
		// generate an error for every 10th element
		if k%10 == 0 {
			err = fmt.Errorf("error %d", k)
		}
		return k, err
	})

	require.Equal(t, length, len(m))

	for k, r := range m {
		if k%10 == 0 {
			require.Error(t, r.Err)
		} else {
			require.NoError(t, r.Err)
		}
		require.Equal(t, k, r.Value)
	}
}

func TestExecutorCollectEmpty(t *testing.T) {
	testSlice := func(slice any) {
		outputs := Collect(context.Background(), 5, slice, func(i int) (any, error) { return nil, nil })

		require.Equal(t, 0, len(outputs))
		require.Equal(t, 0, len(outputs.Values()))
		require.Equal(t, 0, len(outputs.Errors()))
		require.Equal(t, 0, len(outputs.ValuesOnly()))
		require.Equal(t, 0, len(outputs.ErrorsOnly()))

		ok, nok := outputs.Stats()
		require.Equal(t, 0, ok)
		require.Equal(t, 0, nok)
		require.False(t, outputs.HasError())
	}

	testMap := func(slice []any) {
		outputs := CollectMap(context.Background(), 5, slice, func(k any) (any, error) { return nil, nil })

		require.Equal(t, 0, len(outputs))
		require.Equal(t, 0, len(outputs.Values()))
		require.Equal(t, 0, len(outputs.Errors()))
		require.Equal(t, 0, len(outputs.ValuesOnly()))
		require.Equal(t, 0, len(outputs.ErrorsOnly()))

		ok, nok := outputs.Stats()
		require.Equal(t, 0, ok)
		require.Equal(t, 0, nok)
		require.False(t, outputs.HasError())
	}

	testSlice([0]any{})
	testSlice([]any{})
	testSlice(nil)

	testMap(nil)
	testMap([]any{})
}

func TestExecutorCollectPanicsOnBadArg(t *testing.T) {
	test := func(input any) {
		defer func() {
			r := recover()
			require.Equal(t, "expected a slice or array", r)
		}()
		Collect(context.Background(), 5, input, func(i int) (any, error) { return nil, nil })
	}

	test(1)
	test(1.0)
	test("string")
	test(true)
	test(struct{}{})
	test(map[any]any{})
	test(make(chan any))
}
