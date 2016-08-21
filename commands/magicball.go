package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

var responses = []string{
	"It is certain",
	"Reply hazy, try again",
	"Don't count on it",
	"It is decidedly so",
	"Ask again later",
	"My reply is no",
	"Without a doubt",
	"Better not tell you now",
	"My sources say no",
	"Yes, definitely",
	"Cannot predict now",
	"Outlook not so good",
	"You may rely on it",
	"Concentrate and ask again",
	"Very doubtful",
	"As I see it, yes",
	"That question is unclear",
	"Definitely not",
	"Mos defs",
	"Hell na",
}

func MagicBall(s ircx.Sender, m *irc.Message, message string) {
	if len(message) < 1 {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Int31n(int32(len(responses)))
	response := fmt.Sprintf("\u200B%s: %s", m.Name, responses[index])

	messages.QueueMessages(s, &irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})
}
