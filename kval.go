package kval

import (
	"time"
)

// Store is the in memory key value store that holds items for a max duration
// Store is a col
type Store struct {
	LifeTime time.Duration
	cache    map[string]*item
}

// item represents something to be cached in memory
type item struct {
	key       string
	val       interface{}
	createdAt time.Time
}
