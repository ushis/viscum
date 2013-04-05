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
  sums := make(map[string]bool)

  self.Title = node.S(ns, "title")

  for _, entry := range node.SelectNodes(ns, "entry") {
    sum, err := Sha1Sum(entry.String())

    if err != nil {
      Error("[Atom]", err)
      continue
    }

    if self.entries[sum] {
      continue
    }

    e := &db.Entry{
      Sha1:  sum,
      Title: entry.S(ns, "title"),
      Body:  entry.S(ns, "content"),
    }

    if l := entry.SelectNode(ns, "link"); l != nil {
      e.Url = l.As("", "href")
    }

    if len(e.Body) == 0 {
      e.Body = entry.S(ns, "summary")
    }

    if err := self.entryHandler(e); err != nil {
      Error("[Atom]", err)
      continue
    }

    sums[e.Sha1] = true
  }

  self.entries = sums
  return nil
}
