package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andresoro/kval/kval"
	"github.com/gorilla/mux"
)

var (
	store  *kval.Store
	apiURL = "/kval/"
)

// Server for key-value store
type Server struct {
	port   string
	store  *kval.Store
	router *mux.Router
}

// NewHTTP returns an http kval server
func NewHTTP(port string, kval *kval.Store) *Server {
	store = kval
	return &Server{
		port:   port,
		store:  store,
		router: mux.NewRouter(),
	}
}

// Start the server on given port
func (s *Server) Start() {
	s.router.HandleFunc("/kval/{key}", getHandler).Methods("GET")
	s.router.HandleFunc("/kval/{key}", putHandler).Methods("PUT")

	log.Print("Starting HTTP server on port: ", s.port)
	log.Fatal("ListenAndServer", http.ListenAndServe(s.port, s.router))

}

// GET request to url/apiURL/{key}
func getHandler(w http.ResponseWriter, r *http.Request) {

	// get key from URL params
	vars := mux.Vars(r)
	key := vars["key"]

	// get item from store
	item, err := store.Get(key)
	if err != nil {
		w.Write([]byte("Key not found"))
		log.Print(err)
		return
	}

	// write item to client
	json, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("GET request")
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}

// PUT request to url/apiEndpoint/?key=key
func putHandler(w http.ResponseWriter, r *http.Request) {
	// handle url paramters and request body

	vars := mux.Vars(r)
	key := vars["key"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error reading request body"))
		log.Print(err)
		return
	}
	//convert to string
	content := string(body[:])

	// add key-item to cache
	err = store.Add(key, content)
	if err != nil {
		log.Print(err)
		w.Write([]byte("Unable to add key-value pair to the cache"))
		return
	}
	// return success/fail message
	log.Print("PUT request")
	w.Write([]byte("Successfully added key-value pair to the cache"))
}
