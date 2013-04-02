package subscription

import (
  "fmt"
  "viscum/db"
  . "viscum/util"
)

type Service struct {
  db   db.DB      // Database connection
  ctrl chan<- int // Control channel to the fetcher
}

// Returns a new subscription service.
func New(database db.DB, ctrl chan<- int) (string, *Service) {
  return "Subscription", &Service{db: database, ctrl: ctrl}
}

// Subscribes a email to a feed.
func (self *Service) Subscribe(args *Args, reply *Reply) error {
  if _, err := self.db.Subscribe(args.Email, args.Url); err != nil {
    Error("[RPC]", err)
    return err
  }

  reply.Reply = fmt.Sprintf("Subscribed %s to %s", args.Email, args.Url)
  Info("[RPC]", reply.Reply)
  self.ctrl <- CTRL_RELOAD
  return nil
}

// Unsubscribes a email from a feed.
func (self *Service) Unsubscribe(args *Args, reply *Reply) error {
  if _, err := self.db.Unsubscribe(args.Email, args.Url); err != nil {
    Error("[RPC]", err)
    return err
  }

  reply.Reply = fmt.Sprintf("Unsubscribed %s from %s", args.Email, args.Url)
  Info("[RPC]", reply.Reply)
  return nil
}

// Lists subscriptions filtered by email.
func (self *Service) List(args *Args, reply *Reply) (err error) {
  Info("[RPC] Fetch subscriptions for:", args.Email)
  reply.Reply, err = self.db.ListSubscriptions(args.Email)

  if err != nil {
    Error("[RPC]", err)
  }
  return err
}
