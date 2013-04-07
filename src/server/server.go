package main

import (
  "flag"
  "os"
  "os/signal"
  "syscall"
  "viscum/db"
  "viscum/fetcher"
  "viscum/mailer"
  "viscum/rpc"
  "viscum/util"
)

// Path to the config file.
var configFile string

// Init flag parser.
func init() {
  flag.StringVar(&configFile, "config", util.CONF_SERVER, "config file")
}

// Sets up everything and waits for SIGINT.
func main() {
  flag.Parse()

  // Read the config.
  conf, err := util.ReadConfig(configFile)
  if err != nil {
    util.Fatal(err)
  }

  // Connect to the database.
  database := db.New(conf.Get("database", "driver"))

  if err = database.Open(conf.Get("database", "auth")); err != nil {
    util.Fatal(err)
  }
  defer database.Close()

  // Set up the mailer.
  m := mailer.New(database, conf)
  go m.Start()

  // Set up the fetcher.
  f := fetcher.New(database, conf, m.Ctrl)
  go f.Start()

  // Set up the rpc.
  r := rpc.New(database, conf.Get("rpc", "socket"), m.Ctrl, f.Ctrl)
  go r.Start()

  // Listen to SIGINT
  sig := make(chan os.Signal)
  signal.Notify(sig, syscall.SIGINT)

  // Wait for SIGINT and notify goroutines.
  <-sig
  r.Stop()
  f.Stop()
  m.Stop()
}
