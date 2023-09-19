package jarray

import (
	"sync"
)

type Queue struct {
	wg    sync.WaitGroup
	buf   chan struct{}
	tasks []func()
}

func NewQueue(size int) *Queue {
	if size <= 0 {
		return nil
	}
	return &Queue{
		buf: make(chan struct{}, size),
		wg:  sync.WaitGroup{},
	}
}

// AddTask to queue
func (q *Queue) AddTask(f func()) {
	q.tasks = append(q.tasks, f)
}

// Run queue task
func (q *Queue) Run() {
	for _, t := range q.tasks {
		q.wg.Add(1)
		q.buf <- struct{}{}
		go func(f func()) {
			defer q.wg.Done()
			f()
			<-q.buf
		}(t)
	}
	q.wg.Wait()
}
