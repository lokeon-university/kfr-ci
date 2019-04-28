package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
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

var qURLGitHub = "https://sqs.us-east-1.amazonaws.com/492996661514/github_webhook.fifo"

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
		push.Repository.Language = aws.String("None")
	}
	result, err := sqsClient.SendMessage(&sqs.SendMessageInput{
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
		MessageBody:    aws.String("GitHub Webhook " + push.HeadCommit.ID),
		MessageGroupId: aws.String("WebHooks"),
		QueueUrl:       &qURLGitHub,
	})
	if err != nil {
		log.Fatal("Error", err)
	}
	log.Println("Sended Message to Queue", *result.MessageId)
}
