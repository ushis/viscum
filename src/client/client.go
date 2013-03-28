package main

import (
  "flag"
  "fmt"
  "os"
  "viscum/rpc"
  "viscum/util"
)

// Path to the config file.
var configFile string

// Init the flag parser.
func init() {
  flag.Usage = usage
  flag.StringVar(&configFile, "config", util.CONF_CLIENT, "Config file")
}

// Commands we handle.
var commands = map[string]func(*rpc.Client) (rpc.Reply, error){
  "add":     subscribe,
  "rm":      unsubscribe,
  "ls":      subscriptions,
  "queue":   queue,
  "deliver": deliver,
  "mem":     mem,
}

// Lets go.
func main() {
  flag.Parse()

  // Read the config.
  conf, err := util.ReadConfig(configFile)
  if err != nil {
    util.Fatal(err)
  }

  if flag.NArg() < 1 {
    usage()
  }

  // Find the correct function for the command.
  cmd, ok := commands[flag.Arg(0)]
  if !ok {
    usage()
  }

  // Connect to the server.
  client := rpc.NewClient(conf.Get("rpc", "socket"))

  if err := client.Connect(); err != nil {
    util.Fatal(err)
  }
  defer client.Disconnect()

  // Execute the command.
  if reply, err := cmd(client); err != nil {
    util.Fatal(err)
  } else {
    fmt.Println(reply.Text())
  }
}

// Subscribes an email to a feed.
func subscribe(client *rpc.Client) (rpc.Reply, error) {
  if flag.NArg() < 3 {
    usage()
  }
  return client.Subscribe(flag.Arg(1), flag.Arg(2))
}

// Unsubscribes an email from a feed.
func unsubscribe(client *rpc.Client) (rpc.Reply, error) {
  if flag.NArg() < 3 {
    usage()
  }
  return client.Unsubscribe(flag.Arg(1), flag.Arg(2))
}

// Lists all subscription filtered by email.
func subscriptions(client *rpc.Client) (rpc.Reply, error) {
  if flag.NArg() < 2 {
    usage()
  }
  return client.SubscriptionInfo(flag.Arg(1))
}

// Fetches queue info.
func queue(client *rpc.Client) (rpc.Reply, error) {
  return client.QueueInfo()
}

// Attempt to process the queue.
func deliver(client *rpc.Client) (rpc.Reply, error) {
  return client.DeliverQueue()
}

// Fetches the servers mem stats.
func mem(client *rpc.Client) (rpc.Reply, error) {
  return client.MemStats()
}

// Prints the help message and exits.
func usage() {
  fmt.Fprintf(os.Stderr, "Usage: %s <cmd> [options]\n", os.Args[0])

  fmt.Fprintln(os.Stderr, `
Commands:
  add <email> <url>   Subscribe an email to a feed
  rm  <email> <url>   Unsubscribe an email from a feed
  ls  <email>         List subscriptions filtered by email
  queue               Display queue info
  deliver             Attempt to deliver the whole queue
  mem                 Display the servers memory stats

Options:`)

  flag.PrintDefaults()
  os.Exit(1)
}
