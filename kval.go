package kval

import (
	"errors"
	"sync"
	"time"
)

// ErrKeyExists is returned when a key already exists in a Store
var errKeyExists = errors.New("Store: key already exists")

// ErrKeyNotFound is returned when a Store does not have a requested key
var errKeyNotFound = errors.New("Store: key not found in Store")

// Store is the in memory key value store that holds items for a max duration
type Store struct {
	lifeTime time.Duration
	cache    map[string]item
	mu       sync.RWMutex
	frozen   bool
}

// New returns a Store with a lifeTime of 5 minutes
func New() *Store {
	c := make(map[string]item)
	return &Store{
		lifeTime: 5 * time.Minute,
		cache:    c,
		frozen:   false,
	}
}

// Set is a method to set a key-value pair in the cache
func (s *Store) set(key string, val interface{}) {
	item := newItem(key, val)

	s.mu.RLock()
	s.cache[key] = *item
	s.mu.RUnlock()

}

// Add is a method to add an object to the Store
// Add does not replace an item in the Store if the key already exists
func (s *Store) Add(key string, val interface{}) error {

	_, found := s.cache[key]
	if found {
		return errKeyExists
	}

	s.set(key, val)
	return nil

}

// Get is a method to return an item from the Store given a key
// Get should modify an items accessedAt field
func (s *Store) Get(key string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj, found := s.cache[key]
	if !found {
		return nil, errKeyNotFound
	}
	// remove from cache if the item has expired
	if time.Since(obj.accessedAt) > s.lifeTime {
		delete(s.cache, key)

		return nil, errKeyNotFound
	}
	obj.accessed()
	return obj.val, nil
}

// Delete is a method to remove a key-value pair from the Store
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.cache, key)

}

func (s *Store) len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.cache)
}

// Freeze is a function to halt Add/Delete methods
func (s *Store) Freeze() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.frozen = true
}

// Unfreeze allows Add/Delete methods if cache is frozen
func (s *Store) Unfreeze() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.frozen = false
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
