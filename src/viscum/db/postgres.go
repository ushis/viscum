package db

import (
  "database/sql"
  "fmt"
  _ "github.com/ushis/gopgsqldriver"
  "io"
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

// Function that handles rows fetched by PgDB.query().
type RowHandler func(*sql.Rows) error

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

// Updates a feed.
func (self *PgDB) UpdateFeed(f *Feed) (sql.Result, error) {
  return self.Exec("UPDATE feeds SET title = $1 WHERE id = $2", f.Title, f.Id)
}

// Inserts a new entry.
func (self *PgDB) InsertEntry(i int64, e *Entry) (sql.Result, error) {
  return self.Exec("SELECT insert_entry($1, $2, $3, $4)", i, e.Url, e.Title, e.Body)
}

// Dequeues an entry.
func (self *PgDB) Dequeue(e *QueueEntry, success bool) (sql.Result, error) {
  return self.Exec("SELECT dequeue($1, $2)", e.Id, success)
}

// Lists all subscriptions filtered by email.
func (self *PgDB) ListSubscriptions(w io.Writer, email string) error {
  i, err := self.info(w, "SELECT url from subscripts WHERE email = $1", email)

  if err == nil && i == 0 {
    _, err = fmt.Fprint(w, "Couldn't find any subscriptions for: ", email)
  }
  return err
}

// Fetches queue info from the database.
func (self *PgDB) QueueInfo(w io.Writer) error {
  i, err := self.info(w, "SELECT info FROM queue_info")

  if err == nil && i == 0 {
    _, err = fmt.Fprint(w, "The queue is empty.")
  }
  return err
}

// Fetches feeds newer than the provided timestamp and passes them to handler
// in separate goroutines.
func (self *PgDB) FetchNewFeeds(t time.Time, handler func(*Feed)) error {
  return self.query(func(r *sql.Rows) error {
    var title interface{}

    f := new(Feed)

    if err := r.Scan(&f.Id, &f.Url, &title); err != nil {
      return err
    }

    if title != nil {
      f.Title = reflect.ValueOf(title).String()
    }

    go handler(f)

    return nil
  }, "SELECT id, url, title FROM feeds WHERE created_at > $1", t.Format(PG_TIME_FMT))
}

// Fetches unprocessed queue entries and passes them to the handler in
// separate goroutines.
func (self *PgDB) FetchQueue(handler func(*QueueEntry)) error {
  return self.query(func(r *sql.Rows) error {
    var title, fTitle interface{}

    e := &QueueEntry{Entry: new(Entry)}

    if err := r.Scan(&e.Id, &e.Url, &title, &e.Body, &e.Email, &fTitle); err != nil {
      return err
    }

    if title != nil {
      e.Title = reflect.ValueOf(title).String()
    }

    if fTitle != nil {
      e.FeedTitle = reflect.ValueOf(fTitle).String()
    }

    go handler(e)

    return nil
  }, "SELECT id, url, title, body, email, feed_title FROM fetch_queue()")
}

// Queries the database and executes the callback for each row.
func (self *PgDB) query(f RowHandler, q string, args ...interface{}) error {
  rows, err := self.Query(q, args...)

  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    if err := f(rows); err != nil {
      Error("[DB]", err)
    }
  }
  return rows.Err()
}

// Fetches info from the database and writes the results to a writer.
func (self *PgDB) info(w io.Writer, q string, args ...interface{}) (int, error) {
  i := 0

  err := self.query(func(r *sql.Rows) error {
    var info interface{}

    if err := r.Scan(&info); err != nil {
      return err
    }

    if i > 0 {
      if _, err := w.Write([]byte{'\n'}); err != nil {
        return err
      }
    }

    if _, err := fmt.Fprint(w, info); err != nil {
      return err
    }

    i++

    return nil
  }, q, args...)

  return i, err
}
