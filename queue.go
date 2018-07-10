package kval

// Queue represents a FIFO queue used to evict cache values
type Queue struct {
	nodes []*item
}
