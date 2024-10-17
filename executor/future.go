package executor

import (
	"context"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type FutureState int32

const (
	AwaitingExecution FutureState = 0
	Executing         FutureState = 1
	ResultReady       FutureState = 2
	ResultStored      FutureState = 3
)

type Future[T any] struct {
	resultC     chan *Result[T]
	result      atomic.Value
	resultFetch semaphore.Weighted
	state       atomic.Int32
}

func newFuture[T any]() *Future[T] {
	return &Future[T]{
		resultC:     make(chan *Result[T], 1),
		resultFetch: *semaphore.NewWeighted(1),
	}
}

func (f *Future[T]) Get(ctx context.Context) *Result[T] {
	return f.get(ctx, true)
}

func (f *Future[T]) GetNoBlock() (*Result[T], bool) {
	v := f.get(context.Background(), false)
	return v, v != nil
}

func (f *Future[T]) get(ctx context.Context, block bool) *Result[T] {
	result := f.result.Load()
	if result != nil {
		return result.(*Result[T])
	}

	// Only one goroutine can fetch the result from the channel
	// All others wait on its result
	if err := f.resultFetch.Acquire(ctx, 1); err != nil {
		return &Result[T]{Err: err}
	}
	defer f.resultFetch.Release(1)

	// double checked locking
	result = f.result.Load()
	if result != nil {
		return result.(*Result[T])
	}

	var r *Result[T]
	if block {
		select {
		case r = <-f.resultC:
		case <-ctx.Done():
			return &Result[T]{Err: ctx.Err()}
		}
	} else {
		select {
		case r = <-f.resultC:
		default:
			return nil
		}
	}

	f.result.Store(r)
	f.state.Store(int32(ResultStored))
	return r
}

func (f *Future[T]) State() FutureState {
	return FutureState(f.state.Load())
}

func (f *Future[T]) setState(state FutureState) {
	f.state.Store(int32(state))
}

func (f *Future[T]) IsDone() bool {
	state := f.State()
	return state == ResultReady || state == ResultStored
}
