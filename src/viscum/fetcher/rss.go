package fetcher

import (
  "errors"
  "github.com/jteeuwen/go-pkg-xmlx"
  "viscum/db"
  . "viscum/util"
)

type Rss struct {
  *Feed
}

func (self *Rss) process(node *xmlx.Node) error {
  ns := "*"
  sums := make(map[string]bool)
  ch := node.SelectNode(ns, "channel")

  if ch == nil {
    return errors.New("No channels found.")
  }
  self.Title = node.S(ns, "title")

  for _, entry := range node.SelectNodes(ns, "item") {
    sum, err := Sha1Sum(entry.String())

    if err != nil {
      Error("[RSS]", err)
      continue
    }

    if self.entries[sum] {
      continue
    }

    e := &db.Entry{
      Sha1:  sum,
      Url:   entry.S(ns, "link"),
      Title: entry.S(ns, "title"),
      Body:  entry.S(ns, "description"),
    }

    if err := self.entryHandler(e); err != nil {
      Error("[RSS]", err)
      continue
    }

    sums[e.Sha1] = true
  }

  self.entries = sums
  return nil
}
