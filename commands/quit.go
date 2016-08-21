package commands

import (
	"log"
	"os"
	"syscall"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Quit(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}
	messages.QueueMessages(s, &irc.Message{
		Command:  irc.QUIT,
		Params:   []string{" "},
		Trailing: "No one ever asks about Muggy!",
	})
	messages.StopLoop()
	os.Exit(1)
}

func Restart(s ircx.Sender, m *irc.Message, message string) {
	if !isOwner(s, m.Name) {
		return
	}

	gopath := os.Getenv("GOPATH")
	err := syscall.Exec(gopath+"/bin/muggy", []string{"muggy"}, os.Environ())
	log.Print(err)
}
