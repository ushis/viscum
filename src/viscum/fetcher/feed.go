package fetcher

import (
  "errors"
  "github.com/jteeuwen/go-pkg-xmlx"
  "viscum/db"
  . "viscum/util"
)

const (
  NS_ATOM = "http://www.w3.org/2005/Atom" // ATOM namespace
  NS_RSS  = "*"                           // RSS namespace
)

type Feed struct {
  *db.Feed
  entries      map[string]bool       // Set of already processed entries
  entryHandler func(*db.Entry) error // Handles new entries.
}

// Returns a new feed.
func NewFeed(f *db.Feed, h func(*db.Entry) error) *Feed {
  return &Feed{Feed: f, entries: make(map[string]bool), entryHandler: h}
}

// Fetches the feed and sends new entries to the handler.
func (self *Feed) Fetch() error {
  doc := xmlx.New()

  if err := doc.LoadUri(self.Url, nil); err != nil {
    return err
  }

  if n := doc.SelectNode(NS_ATOM, "feed"); n != nil {
    self.processAtom(n)
    return nil
  }

  if n := doc.SelectNode("", "rss"); n != nil {
    self.processRss(n)
    return nil
  }
  return errors.New("Unsupported Format")
}

// Sends new entries to the handler and registers the url.
func (self *Feed) registerEntry(url string, title string, body string) {
  e := &db.Entry{Url: url, Title: Titleize(title), Body: body}

  if err := self.entryHandler(e); err != nil {
    Error("[Feed]", err)
  } else {
    self.entries[e.Url] = true
  }
}

// Processes ATOM feeds.
func (self *Feed) processAtom(node *xmlx.Node) {
  self.Title = Titleize(node.S(NS_ATOM, "title"))

  for _, entry := range node.SelectNodes(NS_ATOM, "entry") {
    url := self.extractAtomUrl(entry)

    if len(url) == 0 || self.entries[url] {
      continue
    }
    body := entry.S(NS_ATOM, "content")

    if len(body) == 0 {
      body = entry.S(NS_ATOM, "summary")
    }
    self.registerEntry(url, entry.S(NS_ATOM, "title"), body)
  }
}

// Extracts the url from ATOM feed items.
func (self *Feed) extractAtomUrl(node *xmlx.Node) string {
  links := node.SelectNodes(NS_ATOM, "link")

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

// Processes RSS feeds.
func (self *Feed) processRss(node *xmlx.Node) {
  ch := node.SelectNode(NS_RSS, "channel")

  if ch == nil {
    return
  }
  self.Title = Titleize(node.S(NS_RSS, "title"))

  for _, entry := range node.SelectNodes(NS_RSS, "item") {
    url := entry.S(NS_RSS, "link")

    if len(url) > 0 && !self.entries[url] {
      self.registerEntry(url, entry.S(NS_RSS, "title"),
        entry.S(NS_RSS, "description"))
    }
  }
}
