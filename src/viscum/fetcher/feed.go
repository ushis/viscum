package fetcher

import (
  "errors"
  "github.com/jteeuwen/go-pkg-xmlx"
  "viscum/db"
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
  doc := xmlx.New()

  if err := doc.LoadUri(self.Url, nil); err != nil {
    return err
  }

  if n := doc.SelectNode("http://www.w3.org/2005/Atom", "feed"); n != nil {
    return (&Atom{self}).process(n)
  }

  if n := doc.SelectNode("", "rss"); n != nil {
    return (&Rss{self}).process(n)
  }

  return errors.New("Unsupported Format")
}
