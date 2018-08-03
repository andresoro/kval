package kval

import (
	"log"
)

// Server is a kval server
type Server struct {
	db *Store
	host string
	port string
}

func (s *Server) Run() {
	log.Printf("Starting server")
	
}