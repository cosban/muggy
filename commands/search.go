package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"github.com/vaughan0/go-ini"
)

var (
	key, cx, site string
)

type Result struct {
	Items []struct {
		Kind, Title, HtmlTitle, Link, DisplayLink, Snippet, HtmlSnippet, CacheId, FormattedUrl, HtmlFormattedUrl string
	}
	Error struct {
		Message string
	}
}

func init() {
	conf, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}
	key, _ = conf.Get("google", "key")
	cx, _ = conf.Get("google", "cx")
	site = ""
}

// Search performs a google search
func Search(s ircx.Sender, m *irc.Message, message string) {
	if len(message) < 1 {
		return
	}
	log.Print("im gay")
	site = parseSite(message)
	q := url.QueryEscape(message)
	request := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?&key=%s&cx=%s&q=%s&siteSearch=%s&fields=items&num=1", key, cx, q, site)
	resp, err := http.Get(request)
	site = ""
	if err != nil {
		log.Print("Issue connecting to google")
		return
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Issue reading json")
		return
	}
	defer resp.Body.Close()

	r := &Result{}
	err = json.Unmarshal(contents, &r)
	if err != nil {
		log.Println("Issue unmartialing json")
		return
	}
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
