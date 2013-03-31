package mailer

import (
  "os/exec"
  . "viscum/util"
)

type Pipe struct {
  cmd string // The command
}

// Registers the handler.
func init() {
  Register("pipe", &Pipe{})
}

// Configures the pipe.
func (self *Pipe) Init(conf *Config) {
  self.cmd = conf.Get("mail", "pipe")
}

// Sends the message.
func (self *Pipe) Send(mail *Mail) error {
  cmd := exec.Command(self.cmd, "-s", mail.GetHeader("Subject"), mail.Email)
  stdin, err := cmd.StdinPipe()

  if err != nil {
    return err
  }
  if err = cmd.Start(); err != nil {
    return err
  }
  if err = mail.WriteBody(stdin); err != nil {
    stdin.Close()
    return err
  }
  if err = stdin.Close(); err != nil {
    return err
  }
  return cmd.Wait()
}
