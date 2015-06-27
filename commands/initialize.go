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
	config  ini.File
	owners  map[string]bool
)

func Configure(t map[string]bool, o map[string]bool, c map[string]IrcCommand, conf *ini.File) {
	trusted = t
	owners = o
	coms = c
	config = *conf
}

func isOwner(user string) bool {
	var b, ok bool
	if b, ok = owners[user]; !b || !ok {
		fmt.Printf("Failed command attempt by: %s, Owner required\n", user)
	}
	return b && ok
}

func isTrusted(user string) bool {
	if b, ok := trusted[user]; ok && b {
		return true
	} else if b, ok := owners[user]; ok && b {
		return true
	}
	fmt.Printf("Failed command attempt by: %s, Trust required\n", user)
	return false
}
