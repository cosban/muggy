package commands

import (
	"fmt"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func AddUser(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
	args := strings.Split(message, " ")
	if len(args[0]) < 1 {
		ListUsers(s, m, message)
		return
	}
	trusted[args[0]] = true
	s.Send(&irc.Message{
		Command:  irc.NOTICE,
		Params:   []string{m.Name},
		Trailing: "I now trust " + args[0],
	})
}

func RemoveUser(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
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

func ListUsers(s ircx.Sender, m *irc.Message, message string) {
	fmt.Printf("Owners: %+v\nTrusted: %+v\n", owners, trusted)
	if !isTrusted(s, m.Name) {
		return
	}
	tlist := ""
	olist := ""
	for k, v := range trusted {
		if v {
			tlist += k + ", "
		}
	}
	if len(tlist) > 0 {
		tlist = tlist[:len(tlist)-2]
	}
	for k, v := range owners {
		if v {
			olist += k + ", "
		}
	}
	if len(olist) > 0 {
		olist = olist[:len(olist)-2]
	}
	s.Send(&irc.Message{
		Command:  irc.NOTICE,
		Params:   []string{m.Name},
		Trailing: "I trust the following users: " + tlist,
	})
	s.Send(&irc.Message{
		Command:  irc.NOTICE,
		Params:   []string{m.Name},
		Trailing: "I obey the following users: " + olist,
	})
	fmt.Printf("I trust the following users: %s\nI obey the following users: %s\n", tlist, olist)
}
