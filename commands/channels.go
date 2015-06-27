package commands

import (
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Join(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(m.Name) {
		return
	}
	if strings.HasPrefix(message, "#") {
		channels := strings.Split(message, " ")
		for i := 0; i < len(channels); i++ {
			s.Send(&irc.Message{
				Command:  irc.JOIN,
				Params:   []string{channels[i]},
				Trailing: "",
			})
		}
		s.Send(&irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: "I have now joined the following channels: " + m.Trailing,
		})
	}
}

func Leave(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(m.Name) {
		return
	}
	s.Send(&irc.Message{
		Command:  irc.PART,
		Params:   m.Params,
		Trailing: "No one ever asks about Muggy!",
	})
}
