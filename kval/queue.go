package kval

import (
	"container/heap"
	"time"
)

// Queue represents a queue used to evict items
// satifies heap interface
type Queue []*Item

// Len returns length of queue
func (q Queue) Len() int { return len(q) }

// Less tests whether item at index i is less than j
// lower time since access implies higher priority
func (q Queue) Less(i, j int) bool {

	return q[i].accessedAt.Before(q[j].accessedAt)
}

// Pop implements heap interface
func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	i := old[n-1]
	i.index = -1
	*q = old[0 : n-1]
	return i
}

// Push implements heap interface
func (q *Queue) Push(x interface{}) {
	tmp := *q
	n := len(tmp)
	i := x.(*Item)
	i.index = n
	tmp = append(tmp, i)
	*q = tmp
}

// Swap implements heap interface
func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

// Peek returns the item at the front of the queue
// must run heap.Init(&q)for Peek to work as intended
func (q *Queue) Peek() *Item {
	if len(*q) != 0 {
		return (*q)[0]
	}

	return nil
}

// Access moves the item to the front of the queue and changes item's
// time accessed field
func (q *Queue) Access(i *Item) {
	heap.Remove(q, i.index)
	i.accessedAt = time.Now()
	heap.Push(q, i)
}
