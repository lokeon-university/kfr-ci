package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lokeon-university/kfr-ci/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *bot) handleOAuth(m *tb.Message) {
	b.bot.Send(m.Sender, "GitHub", &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{{
			Text: "Please, give us permission to see your repositories on GitHub",
			URL:  generateOAuthURL(m),
		}}},
	})
}

func (b *bot) handleStart(m *tb.Message) {
	b.bot.Send(m.Sender, "Welcome to the KFR-CI bot. \n Type /help to see the available commands.")
}

func (b *bot) handleHelp(m *tb.Message) {
	help := `/repos -> Returns a list with the repositories of a previously registered account.
	/auth -> Register a user through your GitHub account.
	/help -> This command.`
	b.bot.Send(m.Sender, help)
}

func (b *bot) handleRepositories(m *tb.Message) {
	token := b.getUserToken(m.Sender)
	if token == "" {
		b.bot.Send(m.Sender, "Please, try /auth")
		return
	}
	inlineKeys := b.getRespositoriesBttns(m.Sender, token)
	b.bot.Send(m.Sender, "Choose a Repository:", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (b *bot) handleRepositoriesResponse(c *tb.Callback) {
	data := strings.Split(c.Data, " ")
	gc := utils.NewGitHubClient(b.ctx, b.getUserToken(c.Sender))
	status, err := gc.SetWebhook(data[1], data[0], strconv.Itoa(c.Sender.ID))
	if err != nil {
		var msg string
		switch status {
		case http.StatusNotFound:
			msg = "Repository cannot be found."
			break
		case http.StatusUnprocessableEntity:
			msg = "The repository was already registered."
			break
		default:
			msg = "Failed to set Webhook."
			break
		}
		b.bot.Respond(c, &tb.CallbackResponse{
			Text:      msg,
			ShowAlert: true,
		})
		return
	}
	b.bot.Respond(c, &tb.CallbackResponse{
		ShowAlert: true,
		Text:      "WebHook created sucesfully",
	})
}
