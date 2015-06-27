package main

import (
	"fmt"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func JoinHandler(s ircx.Sender, m *irc.Message) {
	if m.Name == name {
		for k, _ := range trusted {
			fmt.Printf("Sending for %s\n", k)
			s.Send(&irc.Message{
				Command:  irc.PRIVMSG,
				Params:   []string{"NICKSERV"},
				Trailing: "ACC " + k,
			})
		}
		for k, _ := range owners {
			fmt.Printf("Sending for %s\n", k)
			s.Send(&irc.Message{
				Command:  irc.PRIVMSG,
				Params:   []string{"NICKSERV"},
				Trailing: "ACC " + k,
			})
		}
	} else if _, ok := trusted[m.Name]; ok {
		fmt.Printf("Sending for %s\n", m.Name)
		s.Send(&irc.Message{
			Command:  irc.PRIVMSG,
			Params:   []string{"NICKSERV"},
			Trailing: "ACC " + m.Name,
		})
	} else if _, ok := owners[m.Name]; ok {
		fmt.Printf("Sending for %s\n", m.Name)
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
	} else if v, ok := owners[m.Name]; ok && v {
		owners[m.Trailing] = true
		delete(owners, m.Name)
	}
}
func LeaveHandler(s ircx.Sender, m *irc.Message) {
	if _, ok := trusted[m.Name]; ok {
		trusted[m.Name] = false
	} else if _, ok := owners[m.Name]; ok {
		owners[m.Name] = false
	}
}

func NoticeHandler(s ircx.Sender, m *irc.Message) {
	fmt.Printf("%s sent this notice: %s\n ", m.Name, m.Trailing)
	if m.Name == "NickServ" {
		pieces := strings.Split(m.Trailing, " ")
		if pieces[1] == "ACC" && pieces[2] == "3" {
			if _, ok := trusted[pieces[0]]; ok {
				trusted[pieces[0]] = true
			}
			if _, ok := owners[pieces[0]]; ok {
				owners[pieces[0]] = true
			}
		}
	}
}

func RegisterConnect(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   []string{"NICKSERV"},
		Trailing: "identify " + name + " " + password,
	})
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
