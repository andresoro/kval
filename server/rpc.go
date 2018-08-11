package server

import (
	"log"
	"net"
	"net/rpc"
	"time"

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
func NewRPC(port string, shardNum int, duration time.Duration) *RPCServer {
	return &RPCServer{
		port:  port,
		store: kval.New(shardNum, duration),
	}
}

// Start will init the server, verify necessary methods, and expose them
// to clients
func (r *RPCServer) Start() (err error) {
	log.Print("Initializing rpc server...")
	// register the shared methods
	rpc.Register(&shared.Handler{
		Store: r.store,
	})

	log.Print("Listening on port: ", r.port)
	r.listener, err = net.Listen("tcp", r.port)
	if err != nil {
		return
	}

	rpc.Accept(r.listener)

	return
}
