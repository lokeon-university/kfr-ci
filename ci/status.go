package main

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

type notifier struct {
	cli *pubsub.Client
}

type status struct {
	Branch     string `json:"branch,omitempty"`
	Owner      string `json:"owner,omitempty"`
	RepoName   string `json:"repo_name,omitempty"`
	Status     string `json:"status,omitempty"`
	TelegramID string `json:"telegram_id,omitempty"`
}

func newNotifier() *notifier {
	ctx := context.Background()
	cli, err := pubsub.NewClient(ctx, "kfr-ci")
	if err != nil {
		log.Fatal("Unable to get client for queue")
	}
	return &notifier{cli}
}

func (s *notifier) updateStatus(tgID, sts, repoName, owner, branch string) {
	ctx := context.Background()
	t := s.cli.Topic("status")
	data, _ := json.Marshal(status{
		Branch:     branch,
		Owner:      owner,
		RepoName:   repoName,
		Status:     sts,
		TelegramID: tgID,
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
