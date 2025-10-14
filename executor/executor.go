package executor

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/bcap/go-lib/result"
	"golang.org/x/sync/semaphore"
)

type Executor[T any] struct {
	maxParallelism int
	semaphore      *semaphore.Weighted
	submitted      atomic.Int64
	launched       atomic.Int64
	inFlight       atomic.Int64
	done           atomic.Int64
	pending        atomic.Int64

	waitInactiveLock *sync.Mutex
	waitInactiveCond *sync.Cond

	testSyncCheckpoint chan struct{} // used only for testing manipulation
}

func New[T any](maxParallelism int) *Executor[T] {
	if maxParallelism == 0 {
		maxParallelism = runtime.NumCPU()
	}
	var waitInactiveLock sync.Mutex
	return &Executor[T]{
		maxParallelism:   maxParallelism,
		semaphore:        semaphore.NewWeighted(int64(maxParallelism)),
		waitInactiveCond: sync.NewCond(&waitInactiveLock),
		waitInactiveLock: &waitInactiveLock,
	}
}

// Submits sends a new job to the executor.
// The job will be executed as soon as possible, respecting the maxParallelism setting
func (e *Executor[T]) Submit(ctx context.Context, fn func() (T, error)) *Future[T] {
	e.submitted.Add(1)
	e.pending.Add(1)
	future := newFuture[T]()
	future.state.Store(int32(AwaitingExecution))
	go func() {
		e.launched.Add(1)
		defer e.done.Add(1)
		defer func() {
			pending := e.pending.Add(-1)
			if pending == 0 {
				e.waitInactiveCond.Broadcast()
			}
		}()

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
				future.resultC <- &result.Result[T]{Err: err}
				return
			}
		}

		// Critical Section Start | run the function
		e.inFlight.Add(1)
		future.setState(Executing)
		res, err := fn()

		// Critical Section Stop | allow the next function to run
		e.inFlight.Add(-1)
		if e.maxParallelism > 0 {
			e.semaphore.Release(1)
		}

		// send results
		future.resultC <- &result.Result[T]{Value: res, Err: err}
		future.setState(ResultReady)
	}()
	return future
}

func (e *Executor[T]) MaxParallelism() int {
	return e.maxParallelism
}

func (e *Executor[T]) Submitted() int64 {
	return e.submitted.Load()
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

func (e *Executor[T]) Pending() int64 {
	return e.pending.Load()
}

func (e *Executor[T]) Active() bool {
	return e.Pending() > 0
}

func (e *Executor[T]) Wait() {
	e.waitInactiveLock.Lock()
	e.waitInactiveCond.Wait()
	e.waitInactiveLock.Unlock()
}

func (e *Executor[T]) WaitC(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		e.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
