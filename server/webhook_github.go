package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gopkg.in/go-playground/webhooks.v5/github"
)

type pipeline struct {
	RepositoryID int64
	URL          string
	Repository   string
	Branch       string
	LogFileName  string
	Language     string
}

var qURLGitHub = "https://sqs.us-east-1.amazonaws.com/492996661514/github_webhook.fifo"

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
		responseWebHook(w, r)
		break
	case github.PushPayload:
		sendMessageQueue(payload)
		responseWebHook(w, r)
		break
	}
}

func responseWebHook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(map[string]string{"message": "OK"})
	w.Write(res)
}

func sendMessageQueue(payload interface{}) {
	push := payload.(github.PushPayload)
	result, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"repository": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(push.Repository.Name),
			},
			"repository_id": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(strconv.FormatInt(push.Repository.ID, 10)),
			},
			"branch": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(strings.Split(push.Ref, "/")[2]),
			},
			"language": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: push.Repository.Language,
			},
			"url": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(push.Repository.SSHURL),
			},
			"log": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(push.HeadCommit.ID),
			},
		},
		QueueUrl: &qURLGitHub,
	})
	if err != nil {
		log.Fatal("Error", err)
	}
	log.Println("Sended Message to Queue", *result.MessageId)
}
