package commands

import (
	"log"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Raw(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
	log.Printf("RAW: %s\n%+v", message, irc.ParseMessage(message))
	messages.QueueMessages(s, irc.ParseMessage(message))
}
