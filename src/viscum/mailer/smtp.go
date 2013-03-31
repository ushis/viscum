package mailer

import (
  "bytes"
  "net/smtp"
  . "viscum/util"
)

type Smtp struct {
  addr   string // Server address
  host   string // Server hostname
  user   string // Username
  passwd string // Password
}

// Registers the handler.
func init() {
  Register("smtp", &Smtp{})
}

// Configures the handler.
func (self *Smtp) Init(conf *Config) {
  self.host = conf.Get("mail", "smtp_host")
  self.addr = self.host + ":" + conf.Get("mail", "smtp_port")
  self.user = conf.Get("mail", "smtp_username")
  self.passwd = conf.Get("mail", "smtp_password")
}

// Sends a mail.
func (self *Smtp) Send(mail *Mail) error {
  var buf bytes.Buffer

  if err := mail.WriteHeaders(&buf); err != nil {
    return err
  }

  if err := mail.WriteBody(&buf); err != nil {
    return err
  }

  auth := smtp.PlainAuth("", self.user, self.passwd, self.host)
  return smtp.SendMail(self.addr, auth, mail.From, []string{mail.Email}, buf.Bytes())
}
