package executor

import (
	"context"
	"runtime"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type Executor[T any] struct {
	maxParallelism int
	semaphore      *semaphore.Weighted
	launched       atomic.Int64
	inFlight       atomic.Int64
	done           atomic.Int64

	testSyncCheckpoint chan struct{} // used only for testing manipulation
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

// Submits sends a new job to the executor.
// The job will be executed as soon as possible, respecting the maxParallelism setting
func (e *Executor[T]) Submit(ctx context.Context, fn func() (T, error)) *Future[T] {
	future := newFuture[T]()
	future.state.Store(int32(AwaitingExecution))
	go func() {
		e.launched.Add(1)
		defer e.done.Add(1)

		// inspection point used in testing
		if e.testSyncCheckpoint != nil {
			<-e.testSyncCheckpoint
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

func (e *Executor[T]) MaxParallelism() int {
	return e.maxParallelism
}

func (e *Executor[T]) Launched() int64 {
	return e.launched.Load()
}

func (e *Executor[T]) InFlight() int64 {
	return e.inFlight.Load()
}

func (e *Executor[T]) Done() int64 {
	return e.done.Load()
}

func (e *Executor[T]) Active() bool {
	return e.InFlight() > 0 || e.Launched() > e.Done()
}
