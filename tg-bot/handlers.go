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
			Text: "Confirm Oauth",
			URL:  generateOAuthURL(m),
		}}},
	})
}

func (b *bot) handleStart(m *tb.Message) {
	b.bot.Send(m.Sender, "Bienvenido al bot de KFR-CI. \n Escriba /help para ver los comandos disponibles.")
}

func (b *bot) handleHelp(m *tb.Message) {
	help := `/repos -> Devuelve una lista con los repositorios de una cuenta previamente registrada.
	/auth -> Registra a un usuario mediante su cuenta de Github.`
	b.bot.Send(m.Sender, help)
}

func (b *bot) handleRepositories(m *tb.Message) {
	token := b.getUserToken(m.Sender)
	if token == "" {
		b.bot.Send(m.Sender, "Por favor, escriba /help")
		return
	}
	inlineKeys := b.getRespositoriesBttns(m.Sender, token)
	b.bot.Send(m.Sender, "Seleccione repositorios:", &tb.ReplyMarkup{
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
			msg = "No se ha encontrado el repositorio"
			break
		case http.StatusUnprocessableEntity:
			msg = "El repositorio ya estaba registrado"
			break
		default:
			msg = "Fallo al configurar WebHook"
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
		Text:      "WebHook creado con éxito",
	})
}
