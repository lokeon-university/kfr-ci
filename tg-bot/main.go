package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	tb "gopkg.in/tucnak/telebot.v2"
)

type bot struct {
	bot *tb.Bot
	ctx context.Context
	db  *datastore.Client
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
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "kfr-ci")
	if err != nil {
		return nil, err
	}
	return &bot{b, ctx, client}, nil
}

func (b *bot) start() {
	b.bot.Start()
}

func (b *bot) newHandler(endpoint interface{}, handler interface{}) {
	b.bot.Handle(endpoint, handler)
}

func main() {
	b, err := newBot()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	b.newHandler("/start", b.handleStart)
	b.newHandler("/auth", b.handleOAuth)
	b.newHandler("/help", b.handleHelp)
	b.newHandler("/repo", b.handleRepositories)
	b.start()
}
