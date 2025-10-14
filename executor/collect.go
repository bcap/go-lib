package executor

import (
	"context"
	"reflect"

	"github.com/bcap/go-lib/result"
)

// Collect applies a function to all elements in the given slice or array, returning a slice with the results in the same order of the inputs
//
// The input to the function is the index in the slice. The usage is similar to the sort.Slice function. Example:
//
//	inputs := []int{100, 200, 300}
//	results := Collect(ctx, 0, inputs, func(i int) (int, error) {
//	   return inputs[i] + 10
//	})
//	reflect.DeepEqual(results, []int{110, 210, 310}) // true
func Collect[T any](ctx context.Context, maxParallelism int, slice any, fn func(int) (T, error)) result.Results[T] {
	return CollectE(ctx, New[T](maxParallelism), slice, fn)
}

// CollectE is the same as Collect, but you can pass an existing Executor
func CollectE[T any](ctx context.Context, e *Executor[T], slice any, fn func(int) (T, error)) result.Results[T] {
	if slice == nil {
		return result.Results[T]{}
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
//
// The function is called with each key, and the result is stored in the map.
//
// Example:
//
//	keys := []string{"a", "b", "c"}
//	results := CollectMap(ctx, 0, keys, func(key string) (string, error) {
//		return key + "!", nil
//	})
//	reflect.DeepEqual(results, map[string]string{"a": "a!", "b": "b!", "c": "c!"}) // true
func CollectMap[K comparable, T any](ctx context.Context, maxParallelism int, entries []K, fn func(K) (T, error)) result.ResultsMap[K, T] {
	return CollectMapE(ctx, New[T](maxParallelism), entries, fn)
}

// CollectMapE is the same as CollectMap, but you can pass an existing Executor instance
func CollectMapE[K comparable, T any](ctx context.Context, e *Executor[T], entries []K, fn func(K) (T, error)) result.ResultsMap[K, T] {
	results := map[K]*result.Result[T]{}
	for _, key := range entries {
		results[key] = nil
	}
	CollectMapReplaceE(ctx, e, results, fn)
	return results
}

// CollectMapReplace applies a function to all the keys in the given map, replacing the values with the results
//
// The function is called with each key, and the result is stored back in the passed map.
//
// Example:
//
//	m := map[string]*Result[string]{"a": nil, "b": nil, "c": nil}
//	CollectMapReplace(ctx, 0, m, func(key string) (string, error) {
//		return key + "!", nil
//	})
//	reflect.DeepEqual(m, map[string]*Result[string]{"a": "a!", "b": "b!", "c": "c!"}) // true
func CollectMapReplace[K comparable, T any](ctx context.Context, maxParallelism int, m map[K]*result.Result[T], fn func(K) (T, error)) {
	CollectMapReplaceE(ctx, New[T](maxParallelism), m, fn)
}

// CollectMapReplaceE is the same as CollectMapReplace, but you can pass an existing Executor instance
func CollectMapReplaceE[K comparable, T any](ctx context.Context, e *Executor[T], m map[K]*result.Result[T], fn func(K) (T, error)) {
	futures := map[K]*Future[T]{}
	for key := range m {
		key := key
		futures[key] = e.Submit(ctx, func() (T, error) { return fn(key) })
	}
	for key, future := range futures {
		m[key] = future.Get(ctx)
	}
}

// CollectFutures is a helper function to collect results from a slice of Futures, returning a slice of Results
func CollectFutures[T any](ctx context.Context, futures []*Future[T]) result.Results[T] {
	results := make([]*result.Result[T], len(futures))
	for i, f := range futures {
		results[i] = f.Get(ctx)
	}
	return results
}

// CollectFuturesMap is a helper function to collect results from a map of futures, returning a map of Results.
// Returned map keys are the same as the input map
func CollectFuturesMap[K comparable, T any](ctx context.Context, futures map[K]*Future[T]) result.ResultsMap[K, T] {
	results := map[K]*result.Result[T]{}
	for key, future := range futures {
		results[key] = future.Get(ctx)
	}
	return results
}
