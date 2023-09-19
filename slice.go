package jarray

import "sync"

type rangeFunc[T any] func(k int, v T)
type uniqueFunc[T any] func(v T) (int, int64)
type conditionFunc[T any] func(a T) bool

type StreamSlice[T any] struct {
	slice             []T
	handleFuncList    []interface{}
	distinctFuncCount int
	distinctMap       []map[int64]struct{}
	limit             int
}

// Range func
func (receiver *StreamSlice[T]) Range(f rangeFunc[T]) *StreamSlice[T] {
	receiver.handleFuncList = append(receiver.handleFuncList, f)
	return receiver
}

// Distinct use a unique key func to equal value
func (receiver *StreamSlice[T]) Distinct(f func(T) int64) *StreamSlice[T] {
	receiver.handleFuncList = append(receiver.handleFuncList,
		func(v T) (int, int64) {
			k := receiver.distinctFuncCount
			return k, f(v)
		})
	receiver.distinctFuncCount++
	return receiver
}

// Range use a condition func to filtrate value
func (receiver *StreamSlice[T]) Filter(f conditionFunc[T]) *StreamSlice[T] {
	receiver.handleFuncList = append(receiver.handleFuncList, f)
	return receiver
}

// Limit slice value length at the end of stream
func (receiver *StreamSlice[T]) Limit(limit int) *StreamSlice[T] {
	receiver.limit = limit
	return receiver
}

// Collect run stream node and return result
func (receiver *StreamSlice[T]) Collect() []T {
	result := make([]T, 0, len(receiver.slice))
	if receiver.distinctFuncCount > 0 {
		receiver.distinctMap = make([]map[int64]struct{}, receiver.distinctFuncCount)
		for i := 0; i < receiver.distinctFuncCount; i++ {
			receiver.distinctMap[i] = make(map[int64]struct{})
		}
	}
	for k, v := range receiver.slice {
		invalid := false
		for _, f := range receiver.handleFuncList {
			switch f := f.(type) {
			case rangeFunc[T]:
				f(k, v)
			case uniqueFunc[T]:
				i, u := f(v)
				_, ok := receiver.distinctMap[i][u]
				if ok {
					invalid = true
					break
				}
				receiver.distinctMap[i][u] = struct{}{}
			case conditionFunc[T]:
				if !f(v) {
					invalid = true
					break
				}
			}
		}
		if receiver.limit != 0 &&
			len(result) >= receiver.limit {
			break
		}
		if !invalid {
			result = append(result, v)
		}
	}
	return result
}

func Stream[T any](slice []T) *StreamSlice[T] {
	return &StreamSlice[T]{
		slice: slice,
	}
}

type Slice[T any] []T

func RangeSlice[T any](slice []T) Slice[T] {
	return slice
}

// Range func
func (s Slice[T]) Range(f func(k int, v T)) Slice[T] {
	for i, j := range s {
		f(i, j)
	}
	return s
}

// Range func with goroutine
func (s Slice[T]) ConcurrentRange(f func(k int, v T)) Slice[T] {
	wg := &sync.WaitGroup{}
	for i, j := range s {
		wg.Add(1)
		go func(m int, n T) {
			defer wg.Done()
			f(m, n)
		}(i, j)
	}
	wg.Wait()
	return s
}

// Map a slice type form T to B
func Map[T, B any](slice []T, f func(o T) B) []B {
	result := make([]B, len(slice))
	RangeSlice(slice).ConcurrentRange(
		func(k int, v T) {
			result[k] = f(v)
		})
	return result
}

// Flat slices to slice
func Flat[T any](slice [][]T) []T {
	result := make([]T, 0, 8)
	RangeSlice(slice).Range(func(k int, v []T) {
		result = append(result, v...)
	})
	return result
}
