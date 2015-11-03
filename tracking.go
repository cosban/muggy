package main

import (
	"fmt"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func JoinHandler(s ircx.Sender, m *irc.Message) {
	if m.Name != name {
		userAdd(m.Name, m.Trailing)
	} else {
		fmt.Printf("%+v\n", m.Params)
		s.Send(&irc.Message{
			Command:  irc.WHO,
			Params:   []string{m.Params[0]},
			Trailing: "%na",
		})
	}
}

func AccountHandler(s ircx.Sender, m *irc.Message) {
	userAdd(m.Name, m.Params[0])
}

func WhoisHandler(s ircx.Sender, m *irc.Message) {
	userAdd(m.Params[1], m.Params[2])
}

func NickHandler(s ircx.Sender, m *irc.Message) {
	if v, ok := trusted[m.Name]; ok && v {
		trusted[m.Trailing] = true
		delete(trusted, m.Name)
	} else if v, ok := owners[m.Name]; ok && v {
		owners[m.Trailing] = true
		delete(owners, m.Name)
	} else if v, ok := idents[m.Name]; ok && v {
		idents[m.Trailing] = true
		delete(idents, m.Name)
	}
}

func LeaveHandler(s ircx.Sender, m *irc.Message) {
	if _, ok := trusted[m.Name]; ok {
		trusted[m.Name] = false
	} else if _, ok := owners[m.Name]; ok {
		owners[m.Name] = false
	} else if v, ok := idents[m.Name]; ok && v {
		delete(idents, m.Name)
	}
}

func userAdd(nick string, account string) {
	account = strings.ToLower(account)
	if _, ok := owners[account]; ok {
		owners[nick] = true
	}
	if _, ok := trusted[account]; ok {
		trusted[nick] = true
	}
	if _, ok := idents[account]; !ok {
		idents[nick] = true
	}
}
