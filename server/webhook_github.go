package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/lokeon-university/kfr-ci/utils"
	"gopkg.in/go-playground/webhooks.v5/github"
)

type pipeline struct {
	RepositoryID int64  `json:"repository_id,omitempty"`
	URL          string `json:"url,omitempty"`
	Repository   string `json:"repository,omitempty"`
	Branch       string `json:"branch,omitempty"`
	LogFileName  string `json:"log_file_name,omitempty"`
	Language     string `json:"language,omitempty"`
}

//GitHubWebHookHandler handle GitHub WebHooks events.
func GitHubWebHookHandler(w http.ResponseWriter, r *http.Request) {
	hook, err := github.New(github.Options.Secret(os.Getenv("GH_OSECRET")))
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
		sendMessageQueue(payload)
		break
	}
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(map[string]string{"message": "OK"})
	w.Write(res)
}

func sendMessageQueue(payload interface{}) {
	push := payload.(github.PushPayload)
	if push.Repository.Language == nil {
		push.Repository.Language = utils.String("None")
	}
	t := queueClient.Topic("webhooks")
	data, _ := json.Marshal(pipeline{
		RepositoryID: push.Repository.ID,
		URL:          push.Repository.SSHURL,
		Branch:       strings.Split(push.Ref, "/")[2],
		LogFileName:  push.HeadCommit.ID,
		Language:     *push.Repository.Language,
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
