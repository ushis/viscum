package mailer

import (
  "fmt"
  "io"
  "net/textproto"
  "time"
  "viscum/db"
)

type Mail struct {
  *db.Entry                   // The entry to mail
  From      string            // Sender address
  headers   map[string]string // Mail headers
}

// Returns a new mail.
func NewMail(e *db.Entry, from string) *Mail {
  return &Mail{
    Entry: e,
    From:  from,
    headers: map[string]string{
      "Content-Type": "text/plain; charset=UTF-8",
      "Date":         time.Now().Format(time.RFC1123Z),
      "From":         from,
      "To":           e.Email,
      "Subject":      "[" + e.FeedTitle + "] " + e.Title,
    },
  }
}

// Sets a header. Overrides existing header with the same name.
func (self *Mail) SetHeader(k string, v string) {
  self.headers[textproto.CanonicalMIMEHeaderKey(k)] = v
}

// Returns a already set header.
//
// Returns an empty string, if the header is not already set.
func (self *Mail) GetHeader(k string) string {
  if v, ok := self.headers[textproto.CanonicalMIMEHeaderKey(k)]; ok {
    return v
  }
  return ""
}

// Writes headers to a io.Writer.
func (self *Mail) WriteHeaders(w io.Writer) error {
  for k, v := range self.headers {
    if _, err := w.Write([]byte(k + ": " + v + "\r\n")); err != nil {
      return err
    }
  }
  return nil
}

// Writes the body to a io.Writer.
func (self *Mail) WriteBody(w io.Writer) error {
  if _, err := fmt.Fprintf(w, "%s\n\n%s\n%s\n\n", self.FeedTitle, self.Title, self.Url); err != nil {
    return err
  }

  _, err := w.Write([]byte(self.Body))
  return err
}
