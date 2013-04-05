package rpc

import (
  "net"
  "net/http"
  "net/rpc"
  "os"
  "viscum/db"
  "viscum/rpc/mem"
  "viscum/rpc/queue"
  "viscum/rpc/subscription"
  . "viscum/util"
)

type Server struct {
  db          db.DB      // Database connection
  socket      string     // Rpc socket
  Ctrl        chan int   // Control channel
  MailerCtrl  chan<- int // Control channel to the mailer
  FetcherCtrl chan<- int // Control channel to the fetcher
}

// Returns a new RPC Server
func New(database db.DB, sock string, mc chan<- int, fc chan<- int) *Server {
  registerService(mem.New())
  registerService(queue.New(database, mc))
  registerService(subscription.New(database, fc))
  rpc.HandleHTTP()

  return &Server{
    db:          database,
    socket:      sock,
    Ctrl:        make(chan int),
    MailerCtrl:  mc,
    FetcherCtrl: fc,
  }
}

// Starts the rpc server.
func (self *Server) Start() {
  Info("[RPC] Start...")

  listener, err := self.listen()

  if err != nil {
    Fatal("[RPC]", err)
  }
  defer listener.Close()

  go http.Serve(listener, nil)

  // Wait
  <-self.Ctrl
  Info("[RPC] Stop...")
}

// Commands the rpc to stop.
func (self *Server) Stop() {
  self.Ctrl <- CTRL_STOP
}

// Returns a new socket listener.
func (self *Server) listen() (net.Listener, error) {
  if len(self.socket) > 0 && self.socket[0] != '/' {
    return net.Listen("tcp", self.socket)
  }
  listener, err := net.Listen("unix", self.socket)

  if err != nil {
    return listener, err
  }
  return listener, os.Chmod(self.socket, 0770)
}

// Registers a service.
func registerService(name string, service interface{}) {
  if err := rpc.RegisterName(name, service); err != nil {
    Fatal(err)
  }
}
