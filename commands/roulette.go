package commands

import (
	"math/rand"
	"time"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

var r *rand.Rand
var chamber int
var active int

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	chamber = r.Intn(6)
	spin()
}

func Roulette(s ircx.Sender, m *irc.Message, message string) {
	if chamber == active {
		active = 0
		chamber = r.Intn(6)

		s.Send(&irc.Message{
			Command:  irc.PRIVMSG,
			Params:   m.Params,
			Trailing: "BANG!",
		})
		params := make([]string, len(m.Params))
		for i, v := range m.Params {
			params[i] = v
		}
		params = append(params, m.Name)
		s.Send(&irc.Message{
			Command:  irc.KICK,
			Params:   params,
			Trailing: "BANG!",
		})
	} else {
		s.Send(&irc.Message{
			Command:  irc.PRIVMSG,
			Params:   m.Params,
			Trailing: "Click...",
		})
		active = (active + 1) % 6
	}
}

func spin() {
	active = r.Intn(6)
}

func SpinFire(s ircx.Sender, m *irc.Message, message string) {
	spin()
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: "Spinning the chamber before firing...",
	})
	Roulette(s, m, message)
}
