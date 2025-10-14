package executor

import (
	"context"
	"sync/atomic"

	"github.com/bcap/go-lib/result"
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
	resultC     chan *result.Result[T]
	result      atomic.Value
	resultFetch semaphore.Weighted
	state       atomic.Int32
}

func newFuture[T any]() *Future[T] {
	return &Future[T]{
		resultC:     make(chan *result.Result[T], 1),
		resultFetch: *semaphore.NewWeighted(1),
	}
}

func (f *Future[T]) Get(ctx context.Context) *result.Result[T] {
	return f.get(ctx, true)
}

func (f *Future[T]) GetNoBlock() (*result.Result[T], bool) {
	v := f.get(context.Background(), false)
	return v, v != nil
}

func (f *Future[T]) get(ctx context.Context, block bool) *result.Result[T] {
	res := f.result.Load()
	if res != nil {
		return res.(*result.Result[T])
	}

	// Only one goroutine can fetch the result from the channel
	// All others wait on its result
	if err := f.resultFetch.Acquire(ctx, 1); err != nil {
		return &result.Result[T]{Err: err}
	}
	defer f.resultFetch.Release(1)

	// double checked locking
	res = f.result.Load()
	if res != nil {
		return res.(*result.Result[T])
	}

	var r *result.Result[T]
	if block {
		select {
		case r = <-f.resultC:
		case <-ctx.Done():
			return &result.Result[T]{Err: ctx.Err()}
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

type Futures[T any] []*Future[T]

func (fs *Futures[T]) Add(f *Future[T]) {
	*fs = append(*fs, f)
}

func (fs Futures[T]) Get(ctx context.Context) result.Results[T] {
	results := make([]*result.Result[T], len(fs))
	for i, f := range fs {
		results[i] = f.Get(ctx)
	}
	return results
}
