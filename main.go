package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

type pipeline struct {
	RepositoryID int64  `json:"repository_id,omitempty"`
	URL          string `json:"url,omitempty"`
	Repository   string `json:"repository,omitempty"`
	Branch       string `json:"branch,omitempty"`
	LogFileName  string `json:"log_file_name,omitempty"`
	Language     string `json:"language,omitempty"`
}

func main() {
	ctx := context.Background()
	queueClient, err := pubsub.NewClient(ctx, "kfr-ci")
	if err != nil {
		log.Fatal("Unable to get client for queue")
	}
	var mu sync.Mutex
	received := 0
	sub := queueClient.Subscription("test")
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Printf("Got message: %q\n", string(msg.Data))
		mu.Lock()
		defer mu.Unlock()
		received++
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		fmt.Print(err)
	}
}
