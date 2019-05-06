package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
)

var dbClient *firestore.Client
var queueClient *pubsub.Client
var ctx = context.Background()

func setupQueueDataBase() {
	var err error
	conf := &firebase.Config{ProjectID: "kfr-ci"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	dbClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatal("Unable to get client for database")
	}

	queueClient, err = pubsub.NewClient(ctx, "kfr-ci")
	if err != nil {
		log.Fatal("Unable to get client for queue")
	}
}

func main() {
	setupQueueDataBase()
	setupGitHubOAuth()
	r := mux.NewRouter()
	r.HandleFunc("/", statusOK).Methods("GET")
	r.HandleFunc("/webhooks/{id:[0-9]+}", GitHubWebHookHandler).Methods("POST")
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
