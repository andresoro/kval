package kval

import (
	"hash/fnv"
	"sync"
	"time"
)

// BStore is an in-memory key-value store that uses a max life span
// for items.
type BStore struct {
	cache []*bucket
	sync.RWMutex
	queue    Queue
	frozen   bool
	lifeTime time.Duration
}

// NewBStore returns a new bucket store
func NewBStore() *BStore {
	// init store

	b := &Store{
		cache:    nil,
		frozen:   false,
		lifeTime: 5 * time.Minute,
	}
}

// *TODO* add method to pick/add shards

func (s *BStore) getBucket(key string) *bucket {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	mask := uint32(len(s.cache) - 1)
	return s.cache[hasher.Sum32()&mask]
}
