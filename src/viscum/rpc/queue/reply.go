package queue

type Reply struct {
  Reply string // Reponse text
}

// Returns the response text.
func (self *Reply) Text() string {
  return self.Reply
}
