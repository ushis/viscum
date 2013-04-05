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
  urls := make(map[string]bool)
  ch := node.SelectNode(ns, "channel")

  if ch == nil {
    return errors.New("No channels found.")
  }
  self.Title = node.S(ns, "title")

  for _, entry := range node.SelectNodes(ns, "item") {
    url := entry.S(ns, "link")

    if len(url) == 0 || self.entries[url] {
      urls[url] = true
      continue
    }

    e := &db.Entry{
      Url:   url,
      Title: entry.S(ns, "title"),
      Body:  entry.S(ns, "description"),
    }

    if err := self.entryHandler(e); err != nil {
      Error("[RSS]", err)
      continue
    }

    urls[url] = true
  }

  self.entries = urls
  return nil
}
