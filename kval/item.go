package kval

import (
	"time"
)

// Item represents something to be cached in memory
type Item struct {
	key        string
	val        []byte
	createdAt  time.Time
	accessedAt time.Time
	index      int
}

// NewItem returns a new item with given key-val
func NewItem(key string, val []byte) *Item {
	t := time.Now()
	return &Item{
		key:        key,
		val:        val,
		createdAt:  t,
		accessedAt: t,
		index:      -1,
	}
}

// Size returns the size of the value within an Item
func (i *Item) Size() int {
	return len(i.val)
}

// Less is a function to satisfy google/btree interface
// this creates a strict weak ordering in the cache where items
// are ordered by the time they were accessed. Items accessed more
// recently are greater than items accessed less recently.
// If two items are accessed at the same time the return will default to true
func (i *Item) Less(j *Item) bool {
	timeSinceI := time.Since(i.accessedAt)
	timeSinceJ := time.Since(j.accessedAt)

	// if i was accessed later than j then i < j
	if (timeSinceI - timeSinceJ) > 0 {
		return true
	}
	return false

}
