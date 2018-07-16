package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andresoro/kval"
	"github.com/gorilla/mux"
)

var (
	store  *kval.Store
	apiURL = "/api/cache/"
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

// GET request to url/apiEndpoint/key={key}
func getHandler(w http.ResponseWriter, r *http.Request) {

	// get key from URL params
	query := r.URL.Query()
	key := query.Get("key")

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

	w.Write(json)

}

// PUT request to url/apiEndpoint/{key}?{val}
func putHandler(w http.ResponseWriter, r *http.Request) {
	// handle url paramters and request body

	// add key-item to cache

	// return success/fail message
}
