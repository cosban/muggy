package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cosban/muggy/messages"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type WeatherData struct {
	Weather []struct {
		Id                      int
		Main, Description, Icon string
	}
	Main struct {
		Temp, Pressure, Humidity, Temp_min, Temp_max, Sea_level, Grnd_level float64
	}
	Wind struct {
		Speed, Deg float64
	}
	Name string
}

func unmarshalWeather(message string) (*WeatherData, error) {
	var request string
	zip := parseZip(message)
	if len(zip) > 1 {
		q := url.QueryEscape(zip)
		request = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&zip=%s", weatherkey, q)
	} else {
		q := url.QueryEscape(message)
		request = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&q=%s", weatherkey, q)
	}

	r := &WeatherData{}
	getJSON(request, r)

	fmt.Printf("%+v\n", r)
	if r.Main.Temp == 0 {
		r = nil
	}
	return r, nil
}

func Weather(s ircx.Sender, m *irc.Message, message string) {
	r, err := unmarshalWeather(message)
	response := fmt.Sprintf("\u200B%s: There is obviously no weather at that location, like, ever.", m.Prefix.Name)
	if err == nil && r != nil {
		response = fmt.Sprintf("\u200B%s %s: %s at %s wind at %s %s",
			m.Name,
			r.Name,
			r.Weather[0].Main,
			tempString(r.Main.Temp),
			speedString(r.Wind.Speed),
			directionString(r.Wind.Deg),
		)
	}
	messages.QueueMessages(s, &irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})
}

func Temperature(s ircx.Sender, m *irc.Message, message string) {
	r, err := unmarshalWeather(message)
	response := fmt.Sprintf("\u200B%s: 0.0 Kelvin. Seriously.", m.Prefix.Name)
	if err == nil && r != nil {
		response = fmt.Sprintf("\u200B%s- %s: %s H:%s L:%s ",
			m.Name,
			r.Name,
			tempString(r.Main.Temp),
			tempString(r.Main.Temp_max),
			tempString(r.Main.Temp_min))
	}
	messages.QueueMessages(s, &irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})
}

func tempString(k float64) string {
	return fmt.Sprintf("%.1fC (%.1fF)", kelvinToC(k), kelvinToF(k))
}

func speedString(k float64) string {
	return fmt.Sprintf("%.1f m/s", k)
}

func directionString(i float64) string {
	if i < 34 {
		return "NNE"
	} else if i < 56 {
		return "NE"
	} else if i < 79 {
		return "ENE"
	} else if i < 101 {
		return "E"
	} else if i < 124 {
		return "ESE"
	} else if i < 146 {
		return "SE"
	} else if i < 169 {
		return "SSE"
	} else if i < 191 {
		return "S"
	} else if i < 214 {
		return "SSW"
	} else if i < 236 {
		return "SW"
	} else if i < 259 {
		return "WSW"
	} else if i < 281 {
		return "W"
	} else if i < 304 {
		return "WNW"
	} else if i < 326 {
		return "NW"
	} else if i < 349 {
		return "NNW"
	} else {
		return "N"
	}
}

func kelvinToC(k float64) float64 {
	return (k - 272.15)
}

func kelvinToF(k float64) float64 {
	return (k-273.15)*1.8 + 32
}

func parseZip(s string) string {
	for _, element := range strings.Split(s, " ") {
		if strings.HasPrefix(element, "zip:") {
			return element[len("zip:"):]
		}
	}
	return ""
}
