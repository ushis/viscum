package rpc

import (
  "net/rpc"
  "viscum/rpc/subscription"
)

type Client struct {
  socket     string      // Socket of the rpc server
  connection *rpc.Client // RPC connection
}

// Returns a new client.
func NewClient(socket string) *Client {
  return &Client{socket: socket}
}

// Connects to the rpc server.
func (self *Client) Connect() (err error) {
  self.connection, err = rpc.DialHTTP("unix", self.socket)
  return err
}

// Disconnects from the rpc server.
func (self *Client) Disconnect() error {
  return self.connection.Close()
}

// Subscribes an email to a feed.
func (self *Client) Subscribe(email string, url string) (string, error) {
  var reply subscription.Reply
  args := &subscription.Args{Email: email, Url: url}
  err := self.connection.Call("Service.Subscribe", args, &reply)
  return reply.Reply, err
}

// Unsubscribes an email to a feed.
func (self *Client) Unsubscribe(email string, url string) (string, error) {
  var reply subscription.Reply
  args := &subscription.Args{Email: email, Url: url}
  err := self.connection.Call("Service.Unsubscribe", args, &reply)
  return reply.Reply, err
}
