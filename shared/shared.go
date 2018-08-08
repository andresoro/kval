package shared

import (
	"errors"

	"github.com/andresoro/kval/kval"
)

// holds shared objects for rpc

// Request to the server
type Request struct {
	key string
	val interface{}
}

// Response from the server
type Response struct {
	val interface{}
}

// Handler for rpc methods
type Handler struct {
	Store *kval.Store
}

// Add method to expose
func (h *Handler) Add(req Request, resp *Response) (err error) {
	if req.val == nil {
		err = errors.New("Value must not be empty")
		return
	}

	err = h.Store.Add(req.key, req.val)
	if err != nil {
		return
	}
	resp.val = nil

	return
}

// Get method to expose
func (h *Handler) Get(req Request, resp *Response) (err error) {
	val, err := h.Store.Get(req.key)
	if err != nil {
		return
	}

	resp.val = val

	return

}

// Delete method to expose
func (h *Handler) Delete(req Request, resp *Response) (err error) {
	val, err := h.Store.Delete(req.key)
	if err != nil {
		return
	}
	resp.val = val

	return
}
