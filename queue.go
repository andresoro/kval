package kval

import (
	"sync"
)

// Queue represents a FIFO queue used to evict cache values
type Queue struct {
	nodes []*item
	mu    sync.Mutex
}

// Push inserts a value into the tail end of Queue
func (q *Queue) Push(i *item) {
	q.nodes = append(q.nodes, i)
}
