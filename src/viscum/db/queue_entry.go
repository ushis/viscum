package db

// Defines the queue entry schema.
type QueueEntry struct {
  Id        int64
  Email     string
  Url       string
  Title     string
  Body      string
  FeedTitle string
}
