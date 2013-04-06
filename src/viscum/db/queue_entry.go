package db

type QueueEntry struct {
  Id        int64
  Email     string
  FeedTitle string
  Body      []byte
  *Entry
}
