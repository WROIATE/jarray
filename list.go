package jarray

import (
	"encoding/json"
	"sync"
)

type SyncList[T any] struct {
	slice []T
	lock  *sync.RWMutex
}

// NewSyncList make a new sync list
func NewSyncList[T any]() *SyncList[T] {
	return &SyncList[T]{
		slice: make([]T, 0, 8),
		lock:  &sync.RWMutex{},
	}
}

// ToSyncList package slice to sync list
func ToSyncList[T, B any](slice []B) *SyncList[T] {
	list := make([]T, len(slice))
	copyLock := &sync.Mutex{}
	copyLock.Lock()
	defer copyLock.Unlock()
	for i, v := range slice {
		switch any(v).(type) {
		case T:
			list[i] = any(v).(T)
		case string:
			d, err := unmarshalJson[T]([]byte(any(v).(string)))
			if err != nil {
				return nil
			}
			list[i] = d
		default:
			bytes, err := json.Marshal(v)
			if err != nil {
				return nil
			}
			d, err := unmarshalJson[T](bytes)
			if err != nil {
				return nil
			}
			list[i] = d
		}
	}
	return &SyncList[T]{
		slice: list,
		lock:  &sync.RWMutex{},
	}
}

func unmarshalJson[T any](b []byte) (v T, err error) {
	return v, json.Unmarshal(b, &v)
}

// Add a value into list end
func (s *SyncList[T]) Add(value T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.slice = append(s.slice, value)
}

// Append slice into list end
func (s *SyncList[T]) Append(value []T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.slice = append(s.slice, value...)
}

// Remove index value
func (s *SyncList[T]) Remove(index int) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.slice = append(s.slice[0:index], s.slice[index+1:]...)
	return nil
}

// Get index value
func (s *SyncList[T]) Get(index int) (T, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.slice[index], nil
}

// Set index value
func (s *SyncList[T]) Set(value T, index int) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.slice[index] = value
	return nil
}

// ToSlice return a new slice with list value
func (s *SyncList[T]) ToSlice() []T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	slice := make([]T, len(s.slice))
	copy(slice, s.slice)
	return slice
}

// Size return list length
func (s *SyncList[T]) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.slice)
}
