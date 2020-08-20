package scheduler

import (
	"container/list"
	"sync"
)

type WorkQueue struct {
	ql *list.List
	sync.Mutex
}

func NewWorkQueue() *WorkQueue {
	return &WorkQueue{
		ql: list.New(),
	}
}

func (q *WorkQueue) Add(item interface{}) {
	defer q.Unlock()
	q.Lock()
	q.ql.PushBack(item)
}

func (q *WorkQueue) Front() interface{} {
	return q.ql.Front().Value
}

func (q *WorkQueue) Len() int {
	return q.ql.Len()
}

func (q *WorkQueue) Remove() interface{} {
	defer q.Unlock()
	q.Lock()
	el := q.ql.Front()
	return q.ql.Remove(el)
}
