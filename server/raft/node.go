package raft

import (
	"fmt"
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
	RaftBind string
	Dir      string
}

// New node from given in-mem store
func New(k *kval.Store) *Node {
	return &Node{
		store: k,
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
