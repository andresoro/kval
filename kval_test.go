package kval

import (
	"testing"
)

func TestSet(t *testing.T) {
	store := New()
	if store == nil {
		t.Error("New Store not being returned")
	}

	store.set("test", "data")

	if store.len() != 1 {
		t.Error("Len of store should be 1")
	}

	store.set("test2", "data")

	if store.len() != 2 {
		t.Error("Len of store should be 2")
	}

	store.set("test2", 123)
	if store.len() != 2 {
		t.Error("Len of store should 2")
	}
}

func TestGet(t *testing.T) {
	store := New()
	store.set("test", 154)
	data, err := store.Get("test")
	if err != nil {
		t.Error(err)
	}
	if data != 154 {
		t.Errorf("Value should be %d", 154)
	}
}

func TestAdd(t *testing.T) {
	store := New()

	k, v := "test", 15141

	err := store.Add(k, v)
	if err != nil {
		t.Error(err)
	}

	data, err := store.Get(k)
	if err != nil {
		t.Error(err)
	}

	if data != v {
		t.Errorf("Value returned should be %d, got %d", v, data)
	}

}
