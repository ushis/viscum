package queue

import (
  "viscum/db"
  . "viscum/util"
)

type Service struct {
  db   db.DB    // Database connection
  ctrl chan int // Control channel to the mailer
}

// Returns a new queue service.
func New(database db.DB, ctrl chan int) (string, *Service) {
  return "Queue", &Service{db: database, ctrl: ctrl}
}

// Fetches queue info from the database and sends it to the client.
func (self *Service) List(_ *Args, reply *Reply) (err error) {
  Info("[RPC] Fetch queue info.")
  reply.Reply, err = self.db.QueueInfo()

  if err != nil {
    Error("[RPC]", err)
  }
  return err
}

// Sends the mailer a heads up.
func (self *Service) Deliver(_ *Args, reply *Reply) error {
  Info("[RPC] Hey Mailer! Wake up!")
  self.ctrl <- CTRL_RELOAD
  reply.Reply = "Initiated queue delivery."
  return nil
}
