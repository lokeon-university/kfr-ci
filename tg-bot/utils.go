package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

func generateOAuthURL(m *tb.Message) string {
	baseURL, _ := url.Parse("https://github.com/login/oauth/authorize")
	params := url.Values{}
	params.Add("client_id", os.Getenv("GH_APPID"))
	params.Add("scope", "admin:repo_hook repo read:org")
	params.Add("redirect_uri", os.Getenv("GH_REDIRECT_URI"))
	params.Add("state", fmt.Sprintf("%d", m.Sender.ID))
	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}

func statusOK(w http.ResponseWriter, r *http.Request) {
	log.Println("readed")
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(map[string]string{"status": "OK"})
	w.Write(res)
}
