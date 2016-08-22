package commands

import (
	"fmt"
	"net/url"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type AVWXData struct {
	RawReport      string `json:"Raw-Report"`
	Station, Error string
}

func Metar(s ircx.Sender, m *irc.Message, message string) {
	callAVWX(s, m, message, "metar")
}

func Taf(s ircx.Sender, m *irc.Message, message string) {
	callAVWX(s, m, message, "taf")
}

func callAVWX(s ircx.Sender, m *irc.Message, message, method string) {
	r := &AVWXData{}
	err := getJSON(fmt.Sprintf("https://avwx.rest/api/%s.php?station=%s&format=JSON", method, url.QueryEscape(message)), r)
	response := fmt.Sprintf("\u200B%s: There is obviously no weather at that location, like, ever.", m.Prefix.Name)
	if len(r.Error) > 0 {
		response = fmt.Sprintf("\u200B%s: %s",
			m.Name,
			r.Error,
		)
	} else if err == nil && len(r.RawReport) > 0 {
		response = fmt.Sprintf("\u200B%s: %s",
			m.Name,
			r.RawReport,
		)
	}

	messages.QueueMessages(s, &irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})
}
