package shared

import (
	"errors"
	"log"

	"github.com/andresoro/kval/kval"
)

// holds shared objects for rpc

// Request to the server
type Request struct {
	Key string
	Val interface{}
}

// Response from the server
type Response struct {
	Val interface{}
	Msg string
}

// Handler for rpc methods
type Handler struct {
	Store *kval.Store
}

// Add method to expose
func (h *Handler) Add(req Request, resp *Response) (err error) {
	if req.Val == nil {
		err = errors.New("Value must not be empty")
		return
	}

	err = h.Store.Add(req.Key, req.Val)
	if err != nil {
		log.Print(err)
		return
	}
	resp.Val = nil
	resp.Msg = "Key-Value successfully added"

	log.Printf("Key added to store: %s", req.Key)

	return
}

// Get method to expose
func (h *Handler) Get(req Request, resp *Response) (err error) {
	val, err := h.Store.Get(req.Key)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("Getting key from store: ", req.Key)

	resp.Val = val

	return

}

// Delete method to expose
func (h *Handler) Delete(req Request, resp *Response) (err error) {
	val, err := h.Store.Delete(req.Key)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("Deleting key from store: ", req.Key)
	resp.Val = val
	resp.Msg = "Key successfully deleted"

	return
}

// Flush method to expose
func (h *Handler) Flush(req Request, resp *Response) (err error) {
	h.Store.Flush()
	msg := "Flushing all keys from the cache..."
	log.Print(msg)
	resp.Val = nil
	resp.Msg = msg

	return nil
}
