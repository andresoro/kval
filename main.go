package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/andresoro/kval/server"
	"github.com/andresoro/kval/server/raft"
)

var (
	http     string
	raftBind string
	leader   string
	id       string
)

func init() {
	flag.StringVar(&http, "h", ":8080", "http address to expose")
	flag.StringVar(&raftBind, "r", "", "raft bind address")
	flag.StringVar(&leader, "l", "", "address of leader node to join")
	flag.StringVar(&id, "i", "", "id of this node")
	flag.Parse()
}

func main() {
	if flag.NArg() == 0 {
		fmt.Println("missing raft directory")
	}

	// init node
	node := raft.New()
	// check if this is a leader node
	if leader == "" {
		err := node.Open(true, id)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := node.Open(false, id)
		if err != nil {
			log.Fatalln(err)
		}
		err = node.Join(id, raftBind)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// start server

	srv := server.New(node, ":8080")
	srv.Run()

}
