package executor

import (
	"context"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type Executor[T any] struct {
	maxParallelism int
	semaphore      *semaphore.Weighted
	launched       atomic.Int64
	inFlight       atomic.Int64
	done           atomic.Int64

	testWg *sync.WaitGroup // used only for testing manipulation
}

func New[T any](maxParallelism int) *Executor[T] {
	if maxParallelism == 0 {
		maxParallelism = runtime.NumCPU()
	}
	return &Executor[T]{
		maxParallelism: maxParallelism,
		semaphore:      semaphore.NewWeighted(int64(maxParallelism)),
	}
}

func (e *Executor[T]) Submit(ctx context.Context, fn func() (T, error)) *Future[T] {
	future := newFuture[T]()
	future.state.Store(int32(AwaitingExecution))
	go func() {
		e.launched.Add(1)
		defer e.done.Add(1)

		// inspection point used in testing
		if e.testWg != nil {
			e.testWg.Wait()
		}

		// wait until its our turn to run the function
		if ctx == nil {
			ctx = context.Background()
		}
		if e.maxParallelism > 0 {
			if err := e.semaphore.Acquire(ctx, 1); err != nil {
				future.resultC <- &Result[T]{Err: err}
				return
			}
		}

		// Critical Section Start | run the function
		e.inFlight.Add(1)
		future.setState(Executing)
		result, err := fn()

		// Critical Section Stop | allow the next function to run
		e.inFlight.Add(-1)
		e.semaphore.Release(1)

		// send results
		future.resultC <- &Result[T]{Value: result, Err: err}
		future.setState(ResultReady)
	}()
	return future
}

func (e *Executor[T]) SubmitSlice(ctx context.Context, slice any, fn func(int) (T, error)) []*Future[T] {
	if slice == nil {
		return []*Future[T]{}
	}
	sliceVal := reflect.ValueOf(slice)
	sliceKind := sliceVal.Kind()
	if sliceKind != reflect.Slice && sliceKind != reflect.Array {
		panic("expected a slice or array")
	}
	futures := make([]*Future[T], sliceVal.Len())
	for i := range futures {
		futures[i] = e.Submit(ctx, func() (T, error) {
			return fn(i)
		})
	}
	return futures
}

func (e *Executor[T]) Collect(ctx context.Context, slice any, fn func(int) (T, error)) Results[T] {
	return CollectFutures(ctx, e.SubmitSlice(ctx, slice, fn))
}

func CollectFutures[T any](ctx context.Context, futures []*Future[T]) Results[T] {
	results := make([]*Result[T], len(futures))
	for i, f := range futures {
		results[i] = f.Get(ctx)
	}
	return results
}
