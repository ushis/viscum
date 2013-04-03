package fetcher

import (
  rss "github.com/jteeuwen/go-pkg-rss"
  "time"
  "viscum/db"
  . "viscum/util"
)

type Fetcher struct {
  db         db.DB      // Database connection
  Ctrl       chan int   // Control channel
  MailerCtrl chan<- int // Control channel to the mailer
}

// Returns a new fetcher.
func New(database db.DB, mc chan<- int) *Fetcher {
  return &Fetcher{db: database, Ctrl: make(chan int), MailerCtrl: mc}
}

// Starts fetching.
func (self *Fetcher) Start() {
  var lastUpdate time.Time
  Info("[Fetcher] Start...")

  for {
    err := self.db.FetchNewFeeds(lastUpdate, func(id int64, url string) {
      self.fetch(id, url)
    })

    if err != nil {
      Error("[Fetcher]", err)
    } else {
      lastUpdate = time.Now()
    }

    // Wait for instructions.
    if CTRL_STOP == <-self.Ctrl {
      break
    }
  }

  Info("[Fetcher] Stop...")
}

// Commands the fetcher to stop.
func (self *Fetcher) Stop() {
  self.Ctrl <- CTRL_STOP
}

// Starts fetching a new feed.
func (self *Fetcher) fetch(id int64, url string) {
  Info("[Fetcher] Start fetching:", url)

  f := rss.New(5, true, nil, func(_ *rss.Feed, c *rss.Channel, i []*rss.Item) {
    self.handleNewEntries(id, c, i)
  })

  for {
    if err := f.Fetch(url, nil); err != nil {
      Error("[Fetcher]", err)
    }
    <-time.After(time.Duration(f.SecondsTillUpdate() * 1e9))
  }
}

// Handles new entries.
func (self *Fetcher) handleNewEntries(id int64, ch *rss.Channel, items []*rss.Item) {
  for _, item := range items {
    entry := db.Entry{
      Title: Titleize(item.Title),
      FeedId: id,
      FeedTitle: Titleize(ch.Title),
    }

    // FIXME Find the correct link to the article.
    if len(item.Links) > 0 {
      entry.Url = item.Links[0].Href
    }

    if item.Content != nil {
      entry.Body = item.Content.Text
    }

    if len(entry.Body) == 0 {
      entry.Body = item.Description
    }

    if err := Format(&entry.Body); err != nil {
      Error("[Fetcher]", err)
      return
    }

    if _, err := self.db.InsertEntry(&entry); err != nil {
      Error("[Fetcher]", err)
      return
    }

    self.MailerCtrl <- CTRL_RELOAD
  }
}
