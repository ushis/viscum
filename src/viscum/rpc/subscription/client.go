package subscription

import (
  "net/rpc"
)

// Calls the server.
func call(c *rpc.Client, f string, email string, url string) (*Reply, error) {
  var reply Reply
  err := c.Call(f, &Args{Email: email, Url: url}, &reply)
  return &reply, err
}

// Subscribes an email to a feed.
func Subscribe(con *rpc.Client, email string, url string) (*Reply, error) {
  return call(con, "Subscription.Subscribe", email, url)
}

// Unsubscribes an email from a feed.
func Unsubscribe(con *rpc.Client, email string, url string) (*Reply, error) {
  return call(con, "Subscription.Unsubscribe", email, url)
}

// Fetches info from the server.
func Info(con *rpc.Client, email string) (*Reply, error) {
  return call(con, "Subscription.Info", email, "")
}
