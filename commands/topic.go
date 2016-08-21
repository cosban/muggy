package commands

import (
	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Topic(s ircx.Sender, m *irc.Message, message string) {
	if !isTrusted(s, m.Name) {
		return
	}
	messages.QueueMessages(s, &irc.Message{
		Command:  irc.TOPIC,
		Params:   m.Params,
		Trailing: message,
	})
}
