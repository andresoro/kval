package kval

import (
	"hash/fnv"
	"sync"
	"time"
)

var (
	// placeholder till config is written
	bucketNum = 1
	lifeTime  = 5 * time.Millisecond
)

// BStore is an in-memory key-value store that uses a max life span
// for items.
type BStore struct {
	cache []*bucket
	sync.RWMutex
	frozen   bool
	lifeTime time.Duration
}

// New returns a new bucket store
func New() *BStore {

	b := &BStore{
		cache:    make([]*bucket, bucketNum),
		frozen:   false,
		lifeTime: lifeTime,
	}

	for i := 0; i < bucketNum; i++ {
		b.cache[i] = newBucket(lifeTime)
	}

	go b.janitor()

	return b
}

// Get returns the value of the item with given key
func (b *BStore) Get(key string) (interface{}, error) {
	bucket := b.pickBucket(key)
	item, err := bucket.get(key)
	if err != nil {
		return nil, err
	}
	return item.val, nil

}

// Add method adds a key/val pair to the store and returns an error
// if key already exists
func (b *BStore) Add(key string, val interface{}) error {
	if b.frozen {
		return errStoreIsFrozen
	}

	bucket := b.pickBucket(key)
	err := bucket.set(key, val)
	if err != nil {
		return err
	}

	return nil
}

// Delete method deletes and returns an item with given key from cache
// if item does not exist return an error
func (b *BStore) Delete(key string) (interface{}, error) {
	bucket := b.pickBucket(key)
	i, err := bucket.delete(key)
	if err != nil {
		return nil, err
	}
	return i.val, nil

}

// Freeze a store
func (b *BStore) Freeze() {
	b.frozen = true
}

// Unfreeze a store
func (b *BStore) Unfreeze() {
	b.frozen = false
}

func (b *BStore) clean() {
	var wg sync.WaitGroup

	n := len(b.cache)

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			b.cache[i].clean()
		}(i)
	}

	wg.Wait()
}

func (b *BStore) janitor() {
	for {
		select {
		case <-time.After(b.lifeTime):
			b.clean()
		}
	}
}

// pickBucket is a function to "assign" a key to a bucket
// sharding function is a simple hash(key) % n
// where n is number of buckets
func (b *BStore) pickBucket(key string) *bucket {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	mask := uint32(len(b.cache) - 1)
	return b.cache[hasher.Sum32()&mask]
}
