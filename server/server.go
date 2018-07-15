package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"

	"github.com/andresoro/kval"
	"github.com/gorilla/mux"
)

var (
	store  *kval.Store
	apiURL = "api/cache/"
)

func main() {
	// init store
	store = kval.New()
	log.Printf("Data store initialized")

	// init router
	r := mux.NewRouter()
	r.HandleFunc(apiURL, getHandler).Methods("GET")
	r.HandleFunc(apiURL, putHandler).Methods("POST")

	store.Add("key", "val")

	// run server
	log.Fatal("ListenAndServe", http.ListenAndServe(":8080", r))
}

// GET request to url/apiEndpoint/{key}
func getHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	// if key doesnt exist
	key := r.URL.Path[len(apiURL):]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty request"))
		return
	}

	// get item from store
	item, err := store.Get(key)
	if err != nil {
		w.Write([]byte("Key not found"))
		log.Printf("Error with store.Get")
		return
	}

	// encode item
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error encoding data for GET request")
	}

	w.Write(buf.Bytes())

}

// PUT request to url/apiEndpoint/{key}?{val}
func putHandler(w http.ResponseWriter, r *http.Request) {

}
