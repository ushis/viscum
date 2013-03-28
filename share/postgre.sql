--
-- viscum database init script.
--

--
-- SCHEMA
--

--
-- emails
--
DROP TABLE IF EXISTS emails CASCADE;

CREATE TABLE emails (
  id          serial  PRIMARY KEY,
  email       varchar UNIQUE NOT NULL,
  created_at  timestamp with time zone NOT NULL,
  updated_at  timestamp with time zone NOT NULL,
  CHECK (email ~ '\S')
);

--
-- feeds
--
DROP TABLE IF EXISTS feeds CASCADE;

CREATE TABLE feeds (
  id          serial  PRIMARY KEY,
  url         varchar NOT NULL,
  title       varchar,
  created_at  timestamp with time zone NOT NULL,
  updated_at  timestamp with time zone NOT NULL,
  CHECK (url ~ '\Ahttps?://')
);

--
-- entries
--
DROP TABLE IF EXISTS entries CASCADE;

CREATE TABLE entries (
  id          serial  PRIMARY KEY,
  feed_id     integer REFERENCES feeds (id),
  url         varchar NOT NULL,
  title       varchar,
  body        text,
  created_at  timestamp with time zone NOT NULL,
  updated_at  timestamp with time zone NOT NULL,
  CHECK (url ~ '\Ahttps?://'),
  UNIQUE (feed_id, url)
);

--
-- subscriptions
--
DROP TABLE IF EXISTS subscriptions CASCADE;

CREATE TABLE subscriptions (
  email_id    integer REFERENCES emails (id),
  feed_id     integer REFERENCES feeds (id),
  created_at  timestamp with time zone NOT NULL,
  updated_at  timestamp with time zone NOT NULL,
  UNIQUE (email_id, feed_id)
);

--
-- queue
--
DROP TABLE IF EXISTS queue CASCADE;

CREATE TABLE queue (
  id          serial  PRIMARY KEY,
  email_id    integer REFERENCES emails (id),
  entry_id    integer REFERENCES entries (id),
  pending     boolean DEFAULT false,
  created_at  timestamp with time zone NOT NULL,
  updated_at  timestamp with time zone NOT NULL,
  UNIQUE (email_id, entry_id)
);

--
-- TYPES
--

--
-- queue_entry
--
DROP TYPE IF EXISTS queue_entry CASCADE;

CREATE TYPE queue_entry AS (
  id          integer,
  url         varchar,
  title       varchar,
  body        text,
  email       varchar,
  feed_title  varchar
);

--
-- VIEWS
--

--
-- queue_entries
--
DROP VIEW IF EXISTS queue_entries;

CREATE VIEW queue_entries (id, url, title, body, email, feed_title) AS
  SELECT queue.id       AS id,
         entries.url    AS url,
         entries.title  AS title,
         entries.body   AS body,
         emails.email   AS email,
         feeds.title    AS feed_title
    FROM queue
      INNER JOIN emails ON emails.id = queue.email_id
      INNER JOIN entries ON entries.id = entry_id
      INNER JOIN feeds ON feeds.id = entries.feed_id
        WHERE queue.pending = false
          OR  now() - queue.updated_at > interval '1' hour;

--
-- queue_info
--
DROP VIEW IF EXISTS queue_info;

CREATE VIEW queue_info (info) AS
  SELECT emails.email || ' (' || COUNT(queue.id) || ')'
    FROM queue
      INNER JOIN emails ON emails.id = queue.email_id
      GROUP BY emails.id;

--
-- subscripts
--
DROP VIEW IF EXISTS subscripts;

CREATE VIEW subscripts (email, url) AS
  SELECT emails.email, feeds.url
    FROM subscriptions AS s
      INNER JOIN emails ON emails.id = s.email_id
      INNER JOIN feeds ON feeds.id = s.feed_id;

--
-- TRIGGER
--

--
-- Updates created_at and updated_at.
--
DROP FUNCTION IF EXISTS update_timestamps();

CREATE FUNCTION update_timestamps() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();

  IF (TG_OP = 'INSERT') THEN
    NEW.created_at = NEW.updated_at;
  END IF;

  RETURN NEW;
END
$$ LANGUAGE plpgsql;

--
--
--
DROP TRIGGER IF EXISTS update_timestamps_trig ON emails;

CREATE TRIGGER update_timestamps_trig BEFORE INSERT OR UPDATE ON emails
FOR EACH ROW EXECUTE PROCEDURE update_timestamps();

--
--
--
DROP TRIGGER IF EXISTS update_timestamps_trig ON feeds;

CREATE TRIGGER update_timestamps_trig BEFORE INSERT OR UPDATE ON feeds
FOR EACH ROW EXECUTE PROCEDURE update_timestamps();

--
--
--
DROP TRIGGER IF EXISTS update_timestamps_trig ON entries;

CREATE TRIGGER update_timestamps_trig BEFORE INSERT OR UPDATE ON entries
FOR EACH ROW EXECUTE PROCEDURE update_timestamps();

--
--
--
DROP TRIGGER IF EXISTS update_timestamps_trig ON queue;

CREATE TRIGGER update_timestamps_trig BEFORE INSERT OR UPDATE ON queue
FOR EACH ROW EXECUTE PROCEDURE update_timestamps();

--
--
--
DROP TRIGGER IF EXISTS update_timestamps_trig ON subscriptions;

CREATE TRIGGER update_timestamps_trig BEFORE INSERT OR UPDATE ON subscriptions
FOR EACH ROW EXECUTE PROCEDURE update_timestamps();

--
-- enqueue()
--
DROP FUNCTION IF EXISTS enqueue();

CREATE FUNCTION enqueue() RETURNS trigger AS $$
DECLARE
  s RECORD;
BEGIN
  FOR s IN
    SELECT email_id FROM subscriptions
      WHERE feed_id = NEW.feed_id
        AND now() - created_at  > interval '1' minute
    LOOP
    INSERT INTO queue (email_id, entry_id) VALUES (s.email_id, NEW.id);
  END LOOP;

  RETURN NEW;
END
$$ LANGUAGE plpgsql;

--
-- enqueue_entry_trig
--
DROP TRIGGER IF EXISTS enqueue_trig on entries;

CREATE TRIGGER enqueue_trig AFTER INSERT ON entries
FOR EACH ROW EXECUTE PROCEDURE enqueue();

--
-- FUNCTIONS
--

--
-- fetch_queue()
--
DROP FUNCTION IF EXISTS fetch_queue();

CREATE FUNCTION fetch_queue() RETURNS SETOF queue_entry AS $$
DECLARE
  entry queue_entry%rowtype;
BEGIN
  FOR entry IN SELECT * FROM queue_entries LOOP
    UPDATE queue SET pending = true WHERE id = entry.id;
    RETURN NEXT entry;
  END LOOP;
  RETURN;
END
$$ LANGUAGE plpgsql;

--
-- dequeue()
--
DROP FUNCTION IF EXISTS dequeue(integer, boolean);

CREATE FUNCTION dequeue(id_var integer, processed boolean) RETURNS void AS $$
BEGIN
  IF processed THEN
    DELETE FROM queue WHERE id = id_var;
  ELSE
    UPDATE queue SET pending = false WHERE id = id_var;
  END IF;
END
$$ LANGUAGE plpgsql;

--
-- subscribe(email, url)
--
DROP FUNCTION IF EXISTS subscribe(varchar, varchar);

CREATE FUNCTION subscribe(email_var varchar, url_var varchar) RETURNS void AS $$
DECLARE
  email_id      integer;
  feed_id       integer;
  row_count     integer;
BEGIN
  SELECT id INTO email_id FROM emails WHERE email = email_var LIMIT 1;

  IF NOT FOUND THEN
    INSERT INTO emails (email) VALUES (email_var) RETURNING id into email_id;
  END IF;

  SELECT id INTO feed_id FROM feeds WHERE url LIKE url_var LIMIT 1;

  IF NOT FOUND THEN
    INSERT INTO feeds (url) VALUES (url_var) RETURNING id into feed_id;
  END IF;

  INSERT INTO subscriptions (email_id, feed_id) VALUES (email_id, feed_id);
END
$$ LANGUAGE plpgsql;

--
-- unsubscribe(email, url)
--
DROP FUNCTION IF EXISTS unsubscribe(varchar, varchar);

CREATE FUNCTION unsubscribe(email_var varchar, url_var varchar) RETURNS void AS $$
BEGIN
  DELETE FROM subscriptions USING emails, feeds
    WHERE subscriptions.email_id = emails.id
      AND subscriptions.feed_id = feeds.id
      AND emails.email = email_var
      AND feeds.url LIKE url_var;
END
$$ LANGUAGE plpgsql;

--
--
--
DROP FUNCTION IF EXISTS insert_entry(varchar, varchar, text, integer, varchar);

CREATE FUNCTION insert_entry(
  url_var        varchar,
  title_var      varchar,
  body_var       text,
  feed_id_var    integer,
  feed_title_var varchar
)
RETURNS void AS $$
BEGIN
  PERFORM 1 FROM entries WHERE feed_id = feed_id_var AND url = url_var LIMIT 1;

  IF NOT FOUND THEN
    UPDATE feeds SET title = feed_title_var WHERE id = feed_id_var;

    INSERT INTO entries (feed_id, url, title, body)
    VALUES (feed_id_var, url_var, title_var, body_var);
  END IF;
END
$$ LANGUAGE plpgsql;
