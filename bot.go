package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/cosban/muggy/commands"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"github.com/vaughan0/go-ini"
)

var (
	name, server, password, prefix, channels string

	conf ini.File

	// name of command mapped to the function itself
	coms = make(map[string]commands.CommandStruct)

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

	bot := ircx.Classic(server, name)

	RegisterHandlers(bot)

	fmt.Printf("Attempting to connect...")
	if err := bot.Connect(); err != nil {
		log.Panicln("Unable to dial IRC Server ", err)
	}

	bot.HandleLoop()
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

	fmt.Printf("Owners: %+v\nTrusted: %+v\nBlocked: %+v\n", o, t, b)

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
	registerCommands()
	return config
}

func registerCommand(com commands.CommandStruct) {
	coms[com.Name] = com
}

func registerCommands() {
	registerCommand(commands.CommandStruct{"google", "<query>", "Queries google for search results.", commands.Search})
	registerCommand(commands.CommandStruct{"quit", "", "Forces the bot to quit (SEE ALSO: restart)", commands.Quit})
	registerCommand(commands.CommandStruct{"trust", "<nick>", "Makes the bot trust a nick (SEE ALSO: doubt, trusted)", commands.Trust})
	registerCommand(commands.CommandStruct{"doubt", "<nick>", "Makes the bot no longer trust a nick (SEE ALSO: trust, trusted", commands.Doubt})
	registerCommand(commands.CommandStruct{"alias", "<command>", "Creates an alias which can be used to run a command", commands.Alias})
	registerCommand(commands.CommandStruct{"join", "<channel>", "Makes the bot join a channel (SEE ALSO: leave)", commands.Join})
	registerCommand(commands.CommandStruct{"leave", "<channel>", "Makes the bot leave a channel (SEE ALSO: join)", commands.Leave})
	registerCommand(commands.CommandStruct{"restart", "", "Forces teh bot to restart (SEE ALSO: quit)", commands.Restart})
	registerCommand(commands.CommandStruct{"trusted", "", "Displays a list of trusted users (SEE ALSO: trust, doubt)", commands.Trusted})
	registerCommand(commands.CommandStruct{"mod", "<query>", "Searches Nexus for mods", commands.Mod})
	registerCommand(commands.CommandStruct{"weather", "<city>, [state]", "Queries for weather in a given city (SEE ALSO: temp)", commands.Weather})
	registerCommand(commands.CommandStruct{"temp", "<city>, state", "Queries for temperature in a given city (SEE ALSO: weather)", commands.Temperature})
	registerCommand(commands.CommandStruct{"raw", "<numeric> <recipient>:[message]", "Forces the bot to perform the given raw command", commands.Raw})
	registerCommand(commands.CommandStruct{"topic", "<message>", "Causes the bot to set the topic", commands.Topic})
	registerCommand(commands.CommandStruct{"help", "[command]", "Displays these words", commands.Help})

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
			m[strings.ToLower(list[i])] = false
		}
	}
}

func RegisterHandlers(bot *ircx.Bot) {
	bot.HandleFunc(irc.RPL_WELCOME, RegisterConnect)
	bot.HandleFunc(irc.PING, PingHandler)
	bot.HandleFunc(irc.PRIVMSG, MessageHandler)
	bot.HandleFunc(irc.JOIN, JoinHandler)
	bot.HandleFunc(irc.QUIT, LeaveHandler)
	bot.HandleFunc(irc.PART, LeaveHandler)
	bot.HandleFunc(irc.NICK, NickHandler)
	bot.HandleFunc(irc.ERR_UNKNOWNCOMMAND, UnknownCommandHandler)
	bot.HandleFunc(irc.RPL_WHOREPLY, WhoisHandler)
	// non standard RPL_WHOISACCOUNT
	bot.HandleFunc("330", WhoisHandler)
	bot.HandleFunc("354", WhoisHandler)
	// ACCOUNT-NOTIFY
	bot.HandleFunc("ACCOUNT", AccountHandler)
}

func RegisterConnect(s ircx.Sender, m *irc.Message) {
	fmt.Printf("Connected... now identifying and joining chans")
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   []string{"NICKSERV"},
		Trailing: "identify " + name + " " + password,
	})
	s.Send(&irc.Message{
		Command:  fmt.Sprintf("%s %s", irc.CAP, irc.CAP_REQ),
		Trailing: "account-notify extended-join",
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

func UnknownCommandHandler(s ircx.Sender, m *irc.Message) {
	fmt.Printf("UNKNOWN COMMAND -- %+v\n", m)
}
