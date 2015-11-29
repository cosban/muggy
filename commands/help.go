package commands

import (
	"fmt"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Help(s ircx.Sender, m *irc.Message, message string) {
	if !isTrusted(s, m.Name) {
		return
	}
	args := strings.Split(message, " ")
	if len(args[0]) < 1 {
		s.Send(&irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: "-------- HELP --------",
		})
		for _, v := range coms {
			s.Send(&irc.Message{
				Command:  irc.NOTICE,
				Params:   []string{m.Name},
				Trailing: fmt.Sprintf("%s %s -- %s", v.Name, v.Usage, v.Summary),
			})
		}
		s.Send(&irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: "------ END HELP ------",
		})
	} else {
		if v, ok := coms[args[0]]; ok {
			s.Send(&irc.Message{
				Command:  irc.NOTICE,
				Params:   []string{m.Name},
				Trailing: fmt.Sprintf("%s %s -- %s", v.Name, v.Usage, v.Summary),
			})
		}
	}
}
