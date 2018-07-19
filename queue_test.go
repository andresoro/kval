package kval

import (
	"container/heap"
	"testing"
)

func TestQueue(t *testing.T) {

	itemA := newItem("key", "val")
	itemB := newItem("key2", "val")
	itemC := newItem("key3", "val")

	list := []*item{
		itemA,
		itemB,
		itemC,
	}

	q := make(Queue, len(list))

	for i, item := range list {
		q[i] = item
		q[i].index = i
	}

	heap.Init(&q)

	q.Access(itemA) // push A to the front

	if heap.Pop(&q) != itemB {
		t.Error("Error: after A has been accessed, B should be at the end of queue")
	}

}
