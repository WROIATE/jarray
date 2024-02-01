package jarray

import "sort"

type List[T any] interface {
	Add(value T)
	Remove(index int)
	Get(index int) T
	Set(index int, value T)
	Append(value []T)
	ToSlice() []T
	Size() int
	Clear()
	Contains(value T, equals func(a, b T) bool) bool
	Sort(less func(i, j int, sortSlice []T) bool)
}

type SimpleList[T any] struct {
	slice    []T
	lessFunc func(i, j int, sortSlice []T) bool
}

func (s *SimpleList[T]) Len() int {
	return s.Size()
}

func (s *SimpleList[T]) Less(i, j int) bool {
	if s.lessFunc == nil {
		return true
	}
	return s.lessFunc(i, j, s.slice)
}

func (s *SimpleList[T]) Swap(i, j int) {
	s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
}

func NewSimpleList[T any]() *SimpleList[T] {
	return &SimpleList[T]{
		slice: make([]T, 0, 8),
	}
}

func ToSimpleList[T any](slice []T) *SimpleList[T] {
	if len(slice) == 0 {
		return NewSimpleList[T]()
	}
	copySlice := make([]T, len(slice), len(slice)+8)
	copy(copySlice, slice)
	return &SimpleList[T]{
		slice: copySlice,
	}
}

func (s *SimpleList[T]) Add(value T) {
	s.slice = append(s.slice, value)
}

func (s *SimpleList[T]) Remove(index int) {
	s.slice = append(s.slice[0:index], s.slice[index+1:]...)
}

func (s *SimpleList[T]) Get(index int) T {
	return s.slice[index]
}

func (s *SimpleList[T]) Set(index int, value T) {
	s.slice[index] = value
}

func (s *SimpleList[T]) Append(value []T) {
	s.slice = append(s.slice, value...)
}

func (s *SimpleList[T]) Size() int {
	return len(s.slice)
}

func (s *SimpleList[T]) ToSlice() []T {
	slice := make([]T, len(s.slice))
	copy(slice, s.slice)
	return slice
}

func (s *SimpleList[T]) Clear() {
	s.slice = make([]T, 0, 8)
}

func (s *SimpleList[T]) Contains(value T, equals func(a, b T) bool) bool {
	for _, v := range s.slice {
		if equals(v, value) {
			return true
		}
	}
	return false
}

func (s *SimpleList[T]) Sort(less func(i, j int, sortSlice []T) bool) {
	if less == nil {
		return
	}
	s.lessFunc = less
	sort.Sort(s)
}
