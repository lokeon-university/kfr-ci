package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gorilla/mux"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	Config: aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ID"), os.Getenv("AWS_SECRET"), ""),
	},
}))

// Create DynamoDB client
var dynClient = dynamodb.New(sess)

// Create SQS client
var sqsClient = sqs.New(sess)

func main() {
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
