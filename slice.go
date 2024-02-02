package jarray

import "sync"

type Slice[T any] []T

func EmptySlice[T any]() Slice[T] {
	return make(Slice[T], 0)
}

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

// ConcurrentRange func with goroutine
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
func Map[T, B any](slice []T, f func(index int, value T) B) []B {
	result := make([]B, len(slice))
	RangeSlice(slice).ConcurrentRange(
		func(k int, v T) {
			result[k] = f(k, v)
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
