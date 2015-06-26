package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"github.com/vaughan0/go-ini"
)

var (
	key, cx string
)

type Result struct {
	Items []struct {
		Kind, Title, HtmlTitle, Link, DisplayLink, Snippet, HtmlSnippet, CacheId, FormattedUrl, HtmlFormattedUrl string
	}
}

func init() {
	conf, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}
	key, _ = conf.Get("google", "key")
	cx, _ = conf.Get("google", "cx")
}

func Search(s ircx.Sender, m *irc.Message, message string) {
	if len(message) < 1 {
		return
	}
	q := url.QueryEscape(message)
	request := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=1", key, cx, q)
	resp, err := http.Get(request)
	if err != nil {
		fmt.Println("Issue connecting to google")
		return
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Issue reading json")
		return
	}
	defer resp.Body.Close()

	r := &Result{}
	err = json.Unmarshal(contents, &r)
	if err != nil {
		fmt.Println("Issue unmartialing json")
		return
	}
	response := fmt.Sprintf("%s: %s -> %s", m.Prefix.Name, r.Items[0].Link, r.Items[0].Title)
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: r.Items[0].Snippet,
	})
}
