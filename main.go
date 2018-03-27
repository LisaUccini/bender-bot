package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/develersrl/bender-bot/conf"
	"github.com/develersrl/bender-bot/slackbot"
	"github.com/nlopes/slack"
)

var (
	cfg = flag.String("config", "./config.toml", "configuration file")
)

func file_check(bot *slackbot.Bot, channel string, filename string) {
	for {
		if bot != nil {
			text, err := ioutil.ReadFile(filename)
			if err == nil {
				bot.Message(channel, string(text))
				os.Remove(filename)
			}
		}
		time.Sleep(2500 * time.Millisecond)
	}
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
	embedded := os.Getenv("EMBEDDED_CHANNEL_ID")
	if embedded == "" {
		embedded = conf.Bot.ChannelID
	}
	fmt.Printf("Embedded Channel: %s\n\r", embedded)
	test_log_file := os.Getenv("TEST_LOG_FILE")
	if test_log_file == "" {
		test_log_file = conf.Bot.TestLogFile
	}
	fmt.Printf("Test Log File: %s\n\r", test_log_file)

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

	go file_check(bot, embedded, test_log_file)

	if err := bot.Start(); err != nil {
		log.Fatalln(err)
	}
}
