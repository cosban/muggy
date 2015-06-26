package commands

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type IrcCommand func(ircx.Sender, *irc.Message, string)

var (
	coms    map[string]IrcCommand
	trusted map[string]bool
)

func Configure(t map[string]bool, c map[string]IrcCommand) {
	trusted = t
	coms = c
}
