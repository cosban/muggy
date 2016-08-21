package messages

import "github.com/sorcix/irc"

type queuedMessage struct {
	Sender  irc.Sender
	Message *irc.Message
}

type queue []queuedMessage

// Remove removes an element from the queue and then returns the element it removed
func (q *queue) Remove(i int) queuedMessage {
	s := *q
	m := s[0]
	s = append(s[:i], s[i+1:]...)
	*q = s
	return m
}

// IsEmpty returns true if the queue has no elements
func (q *queue) IsEmpty() bool {
	return len(*q) == 0
}
