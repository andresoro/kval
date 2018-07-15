package kval

// Queue represents a FIFO queue used to evict cache values
type Queue []*item

// Push inserts a value into the tail end of Queue
func (q *Queue) Push(i *item) {
	*q = append(*q, i)
}

// Pop returns the last item in the queue
func (q *Queue) Pop() (i *item) {
	n := (*q)[0]
	*q = (*q)[1:]
	return n
}
