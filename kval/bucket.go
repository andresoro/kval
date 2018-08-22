package kval

import (
	"container/heap"
	"sync"
	"time"
)

// buckets will be used to free up lock contention over a cache.
// buckets are essentially thread safe maps acting as shards.
type bucket struct {
	cache map[string]*Item
	sync.RWMutex
	queue      Queue
	timeToLive time.Duration
}

func newBucket(ttl time.Duration) *bucket {
	c := make(map[string]*Item)
	q := make(Queue, 0)
	heap.Init(&q)

	b := &bucket{
		cache:      c,
		queue:      q,
		timeToLive: ttl,
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
	i := NewItem(key, val)
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

func (b *bucket) clean() {
	b.Lock()
	defer b.Unlock()

	clean := false

	for !clean {
		item := b.queue.Peek()
		if item != nil {
			if time.Since(item.accessedAt) > b.timeToLive {
				heap.Pop(&b.queue)
				delete(b.cache, item.key)
			} else {
				clean = true
			}
		} else {
			clean = true
		}
	}

	return
}
