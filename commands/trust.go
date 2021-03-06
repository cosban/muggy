package commands

import (
	"log"
	"strings"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Trust(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
	args := strings.Split(message, " ")
	if len(args[0]) < 1 {
		Trusted(s, m, message)
		return
	}
	trusted[args[0]] = true
	messages.QueueMessages(s, &irc.Message{
		Command:  irc.NOTICE,
		Params:   []string{m.Name},
		Trailing: "I now trust " + args[0],
	})
}

func Doubt(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
	args := strings.Split(message, " ")
	if _, ok := trusted[args[0]]; ok {
		delete(trusted, args[0])
		messages.QueueMessages(s, &irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: "I no longer trust " + args[0],
		})
	}
}

func Block(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
	args := strings.Split(message, " ")
	if _, ok := trusted[args[0]]; ok {
		idents[args[0]] = false
		messages.QueueMessages(s, &irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: args[0] + " is now blocked.",
		})
	}

}

func Trusted(s ircx.Sender, m *irc.Message, message string) {
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
	messages.QueueMessages(s,
		&irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: "I trust the following users: " + tlist,
		},
		&irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: "I obey the following users: " + olist,
		},
	)
	log.Printf("I trust the following users: %s\nI obey the following users: %s\n", tlist, olist)
}
