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

type ResultWeather struct {
	Weather []struct {
		Id                      int
		Main, Description, Icon string
	}
	Main struct {
		Temp, Pressure, Humidity, Temp_min, Temp_max float64
	}
	Name string
}

var weatherkey string

func init() {
	conf, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}
	weatherkey, _ = conf.Get("weather", "key")
	site = ""
}

func unmarshalWeather(message string) (*ResultWeather, error) {
	q := url.QueryEscape(message)
	request := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&q=%s", weatherkey, q)
	fmt.Printf("Weather: %s\n", request)
	resp, err := http.Get(request)
	if err != nil {
		fmt.Println("Issue connecting to Weather")
		return nil, err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Issue reading json")
		return nil, err
	}
	defer resp.Body.Close()

	r := &ResultWeather{}
	err = json.Unmarshal(contents, &r)
	if err != nil {
		fmt.Println("Issue unmartialing json")
		return nil, err
	}
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
		response = fmt.Sprintf("\u200B%s: %s - %s at %s",
			m.Name,
			r.Name,
			r.Weather[0].Main,
			tempString(r.Main.Temp))
	}
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})
}

func Tempurature(s ircx.Sender, m *irc.Message, message string) {
	r, err := unmarshalWeather(message)
	response := fmt.Sprintf("\u200B%s: 0.0 Kelvin. Seriously.", m.Prefix.Name)
	if err == nil && r != nil {
		response = fmt.Sprintf("\u200B%s: %s - %s H:%s L:%s ",
			m.Name,
			r.Name,
			tempString(r.Main.Temp),
			tempString(r.Main.Temp_max),
			tempString(r.Main.Temp_min))
	}
	s.Send(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   m.Params,
		Trailing: response,
	})
}

func tempString(k float64) string {
	return fmt.Sprintf("%.1fC (%.1fF)", kelvinToC(k), kelvinToF(k))
}

func kelvinToC(k float64) float64 {
	return (k - 272.15)
}

func kelvinToF(k float64) float64 {
	return (k-273.15)*1.8 + 32
}
