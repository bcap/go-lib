package pool

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type myObject struct {
	resetCalls int
}

func (m *myObject) Reset() {
	m.resetCalls++
}

func TestPool(t *testing.T) {
	var newCalls int
	newFn := func() *myObject {
		newCalls++
		return &myObject{}
	}
	p := New[*myObject](newFn, WithMaxSize[*myObject](3))
	seen := map[*myObject]struct{}{}

	assert.Equal(t, 0, newCalls, "Expected no new objects to be created yet")

	obj1 := p.MustGet()
	assert.NotNil(t, obj1, "Expected to get an object from the pool")
	assert.Equal(t, 1, newCalls, "Expected one new object to be created when calling MustGet")
	assert.Equal(t, obj1.resetCalls, 0, "Expected the reset method to be not called yet")
	assert.Equal(t, p.lent.Load(), int64(1), "Expected the lent count to be zero after returning all objects to the pool")
	seen[obj1] = struct{}{}

	obj2 := p.MustGet()
	assert.NotNil(t, obj2, "Expected to get an object from the pool")
	assert.Equal(t, 2, newCalls, "Expected one new object to be created when calling MustGet")
	assert.NotSame(t, obj1, obj2, "Expected to get a different object from the pool")
	assert.Equal(t, p.lent.Load(), int64(2), "Expected the lent count to be zero after returning all objects to the pool")
	seen[obj2] = struct{}{}

	p.Return(obj1)
	assert.Equal(t, p.lent.Load(), int64(1), "Expected the lent count to be zero after returning all objects to the pool")
	assert.Equal(t, obj1.resetCalls, 1, "Expected the reset method to be called on the object when returned to the pool")

	obj3 := p.MustGet()
	assert.NotNil(t, obj3, "Expected to get an object from the pool")
	assert.Equal(t, 2, newCalls, "No new calls to NewFn expected")
	assert.Equal(t, obj3, obj1, "Expected to get the same object back from the pool after returning it")
	assert.Equal(t, p.lent.Load(), int64(2), "Expected the lent count to be zero after returning all objects to the pool")

	obj4 := p.MustGet()
	assert.NotNil(t, obj4, "Expected to get an object from the pool")
	assert.Equal(t, 3, newCalls)
	assert.Equal(t, p.lent.Load(), int64(3), "Expected the lent count to be zero after returning all objects to the pool")

	// At this point the pool is full. Any further Get calls will block

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	obj5, err := p.Get(ctx)
	assert.Nil(t, obj5)
	assert.NotNil(t, err)

	p.Return(obj2)
	assert.Equal(t, p.lent.Load(), int64(2), "Expected the lent count to be zero after returning all objects to the pool")
	assert.Equal(t, obj2.resetCalls, 1, "Expected the reset method to be called on the second object when returned to the pool")

	p.Return(obj3)
	assert.Equal(t, p.lent.Load(), int64(1), "Expected the lent count to be zero after returning all objects to the pool")
	assert.Equal(t, obj3.resetCalls, 2, "Expected the reset method to be called on the first object again when returned to the pool")

	p.Return(obj4)
	assert.Equal(t, p.lent.Load(), int64(0), "Expected the lent count to be zero after returning all objects to the pool")
	assert.Equal(t, obj4.resetCalls, 1, "Expected the reset method to be called on the first object again when returned to the pool")
}
