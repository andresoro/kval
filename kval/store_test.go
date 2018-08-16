package kval

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestParrallel(t *testing.T) {
	t.Parallel()

	store, _ := New(4, 1*time.Minute)
	testVal := "test value"

	var val interface{}
	var wg sync.WaitGroup
	wg.Add(3)
	keys := 1000

	go func() {
		defer wg.Done()
		for i := 0; i < keys; i++ {
			store.Add(fmt.Sprintf("key%d", i), testVal)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < keys; i++ {
			val, _ = store.Get(fmt.Sprintf("key%d", i))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < keys; i++ {
			store.Delete(fmt.Sprintf("key%d", i))
		}
	}()

	wg.Wait()

}

func TestGet(t *testing.T) {
	store, _ := New(4, 5*time.Minute)

	store.Add("test", 154)
	data, err := store.Get("test")

	if err != nil {
		t.Error(err)
	}

	if data != 154 {
		t.Errorf("Value should be %d", 154)
	}
}

func TestAdd(t *testing.T) {
	store, _ := New(4, 5*time.Minute)

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

	err2 := store.Add(k, "data")
	if err2 == nil {
		t.Error("Store should return an error when Adding an existing key")
	}

}

func TestDelete(t *testing.T) {
	store, _ := New(4, 5*time.Minute)

	store.Add("test", "data")
	_, err := store.Get("test")
	if err != nil {
		t.Error("Not adding value to store")
	}

	store.Delete("test")
	_, err = store.Get("test")
	if err != errKeyNotFound {
		t.Error("Key found in store when it should have been deleted")
	}
}

func TestClean(t *testing.T) {
	store, _ := New(4, 5*time.Millisecond)

	store.Add("key", "val")
	time.Sleep(6 * time.Millisecond)

	i, err := store.Get("key")
	fmt.Println(i)
	if err == nil {
		t.Error("Item should not be in cache after lifetime duration, clean or janitor func not working")
	}

}

func TestFreeze(t *testing.T) {
	store, _ := New(4, 5*time.Minute)

	store.Add("key", 981093813)
	store.Freeze()
	store.Add("key2", 1313414)

	_, err := store.Get("key2")
	if err != errKeyNotFound {
		t.Error("When frozen, store should not add values")
	}

	store.Unfreeze()
	store.Add("key2", 1313414)
	_, err = store.Get("key2")
	if err != nil {
		t.Error("Not adding to store after unfreeze")
	}

}

func TestLess(t *testing.T) {
	a := newItem("key", "val")
	time.Sleep(5 * time.Millisecond)
	b := newItem("key2", "val")

	if a.Less(b) != true {
		t.Error("The item added later should be Less than b")
	}

	if b.Less(a) != false {
		t.Error("The item added earlier should be greater than later item")
	}
}
