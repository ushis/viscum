package mailer

import (
  "html/template"
  "io"
  "net/textproto"
  "time"
  "viscum/db"
)

var tmpl = template.Must(template.New("mail").Parse(
  `{{.FeedTitle}}

{{.Title}}
{{.Url}}

{{.Body}}

Proudly delivered by viscum.
`))

type Mail struct {
  *db.Entry
  From    string
  headers map[string]string
}

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

func (self *Mail) SetHeader(k string, v string) {
  self.headers[textproto.CanonicalMIMEHeaderKey(k)] = v
}

func (self *Mail) GetHeader(k string) string {
  if v, ok := self.headers[textproto.CanonicalMIMEHeaderKey(k)]; ok {
    return v
  }
  return ""
}

func (self *Mail) WriteHeaders(w io.Writer) error {
  for k, v := range self.headers {
    if _, err := w.Write([]byte(k + ": " + v + "\r\n")); err != nil {
      return err
    }
  }
  return nil
}

func (self *Mail) WriteBody(w io.Writer) error {
  return tmpl.Execute(w, self)
}
