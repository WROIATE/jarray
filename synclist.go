package jarray

import (
	"sync"
)

type SyncList[T any] struct {
	*SimpleList[T]
	lock *sync.RWMutex
}

// NewSyncList make a new sync list
func NewSyncList[T any]() *SyncList[T] {
	return &SyncList[T]{
		SimpleList: NewSimpleList[T](),
		lock:       &sync.RWMutex{},
	}
}

// ToSyncList package slice to sync list
func ToSyncList[T any](slice []T) *SyncList[T] {
	copyLock := &sync.Mutex{}
	copyLock.Lock()
	defer copyLock.Unlock()
	return &SyncList[T]{
		SimpleList: ToSimpleList(slice),
		lock:       &sync.RWMutex{},
	}
}

// Add a value into list end
func (s *SyncList[T]) Add(value T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.SimpleList.Add(value)
}

// Append slice into list end
func (s *SyncList[T]) Append(value []T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.SimpleList.Append(value)
}

// Remove index value
func (s *SyncList[T]) Remove(index int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.SimpleList.Remove(index)
}

// Get index value
func (s *SyncList[T]) Get(index int) T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.SimpleList.Get(index)
}

// Set index value
func (s *SyncList[T]) Set(index int, value T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.SimpleList.Set(index, value)
}

// ToSlice return a new slice with list value
func (s *SyncList[T]) ToSlice() []T {
	defer s.lock.RUnlock()
	return s.SimpleList.ToSlice()
}

// Size return list length
func (s *SyncList[T]) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.SimpleList.Size()
}

func (s *SyncList[T]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.SimpleList.Clear()
}

func (s *SyncList[T]) Contains(value T, equals func(a T, b T) bool) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.SimpleList.Contains(value, equals)
}

func (s *SyncList[T]) Sort(less func(i, j int, sortSlice []T) bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.SimpleList.Sort(less)
}
