package fetcher

import (
  "time"
  "viscum/db"
  . "viscum/util"
)

type Fetcher struct {
  db         db.DB         // Database connection
  poll       time.Duration // Poll interval
  Ctrl       chan int      // Control channel
  MailerCtrl chan<- int    // Control channel to the mailer
}

// Returns a new fetcher.
func New(database db.DB, conf *Config, mc chan<- int) *Fetcher {
  return &Fetcher{
    db:         database,
    poll:       conf.GetDuration("feed", "poll"),
    Ctrl:       make(chan int),
    MailerCtrl: mc,
  }
}

// Starts fetching.
func (self *Fetcher) Start() {
  Info("[Fetcher] Start...")
  var lastUpdate time.Time

  for {
    err := self.db.FetchNewFeeds(lastUpdate, func(feed *db.Feed) {
      self.fetch(feed)
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
func (self *Fetcher) fetch(f *db.Feed) {
  Info("[Fetcher] Start fetching", f.Url, "every", self.poll)

  feed := NewFeed(f, func(entry *db.Entry) error {
    return self.handleEntry(f.Id, entry)
  })

  for {
    if err := feed.Fetch(); err != nil {
      Error("[Fetcher]", err)
    } else if _, err := self.db.UpdateFeed(feed.Feed); err != nil {
      Error("[Fetcher]", err)
    }
    <-time.After(self.poll)
  }
}

// Handles a new entry.
func (self *Fetcher) handleEntry(feedId int64, entry *db.Entry) (err error) {
  if _, err = self.db.InsertEntry(feedId, entry); err == nil {
    self.MailerCtrl <- CTRL_RELOAD
  }
  return
}
