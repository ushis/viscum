package rpc

import (
  "net/rpc"
)

type Client struct {
  *rpc.Client        // RPC connection
  socket      string // Socket of the rpc server
}

// Returns a new client.
func NewClient(socket string) *Client {
  return &Client{socket: socket}
}

// Connects to the rpc server.
func (self *Client) Connect() (err error) {
  if len(self.socket) > 0 && self.socket[0] != '/' {
    self.Client, err = rpc.DialHTTP("tcp", self.socket)
  } else {
    self.Client, err = rpc.DialHTTP("unix", self.socket)
  }
  return
}

// Subscribes a email to a feed.
func (self *Client) Subscribe(email string, url string) (*Reply, error) {
  return self.call("S.Subscribe", email, url)
}

// Unsubscribes an email from a feed.
func (self *Client) Unsubscribe(email string, url string) (*Reply, error) {
  return self.call("S.Unsubscribe", email, url)
}

// Fetches subscription info from the server.
func (self *Client) ListSubscriptions(email string) (*Reply, error) {
  return self.call("S.ListSubscriptions", email, "")
}

// Fetches queue info from the server.
func (self *Client) QueueInfo() (*Reply, error) {
  return self.call("S.QueueInfo", "", "")
}

// Attempt to send all remaining mail.
func (self *Client) Deliver() (*Reply, error) {
  return self.call("S.Deliver", "", "")
}

// Fetches mem stats from the server.
func (self *Client) MemStats() (*Reply, error) {
  return self.call("S.MemStats", "", "")
}

// Calls the server.
func (self *Client) call(f string, email string, url string) (*Reply, error) {
  var reply Reply
  err := self.Call(f, &Args{Email: email, Url: url}, &reply)
  return &reply, err
}
