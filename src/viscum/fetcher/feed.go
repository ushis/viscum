package fetcher

import (
  "errors"
  "github.com/jteeuwen/go-pkg-xmlx"
  "io/ioutil"
  "net/http"
  "viscum/db"
  . "viscum/util"
)

type Feed struct {
  *db.Feed
  entries      map[string]bool
  entryHandler func(*db.Entry) error
}

//
func NewFeed(f *db.Feed, h func(*db.Entry) error) *Feed {
  return &Feed{Feed: f, entries: make(map[string]bool), entryHandler: h}
}

//
func (self *Feed) Fetch() error {
  resp, err := http.Get(self.Url)

  if err != nil {
    return err
  }
  defer resp.Body.Close()

  content, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    return err
  }
  sum, err := Sha1Sum(content)

  if err != nil {
    return err
  }

  if sum == self.Sha1 {
    return nil
  }
  self.Sha1 = sum

  doc := xmlx.New()

  if err := doc.LoadBytes(content, nil); err != nil {
    return err
  }
  return self.Process(doc)
}

//
func (self *Feed) Process(doc *xmlx.Document) error {
  if n := doc.SelectNode("http://www.w3.org/2005/Atom", "feed"); n != nil {
    return (&Atom{self}).process(n)
  }
  if n := doc.SelectNode("", "rss"); n != nil {
    return (&Rss{self}).process(n)
  }
  return errors.New("Unsupported Format")
}
