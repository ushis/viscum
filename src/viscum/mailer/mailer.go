package mailer

import (
  "viscum/db"
  . "viscum/util"
)

type Mailer struct {
  db      db.DB    // Database connection
  conf    *Config  // Config
  Ctrl    chan int // Control channel
  handler Handler  // Mail handler
  from    string   // Sender address
}

// Available handlers.
var handlers = make(map[string]Handler)

// Registers a new handler.
func Register(name string, handler Handler) {
  if _, dup := handlers[name]; dup {
    Fatal("Handler registered twice:", name)
  }
  handlers[name] = handler
}

// Returns a new mailer.
func New(database db.DB, conf *Config) *Mailer {
  name := conf.Get("mail", "mailer")
  handler, ok := handlers[name]

  if !ok {
    Fatal("Unknown mail handler:", name)
  }
  handler.Init(conf)

  return &Mailer{
    db:      database,
    conf:    conf,
    Ctrl:    make(chan int),
    handler: handler,
    from:    conf.Get("mail", "from"),
  }
}

// Starts the mailer.
func (self *Mailer) Start() {
  Info("[Mailer] Start...")

  for {
    err := self.db.FetchQueue(func(entry *db.Entry) {
      self.handleEntry(entry)
    })

    if err != nil {
      Error(err)
    }

    // Wait for instructions.
    if CTRL_STOP == <-self.Ctrl {
      break
    }
  }

  Info("[Mailer] Stop...")
}

// Commands the mailer to stop.
func (self *Mailer) Stop() {
  self.Ctrl <- CTRL_STOP
}

// Handles a queue entry.
func (self *Mailer) handleEntry(e *db.Entry) {
  Info("[Mailer] Send:", e.Id, e.Email, e.Title)
  success := true

  if err := self.handler.Send(NewMail(e, self.from)); err != nil {
    success = false
    Error("[Mailer]", err)
  }

  if _, err := self.db.Dequeue(e, success); err != nil {
    Error("[Mailer]", err)
  }
}
