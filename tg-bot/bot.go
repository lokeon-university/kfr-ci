package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/lokeon-university/kfr-ci/utils"
	"google.golang.org/api/iterator"
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

type callBackData struct {
	Owner string `json:"owner,omitempty"`
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}

func (b *bot) getUserToken(u *tb.User) (string, error) {
	q := datastore.NewQuery("users").Filter("ID =", u.ID)
	it := b.db.Run(b.ctx, q)
	var user utils.User
	for {
		_, err := it.Next(&user)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return user.Token, nil
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
		data, _ := json.Marshal(callBackData{*repo.Owner.Name, *repo.Name, token})
		inlineBtn := tb.InlineButton{
			Unique: strconv.FormatInt(*repo.ID, 10),
			Text:   *repo.FullName,
			Data:   string(data),
		}
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
		b.bot.Handle(&inlineBtn, b.handleRepositoriesResponse)
	}
	return inlineKeys
}
