package kval

import (
	"sync"
)

// buckets will be used to free up lock contention over a cache.
// buckets are essentially thread safe maps acting as shards.
type bucket struct {
	cache map[string]*Item
	sync.RWMutex
}

func (b *bucket) set(key string, val interface{}) error {
	b.Lock()
	defer b.Unlock()
	exists := b.cache[key]

	if exists != nil {
		return errKeyExists
	}

	i := newItem(key, val)

	b.cache[key] = i
	return nil

}

func (b *bucket) get(key string) *Item {
	b.RLock()
	defer b.RUnlock()
	return b.cache[key]
}

func (b *bucket) delete(key string) *Item {
	b.Lock()
	defer b.Unlock()
	i := b.cache[key]
	delete(b.cache, key)
	return i
}
