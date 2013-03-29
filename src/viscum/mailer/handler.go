package mailer

import (
  "viscum/util"
)

type Handler interface {
  Init(conf *util.Config) // Configures the handler
  Send(entry *Mail) error // Sends the mail.
}
