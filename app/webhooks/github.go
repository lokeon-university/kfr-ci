package webhooks

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/go-playground/webhooks.v5/github"
)

//GitHubWebHookHandler handle GitHub WebHooks events.
func GitHubWebHookHandler(w http.ResponseWriter, r *http.Request) {
	hook, err := github.New(github.Options.Secret("GH_OSECRET"))
	if err != nil {
		log.Println("Error connecting to GitHub")
	}
	payload, err := hook.Parse(r, github.PushEvent, github.PingEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			log.Println("event wasn't one of the ones asked to be parsed")
		}
	}
	switch payload.(type) {
	case github.PingPayload:
		w.Header().Set("Content-Type", "application/json")
		res, _ := json.Marshal(map[string]string{"message": "OK"})
		w.Write(res)
		break
	case github.PushPayload:
		//TODO sent event to sqs
		break
	}
}
