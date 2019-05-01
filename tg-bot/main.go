package main

import (
	"log"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

type bot struct {
	bot *tb.Bot
}

func newBot() (*bot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token: os.Getenv("TG_TOKEN"),
		Poller: &tb.Webhook{
			Listen: ":" + os.Getenv("PORT"),
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: os.Getenv("TG_WEBHOOK"),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &bot{b}, nil
}

func (b *bot) start() {
	b.bot.Start()
}

func (b *bot) newHandler(endpoint interface{}, handler func(m *tb.Message)) {
	b.bot.Handle(endpoint, handler)
}

func main() {
	b, err := newBot()
	if err != nil {
		log.Fatal("")
	}
	b.newHandler("/start", b.handleStart)
	b.newHandler("/auth", b.handleOAuth)
	b.newHandler("/help", b.handleHelp)
	bot.handleRepos()
	bot.handleRepoResponse()
	b.start()
}
