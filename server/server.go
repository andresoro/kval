package server

import (
	"log"
	"net/http"
	"time"

	"github.com/andresoro/kval/kval"
	"github.com/gorilla/mux"
)

var (
	store  *kval.Store
	apiURL = "/api/cache/"
)

// Server for key-value store
type Server struct {
	store  *kval.Store
	router *mux.Router
}

// New returns a kval server
func New() *Server {
	return &Server{
		store:  kval.New(4, 5*time.Minute),
		router: mux.NewRouter(),
	}
}

// Start the server on given port
func (s *Server) Start(port string) {
	s.router.HandleFunc("/api/cache/{key}", getHandler).Methods("GET")
	s.router.HandleFunc(apiURL, putHandler).Methods("POST")

	log.Fatal("ListenAndServer", http.ListenAndServe(port, s.router))

}
