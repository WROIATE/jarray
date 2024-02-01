package jarray

type rangeFunc[T any] func(k int, v T)
type uniqueFunc[T any] func(v T) (int, int64)
type conditionFunc[T any] func(a T) bool

// NewStream function is used to create a new StreamFlow.
// The parameter slice is a generic slice that represents the data source for streaming processing.
func NewStream[T any](slice []T) StreamFlow[T, T] {
	return &streamSlice[T, T]{
		slice:             slice,
		handleFuncList:    make([]interface{}, 0),
		distinctFuncCount: 0,
		distinctMap:       make([]map[int64]struct{}, 0),
		limit:             0,
		mapFunc: func(t T) T {
			return t
		},
	}
}

// NewMapStream function is used to create a new StreamStart
// allow to use mapFunc to transfer struct type in collect step
func NewMapStream[T, B any](slice []T) StreamStart[T, B] {
	return &streamSlice[T, B]{
		slice:             slice,
		handleFuncList:    make([]interface{}, 0),
		distinctFuncCount: 0,
		distinctMap:       make([]map[int64]struct{}, 0),
		limit:             0,
		mapFunc:           nil,
	}
}

type Stream[T, B any] interface {
	StreamStart[T, B]
	StreamFlow[T, B]
	StreamEnd[B]
}

type StreamStart[T, B any] interface {
	Map(func(T) B) StreamFlow[T, B]
}

type StreamFlow[T, B any] interface {
	Range(rangeFunc[T]) StreamFlow[T, B]
	Distinct(f func(T) int64) StreamFlow[T, B]
	Filter(f conditionFunc[T]) StreamFlow[T, B]
	Limit(limit int) StreamEnd[B]
}

type StreamEnd[T any] interface {
	Collect() []T
}

type streamSlice[T, B any] struct {
	slice             []T
	handleFuncList    []interface{}
	distinctFuncCount int
	distinctMap       []map[int64]struct{}
	limit             int
	mapFunc           func(T) B
}

// Map will transfer struct type at the end of collect
func (s *streamSlice[T, B]) Map(f func(T) B) StreamFlow[T, B] {
	s.mapFunc = f
	return s
}

// Range func
func (s *streamSlice[T, B]) Range(f rangeFunc[T]) StreamFlow[T, B] {
	s.handleFuncList = append(s.handleFuncList, f)
	return s
}

// Distinct use a unique key func to equal value
func (s *streamSlice[T, B]) Distinct(f func(T) int64) StreamFlow[T, B] {
	s.handleFuncList = append(s.handleFuncList,
		func(v T) (int, int64) {
			k := s.distinctFuncCount
			return k, f(v)
		})
	s.distinctFuncCount++
	return s
}

// Filter use a condition func to filtrate value
func (s *streamSlice[T, B]) Filter(f conditionFunc[T]) StreamFlow[T, B] {
	s.handleFuncList = append(s.handleFuncList, f)
	return s
}

// Limit slice value length at the end of StreamFlow
func (s *streamSlice[T, B]) Limit(limit int) StreamEnd[B] {
	s.limit = limit
	return s
}

// Collect run StreamFlow node and return result
func (s *streamSlice[T, B]) Collect() []B {
	result := make([]B, 0, len(s.slice))
	if s.distinctFuncCount > 0 {
		s.distinctMap = make([]map[int64]struct{}, s.distinctFuncCount)
		for i := 0; i < s.distinctFuncCount; i++ {
			s.distinctMap[i] = make(map[int64]struct{})
		}
	}
	for k, v := range s.slice {
		invalid := false
		for _, f := range s.handleFuncList {
			switch f := f.(type) {
			case rangeFunc[T]:
				f(k, v)
			case uniqueFunc[T]:
				i, u := f(v)
				_, ok := s.distinctMap[i][u]
				if ok {
					invalid = true
					break
				}
				s.distinctMap[i][u] = struct{}{}
			case conditionFunc[T]:
				if !f(v) {
					invalid = true
					break
				}
			}
		}
		if s.limit != 0 &&
			len(result) >= s.limit {
			break
		}
		if !invalid && s.mapFunc != nil {
			result = append(result, s.mapFunc(v))
		}
	}
	return result
}
