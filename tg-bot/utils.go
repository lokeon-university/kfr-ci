package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/lokeon-university/kfr-ci/utils"
	"google.golang.org/api/iterator"
	tb "gopkg.in/tucnak/telebot.v2"
)

type callBackData struct {
	Owner string `json:"owner,omitempty"`
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}

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
