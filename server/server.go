package server

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Service is the key-value node interface we want to expose
type Service interface {
	Get(key string) ([]byte, error)
	Add(key string, val []byte) error
	Delete(key string) error
}

// Server is the http server for the raft node
type Server struct {
	store  Service
	addr   string
	router *mux.Router
}

// New returns a new server for the service
func New(serv Service, addr string) *Server {

	s := &Server{
		store:  serv,
		router: mux.NewRouter(),
		addr:   addr,
	}

	// handle the add, get, delete method
	s.router.HandleFunc("/key/{key}", s.handle)
	// handle a new node joining
	s.router.HandleFunc("/join", s.join)

	return s
}

// run the server on given addr
func (s *Server) Run() {
	http.ListenAndServe(s.addr, s.router)
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	switch r.Method {
	case "GET":
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		val, err := s.store.Get(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write(val)
		return

	case "POST":
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = s.store.Add(key, body)
		if err != nil {
			w.Write([]byte("Key already exists"))
		}

	}

}

func (s *Server) join(w http.ResponseWriter, r *http.Request) {

}
