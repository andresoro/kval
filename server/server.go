package main

import (
	"log"
	"net/http"

	"github.com/andresoro/kval"
	"github.com/gorilla/mux"
)

var (
	store *kval.Store
)

func main() {
	// init store
	store = kval.New()
	log.Printf("Data store initialized")

	// init router
	r := mux.NewRouter()
	r.HandleFunc("/api/", getHandler).Methods("GET")
	r.HandleFunc("/api/", putHandler).Methods("POST")

	// run server
	log.Fatal("ListenAndServe", http.ListenAndServe(":8080", r))
}

// GET request to url/apiEndpoint/{key}
func getHandler(w http.ResponseWriter, r *http.Request) {

}

// PUT request to url/apiEndpoint/{key}?{val}
func putHandler(w http.ResponseWriter, r *http.Request) {

}
