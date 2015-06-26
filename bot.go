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
	name, server, password, prefix, channels, owner string
)

type CommandFunc func(ircx.Sender, *irc.Message, string)

// name of command mapped to the function itself
var coms = make(map[string]CommandFunc)

// regex mapped to its reply
var replies = make(map[*regexp.Regexp]string)

// the nickname of the user mapped to a boolean of whether they are logged in
var trusted = make(map[string]bool)

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
	owner, _ = conf.Get("bot", "owner")
	t, _ := conf.Get("bot", "trusted")
	t = strings.Replace(t, " ", "", -1)
	tusers := strings.Split(t, ",")

	trusted[owner] = false
	for i := 0; i < len(tusers); i++ {
		trusted[tusers[i]] = false
	}

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
	bot.AddCallback(irc.JOIN, ircx.Callback{Handler: ircx.HandlerFunc(JoinHandler)})
	bot.AddCallback(irc.QUIT, ircx.Callback{Handler: ircx.HandlerFunc(LeaveHandler)})
	bot.AddCallback(irc.PART, ircx.Callback{Handler: ircx.HandlerFunc(LeaveHandler)})
	bot.AddCallback(irc.NOTICE, ircx.Callback{Handler: ircx.HandlerFunc(NoticeHandler)})
	bot.AddCallback(irc.NICK, ircx.Callback{Handler: ircx.HandlerFunc(NickHandler)})
}
