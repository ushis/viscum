package rpc

import (
  "net/rpc"
  "viscum/rpc/mem"
  "viscum/rpc/queue"
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

// Subscribes a email to a feed.
func (self *Client) Subscribe(email string, url string) (Reply, error) {
  return subscription.Subscribe(self.connection, email, url)
}

// Unsubscribes an email from a feed.
func (self *Client) Unsubscribe(email string, url string) (Reply, error) {
  return subscription.Unsubscribe(self.connection, email, url)
}

// Fetches subscription info from the server.
func (self *Client) ListSubscriptions(email string) (Reply, error) {
  return subscription.List(self.connection, email)
}

// Fetches queue info from the server.
func (self *Client) ListQueue() (Reply, error) {
  return queue.List(self.connection)
}

// Attempt to send all remaining mail.
func (self *Client) DeliverQueue() (Reply, error) {
  return queue.Deliver(self.connection)
}

// Fetches mem stats from the server.
func (self *Client) MemStats() (Reply, error) {
  return mem.Stats(self.connection)
}
