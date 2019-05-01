package main

import (
	"context"
	"encoding/json"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

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
	token, err := getToken(m.Sender)
	if err != nil {
		b.bot.Send(m.Sender, "First call /auth")
	} else {
		gc := createClient(context.Background(), token)
		buttons, err := gc.getRepos()
		if err == nil {
			b.bot.Send(m.Sender, "Repositories", &tb.ReplyMarkup{
				InlineKeyboard: buttons,
			})
		} else {
			log.Println(err)
			b.bot.Send(m.Sender, "Unable to read your repositories")
		}
	}
}

func (b *bot) handleRepositoriesResponse(c *tb.Callback) {
	buttons := tb.InlineButton{Unique: "repos"}
	b.bot.Handle(&buttons, func(c *tb.Callback) {
		var payload map[string]string
		er := json.Unmarshal([]byte(c.Data), &payload)
		if err != nil {
			log.Println(err)
		}
		b.bot.Respond(c, &tb.CallbackResponse{})
	})
}
