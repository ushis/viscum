package subscription

import (
  "fmt"
  "viscum/db"
  "viscum/util"
)

type Service struct {
  db          db.DB    // Database connection
  fetcherCtrl chan int // Control channel to the fetcher
}

type Args struct {
  Email string // Email address
  Url   string // Feed url
}

type Reply struct {
  Reply string // Reply text
}

// Returns a new subscription service.
func New(database db.DB, ctrl chan int) *Service {
  return &Service{db: database, fetcherCtrl: ctrl}
}

// Subscribes a email to a feed.
func (self *Service) Subscribe(args *Args, reply *Reply) error {
  if _, err := self.db.Subscribe(args.Email, args.Url); err != nil {
    util.Error("[RPC]", err)
    return err
  }

  reply.Reply = fmt.Sprintf("Subscribed %s to %s", args.Email, args.Url)
  util.Info("[RPC]", reply.Reply)
  self.fetcherCtrl <- util.CTRL_RELOAD
  return nil
}

// Unsubscribes a email from a feed.
func (self *Service) Unsubscribe(args *Args, reply *Reply) error {
  if _, err := self.db.Unsubscribe(args.Email, args.Url); err != nil {
    util.Error("[RPC]", err)
    return err
  }

  reply.Reply = fmt.Sprintf("Unsubscribed %s from %s", args.Email, args.Url)
  util.Info("[RPC]", reply.Reply)
  return nil
}
