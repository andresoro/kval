package kval

import (
	"errors"
	"sync"
	"time"
)

// ErrKeyExists is returned when a key already exists in a Store
var ErrKeyExists = errors.New("Store: key already exists")

// ErrKeyNotFound is returned when a Store does not have a requested key
var ErrKeyNotFound = errors.New("Store: key not found in Store")

// Store is the in memory key value store that holds items for a max duration
type Store struct {
	lifeTime time.Duration
	cache    map[string]item
	mu       sync.RWMutex
}

// New returns a Store with a lifeTime of 5 minutes
func New() *Store {
	c := make(map[string]item)
	return &Store{
		lifeTime: 5 * time.Minute,
		cache:    c,
	}
}

// Set is a method to set a key-value pair in the cache
func (s *Store) Set(key string, val interface{}) {
	item := item{
		key: key,
		val: val,
	}

	s.mu.Lock()
	s.cache[key] = item
	s.mu.Unlock()

}

// Add is a method to add an object to the Store
// Add does not replace an item in the Store if the key already exists
func (s *Store) Add(key string, val interface{}) error {

	_, found := s.cache[key]
	if found {
		return ErrKeyExists
	}

	s.mu.Lock()
	s.Set(key, val)
	s.mu.Unlock()
	return nil

}

// Get is a method to return an item from the Store given a key
// Get should modify an items accessedAt field
func (s *Store) Get(key string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj, found := s.cache[key]
	if !found {
		return nil, ErrKeyNotFound
	}
	// remove from cache if the item has expired
	if time.Since(obj.accessedAt) > s.lifeTime {
		delete(s.cache, key)

		return nil, ErrKeyNotFound
	}
	obj.accessed()
	return obj.val, nil
}

func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.cache)
}

// item represents something to be cached in memory
type item struct {
	key        string
	val        interface{}
	createdAt  time.Time
	accessedAt time.Time
}

func newItem(key string, val interface{}) *item {
	t := time.Now()
	return &item{
		key:        key,
		val:        val,
		createdAt:  t,
		accessedAt: t,
	}
}

func (i *item) accessed() {
	i.accessedAt = time.Now()
}
