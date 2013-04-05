package fetcher

import (
  "github.com/jteeuwen/go-pkg-xmlx"
  "viscum/db"
  . "viscum/util"
)

type Atom struct {
  *Feed
}

func (self *Atom) process(node *xmlx.Node) error {
  ns := "http://www.w3.org/2005/Atom"
  urls := make(map[string]bool)

  self.Title = node.S(ns, "title")

  for _, entry := range node.SelectNodes(ns, "entry") {
    var url string

    if l := entry.SelectNode(ns, "link"); l != nil {
      url = l.As("", "href")
    }

    if len(url) == 0 || self.entries[url] {
      urls[url] = true
      continue
    }

    e := &db.Entry{
      Url:   url,
      Title: entry.S(ns, "title"),
      Body:  entry.S(ns, "content"),
    }

    if len(e.Body) == 0 {
      e.Body = entry.S(ns, "summary")
    }

    if err := self.entryHandler(e); err != nil {
      Error("[Atom]", err)
      continue
    }

    urls[url] = true
  }

  self.entries = urls
  return nil
}
