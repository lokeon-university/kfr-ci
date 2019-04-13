package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/lokeon-university/kfr-ci/app/webhooks"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

// Create DynamoDB client
var svc = dynamodb.New(sess)

func main() {
	setupGitHubOAuth()
	r := mux.NewRouter()
	r.HandleFunc("/", HandleMain).Methods("GET")
	r.HandleFunc("/webhooks", webhooks.GitHubWebHookHandler).Methods("POST")
	r.HandleFunc("/auth", GitHubOAuthHandler).Methods("GET")
	srv := &http.Server{
		Addr:         "0.0.0.0:5000",
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
