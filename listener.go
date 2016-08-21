package main

import (
	"log"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"github.com/sorcix/irc/ctcp"
)

func MessageHandler(s ircx.Sender, m *irc.Message) {
	msg := m.Trailing
	var command string
	if m.Params[0] == name {
		m.Params = []string{m.Name}
	}

	if strings.HasPrefix(strings.ToLower(msg), strings.ToLower(name)) {
		pieces := strings.Split(msg, " ")
		if len(pieces) >= 2 {
			command = pieces[1]
			runCommand(command, pieces, msg, s, m)
			return
		}
	} else if strings.HasPrefix(msg, prefix) {
		pieces := strings.Split(msg, " ")
		if len(pieces) >= 1 {
			command = pieces[0][1:]
			runCommand(command, pieces, msg, s, m)
			return
		}
	} else if strings.HasPrefix(msg, "\x01VERSION") {
		log.Println(ctcp.VersionReply())
		s.Send(&irc.Message{
			Command:  irc.PRIVMSG,
			Params:   m.Params,
			Trailing: ctcp.VersionReply(),
		})
		return
	} else {
		for k := range replies {
			if ok := k.FindAllString(msg, 1); ok != nil {
				s.Send(&irc.Message{
					Command:  irc.PRIVMSG,
					Params:   m.Params,
					Trailing: replies[k],
				})
				return
			}
		}
	}
}

func runCommand(command string, pieces []string, msg string, s ircx.Sender, m *irc.Message) {
	var params string
	if len(pieces) >= 2 {
		params = strings.SplitN(msg, command, 2)[1]
		params = strings.Trim(params, " ")
	} else {
		params = ""
	}
	if c, ok := coms[command]; ok {
		log.Printf("User %s issued command: %s\n", m.Name, command)
		c.Handle(s, m, params)
	}
}
