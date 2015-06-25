package commands

import (
	"os"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Quit(s ircx.Sender, m *irc.Message, message string) {
	s.Send(&irc.Message{
		Command:  irc.QUIT,
		Params:   m.Params,
		Trailing: "No one ever asks about Muggy!",
	})
	os.Exit(1)
}
