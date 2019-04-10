package webhooks

import (
	"log"
	"net/http"

	"gopkg.in/go-playground/webhooks.v5/github"
)

//GitHubWebHookHandler handle GitHub WebHooks events.
func GitHubWebHookHandler(w http.ResponseWriter, r *http.Request) {
	hook, err := github.New(github.Options.Secret(""))
	if err != nil {
		log.Printf("")
	}
	payload, err := hook.Parse(r, github.PushEvent, github.PingEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			log.Println("event wasn't one of the ones asked to be parsed")
		}
	}
	switch payload.(type) {
	case github.PingPayload:
		// TODO write back
		break
	case github.PushPayload:
		//TODO sent event to sqs
		break
	}
}
