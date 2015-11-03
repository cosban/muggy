package commands

import (
	"fmt"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"github.com/vaughan0/go-ini"
)

type IrcCommand func(ircx.Sender, *irc.Message, string)

var (
	coms    map[string]IrcCommand
	trusted map[string]bool
	idents  map[string]bool
	owners  map[string]bool
	config  ini.File
)

func Configure(t map[string]bool, o map[string]bool, i map[string]bool, c map[string]IrcCommand, conf *ini.File) {
	trusted = t
	owners = o
	idents = i
	coms = c
	config = *conf
}

func identRequest(s ircx.Sender, user string) {
	s.Send(&irc.Message{
		Command:  irc.NOTICE,
		Params:   []string{user},
		Trailing: "I don't recognize you! Please say identify with nickserv then give me a mug :D",
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
