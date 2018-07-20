package kval

import (
	"sync"
	"time"
)

type bucket struct {
	lifetime time.Duration
	cache    map[string]*Item
	sync.RWMutex
}

func (b *bucket) set(key string, item *Item) {
	b.Lock()
	defer b.Unlock()
	b.cache[key] = item
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
}
