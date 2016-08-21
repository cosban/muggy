package commands

import (
	"fmt"
	"strings"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func Help(s ircx.Sender, m *irc.Message, message string) {
	if !isTrusted(s, m.Name) {
		return
	}
	args := strings.Split(message, " ")
	if len(args[0]) < 1 {
		var msgs []string
		comstr := ""
		for _, v := range coms {
			comstr = fmt.Sprintf("%s, %s", comstr, v.Name)
		}
		msgs = append(msgs, "-------- HELP --------", "NOTE: To view usage of a command, please use .help <command>", comstr[2:], "------ END HELP ------")
		messages.BuildMessages(s, irc.NOTICE, []string{m.Name}, msgs...)
	} else if v, ok := coms[args[0]]; ok {
		messages.BuildMessages(s, irc.NOTICE, []string{m.Name}, fmt.Sprintf("%s %s -- %s", v.Name, v.Usage, v.Summary))
	}
}
