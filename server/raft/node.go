package raft

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/andresoro/kval/kval"
	"github.com/hashicorp/raft"
)

var (
	timeout   = 10 * time.Second
	snapCount = 2
)

// Node represents a raft node with an inmemory datastore
type Node struct {
	raft     *raft.Raft
	store    *kval.Store
	mu       sync.Mutex
	log      *log.Logger
	RaftBind string
	Dir      string
}

// New node from given in-mem store
func New() *Node {
	s, err := kval.New(2, 5*time.Minute)
	if err != nil {
		return nil
	}
	return &Node{
		store: s,
	}
}

// Open inititalizes the raft cluster system. If single is set to true
// and there are no other nodes in the cluster, this node becomes the leader.
// id is the server identifier for this node
func (n *Node) Open(single bool, id string) error {
	// raft config
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(id)

	// raft communication
	addr, err := net.ResolveTCPAddr("tcp", n.RaftBind)
	if err != nil {
		return err
	}
	trans, err := raft.NewTCPTransport(n.RaftBind, addr, 3, timeout, os.Stderr)
	if err != nil {
		return err
	}

	// snapshot store
	snaps, err := raft.NewFileSnapshotStore(n.Dir, snapCount, os.Stderr)
	if err != nil {
		return err
	}

	// log store, stable store
	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()

	// initiate raft system
	ra, err := raft.NewRaft(conf, (*fsm)(n), logStore, stableStore, snaps, trans)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}

	n.raft = ra

	if single {
		c := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      conf.LocalID,
					Address: trans.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(c)
	}

	return nil
}

// Get returns a value for the given key
// TODO: add raft method to update TTL in all nodes
func (n *Node) Get(key string) ([]byte, error) {
	return n.store.Get(key)
}

// Add a key to the node and then pass to the other nodes in the cluster
func (n *Node) Add(key string, val []byte) error {
	if n.raft.State() != raft.Leader {
		return errors.New("Not cluster leader")
	}

	c := &command{
		Op:  "add",
		Key: key,
		Val: val,
	}

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := n.raft.Apply(b, timeout)
	return f.Error()
}

// Delete a key from the store and tell other nodes
func (n *Node) Delete(key string) error {
	if n.raft.State() != raft.Leader {
		return errors.New("Not cluster leader")
	}

	c := &command{
		Op:  "delete",
		Key: "key",
	}

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := n.raft.Apply(b, timeout)
	return f.Error()
}

// Join a node to this store
func (n *Node) Join(id, addr string) error {
	confFuture := n.raft.GetConfiguration()
	err := confFuture.Error()
	if err != nil {
		return err
	}

	for _, srv := range confFuture.Configuration().Servers {

		//if a node exists with the joining server ID or addr, it may need to be removed
		if srv.ID == raft.ServerID(id) || srv.Address == raft.ServerAddress(addr) {
			// if ID and addr are the same, do nothing
			if srv.ID == raft.ServerID(id) && srv.Address == raft.ServerAddress(addr) {
				n.log.Printf("Node already in the cluster. id: %s, addr: %s \n", id, addr)
				return nil
			}
			//remove
			future := n.raft.RemoveServer(srv.ID, 0, 0)
			err = future.Error()
			if err != nil {
				return err
			}
		}
	}

	f := n.raft.AddVoter(raft.ServerID(id), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return err
	}
	n.log.Printf("node added with addr: %s, id: %s", id, addr)
	return nil
}
