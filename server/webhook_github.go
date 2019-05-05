package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
	"github.com/lokeon-university/kfr-ci/utils"
	"gopkg.in/go-playground/webhooks.v5/github"
)

type pipeline struct {
	Branch      string `json:"branch,omitempty"`
	Language    string `json:"language,omitempty"`
	LogFileName string `json:"log_file_name,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Repository  string `json:"repository,omitempty"`
	TelegramID  string `json:"telegram_id,omitempty"`
	URL         string `json:"url,omitempty"`
}

//GitHubWebHookHandler handle GitHub WebHooks events.
func GitHubWebHookHandler(w http.ResponseWriter, r *http.Request) {
	hook, err := github.New(github.Options.Secret(os.Getenv("GH_APPSECRET")))
	if err != nil {
		log.Println("Error connecting to GitHub")
	}
	payload, err := hook.Parse(r, github.PushEvent, github.PingEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			log.Println("event wasn't one of the ones asked to be parsed")
		} else {
			log.Println(err)
		}
	}
	switch payload.(type) {
	case github.PingPayload:
		break
	case github.PushPayload:
		sendMessageQueue(payload, mux.Vars(r))
		break
	}
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(map[string]string{"message": "OK"})
	w.Write(res)
}

func sendMessageQueue(payload interface{}, vars map[string]string) {
	push := payload.(github.PushPayload)
	tgID, _ := vars["id"]
	if push.Repository.Language == nil {
		push.Repository.Language = utils.String("none")
	}
	t := queueClient.Topic("webhooks")
	data, _ := json.Marshal(pipeline{
		Branch:      strings.Split(push.Ref, "/")[2],
		Language:    strings.ToLower(*push.Repository.Language),
		LogFileName: push.HeadCommit.ID,
		Owner:       push.Repository.Owner.Login,
		Repository:  push.Repository.Name,
		TelegramID:  tgID,
		URL:         push.Repository.CloneURL,
	})
	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	id, err := result.Get(ctx)
	if err != nil {
		log.Print(err)
	}
	log.Printf("Published a message; msg ID: %v\n", id)
}
