package main

import (
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func JoinHandler(s ircx.Sender, m *irc.Message) {
	if m.Name == name {
		for k, _ := range trusted {
			s.Send(&irc.Message{
				Command:  irc.PRIVMSG,
				Params:   []string{"NICKSERV"},
				Trailing: "ACC " + k,
			})
		}
	} else if _, ok := trusted[m.Name]; ok {
		s.Send(&irc.Message{
			Command:  irc.PRIVMSG,
			Params:   []string{"NICKSERV"},
			Trailing: "ACC " + m.Name,
		})
	}
}

func NickHandler(s ircx.Sender, m *irc.Message) {
	if v, ok := trusted[m.Name]; ok && v {
		trusted[m.Trailing] = true
		delete(trusted, m.Name)
		if m.Name == owner {
			owner = m.Trailing
		}
	}
}

func LeaveHandler(s ircx.Sender, m *irc.Message) {
	if _, ok := trusted[m.Name]; ok {
		trusted[m.Name] = false
	}
}

func NoticeHandler(s ircx.Sender, m *irc.Message) {
	if m.Name == "NickServ" {
		pieces := strings.Split(m.Trailing, " ")
		if pieces[1] == "ACC" {
			trusted[pieces[0]] = (pieces[2] == "3")
		}
	}
}

func RegisterConnect(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{channels},
	})
}

func PingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}
