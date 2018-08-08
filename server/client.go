package server

import (
	"context"
	"net/rpc"

	"github.com/andresoro/kval/shared"
)

// Client represents an RPC client
type Client struct {
	port   string
	client *rpc.Client
}

// Init client
func (c *Client) Init() (err error) {
	c.client, err = rpc.Dial("tcp", "127.0.0.1:"+c.port)
	if err != nil {
		return
	}
	return
}

// Close method
func (c *Client) Close() (err error) {
	if c.client != nil {
		err = c.client.Close()
		return
	}
	return

}

// Add client method
func (c *Client) Add(ctx context.Context, key string, val interface{}) (msg string, err error) {
	var (
		request = &shared.Request{
			Key: key,
			Val: val,
		}
		response = new(shared.Response)
	)

	err = c.client.Call("Handler.Add", request, response)
	if err != nil {
		return
	}

	msg = response.Msg
	return
}

// Get client method
func (c *Client) Get(ctx context.Context, key string) (msg string, err error) {
	var (
		request = &shared.Request{
			Key: key,
			Val: nil,
		}
		response = new(shared.Response)
	)

	err = c.client.Call("Handler.Get", request, response)
	if err != nil {
		return
	}
	msg = response.Msg
	return
}

// Delete client method
func (c *Client) Delete(ctx context.Context, key string) (msg string, err error) {
	var (
		request = &shared.Request{
			Key: key,
			Val: nil,
		}
		response = new(shared.Response)
	)

	err = c.client.Call("Handler.Delete", request, response)
	if err != nil {
		return
	}
	msg = response.Msg
	return
}
