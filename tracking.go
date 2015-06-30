package main

import (
	"fmt"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

var acc = make(map[string]bool)

func JoinHandler(s ircx.Sender, m *irc.Message) {
	fmt.Printf("JOIN %+v\n", m)
	if m.Name != name {
		if _, ok := trusted[m.Name]; ok {
			CheckACCFor(s, m.Name)
		} else if _, ok := owners[m.Name]; ok {
			CheckACCFor(s, m.Name)
		} else if _, ok := idents[m.Name]; !ok {
			CheckACCFor(s, m.Name)
		}
	}
}

func NamesHandler(s ircx.Sender, m *irc.Message) {
	fmt.Printf("NAMES: %s\n", m.Trailing)
	users := strings.Split(m.Trailing, " ")
	for i := 0; i < len(users); i++ {
		if users[i] != name {
			if strings.HasPrefix(users[i], "@") || strings.HasPrefix(users[i], "+") || strings.HasPrefix(users[i], "!") {
				users[i] = users[i][1:]
			}
			acc[users[i]] = true
			if len(acc) == 1 {
				CheckACC(s)
			}
		}
	}
}

func CheckACC(s ircx.Sender) {
	for user := range acc {
		if acc[user] {
			acc[user] = false
			CheckACCFor(s, user)
			break
		} else {
			delete(acc, user)
		}
	}
}

func CheckACCFor(s ircx.Sender, user string) {
	if len(user) > 0 {
		s.Send(&irc.Message{
			Command:  irc.PRIVMSG,
			Params:   []string{"NICKSERV"},
			Trailing: "ACC " + user,
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

func NoticeHandler(s ircx.Sender, m *irc.Message) {
	fmt.Printf("NOTICE %+v\n", m)
	if m.Name == "NickServ" {
		pieces := strings.Split(m.Trailing, " ")
		if pieces[1] == "ACC" {
			var user string
			if pieces[2] == "3" {
				user = pieces[0]
				if _, ok := trusted[pieces[0]]; ok {
					trusted[user] = true
				}
				if _, ok := owners[pieces[0]]; ok {
					owners[user] = true
				}
				if _, ok := idents[pieces[0]]; !ok {
					idents[user] = true
				}
			}
			if len(acc) > 0 {
				CheckACC(s)
			}
		}
	}
}
