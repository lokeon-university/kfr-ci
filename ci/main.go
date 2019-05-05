package main

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

const workers = 4

var agnt *agent
var notify *notifier

func worker(ctx context.Context, msg *pubsub.Message) {
	//TODO call function to run docker container
	msg.Ack()
	var pipe pipeline
	_ = json.Unmarshal(msg.Data, &pipe)
	pipe.Status = notify.updateStatus
	agnt.buildPipeline(&pipe)
}

func main() {
	ctx := context.Background()
	agnt = newAgent()
	notify = newNotifier()
	queueClient, err := pubsub.NewClient(ctx, "kfr-ci")
	if err != nil {
		log.Fatal("Unable to get client for queue")
	}
	sub := queueClient.Subscription("webhooks")
	sub.ReceiveSettings.NumGoroutines = workers
	sub.ReceiveSettings.MaxOutstandingMessages = workers
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, worker)
	if err != nil {
		log.Println(err, "Unable to process more messages")
		cancel()
	}
}
