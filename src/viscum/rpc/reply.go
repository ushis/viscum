package rpc

type Reply struct {
  Reply []string // Response
}

// Returns the response.
func (self *Reply) Response() []string {
  return self.Reply
}

// Appends reponse text.
func (self *Reply) Append(s string) {
  self.Reply = append(self.Reply, s)
}
