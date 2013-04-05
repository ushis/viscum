package db

import (
  "bytes"
  "database/sql"
  _ "github.com/jbarham/gopgsqldriver"
  "reflect"
  "time"
  . "viscum/util"
)

const (
  // Postgres timestamp layout.
  //
  // See http://golang.org/src/pkg/time/format.go#L9
  PG_TIME_FMT = "2006-01-02 15:04:05.000000Z0700"
)

// Postgres database.
type PgDB struct {
  *sql.DB
}

// Register init function.
func init() {
  Register("postgres", &PgDB{})
}

// Opens the connection using the postgres driver.
func (self *PgDB) Open(auth string) (err error) {
  self.DB, err = sql.Open("postgres", auth)
  return err
}

// Inserts a new subscription.
func (self *PgDB) Subscribe(email string, url string) (sql.Result, error) {
  return self.Exec("SELECT subscribe($1, $2)", email, url)
}

// Removes a subscription.
func (self *PgDB) Unsubscribe(email string, url string) (sql.Result, error) {
  return self.Exec("SELECT unsubscribe($1, $2)", email, url)
}

//
func (self *PgDB) UpdateFeed(f *Feed) (sql.Result, error) {
  return self.Exec("UPDATE feeds SET sha1 = $1, title = $2 WHERE id = $3",
    f.Sha1, f.Title, f.Id)
}

// Inserts a new entry.
func (self *PgDB) InsertEntry(feedId int64, e *Entry) (sql.Result, error) {
  return self.Exec("SELECT insert_entry($1, $2, $3, $4, $5)",
    feedId, e.Sha1, e.Url, e.Title, e.Body)
}

// Dequeues an entry.
func (self *PgDB) Dequeue(e *QueueEntry, success bool) (sql.Result, error) {
  return self.Exec("SELECT dequeue($1, $2)", e.Id, success)
}

//
func (self *PgDB) ListSubscriptions(email string) (string, error) {
  s, i, err := self.info("SELECT url from subscripts WHERE email = $1", email)

  if err == nil && i == 0 {
    return "Couldn't find any subscriptions for: " + email, nil
  }
  return s, err
}

//
func (self *PgDB) QueueInfo() (string, error) {
  s, i, err := self.info("SELECT info FROM queue_info")

  if err == nil && i == 0 {
    return "The queue is empty.", nil
  }
  return s, err
}

func (self *PgDB) info(q string, args ...interface{}) (string, int, error) {
  rows, err := self.Query(q, args...)

  if err != nil {
    return "", 0, err
  }
  defer rows.Close()

  var buffer bytes.Buffer
  i := 0

  for rows.Next() {
    var info string

    if err := rows.Scan(&info); err != nil {
      Error("[DB]", err)
      continue
    }

    if i > 0 {
      buffer.WriteByte('\n')
    }

    buffer.WriteString(info)
    i++
  }

  return buffer.String(), i, rows.Err()
}

//
func (self *PgDB) FetchNewFeeds(t time.Time, handler func(*Feed)) error {
  rows, err := self.Query(
    "SELECT id, url, sha1, title FROM feeds WHERE created_at > $1",
    t.Format(PG_TIME_FMT))

  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    var f Feed
    var sha1, title interface{}

    if err := rows.Scan(&f.Id, &f.Url, &sha1, &title); err != nil {
      Error("[DB]", err)
      continue
    }
    if sha1 != nil {
      f.Sha1 = reflect.ValueOf(sha1).String()
    }
    if title != nil {
      f.Title = reflect.ValueOf(title).String()
    }
    go handler(&f)
  }
  return rows.Err()
}

func (self *PgDB) FetchQueue(handler func(*QueueEntry)) error {
  rows, err := self.Query("SELECT id, url, title, body, email, feed_title FROM fetch_queue()")

  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    e := &QueueEntry{Entry: new(Entry)}
    err := rows.Scan(&e.Id, &e.Url, &e.Title, &e.Body, &e.Email, &e.FeedTitle)

    if err != nil {
      Error("[DB]", err)
      continue
    }

    go handler(e)
  }
  return rows.Err()
}
