package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/lokeon-university/kfr-ci/utils"
	"google.golang.org/api/iterator"
	tb "gopkg.in/tucnak/telebot.v2"
)

type bot struct {
	bot *tb.Bot
	ctx context.Context
	db  *firestore.Client
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
	conf := &firebase.Config{ProjectID: "kfr-ci"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	return &bot{b, ctx, client}, nil
}

func (b *bot) start() {
	b.bot.Start()
}

func (b *bot) newHandler(endpoint interface{}, handler interface{}) {
	b.bot.Handle(endpoint, handler)
}

type callBackData struct {
	Owner string `json:"owner,omitempty"`
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}

func (b *bot) getUserToken(u *tb.User) string {
	iter := b.db.Collection("users").Where("ID", "==", u.ID).Documents(b.ctx)
	var user utils.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		doc.DataTo(&user)
	}
	return user.Token
}

func (b *bot) getRespositoriesBttns(u *tb.User, token string) [][]tb.InlineButton {
	inlineKeys := [][]tb.InlineButton{}
	gc := utils.NewGitHubClient(b.ctx, token)
	repos, err := gc.GetRespositories()
	if err != nil {
		b.bot.Send(u, "Unable to get your repositories")
		return inlineKeys
	}
	for _, repo := range repos {
		inlineBtn := tb.InlineButton{
			Unique: strconv.FormatInt(*repo.ID, 10),
			Text:   *repo.FullName,
			Data:   fmt.Sprintf("%s %s", *repo.Owner.Login, *repo.Name),
		}
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
		b.bot.Handle(&inlineBtn, b.handleRepositoriesResponse)
	}
	return inlineKeys
}
