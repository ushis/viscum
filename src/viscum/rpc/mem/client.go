package mem

import (
  "net/rpc"
)

// Fetches memory stats from the server.
func Stats(c *rpc.Client) (*Reply, error) {
  reply := &Reply{}
  err := c.Call("Mem.Stats", &Args{}, reply)
  return reply, err
}
