package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/kyokomi/emoji"
	"github.com/lokeon-university/kfr-ci/utils"
	"google.golang.org/api/iterator"
	tb "gopkg.in/tucnak/telebot.v2"
)

type bot struct {
	bot *tb.Bot
	ctx context.Context
	db  *firestore.Client
}

type status struct {
	Branch     string `json:"branch,omitempty"`
	Owner      string `json:"owner,omitempty"`
	RepoName   string `json:"repo_name,omitempty"`
	Status     string `json:"status,omitempty"`
	TelegramID string `json:"telegram_id,omitempty"`
}

func (s *status) getLogURL() string {
	return fmt.Sprintf("https://kfr-ci.tk/logs/%s", s.Status)
}

func (s *status) getMessage() (interface{}, bool) {
	if strings.HasSuffix(s.Status, ".log") {
		return &tb.ReplyMarkup{
			InlineKeyboard: [][]tb.InlineButton{{{
				Text: "Click here.",
				URL:  s.getLogURL(),
			}}},
		}, true
	}
	return emoji.Sprintf("Pipeline Status:\n%s/%s %s branch\n%s", s.Owner, s.RepoName, s.Branch, s.Status), false
}

type updateStatus struct {
	Message struct {
		Attributes struct {
			Key string `json:"key,omitempty"`
		} `json:"attributes,omitempty"`
		Data      status `json:"data,omitempty"`
		MessageID string `json:"messageId,omitempty"`
	} `json:"message,omitempty"`
	Subscription string `json:"subscription,omitempty"`
}

func (u *updateStatus) UnmarshalJSON(data []byte) error {
	type Alias updateStatus
	aux := &struct {
		Message struct {
			Data string `json:"data"`
		} `json:"message,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(aux.Message.Data)
	if err = json.Unmarshal(data, &u.Message.Data); err != nil {
		return err
	}
	return nil
}

func newBot(p *tb.Webhook) (*bot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: p,
	})
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "kfr-ci"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	return &bot{b, ctx, client}, nil
}

func (b *bot) start() {
	log.Println("bot started")
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

func (b *bot) getUserToken(u *tb.User) string {
	iter := b.db.Collection("users").Where("ID", "==", u.ID).Documents(b.ctx)
	var user utils.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		doc.DataTo(&user)
	}
	return user.Token
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
		inlineBtn := tb.InlineButton{
			Unique: strconv.FormatInt(*repo.ID, 10),
			Text:   *repo.FullName,
			Data:   fmt.Sprintf("%s %s", *repo.Owner.Login, *repo.Name),
		}
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
		b.bot.Handle(&inlineBtn, b.handleRepositoriesResponse)
	}
	return inlineKeys
}

func (b *bot) updateStatus(w http.ResponseWriter, r *http.Request) {
	var status updateStatus
	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, _ := strconv.Atoi(status.Message.Data.TelegramID)
	msg, hasLog := status.Message.Data.getMessage()
	user := &tb.User{ID: id}
	if hasLog {
		b.bot.Send(user, "If you want to see the log file.", msg)
	} else {
		b.bot.Send(user, msg)
	}
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(map[string]string{"status": "OK"})
	w.Write(res)
}
