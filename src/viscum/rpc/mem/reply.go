package mem

type Reply struct {
  Reply string // Response text
}

// Returns the response text.
func (self *Reply) Text() string {
  return self.Reply
}
