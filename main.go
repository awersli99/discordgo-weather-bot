package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

type locationOF struct {
	country string
	lat     float64
	lon     float64
	tzID    string
	name    string
	region  string
}

type weather struct {
	pressureMb float64
	tempC      float64
	windMPH    float64
	windDegree float64
	humidity   float64
	feelsLikeF float64
	gustMPH    float64
	condition  map[string]interface{}
}

func init() {

	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func getWeather(location string) (weather, locationOF, string) {
	url := "http://api.apixu.com/v1/current.json?key=43c0be7ca1624d31b23124340190907&q=" + location
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	weatherMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &weatherMap)
	if err != nil {
		return weather{}, locationOF{}, "error"
	}
	if weatherMap["current"] == nil {
		return weather{}, locationOF{}, "error"
	}
	currentWeather := weatherMap["current"].(map[string]interface{})
	weatherStruct := weather{
		currentWeather["pressure_mb"].(float64),
		currentWeather["temp_c"].(float64),
		currentWeather["wind_mph"].(float64),
		currentWeather["wind_degree"].(float64),
		currentWeather["humidity"].(float64),
		currentWeather["feelslike_f"].(float64),
		currentWeather["gust_mph"].(float64),
		currentWeather["condition"].(map[string]interface{}),
	}
	locationx := weatherMap["location"].(map[string]interface{})
	locationStruct := locationOF{
		locationx["country"].(string),
		locationx["lat"].(float64),
		locationx["lon"].(float64),
		locationx["tz_id"].(string),
		locationx["name"].(string),
		locationx["region"].(string),
	}
	return weatherStruct, locationStruct, ""
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandContent := strings.Split(m.Content, " ")
	if commandContent[0] == "?weather" {
		var location string
		for i := range commandContent {
			if i == 0 {
				continue
			} else if i == 1 {
				location += commandContent[i]
			} else {
				location += "-" + commandContent[i]
			}
		}
		weatherStruct, locationStruct, errorx := getWeather(location)
		if errorx != "" {
			embed := &discordgo.MessageEmbed{
				Title:       "Error!",
				Description: "Location " + location + " not found.",
				Color:       0x7ec0ee,
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}
		icon := "https:" + weatherStruct.condition["icon"].(string)
		embed := &discordgo.MessageEmbed{
			Title:       locationStruct.name,
			Description: "Weather statistics:",
			Color:       0x7ec0ee,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: icon,
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Location: ",
					Value:  locationStruct.name,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Region: ",
					Value:  locationStruct.region + ", " + locationStruct.country,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Latitude: ",
					Value:  fmt.Sprintf("%.2f", locationStruct.lat),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Longitude: ",
					Value:  fmt.Sprintf("%.2f", locationStruct.lon),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Temperature (Celcius): ",
					Value:  fmt.Sprintf("%.0f", weatherStruct.tempC),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Temperature (Fahrenheit): ",
					Value:  fmt.Sprintf("%.0f", weatherStruct.feelsLikeF),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Wind MPH: ",
					Value:  fmt.Sprintf("%.0f", weatherStruct.windMPH),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Wind Degree: ",
					Value:  fmt.Sprintf("%.0f", weatherStruct.windDegree),
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Weather bot made by Awersli99 using: https://www.apixu.com/",
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		// http://api.apixu.com/v1/current.json?key=43c0be7ca1624d31b23124340190907&q=CITYNAME
	}
}

func main() {

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
