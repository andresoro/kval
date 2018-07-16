package main

import (
	"encoding/json"
	"io/ioutil"
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
	r.HandleFunc("/api/cache/{key}", getHandler).Methods("GET")
	r.HandleFunc(apiURL, putHandler).Methods("POST")

	// test val with an array of int
	store.Add("key", []int{1, 5141, 13131})

	// run server
	log.Fatal("ListenAndServe", http.ListenAndServe(":8080", r))
}

// GET request to url/apiURL/?key={key}
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}

// PUT request to url/apiEndpoint/?key=key
func putHandler(w http.ResponseWriter, r *http.Request) {
	// handle url paramters and request body
	query := r.URL.Query()
	key := query.Get("key")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	// add key-item to cache
	err = store.Add(key, body)
	if err != nil {
		log.Print(err)
		w.Write([]byte("Unable to add key-value pair to the cache"))
		return
	}
	// return success/fail message
	w.Write([]byte("Successfully added key-value pair to the cache"))
}
