package main

import (
  "flag"
  "fmt"
  "os"
  "viscum/rpc"
  "viscum/util"
)

var config string

func init() {
  flag.StringVar(&config, "config", util.CONF_FILE, "config file")
}

func main() {
  flag.Parse()

  if flag.NArg() < 2 {
    fmt.Println(config)
    os.Exit(1)
  }

  client := rpc.NewClient("/tmp/viscum.sock")

  if err := client.Connect(); err != nil {
    util.Fatal(err)
  }
  defer client.Disconnect()

  reply, err := client.Subscribe(flag.Arg(0), flag.Arg(1))

  if err != nil {
    util.Fatal(err)
  }

  fmt.Println(reply)
}
