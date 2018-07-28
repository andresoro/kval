package kval

import (
	"container/heap"
	"sync"
)

// buckets will be used to free up lock contention over a cache.
// buckets are essentially thread safe maps acting as shards.
type bucket struct {
	cache map[string]*Item
	sync.RWMutex
	queue Queue
}

func newBucket() *bucket {
	c := make(map[string]*Item)
	q := make(Queue, 0)
	heap.Init(&q)

	b := &bucket{
		cache: c,
		queue: q,
	}

	return b

}

// set will return an error if key already exists
func (b *bucket) set(key string, val interface{}) error {
	b.Lock()
	defer b.Unlock()
	exists := b.cache[key]

	if exists != nil {
		return errKeyExists
	}

	// init item and add to queue
	i := newItem(key, val)
	heap.Push(&b.queue, i)

	b.cache[key] = i
	return nil

}

func (b *bucket) get(key string) (*Item, error) {
	b.RLock()
	defer b.RUnlock()
	i, found := b.cache[key]
	if !found {
		return nil, errKeyNotFound
	}
	b.queue.Access(i)
	return i, nil
}

func (b *bucket) delete(key string) (*Item, error) {
	b.Lock()
	defer b.Unlock()

	i, found := b.cache[key]
	if !found {
		return nil, errKeyNotFound
	}

	//delete from cache and queue
	delete(b.cache, key)
	heap.Remove(&b.queue, i.index)
	return i, nil
}
