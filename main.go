package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/develersrl/bender-bot/conf"
	"github.com/develersrl/bender-bot/slackbot"
	"github.com/nlopes/slack"
)

var (
	cfg = flag.String("config", "./config.toml", "configuration file")
)

type BenderBot struct {
	Bot     *slackbot.Bot
	Channel string
}

type BenderMsg struct {
	Channel string `json:"channelid"`
	Name    string `json:"name"`
	Msg     string `json:"msg"`
}

func (b *BenderBot) postMsg(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var t BenderMsg
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}
	fmt.Println(t.Name, t.Msg, t.Channel)
	channelId := b.Channel
	if t.Channel != "" {
		channelId = t.Channel
	}
	b.Bot.Message(channelId, string(t.Name+": "+t.Msg))
}

func main() {
	flag.Parse()

	conf, err := config.Load(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		token = conf.Bot.SlackToken
	}
	fmt.Printf("Token: %s\n\r", token)
	defaultChannel := os.Getenv("CHANNEL_ID")
	if defaultChannel == "" {
		defaultChannel = conf.Bot.ChannelID
	}

	// Slack Bot filter
	var opts slackbot.Config
	bot := slackbot.New(token, opts)

	bot.DefaultResponse(func(b *slackbot.Bot, msg *slack.Msg) {
		fmt.Printf("Message from channel (%s): %s", msg.Channel, msg.Text)
		bot.Message(msg.Channel, "Non ho capito")
	})

	//bot.RespondTo("^(.*)$", func(b *slackbot.Bot, msg *slack.Msg, args ...string) {
	//    fmt.Printf("Message from channel (%s): %s", msg.Channel, msg.Text)
	//    bot.Message(msg.Channel, "Antani la supercazzola, con scappellamento a destra!")
	//})

	fmt.Printf("Run Bot server\n\r")
	go func(b *slackbot.Bot) {
		if err := b.Start(); err != nil {
			log.Fatalln(err)
		}
	}(bot)

	bender := &BenderBot{Bot: bot, Channel: defaultChannel}
	// Routing
	http.HandleFunc("/show", bender.postMsg)

	fmt.Printf("Run HTTP server\n\r")
	http.ListenAndServe(":8080", nil)
}
