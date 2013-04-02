package db

import (
  "bytes"
  "database/sql"
  _ "github.com/jbarham/gopgsqldriver"
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

// Inserts a new entry.
func (self *PgDB) InsertEntry(e *Entry) (sql.Result, error) {
  return self.Exec("SELECT insert_entry($1, $2, $3, $4, $5)",
    e.Url, e.Title, e.Body, e.FeedId, e.FeedTitle)
}

// Dequeues an entry.
func (self *PgDB) Dequeue(e *Entry, success bool) (sql.Result, error) {
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
func (self *PgDB) FetchNewFeeds(t time.Time, handler func(int64, string)) error {
  rows, err := self.Query("SELECT id, url FROM feeds WHERE created_at > $1", t.Format(PG_TIME_FMT))

  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    var id int64
    var url string

    if err := rows.Scan(&id, &url); err != nil {
      Error("[DB]", err)
      continue
    }

    go handler(id, url)
  }
  return rows.Err()
}

func (self *PgDB) FetchQueue(handler func(*Entry)) error {
  rows, err := self.Query("SELECT id, url, title, body, email, feed_title FROM fetch_queue()")

  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    var e Entry
    err := rows.Scan(&e.Id, &e.Url, &e.Title, &e.Body, &e.Email, &e.FeedTitle)

    if err != nil {
      Error("[DB]", err)
      continue
    }

    go handler(&e)
  }
  return rows.Err()
}
