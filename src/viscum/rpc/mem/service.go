package mem

import (
  "fmt"
  "runtime"
)

type Service struct{}

// Returns a new service.
func New() (string, *Service) {
  return "Mem", new(Service)
}

// Serves memory stats.
func (self *Service) Stats(_ *Args, reply *Reply) error {
  m := new(runtime.MemStats)
  runtime.ReadMemStats(m)
  reply.Reply = fmt.Sprintf("Sys: %d Heap: %d", m.Sys, m.HeapAlloc)
  return nil
}
