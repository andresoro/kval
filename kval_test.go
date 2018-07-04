package kval

import (
	"testing"
)

func TestSet(t *testing.T) {
	store := New()
	if store == nil {
		t.Error("New Store not being returned")
	}

	store.Set("test", "data")

	if store.Len() != 1 {
		t.Error("Len of store should be 1")
	}
}

func TestGet(t *testing.T) {
	store := New()
	store.Set("test", 154)
	data, err := store.Get("test")
	if err != nil {
		t.Error(err)
	}
	if data != 154 {
		t.Errorf("Value should be %d", 154)
	}
}
