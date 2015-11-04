package commands

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Topic(s ircx.Sender, m *irc.Message, message string) {
	if !isTrusted(s, m.Name) {
		return
	}
	s.Send(&irc.Message{
		Command:  irc.TOPIC,
		Params:   m.Params,
		Trailing: message,
	})
}
