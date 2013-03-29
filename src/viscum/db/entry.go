package db

// Defines the entry schema.
type Entry struct {
  Id        int64
  Url       string
  Title     string
  Body      string
  FeedId    int64
  FeedTitle string
  Email     string
}
