package rpc

import (
	"encoding/gob"
	"log"
	"net"
	"net/rpc"

	"github.com/andresoro/kval/kval"
	"github.com/andresoro/kval/shared"
)

// RPCServer is an rpc to export methods over a network
type RPCServer struct {
	port     string
	listener net.Listener
	store    *kval.Store
}

// NewRPC returns an rpc server with a cache shard size and item time duration
func New(port string, kval *kval.Store) *RPCServer {
	return &RPCServer{
		port:  port,
		store: kval,
	}
}

// Start will init the server, verify necessary methods, and expose them
// to clients
func (r *RPCServer) Start() (err error) {
	// register the shared methods
	rpc.Register(&shared.Handler{
		Store: r.store,
	})

	log.Print("Starting RPC server on port: ", r.port)
	r.listener, err = net.Listen("tcp", r.port)
	if err != nil {
		return
	}

	rpc.Accept(r.listener)

	return
}

// Register will register an obj to encode over rpc
func (r *RPCServer) Register(obj interface{}) {
	gob.Register(obj)
}
