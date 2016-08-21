package commands

import (
	"fmt"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"github.com/vaughan0/go-ini"
)

type CommandHandler func(ircx.Sender, *irc.Message, string)

type CommandStruct struct {
	Name, Usage, Summary string
	Handle               CommandHandler
}

var (
	coms    map[string]CommandStruct
	trusted map[string]bool
	idents  map[string]bool
	owners  map[string]bool
	config  ini.File
)

func Configure(t map[string]bool, o map[string]bool, i map[string]bool, c map[string]CommandStruct, conf *ini.File) {
	trusted = t
	owners = o
	idents = i
	coms = c
	config = *conf
}

func identRequest(s ircx.Sender, user string) {
	messages.QueueMessages(s, &irc.Message{
		Command:  irc.NOTICE,
		Params:   []string{user},
		Trailing: "I don't recognize you! Please identify with services.",
	})
}

func isOwner(s ircx.Sender, user string) bool {
	var b, ok bool
	if b, ok = owners[user]; !b || !ok {
		fmt.Printf("Failed command attempt by: %s, Owner required\n", user)
		identRequest(s, user)
	}
	return b && ok
}

func isTrusted(s ircx.Sender, user string) bool {
	if b, ok := owners[user]; ok && b {
		return true
	} else if b, ok := trusted[user]; ok && b {
		return true
	} else if !b {
		return false
	}
	fmt.Printf("Failed command attempt by: %s, Trust required\n", user)
	identRequest(s, user)
	return false
}
