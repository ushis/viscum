package rpc

import (
  "net/rpc"
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
func (self *Client) SubscriptionInfo(email string) (Reply, error) {
  return subscription.Info(self.connection, email)
}

// Fetches queue info from the server.
func (self *Client) QueueInfo() (Reply, error) {
  return queue.Info(self.connection)
}

// Attempt to send all remaining mail.
func (self *Client) DeliverQueue() (Reply, error) {
  return queue.Deliver(self.connection)
}
