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

type YTResult struct {
	Items []struct {
		Id struct {
			VideoId string
		}
		Snippet struct {
			Title, Description, ChannelTitle string
		}
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

func SearchYoutube(s ircx.Sender, m *irc.Message, message string) {
	if len(message) < 1 {
		return
	}
	site = parseSite(message)
	q := url.QueryEscape(message)
	request := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?&key=%s&part=id,snippet&maxResults=1&q=%s", key, q)
	fmt.Println(request)
	resp, err := http.Get(request)
	site = ""
	if err != nil {
		fmt.Println("Issue connecting to youtube")
		return
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Issue reading json")
		return
	}
	defer resp.Body.Close()

	r := &YTResult{}
	err = json.Unmarshal(contents, &r)
	fmt.Println(r)
	if err != nil {
		fmt.Println("Issue unmartialing json")
		return
	}
	response := fmt.Sprintf("\u200B%s: No results found :(", m.Prefix.Name)
	if len(r.Items) > 0 {
		response = fmt.Sprintf("\u200B%s: https://youtube.com/watch?v=%s -- \u0002%s by %s\u0002: \"%s\" ", m.Prefix.Name, r.Items[0].Id.VideoId, r.Items[0].Snippet.Title, r.Items[0].Snippet.ChannelTitle, r.Items[0].Snippet.Description)
	}
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})

}
