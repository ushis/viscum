package fetcher

import (
  "errors"
  "github.com/jteeuwen/go-pkg-xmlx"
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
  doc := xmlx.New()

  if err := doc.LoadUri(self.Url, nil); err != nil {
    return err
  }

  if n := doc.SelectNode("http://www.w3.org/2005/Atom", "feed"); n != nil {
    self.processAtom(n)
    return nil
  }

  if n := doc.SelectNode("", "rss"); n != nil {
    self.processRss(n)
    return nil
  }

  return errors.New("Unsupported Format")
}

//
func (self *Feed) processAtom(node *xmlx.Node) {
  ns := "http://www.w3.org/2005/Atom"
  self.Title = Titleize(node.S(ns, "title"))

  for _, entry := range node.SelectNodes(ns, "entry") {
    url := self.extractAtomUrl(ns, entry)

    if len(url) == 0 || self.entries[url] {
      continue
    }

    e := &db.Entry{
      Url:   url,
      Title: Titleize(entry.S(ns, "title")),
      Body:  entry.S(ns, "content"),
    }

    if len(e.Body) == 0 {
      e.Body = entry.S(ns, "summary")
    }

    if err := self.entryHandler(e); err != nil {
      Error("[Feed]", err)
      continue
    }

    self.entries[url] = true
  }
}

//
func (self *Feed) extractAtomUrl(ns string, node *xmlx.Node) string {
  links := node.SelectNodes(ns, "link")

  for _, link := range links {
    if link.As("", "rel") == "alternate" {
      return link.As("", "href")
    }
  }

  if len(links) > 0 {
    return links[0].As("", "href")
  }

  return ""
}

//
func (self *Feed) processRss(node *xmlx.Node) {
  ns := "*"
  ch := node.SelectNode(ns, "channel")

  if ch == nil {
    return
  }
  self.Title = Titleize(node.S(ns, "title"))

  for _, entry := range node.SelectNodes(ns, "item") {
    url := entry.S(ns, "link")

    if len(url) == 0 || self.entries[url] {
      continue
    }

    e := &db.Entry{
      Url:   url,
      Title: Titleize(entry.S(ns, "title")),
      Body:  entry.S(ns, "description"),
    }

    if err := self.entryHandler(e); err != nil {
      Error("[Feed]", err)
      continue
    }

    self.entries[url] = true
  }
}
