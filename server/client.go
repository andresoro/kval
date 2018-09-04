package server

import (
	"fmt"
	"net/rpc"

	"github.com/andresoro/kval/shared"
)

// Client represents an RPC client
type Client struct {
	Port   string
	client *rpc.Client
}

// Init client
func (c *Client) Init() (err error) {
	c.client, err = rpc.Dial("tcp", "127.0.0.1"+c.Port)
	if err != nil {
		fmt.Println(err)
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
func (c *Client) Add(key string, val interface{}) (msg string, err error) {
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
func (c *Client) Get(key string) (msg interface{}, err error) {
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
	msg = response.Val
	return
}

// Delete client method
func (c *Client) Delete(key string) (msg interface{}, err error) {
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
	msg = response.Val
	return
}

// Flush client method
func (c *Client) Flush() (msg interface{}, err error) {
	var (
		request = &shared.Request{
			Key: "",
			Val: nil,
		}
		response = new(shared.Response)
	)

	c.client.Call("Handler.Flush", request, response)
	msg = response.Val
	return
}
