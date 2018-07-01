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
//
type Store struct {
	lifeTime time.Duration
	cache    map[string]*item
	sync.RWMutex
}

// Add is a method to add an object to the Store
// Add does not replace an item in the Store if the key already exists
func (s *Store) Add(key string, val interface{}) error {
	s.Lock()
	defer s.Unlock()
	i := newItem(key, val)

	// if key exists then return ErrKeyExists
	_, err := s.Get(key)
	if err != nil {
		return ErrKeyExists
	}

	s.cache[key] = i
	return nil
}

// Get is a method to return an item from the Store given a key
// Get should modify an items accessedAt field
func (s *Store) Get(key string) (interface{}, error) {
	s.Lock()
	defer s.Unlock()

	obj, found := s.cache[key]
	if !found {
		return nil, ErrKeyNotFound
	}
	if time.Since(obj.accessedAt) > s.lifeTime {
		delete(s.cache, key)

		return nil, ErrKeyNotFound
	}

	return obj.val, nil
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
