package commands

import (
	"fmt"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Raw(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
	fmt.Printf("RAW: %s\n%+v", message, irc.ParseMessage(message))
	s.Send(irc.ParseMessage(message))
}
