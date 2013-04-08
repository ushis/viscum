package rpc

import (
  "fmt"
  "runtime"
  "viscum/db"
  . "viscum/util"
)

type Service struct {
  db          db.DB      // Database connection
  mailerCtrl  chan<- int // Control channel to the mailer
  fetcherCtrl chan<- int //Control channel to the fetcher
}

// Returns a new service.
func NewService(db db.DB, mc chan<- int, fc chan<- int) *Service {
  return &Service{db: db, mailerCtrl: mc, fetcherCtrl: fc}
}

// Serves memory stats.
func (self *Service) MemStats(_ *Args, r *Reply) (err error) {
  m := new(runtime.MemStats)
  runtime.ReadMemStats(m)
  r.Append(fmt.Sprint("Sys: ", m.Sys, " Heap: ", m.HeapAlloc))
  return
}

// Fetches queue info from the database and sends it to the client.
func (self *Service) QueueInfo(_ *Args, r *Reply) (err error) {
  Info("[RPC] Fetch queue info.")

  if r.Reply, err = self.db.QueueInfo(); err != nil {
    Error("[RPC]", err)
  }
  return
}

// Sends the mailer a heads up.
func (self *Service) Deliver(_ *Args, r *Reply) (err error) {
  Info("[RPC] Hey Mailer! Wake up!")
  self.mailerCtrl <- CTRL_RELOAD
  r.Append("Initiated queue delivery.")
  return
}

// Subscribes a email to a feed.
func (self *Service) Subscribe(args *Args, r *Reply) (err error) {
  if _, err = self.db.Subscribe(args.Email, args.Url); err != nil {
    Error("[RPC]", err)
    return
  }
  // Notify the fetcher about new subscriptions.
  self.fetcherCtrl <- CTRL_RELOAD
  r.Append(fmt.Sprint("Subscribed ", args.Email, " to ", args.Url))
  Info("[RPC] Subscribed", args.Email, "to", args.Url)
  return
}

// Unsubscribes a email from a feed.
func (self *Service) Unsubscribe(args *Args, r *Reply) (err error) {
  if _, err = self.db.Unsubscribe(args.Email, args.Url); err != nil {
    Error("[RPC]", err)
    return
  }
  r.Append(fmt.Sprint("Unsubscribed ", args.Email, " from ", args.Url))
  Info("[RPC] Unsubscribed", args.Email, "from", args.Url)
  return
}

// Lists subscriptions filtered by email.
func (self *Service) ListSubscriptions(args *Args, r *Reply) (err error) {
  Info("[RPC] Fetch subscriptions for:", args.Email)

  if r.Reply, err = self.db.ListSubscriptions(args.Email); err != nil {
    Error("[RPC]", err)
  }
  return
}
