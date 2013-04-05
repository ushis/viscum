package db

import (
  "database/sql"
  "time"
  . "viscum/util"
)

// Database interface.
type DB interface {

  // Opens the connection.
  Open(auth string) (err error)

  // Closes the connection.
  Close() (err error)

  // Adds a new subscription.
  Subscribe(email string, url string) (r sql.Result, err error)

  // Removes a subscription.
  Unsubscribe(email string, url string) (r sql.Result, err error)

  // Lists subscriptions filtered by email.
  ListSubscriptions(email string) (info string, err error)

  // Updates a feed.
  UpdateFeed(f *Feed) (r sql.Result, err error)

  // Adds a new entry, checks the subscriptions and enqueues them.
  InsertEntry(feedId int64, e *Entry) (r sql.Result, err error)

  // Dequeues an entry. It removes the entry from queue, if processed is true.
  // If it is false it removes the pending flag from the entry.
  Dequeue(e *QueueEntry, processed bool) (r sql.Result, err error)

  // Returns an array of info strings about all queue entries.
  QueueInfo() (info string, err error)

  // Fetches new feeds and passes them to a handler function.
  FetchNewFeeds(age time.Time, handler func(feed *Feed)) (err error)

  // Fetches queue entries and passes them to a handler function.
  FetchQueue(handler func(entry *QueueEntry)) (err error)
}

// Map of all registered databases.
var databases = make(map[string]DB)

// Registers a database.
func Register(name string, db DB) {
  if _, dup := databases[name]; dup {
    Fatal("[DB] Database registered twice:", name)
  }
  databases[name] = db
}

// Returns a new database.
func New(name string) DB {
  db, ok := databases[name]

  if !ok {
    Fatal("[DB] Couldn't find database:", name)
  }
  return db
}
