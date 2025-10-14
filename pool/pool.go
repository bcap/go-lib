package pool

import (
	"context"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type Poolable interface {
	Reset()
}

type Pool[P Poolable] struct {
	pool    *sync.Pool
	sem     *semaphore.Weighted
	lent    atomic.Int64
	size    atomic.Int64
	maxSize int64
}

type Option[P Poolable] func(*Pool[P])

func WithMinSize[P Poolable](minSize int64) Option[P] {
	return func(p *Pool[P]) {
		for i := int64(0); i < minSize; i++ {
			p.Return(p.MustGet())
		}
	}
}

func WithMaxSize[P Poolable](maxSize int64) Option[P] {
	return func(p *Pool[P]) {
		p.maxSize = maxSize
		p.sem = semaphore.NewWeighted(maxSize)
	}
}

func New[P Poolable](newFn func() P, opts ...Option[P]) *Pool[P] {
	var p Pool[P]
	pool := sync.Pool{New: func() any {
		defer p.size.Add(1)
		return newFn()
	}}
	p.pool = &pool
	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

func (p *Pool[P]) Get(ctx context.Context) (P, error) {
	if p.sem != nil {
		if err := p.sem.Acquire(ctx, 1); err != nil {
			var zeroVal P
			return zeroVal, err
		}
	}
	obj := p.pool.Get().(P)
	p.lent.Add(1)
	return obj, nil
}

func (p *Pool[P]) MustGet() P {
	obj, err := p.Get(context.Background())
	if err != nil {
		panic(err)
	}
	return obj
}

func (p *Pool[T]) Return(t T) {
	t.Reset()
	p.pool.Put(t)
	p.lent.Add(-1)
	if p.sem != nil {
		p.sem.Release(1)
	}
}
