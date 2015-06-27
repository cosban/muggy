package commands

import (
	"log"
	"os"
	"os/exec"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Quit(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(m.Name) {
		return
	}
	msg := &irc.Message{
		Command:  irc.QUIT,
		Params:   []string{" "},
		Trailing: "No one ever asks about Muggy!",
	}
	s.Send(msg)

	log.Println(msg.Trailing)
	os.Exit(1)
}

func Restart(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(m.Name) {
		return
	}
	msg := &irc.Message{
		Command:  irc.QUIT,
		Params:   []string{" "},
		Trailing: "No one ever asks about Muggy!",
	}

	s.Send(msg)

	script, _ := config.Get("bot", "script")
	exec.Command(script)
}
