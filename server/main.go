package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
)

var dbClient *datastore.Client
var queueClient *pubsub.Client
var ctx = context.Background()

func setupQueueDataBase() {
	var err error
	dbClient, err = datastore.NewClient(ctx, "")
	if err != nil {
		log.Fatal("Unable to get client for database")
	}
	queueClient, err = pubsub.NewClient(ctx, "")
	if err != nil {
		log.Fatal("Unable to get client for queue")
	}
}

func main() {
	setupQueueDataBase()
	setupGitHubOAuth()
	r := mux.NewRouter()
	r.HandleFunc("/", HandleMain).Methods("GET")
	r.HandleFunc("/webhooks", GitHubWebHookHandler).Methods("POST")
	r.HandleFunc("/auth", GitHubOAuthHandler).Methods("GET")
	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	log.Println(srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
