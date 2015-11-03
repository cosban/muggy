package commands

import (
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func AddAlias(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
	args := strings.Split(message, " ")
	if len(args) == 2 {
		if v, ok := coms[args[0]]; ok {
			coms[args[1]] = v
			s.Send(&irc.Message{
				Command:  irc.NOTICE,
				Params:   []string{m.Name},
				Trailing: args[1] + " is now an alias of " + args[0],
			})
		}
	}
}
