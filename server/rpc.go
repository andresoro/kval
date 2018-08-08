package server

import (
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

// NewRPC returns an rpc server
func NewRPC(port string) *RPCServer {
	return &RPCServer{
		port:  port,
		store: kval.New(4, 3*time.Minute),
	}
}

// Start will init the server, verify necessary methods, and expose them
// to clients
func (r *RPCServer) Start() (err error) {

	// register the shared methods
	rpc.Register(&shared.Handler{
		Store: r.store,
	})

	r.listener, err = net.Listen("tcp", r.port)
	if err != nil {
		return
	}

	rpc.Accept(r.listener)

	return
}
