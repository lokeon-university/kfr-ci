package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type kfrBot struct {
	bot *tb.Bot
}

func initBot() kfrBot {
	bot, err := tb.NewBot(tb.Settings{
		Token: os.Getenv("KFR_TELEGRAM"),
		Poller: &tb.Webhook{
			Listen: ":" + os.Getenv("PORT"),
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: fmt.Sprintf("https://%s.herokuapp.com/%s", os.Getenv("HEROKU_APP_NAME"), os.Getenv("KFR_TELEGRAM")),
			},
		},
	})
	checkError(err)
	return kfrBot{bot}
}

func (kfr *kfrBot) start() {
	kfr.bot.Start()
}

func generateOauthURL(m *tb.Message) string {
	baseURL, err := url.Parse("https://github.com/login/oauth/authorize")
	checkError(err)
	params := url.Values{}
	params.Add("client_id", os.Getenv("GH_OCID"))
	params.Add("scope", "admin:repo_hook repo read:org")
	params.Add("redirect_uri", os.Getenv("GH_REDIRECT_URI"))
	params.Add("state", fmt.Sprintf("%d %d", m.Sender.ID, m.Chat.ID))
	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}

func (kfr *kfrBot) hadleAuth() {
	kfr.bot.Handle("/auth", func(m *tb.Message) {
		kfr.bot.Send(m.Sender, "GitHub", &tb.ReplyMarkup{
			InlineKeyboard: [][]tb.InlineButton{{{
				Text: "Confirm OAuth",
				URL:  generateOauthURL(m)},
			},
			},
		})
	})
	log.Println("Handled Auth")
}

func (kfr *kfrBot) handleStart() {
	kfr.bot.Handle("/start", func(m *tb.Message) {
		kfr.bot.Send(m.Sender, "Welcome to KFR CI")
	})
	log.Println("Handled Start")
}

func main() {
	bot := initBot()
	bot.hadleAuth()
	bot.handleRepos()
	bot.handleHelp()
	bot.handleStart()
	bot.handleRepoResponse()
	bot.start()
}
