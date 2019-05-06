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
<<<<<<< HEAD
	b.bot.Send(m.Sender, "Welcome to the KFR-CI bot. \n Type /help to see the available commands.")
=======
	b.bot.Send(m.Sender, "Bienvenido al bot de KFR-CI. \n Escriba /help para ver los comandos disponibles.")
>>>>>>> 317831a1bce3ead9dc260b4d35d28616567347c7
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
<<<<<<< HEAD
		b.bot.Send(m.Sender, "Please, try /auth")
		return
	}
	inlineKeys := b.getRespositoriesBttns(m.Sender, token)
	b.bot.Send(m.Sender, "Choose a Repository:", &tb.ReplyMarkup{
=======
		b.bot.Send(m.Sender, "Por favor, escriba /help")
		return
	}
	inlineKeys := b.getRespositoriesBttns(m.Sender, token)
	b.bot.Send(m.Sender, "Seleccione repositorios:", &tb.ReplyMarkup{
>>>>>>> 317831a1bce3ead9dc260b4d35d28616567347c7
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
<<<<<<< HEAD
			msg = "Repository cannot be found."
			break
		case http.StatusUnprocessableEntity:
			msg = "The repository was already registered."
			break
		default:
			msg = "Failed to set Webhook."
=======
			msg = "No se ha encontrado el repositorio"
			break
		case http.StatusUnprocessableEntity:
			msg = "El repositorio ya estaba registrado"
			break
		default:
			msg = "Fallo al configurar WebHook"
>>>>>>> 317831a1bce3ead9dc260b4d35d28616567347c7
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
<<<<<<< HEAD
		Text:      "WebHook created sucesfully",
=======
		Text:      "WebHook creado con éxito",
>>>>>>> 317831a1bce3ead9dc260b4d35d28616567347c7
	})
}
