package queue

import (
  "net/rpc"
)

// Calls the server.
func call(c *rpc.Client, name string) (*Reply, error) {
  var reply Reply
  err := c.Call(name, &Args{}, &reply)
  return &reply, err
}

// Commands the server to attempt to deliver all queue items.
func Deliver(c *rpc.Client) (*Reply, error) {
  return call(c, "Queue.Deliver")
}

// Fetches queue info from the server.
func Info(c *rpc.Client) (*Reply, error) {
  return call(c, "Queue.Info")
}
