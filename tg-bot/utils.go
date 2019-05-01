package main

import (
	"fmt"
	"net/url"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

func generateOAuthURL(m *tb.Message) string {
	baseURL, _ := url.Parse("https://github.com/login/oauth/authorize")
	params := url.Values{}
	params.Add("client_id", os.Getenv("GH_OCID"))
	params.Add("scope", "admin:repo_hook repo read:org")
	params.Add("redirect_uri", os.Getenv("GH_REDIRECT_URI"))
	params.Add("state", fmt.Sprintf("%d %d", m.Sender.ID, m.Chat.ID))
	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}
