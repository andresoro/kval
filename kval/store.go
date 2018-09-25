package kval

import (
	"errors"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errStoreIsFrozen = errors.New("Error: Store is frozen")
	errKeyExists     = errors.New("Error: Key already exists in store")
	errKeyNotFound   = errors.New("Error: Key not found in store")
)

// Store is an in-memory key-value store that uses a max life span
// for items.
type Store struct {
	cache     []*bucket
	frozen    bool
	cacheSize uint64
	lifeTime  time.Duration
}

// New returns a new bucket store
func New(shardNum int, timeToLive time.Duration) (*Store, error) {

	if !powerOfTwo(shardNum) {
		return nil, errors.New("Number of shards should be power of two")
	}

	s := &Store{
		cache:     make([]*bucket, shardNum),
		frozen:    false,
		lifeTime:  timeToLive,
		cacheSize: 0,
	}

	for i := 0; i < shardNum; i++ {
		s.cache[i] = newBucket(timeToLive)
	}

	go s.janitor()

	return s, nil
}

// Get returns the value of the item with given key
func (s *Store) Get(key string) (interface{}, error) {
	bucket := s.pickBucket(key)
	item, err := bucket.get(key)
	if err != nil {
		return nil, err
	}
	return item.val, nil

}

// Add method adds a key/val pair to the store and returns an error
// if key already exists
func (s *Store) Add(key string, val interface{}) error {
	if s.frozen {
		return errStoreIsFrozen
	}

	atomic.AddUint64(&s.cacheSize, 1)

	bucket := s.pickBucket(key)
	err := bucket.set(key, val)
	if err != nil {
		return err
	}

	return nil
}

// Delete method deletes and returns an item with given key from cache
// if item does not exist return an error
func (s *Store) Delete(key string) (interface{}, error) {
	bucket := s.pickBucket(key)
	i, err := bucket.delete(key)
	if err != nil {
		return nil, err
	}
	atomic.AddUint64(&s.cacheSize, ^uint64(0))
	return i.val, nil

}

// Flush method will delete every key from every bucket
func (s *Store) Flush() {

	for _, bucket := range s.cache {
		bucket.flush()
	}

}

// Size returns number of entries in cache
func (s *Store) Size() uint64 {
	return atomic.LoadUint64(&s.cacheSize)
}

// Freeze a store
func (s *Store) Freeze() {
	s.frozen = true
}

// Unfreeze a store
func (s *Store) Unfreeze() {
	s.frozen = false
}

func (s *Store) clean() {
	var wg sync.WaitGroup

	n := len(s.cache)

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			s.cache[i].clean()
		}(i)
	}

	wg.Wait()
}

func (s *Store) janitor() {
	for {
		select {
		case <-time.After(s.lifeTime):
			s.clean()
		}
	}
}

// pickBucket is a function to "assign" a key to a bucket
// sharding function is a simple hash(key) % n
// where n is number of buckets
func (s *Store) pickBucket(key string) *bucket {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	mask := uint32(len(s.cache) - 1)
	return s.cache[hasher.Sum32()&mask]
}

func powerOfTwo(i int) bool {
	return (i & (i - 1)) == 0
}
