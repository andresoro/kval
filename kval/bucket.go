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
	bsize int64
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
		bsize:      0,
	}

	return b

}

// set will return an error if key already exists
func (b *bucket) set(key string, val []byte) error {
	b.Lock()
	exists := b.cache[key]

	if exists != nil {
		b.Unlock()
		return errKeyExists
	}

	// init item and add to queue
	i := NewItem(key, val)
	heap.Push(&b.queue, i)
	b.bsize += i.Size()
	b.cache[key] = i

	b.Unlock()
	return nil

}

func (b *bucket) get(key string) (*Item, error) {
	b.RLock()

	i, found := b.cache[key]
	if !found {
		b.RUnlock()
		return nil, errKeyNotFound
	}
	b.queue.Access(i)

	b.RUnlock()
	return i, nil
}

func (b *bucket) delete(key string) (*Item, error) {
	b.Lock()

	i, found := b.cache[key]
	if !found {
		b.Unlock()
		return nil, errKeyNotFound
	}
	b.bsize -= i.Size()

	//delete from cache and queue
	delete(b.cache, key)
	heap.Remove(&b.queue, i.index)

	b.Unlock()
	return i, nil
}

func (b *bucket) clean() {
	b.Lock()

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
	b.Unlock()
	return
}

// delete all entries of the bucket
func (b *bucket) flush() {
	b.Lock()

	for key := range b.cache {
		delete(b.cache, key)
	}

	b.Unlock()
}

func (b *bucket) size() int64 {
	return b.bsize
}
