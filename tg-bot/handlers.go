package main

import (
	"encoding/json"
	"os"

	"github.com/google/go-github/github"
	"github.com/lokeon-university/kfr-ci/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

var githubWebhook = github.Hook{
	Active: github.Bool(true),
	Events: []string{"push"},
	Config: map[string]interface{}{
		"content_type": "json",
		"url":          os.Getenv("GH_WH"),
		"secret":       os.Getenv("GH_OSECRET"),
	}}

func (b *bot) handleOAuth(m *tb.Message) {
	b.bot.Send(m.Sender, "GitHub", &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{{
			Text: "Confirm Oauth",
			URL:  generateOAuthURL(m),
		}}},
	})
}

func (b *bot) handleStart(m *tb.Message) {
	b.bot.Send(m.Sender, "Welcome to kfr-ci")
}

func (b *bot) handleHelp(m *tb.Message) {
	help := `/repos -> Devuelve una lista con los repositorios de una cuenta previamente registrada.
	/auth -> Registra a un usuario mediante su cuenta de Github.`
	b.bot.Send(m.Sender, help)
}

func (b *bot) handleRepositories(m *tb.Message) {
	token, err := b.getUserToken(m.Sender)
	if err != nil {
		b.bot.Send(m.Sender, "Please call /help")
		return
	}
	inlineKeys := b.getRespositoriesBttns(m.Sender, token)
	b.bot.Send(m.Sender, "Choose Repositorie:", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (b *bot) handleRepositoriesResponse(c *tb.Callback) {
	var data callBackData
	_ = json.Unmarshal([]byte(c.Data), &data)
	gc := utils.NewGitHubClient(b.ctx, data.Token)
	err := gc.SetWebhook(data.Name, data.Owner, &githubWebhook)
	if err != nil {
		b.bot.Respond(c, &tb.CallbackResponse{})
		return
	}
	b.bot.Respond(c, &tb.CallbackResponse{})
}
