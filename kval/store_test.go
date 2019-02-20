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

	var wg sync.WaitGroup
	wg.Add(3)
	keys := 1000

	go func() {
		defer wg.Done()
		for i := 0; i < keys; i++ {
			store.Add(fmt.Sprintf("key%d", i), []byte(testVal))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < keys; i++ {
			store.Get(fmt.Sprintf("key%d", i))
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

func TestNew(t *testing.T) {

	_, err := New(3, time.Minute)
	if err == nil {
		t.Error("Store must only take power of two for shard num")
	}

}
func TestGet(t *testing.T) {
	store, _ := New(4, 5*time.Millisecond)

	store.Add("test", []byte("154"))
	data, err := store.Get("test")

	if err != nil {
		t.Error(err)
	}

	if string(data) != "154" {
		t.Errorf("Value should be %d", 154)
	}
}

func TestAdd(t *testing.T) {
	store, _ := New(4, 5*time.Millisecond)

	k, v := "test", "data"

	err := store.Add(k, []byte(v))
	if err != nil {
		t.Error(err)
	}

	data, err := store.Get(k)
	if err != nil {
		t.Error(err)
	}

	if string(data) != v {
		t.Errorf("Value returned should be %s, got %s", v, string(data))
	}

	err2 := store.Add(k, []byte("data"))
	if err2 == nil {
		t.Error("Store should return an error when Adding an existing key")
	}

}

func TestDelete(t *testing.T) {
	store, _ := New(4, 5*time.Millisecond)

	store.Add("test", []byte("data"))
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

	store.Add("key", []byte("val"))
	time.Sleep(10 * time.Millisecond)

	i, err := store.Get("key")
	fmt.Println(i)
	if err == nil {
		t.Error("Item should not be in cache after lifetime duration, clean or janitor func not working")
	}

}

func TestFlush(t *testing.T) {
	store, _ := New(4, time.Second)

	store.Add("key", []byte("val"))
	store.Add("key2", []byte("val"))

	store.Flush()
	_, err := store.Get("key")
	if err == nil {
		t.Error("Flush should delete all keys")
	}

	_, err = store.Get("key2")
	if err == nil {
		t.Error("Flush should delete all keys")
	}

}

func TestAtomic(t *testing.T) {
	store, _ := New(4, time.Minute)

	store.Add("key1", []byte("fhjjajfa"))
	store.Add("key2", []byte("jahdada"))

	if store.Size() != 2 {
		t.Error("Atomic size of cache should be 2")
	}

	store.Delete("key2")

	if store.Size() != 1 {
		t.Error("Atomic size of cache should be 1")
	}

	store.Delete("key1")

	if store.Size() != 0 {
		t.Error("Atomic size of cache should be 0")
	}

}

func TestFreeze(t *testing.T) {
	store, _ := New(4, 5*time.Minute)

	store.Add("key", []byte("Test"))
	store.Freeze()
	store.Add("key2", []byte("Test"))

	_, err := store.Get("key2")
	if err != errKeyNotFound {
		t.Error("When frozen, store should not add values")
	}

	store.Unfreeze()
	store.Add("key2", []byte("1313414"))
	_, err = store.Get("key2")
	if err != nil {
		t.Error("Not adding to store after unfreeze")
	}

}

func TestLess(t *testing.T) {
	a := NewItem("key", []byte("val"))
	time.Sleep(5 * time.Millisecond)
	b := NewItem("key2", []byte("val"))

	if a.Less(b) != true {
		t.Error("The item added later should be Less than b")
	}

	if b.Less(a) != false {
		t.Error("The item added earlier should be greater than later item")
	}
}
