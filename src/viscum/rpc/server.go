package rpc

import (
  "net"
  "net/http"
  "net/rpc"
  "os"
  "viscum/db"
  . "viscum/util"
)

type Server struct {
  db     db.DB    // Database connection
  socket string   // Rpc socket
  Ctrl   chan int // Control channel
}

// Returns a new RPC Server
func New(database db.DB, sock string, mc chan<- int, fc chan<- int) *Server {
  if err := rpc.RegisterName("S", NewService(database, mc, fc)); err != nil {
    Fatal("[RPC]", err)
  }
  rpc.HandleHTTP()

  return &Server{db: database, socket: sock, Ctrl: make(chan int)}
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
