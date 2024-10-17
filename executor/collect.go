package executor

import (
	"context"
	"reflect"
)

// Collect applies a function to all elements in the given slice or array, returning a slice with the results in the same order of the inputs
func Collect[T any](ctx context.Context, maxParallelism int, slice any, fn func(int) (T, error)) Results[T] {
	return CollectE(ctx, New[T](maxParallelism), slice, fn)
}

// CollectE is the same as Collect, but you can pass an existing Executor. This is the same as calling Executor.Collect
func CollectE[T any](ctx context.Context, e *Executor[T], slice any, fn func(int) (T, error)) Results[T] {
	if slice == nil {
		return Results[T]{}
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
	return CollectFutures(ctx, futures)
}

// CollectMap applies a function to all given keys, returning a map of key -> results
func CollectMap[K comparable, T any](ctx context.Context, maxParallelism int, entries []K, fn func(K) (T, error)) ResultsMap[K, T] {
	return CollectMapE(ctx, New[T](maxParallelism), entries, fn)
}

// CollectMapE is the same as CollectMap, but you can pass an existing Executor instance
func CollectMapE[K comparable, T any](ctx context.Context, e *Executor[T], entries []K, fn func(K) (T, error)) ResultsMap[K, T] {
	results := map[K]*Result[T]{}
	for _, key := range entries {
		results[key] = nil
	}
	CollectMapReplaceE(ctx, e, results, fn)
	return results
}

// CollectMapReplace applies a function to all the keys in the given map, replacing the values with the results
func CollectMapReplace[K comparable, T any](ctx context.Context, maxParallelism int, m map[K]*Result[T], fn func(K) (T, error)) {
	CollectMapReplaceE(ctx, New[T](maxParallelism), m, fn)
}

// CollectMapReplaceE is the same as CollectMapReplace, but you can pass an existing Executor instance
func CollectMapReplaceE[K comparable, T any](ctx context.Context, e *Executor[T], m map[K]*Result[T], fn func(K) (T, error)) {
	futures := map[K]*Future[T]{}
	for key := range m {
		key := key
		futures[key] = e.Submit(ctx, func() (T, error) { return fn(key) })
	}
	for key, future := range futures {
		m[key] = future.Get(ctx)
	}
}

// CollectFutures is a helper function to collect results from a slice of futures
func CollectFutures[T any](ctx context.Context, futures []*Future[T]) Results[T] {
	results := make([]*Result[T], len(futures))
	for i, f := range futures {
		results[i] = f.Get(ctx)
	}
	return results
}

// CollectFuturesMap is a helper function to collect results from a map of futures
func CollectFuturesMap[K comparable, T any](ctx context.Context, futures map[K]*Future[T]) ResultsMap[K, T] {
	results := map[K]*Result[T]{}
	for key, future := range futures {
		results[key] = future.Get(ctx)
	}
	return results
}
