package kval

import (
	"container/heap"
	"testing"
)

func TestQueue(t *testing.T) {

	itemA := newItem("key", "val")
	itemB := newItem("key2", "val")
	itemC := newItem("key3", "val")

	list := []*Item{
		itemA,
		itemB,
		itemC,
	}

	q := make(Queue, 0)

	heap.Init(&q)

	for _, item := range list {
		heap.Push(&q, item)
	}

	q.Access(itemA) // push A to the front

	if q.Peek() != itemB {
		t.Error("Error: Peek function should return B since it is now at front")
	}

	if heap.Pop(&q) != itemB {
		t.Error("Error: after A has been accessed, B should be at the front of queue")
	}

}
