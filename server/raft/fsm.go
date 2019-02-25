package raft

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"
)

type command struct {
	Op  string `json:"op,omitempty"`
	Key string `json:"key,omitempty"`
	Val []byte `json:"val,omitempty"`
}

// implements the FiniteStateMachine interface for the raft consensus
type fsm Node

// Apply raft log entry to the store
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command

	err := json.Unmarshal(l.Data, &c)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case "add":
		return f.applySet(c.Key, c.Val)
	case "delete":
		return f.applyDelete(c.Key)
	case "get":
		return f.applyGet(c.Key)
	default:
		panic(fmt.Sprintf("unrecognized command %s", c.Op))
	}
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return &fsmSnapshot{store: f.store.Snap()}, nil
}

func (f *fsm) Restore(rc io.ReadCloser) error {
	o := make(map[string]string)
	if err := json.NewDecoder(rc).Decode(&o); err != nil {
		return err
	}

	f.store.UnSnap(o)
	return nil
}

func (f *fsm) applyGet(key string) interface{} {
	f.store.Get(key)
	return nil
}

func (f *fsm) applySet(key string, value []byte) interface{} {
	err := f.store.Add(key, value)
	return err
}

func (f *fsm) applyDelete(key string) interface{} {
	_, err := f.store.Delete(key)
	return err
}

type fsmSnapshot struct {
	store map[string]string
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}
		if _, err := sink.Write(b); err != nil {
			return err
		}
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {}
