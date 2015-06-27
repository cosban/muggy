package main

import (
	"fmt"
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
	o = strings.Replace(o, " ", "", -1)
	ousers := strings.Split(o, ",")
	fmt.Printf("Owners: ")
	for i := 0; i < len(ousers); i++ {
		owners[ousers[i]] = false
		fmt.Printf("%s ", ousers[i])
	}

	fmt.Printf("\nTrusted: ")
	t, _ := config.Get("bot", "trusted")
	t = strings.Replace(t, " ", "", -1)
	tusers := strings.Split(t, ",")
	for i := 0; i < len(tusers); i++ {
		trusted[tusers[i]] = false
		fmt.Printf("%s ", tusers[i])
	}
	fmt.Printf("\n")

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
		commands.Configure(trusted, owners, coms, &conf)
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

	log.Printf("%+v", conf)

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
