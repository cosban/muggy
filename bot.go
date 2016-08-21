package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/cosban/muggy/commands"
	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
	"github.com/vaughan0/go-ini"
)

var (
	name, server, password, prefix, channels, serverpass string

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
	var bot *ircx.Bot
	if len(serverpass) > 0 {
		bot = ircx.WithLoginTLS(server, name, name, serverpass, nil)
	} else {
		bot = ircx.Classic(server, name)
	}

	performRegistrations(bot)

	if err := bot.Connect(); err != nil {
		log.Fatal(err)
	}
	messages.WriteLoop()
	fmt.Printf("Connecting...\n")
	bot.HandleLoop()
}

func configure() ini.File {
	config, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}

	name, _ = config.Get("bot", "name")
	password, _ = config.Get("bot", "password")
	serverpass, _ = config.Get("bot", "serverpass")
	server, _ = config.Get("bot", "server")
	ch, _ := config.Get("bot", "channels")
	channels = strings.Replace(ch, " ", "", -1)
	prefix, _ = config.Get("bot", "prefix")

	ownerMap, _ := config.Get("bot", "owners")
	PopulateMap(ownerMap, owners)
	trustedMap, _ := config.Get("bot", "trusted")
	PopulateMap(trustedMap, trusted)
	blockedMap, _ := config.Get("bot", "blocked")
	PopulateMap(blockedMap, idents)

	log.Printf("User permissions are as follows:\nOwners: %+v\nTrusted: %+v\nBlocked: %+v\n", ownerMap, trustedMap, blockedMap)

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

func registerCommand(com commands.CommandStruct) {
	coms[com.Name] = com
}

// PopulateMap takes in a stringified list and a map and then populates the map with keys taken from the list
func PopulateMap(s string, m map[string]bool) {
	s = strings.Replace(s, " ", "", -1)
	list := strings.Split(s, ",")
	for i := 0; i < len(list); i++ {
		if len(list[i]) > 0 {
			m[strings.ToLower(list[i])] = false
		}
	}
}

func registerConnect(s ircx.Sender, m *irc.Message) {
	log.Print("Identifying and joining channels...")
	if len(password) > 0 {
		messages.QueueMessages(s,
			&irc.Message{
				Command:  irc.PRIVMSG,
				Params:   []string{"NICKSERV"},
				Trailing: "identify " + name + " " + password,
			},
		)
	}

	messages.QueueMessages(s,
		&irc.Message{
			Command:  fmt.Sprintf("%s %s", irc.CAP, irc.CAP_REQ),
			Trailing: "account-notify extended-join",
		},
	)
	messages.QueueMessages(s,
		&irc.Message{
			Command: irc.JOIN,
			Params:  []string{channels},
		},
	)
}

func pingHandler(s ircx.Sender, m *irc.Message) {
	// special case, don't use message queue
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}

func loginErrHandler(s ircx.Sender, m *irc.Message) {
	log.Println(m.Trailing)
}

func unknownCommandHandler(s ircx.Sender, m *irc.Message) {
	fmt.Printf("UNKNOWN COMMAND -- %+v\n", m)
}

func performRegistrations(bot *ircx.Bot) {
	bot.HandleFunc(irc.RPL_WELCOME, registerConnect)
	bot.HandleFunc(irc.PING, pingHandler)
	bot.HandleFunc(irc.PRIVMSG, MessageHandler)
	bot.HandleFunc(irc.NOTICE, MessageHandler)
	bot.HandleFunc(irc.JOIN, JoinHandler)
	bot.HandleFunc(irc.QUIT, LeaveHandler)
	bot.HandleFunc(irc.PART, LeaveHandler)
	bot.HandleFunc(irc.NICK, NickHandler)
	bot.HandleFunc(irc.ERR_UNKNOWNCOMMAND, unknownCommandHandler)
	bot.HandleFunc(irc.RPL_WHOREPLY, WhoisHandler)
	// non standard RPL_WHOISACCOUNT
	bot.HandleFunc("330", WhoisHandler)
	bot.HandleFunc("354", WhoisHandler)
	bot.HandleFunc("353", WhoisHandler)
	bot.HandleFunc(irc.ERR_NOLOGIN, loginErrHandler)
	// ACCOUNT-NOTIFY
	bot.HandleFunc("ACCOUNT", AccountHandler)

	registerCommand(commands.CommandStruct{
		Name:    "google",
		Usage:   "<query>",
		Summary: "Queries google for search results.",
		Handle:  commands.Search,
	})
	registerCommand(commands.CommandStruct{
		Name:    "yt",
		Usage:   "<query>",
		Summary: "Queries Youtube for videos.",
		Handle:  commands.SearchYoutube,
	})
	registerCommand(commands.CommandStruct{
		Name:    "quit",
		Usage:   "",
		Summary: "Forces the bot to quit (SEE ALSO: restart)",
		Handle:  commands.Quit,
	})
	registerCommand(commands.CommandStruct{
		Name:    "trust",
		Usage:   "<nick>",
		Summary: "Makes the bot trust a nick (SEE ALSO: doubt, trusted)",
		Handle:  commands.Trust,
	})
	registerCommand(commands.CommandStruct{
		Name:    "doubt",
		Usage:   "<nick>",
		Summary: "Makes the bot no longer trust a nick (SEE ALSO: trust, trusted",
		Handle:  commands.Doubt,
	})
	registerCommand(commands.CommandStruct{
		Name:    "alias",
		Usage:   "<command>",
		Summary: "Creates an alias which can be used to run a command",
		Handle:  commands.Alias,
	})
	registerCommand(commands.CommandStruct{
		Name:    "join",
		Usage:   "<channel>",
		Summary: "Makes the bot join a channel (SEE ALSO: leave)",
		Handle:  commands.Join,
	})
	registerCommand(commands.CommandStruct{
		Name:    "leave",
		Usage:   "<channel>",
		Summary: "Makes the bot leave a channel (SEE ALSO: join)",
		Handle:  commands.Leave,
	})
	registerCommand(commands.CommandStruct{
		Name:    "restart",
		Usage:   "",
		Summary: "Forces the bot to restart (SEE ALSO: quit)",
		Handle:  commands.Restart,
	})
	registerCommand(commands.CommandStruct{
		Name:    "trusted",
		Usage:   "",
		Summary: "Displays a list of trusted users (SEE ALSO: trust, doubt)",
		Handle:  commands.Trusted,
	})
	registerCommand(commands.CommandStruct{
		Name:    "weather",
		Usage:   "<city>, [state]",
		Summary: "Queries for weather in a given city (SEE ALSO: temp)",
		Handle:  commands.Weather,
	})
	registerCommand(commands.CommandStruct{
		Name:    "temperature",
		Usage:   "<city>, state",
		Summary: "Queries for temperature in a given city (SEE ALSO: weather)",
		Handle:  commands.Temperature,
	})
	registerCommand(commands.CommandStruct{
		Name:    "raw",
		Usage:   "<numeric> <recipient>:[message]",
		Summary: "Forces the bot to perform the given raw command",
		Handle:  commands.Raw,
	})
	registerCommand(commands.CommandStruct{
		Name:    "topic",
		Usage:   "<message>",
		Summary: "Causes the bot to set the topic",
		Handle:  commands.Topic,
	})
	registerCommand(commands.CommandStruct{
		Name:    "help",
		Usage:   "[command]",
		Summary: "Displays these words",
		Handle:  commands.Help,
	})
	registerCommand(commands.CommandStruct{
		Name:    "8ball",
		Usage:   "[question]",
		Summary: "Summons the power of the magic 8ball to answer your question",
		Handle:  commands.MagicBall,
	})
	registerCommand(commands.CommandStruct{
		Name:    "roulette",
		Usage:   "",
		Summary: "Russian roulette with a six shooter",
		Handle:  commands.Roulette,
	})
	registerCommand(commands.CommandStruct{
		Name:    "spin",
		Usage:   "",
		Summary: "Spins the six shooter before firing",
		Handle:  commands.SpinFire,
	})

	for k := range coms {
		log.Printf("Command Registered: %s", k)
	}

	for k, v := range conf["aliases"] {
		k = strings.Trim(k, " ")
		v = strings.Trim(v, " ")

		if com, ok := coms[v]; ok {
			log.Printf("Added alias %s from %s", k, v)
			registerCommand(commands.CommandStruct{
				Name:    k,
				Usage:   com.Usage,
				Summary: com.Summary,
				Handle:  com.Handle,
			})
		}
	}
}
