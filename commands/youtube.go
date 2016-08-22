package commands

import (
	"fmt"
	"net/url"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type YoutubeData struct {
	Items []struct {
		Id struct {
			VideoID string
		}
		Snippet struct {
			Title, Description, ChannelTitle string
		}
	}
}

func SearchYoutube(s ircx.Sender, m *irc.Message, message string) {
	if len(message) < 1 {
		return
	}
	q := url.QueryEscape(message)
	request := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?&key=%s&part=id,snippet&maxResults=1&q=%s", key, q)

	r := &YoutubeData{}
	getJSON(request, r)

	response := fmt.Sprintf("\u200B%s: No results found :(", m.Prefix.Name)
	if len(r.Items) > 0 {
		response = fmt.Sprintf(
			"\u200B%s: https://youtube.com/watch?v=%s -- \u0002%s by %s\u0002: \"%s\" ",
			m.Prefix.Name,
			r.Items[0].Id.VideoID,
			r.Items[0].Snippet.Title,
			r.Items[0].Snippet.ChannelTitle,
			r.Items[0].Snippet.Description,
		)
	}

	messages.QueueMessages(s, &irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})

}
