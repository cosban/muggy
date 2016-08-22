package commands

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type SearchData struct {
	Items []struct {
		Kind, Title, HTMLTitle, Link, DisplayLink, Snippet, HTMLSnippet, CacheID, FormattedURL, HTMLFormattedURL string
	}
	Error struct {
		Message string
	}
}

// Search performs a google search
func Search(s ircx.Sender, m *irc.Message, message string) {
	if len(message) < 1 {
		return
	}
	log.Print("im gay")
	site := parseSite(message)
	q := url.QueryEscape(message)

	request := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?&key=%s&cx=%s&q=%s&siteSearch=%s&fields=items&num=1", key, cx, q, site)
	r := &SearchData{}
	getJSON(request, r)

	response := fmt.Sprintf("\u200B%s: No results found :(", m.Prefix.Name)
	if len(r.Items) > 0 {
		response = fmt.Sprintf("\u200B%s: %s -- \u0002%s\u0002: \"%s\" ", m.Prefix.Name, r.Items[0].Link, r.Items[0].Title, r.Items[0].Snippet)
	} else if len(r.Error.Message) > 0 {
		log.Print(r.Error.Message)
	}
	messages.QueueMessages(s, &irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})
}

func parseSite(s string) string {
	for _, element := range strings.Split(s, " ") {
		if strings.HasPrefix(element, "site:") {
			return element[len("site:"):]
		}
	}
	return ""
}
