package main

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

const workers = 4

func worker(ctx context.Context, msg *pubsub.Message) {
	// TODO call function to run docker container
	msg.Ack()
	var pipe pipeline
	_ = json.Unmarshal(msg.Data, &pipe)
}

func main() {
	ctx := context.Background()
	queueClient, err := pubsub.NewClient(ctx, "kfr-ci")
	if err != nil {
		log.Fatal("Unable to get client for queue")
	}
	sub := queueClient.Subscription("test")
	sub.ReceiveSettings.NumGoroutines = 4
	sub.ReceiveSettings.MaxOutstandingMessages = 4
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, worker)
	if err != nil {
		log.Println("Unable to process more messages")
		cancel()
	}
}
