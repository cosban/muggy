package commands

import (
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func AddUser(s ircx.Sender, m *irc.Message, message string) {
	args := strings.Split(message, " ")
	trusted[args[0]] = true
	s.Send(&irc.Message{
		Command:  irc.NOTICE,
		Params:   []string{m.Name},
		Trailing: "I now trust " + args[0],
	})
}

func RemoveUser(s ircx.Sender, m *irc.Message, message string) {
	args := strings.Split(message, " ")
	if _, ok := trusted[args[0]]; ok {
		delete(trusted, args[0])
		s.Send(&irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: "I no longer trust " + args[0],
		})
	}
}
