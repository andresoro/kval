package shared

import (
	"errors"

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
		return
	}
	resp.Val = nil
	resp.Msg = "Key-Value successfully added"

	return
}

// Get method to expose
func (h *Handler) Get(req Request, resp *Response) (err error) {
	val, err := h.Store.Get(req.Key)
	if err != nil {
		return
	}

	resp.Val = val

	return

}

// Delete method to expose
func (h *Handler) Delete(req Request, resp *Response) (err error) {
	val, err := h.Store.Delete(req.Key)
	if err != nil {
		return
	}
	resp.Val = val
	resp.Msg = "Key successfully deleted"

	return
}
