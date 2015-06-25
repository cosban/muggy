package commands

import (
	"os"
	"time"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Quit(s ircx.Sender, m *irc.Message, message string) {
	s.Send(&irc.Message{
		Command:  irc.QUIT,
		Params:   m.Params,
		Trailing: "No one ever asks about Muggy!",
	})

	time.Sleep(1000)
	os.Exit(1)
}
