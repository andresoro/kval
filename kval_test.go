package kval

import (
	"testing"
)

func TestAdd(t *testing.T) {
	store := New()
	err := store.Add("test", "data")
	if err != nil {
		t.Fatal(err)
	}

	c := store.Len()
	if c != 1 {
		t.Fatalf("Expected len of store to be 1, got %d", c)
	}
}
