package main

import (
	"fmt"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func MessageHandler(s ircx.Sender, m *irc.Message) {
	msg := m.Trailing
	var command string
	if strings.HasPrefix(strings.ToLower(msg), strings.ToLower(name)) {
		pieces := strings.Split(msg, " ")
		if len(pieces) >= 2 {
			command = pieces[1]
			runCommand(command, pieces, msg, s, m)
		}
	} else if strings.HasPrefix(msg, prefix) {
		pieces := strings.Split(msg, " ")
		if len(pieces) >= 1 {
			command = pieces[0][1:]
			runCommand(command, pieces, msg, s, m)
		}
	} else {
		for k, _ := range replies {
			if ok := k.FindAllString(msg, 1); ok != nil {
				s.Send(&irc.Message{
					Command:  irc.PRIVMSG,
					Params:   m.Params,
					Trailing: replies[k],
				})
				break
			}
		}
	}
}

func runCommand(command string, pieces []string, msg string, s ircx.Sender, m *irc.Message) {
	fmt.Printf("User %s issued command: %s\n", m.Name, command)
	var params string
	if len(pieces) >= 2 {
		params = strings.SplitN(msg, command, 2)[1]
		params = strings.Trim(params, " ")
	} else {
		params = ""
	}
	if c, ok := coms[command]; ok {
		c(s, m, params)
	}
}
