package kval

import (
	"container/heap"
	"testing"
)

func TestQueue(t *testing.T) {
	itemA := newItem("key", "val")
	itemB := newItem("key2", "val")
	itemC := newItem("key3", "val")

	q := Queue{
		itemA,
		itemB,
		itemC,
	}

	heap.Init(&q)

	// itemA should be last in the queue
	if heap.Pop(&q) != itemA {
		t.Error("itemC should be the last item")
	}

}
