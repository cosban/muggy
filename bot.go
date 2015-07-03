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

	conf ini.File

	// name of command mapped to the function itself
	coms = make(map[string]commands.IrcCommand)

	// regex mapped to its reply
	replies = make(map[*regexp.Regexp]string)

	// the nickname of the user mapped to a boolean of whether they are logged in
	trusted = make(map[string]bool)
	owners  = make(map[string]bool)
	// this one is different, we will only allow registered users to commands, the
	// boolean is turned to false if they are blocked
	idents = make(map[string]bool)
)

func main() {
	conf = configure()

	RegisterCommands()

	bot := ircx.Classic(server, name)
	if err := bot.Connect(); err != nil {
		log.Panicln("Unable to dial IRC Server ", err)
	}

	RegisterHandlers(bot)
	bot.CallbackLoop()
}

func configure() ini.File {
	config, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}

	name, _ = config.Get("bot", "name")
	password, _ = config.Get("bot", "password")
	server, _ = config.Get("bot", "server")
	ch, _ := config.Get("bot", "channels")
	channels = strings.Replace(ch, " ", "", -1)
	prefix, _ = config.Get("bot", "prefix")

	o, _ := config.Get("bot", "owners")
	PopulateMap(o, owners)
	t, _ := config.Get("bot", "trusted")
	PopulateMap(t, trusted)
	b, _ := config.Get("bot", "blocked")
	PopulateMap(b, idents)

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
		commands.Configure(trusted, owners, idents, coms, &conf)
	}
	return config
}

func RegisterCommands() {
	coms["g"] = commands.Search
	coms["suicide"] = commands.Quit
	coms["trust"] = commands.AddUser
	coms["doubt"] = commands.RemoveUser
	coms["alias"] = commands.AddAlias
	coms["join"] = commands.Join
	coms["leave"] = commands.Leave
	coms["restart"] = commands.Restart
	coms["trusted"] = commands.ListUsers
	coms["mod"] = commands.ModSearch

	for k, v := range conf["aliases"] {
		k = strings.Trim(k, " ")
		v = strings.Trim(v, " ")

		if com, ok := coms[v]; ok {
			log.Printf("Added alias %s from %s", k, v)
			coms[k] = com
		}
	}

	for k, _ := range coms {
		log.Printf("Command Registered: %s", k)
	}
}

func PopulateMap(s string, m map[string]bool) {
	s = strings.Replace(s, " ", "", -1)
	list := strings.Split(s, ",")
	for i := 0; i < len(list); i++ {
		if len(list[i]) > 0 {
			m[list[i]] = false
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
	bot.AddCallback(irc.RPL_NAMREPLY, ircx.Callback{Handler: ircx.HandlerFunc(NamesHandler)})
}

func RegisterConnect(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   []string{"NICKSERV"},
		Trailing: "identify " + name + " " + password,
	})
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
