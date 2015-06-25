package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/cosban/harkness/commands"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"github.com/vaughan0/go-ini"
)

var (
	name, server, password, prefix, channels string
)

type CommandFunc func(ircx.Sender, *irc.Message, string)

var coms = make(map[string]CommandFunc)
var replies = make(map[*regexp.Regexp]string)

func main() {
	configure()

	coms["g"] = commands.Search
	coms["suicide"] = commands.Quit

	bot := ircx.WithLogin(server, name, name, password)
	if err := bot.Connect(); err != nil {
		log.Panicln("Unable to dial IRC Server ", err)
	}

	RegisterHandlers(bot)
	bot.CallbackLoop()
}

func configure() {
	conf, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}

	name, _ = conf.Get("bot", "name")
	password, _ = conf.Get("bot", "password")
	server, _ = conf.Get("bot", "server")
	ch, _ := conf.Get("bot", "channels")
	channels = strings.Replace(ch, " ", "", -1)
	prefix, _ = conf.Get("bot", "prefix")

	body, err := ioutil.ReadFile("./replies")
	if err != nil {
		log.Panicln("Could not read replies!")
	} else {
		lines := strings.Split(string(body), "\n")
		for i := 0; i < len(lines)-1; i++ {
			kvline := strings.Split(lines[i], ":=:")
			key := regexp.MustCompile(strings.Trim(kvline[0], " "))
			replies[key] = strings.Trim(kvline[1], " ")
		}
	}
}

func RegisterHandlers(bot *ircx.Bot) {
	bot.AddCallback(irc.RPL_WELCOME, ircx.Callback{Handler: ircx.HandlerFunc(RegisterConnect)})
	bot.AddCallback(irc.PING, ircx.Callback{Handler: ircx.HandlerFunc(PingHandler)})
	bot.AddCallback(irc.PRIVMSG, ircx.Callback{Handler: ircx.HandlerFunc(MessageHandler)})
}

func RegisterConnect(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{channels},
	})
}

func PingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}

func MessageHandler(s ircx.Sender, m *irc.Message) {
	msg := m.Trailing
	var command string
	var params string
	if strings.HasPrefix(msg, name) {
		log.Println("Found Name Match")
		pieces := strings.Split(msg, " ")
		if len(pieces) >= 2 {
			log.Println("Correct length " + pieces[0] + " " + pieces[1] + " " + pieces[2])
			command = pieces[1]
			if len(pieces) > 2 {
				params = strings.SplitN(msg, command, 2)[1]
				params = strings.Trim(params, " ")
			} else {
				params = ""
			}
			if c, ok := coms[command]; ok {
				log.Println("Valid Command " + params)
				c(s, m, params)
			}
		}
	} else if strings.HasPrefix(msg, prefix) {
		log.Println("Found prefix match")
		pieces := strings.Split(msg, " ")
		if len(pieces) >= 1 {
			log.Println("Correct length")
			command = pieces[0][1:]
			if len(pieces) >= 2 {
				params = strings.SplitN(msg, command, 2)[1]
				params = strings.Trim(params, " ")
			} else {
				params = ""
			}
			if c, ok := coms[command]; ok {
				log.Println("Valid Command " + params)
				c(s, m, params)
			}
		}
	} else {
		for k, _ := range replies {
			//TODO: rather than compiling each time, compile initially and store in the map
			if ok := k.FindAllString(msg, 1); ok != nil {
				s.Send(&irc.Message{
					Command:  irc.PRIVMSG,
					Params:   m.Params,
					Trailing: replies[k],
				})
				break
			}
		}
	}
}
