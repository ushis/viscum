package mailer

import (
  "fmt"
  "os/exec"
  "viscum/db"
  "viscum/util"
)

type Mailer struct {
  db   db.DB    // Database connection
  cmd  string   // Mail command
  Ctrl chan int // Control channel
}

// Returns a new mailer.
func New(database db.DB, cmd string) *Mailer {
  return &Mailer{db: database, cmd: cmd, Ctrl: make(chan int)}
}

// Starts the mailer.
func (self *Mailer) Start() {
  util.Info("[Mailer] Start...")

  for {
    err := self.db.FetchQueue(func(entry *db.QueueEntry) {
      self.handleQueueEntry(entry)
    })

    if err != nil {
      util.Error(err)
    }

    // Wait for instructions.
    if ctrl := <-self.Ctrl; ctrl == util.CTRL_STOP {
      break
    }
  }

  util.Info("[Mailer] Stop...")
}

// Commands the mailer to stop.
func (self *Mailer) Stop() {
  self.Ctrl <- util.CTRL_STOP
}

// Handles a queue entry.
func (self *Mailer) handleQueueEntry(e *db.QueueEntry) {
  success := true

  if err := self.send(e); err != nil {
    success = false
    util.Error("[Mailer]", err)
  }

  self.db.Dequeue(e, success)
}

// Sends the message.
func (self *Mailer) send(e *db.QueueEntry) error {
  subject := fmt.Sprintf("[%s] %s", e.FeedTitle, e.Title)
  util.Info("[Mailer] Send:", e.Id, e.Email, subject)

  cmd := exec.Command(self.cmd, "-s", subject, e.Email)
  stdin, err := cmd.StdinPipe()

  if err != nil {
    return err
  }
  if err = cmd.Start(); err != nil {
    return err
  }
  if _, err = fmt.Fprintf(stdin, "%s\n\n%s\n%s\n\n", e.FeedTitle, e.Title, e.Url); err != nil {
    stdin.Close()
    return err
  }
  if _, err = fmt.Fprint(stdin, e.Body); err != nil {
    stdin.Close()
    return err
  }
  if err = stdin.Close(); err != nil {
    return err
  }
  return cmd.Wait()
}
