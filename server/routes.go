package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
)

var ghOAuth *oauth2.Config

func setupGitHubOAuth() {
	ghOAuth = &oauth2.Config{
		ClientID:     os.Getenv("GH_OCID"),
		ClientSecret: os.Getenv("GH_OCIDS"),
		Endpoint:     gh.Endpoint,
		Scopes:       []string{"repo", "admin:repo_hook", "read:org"},
	}
}

//User GitHub user data.
type User struct {
	Token  string `json:"token,omitempty"`
	ChatID string `json:"chat_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

//GitHubOAuthHandler handle GitHub OAuth event.
func GitHubOAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	log.Println(url.Get("code"), url.Get("state"))
	token, err := ghOAuth.Exchange(oauth2.NoContext, url.Get("code"))
	if err != nil {
		log.Fatal(err)
	}
	ids := strings.Split(url.Get("state"), " ")
	user := User{
		Token:  token.AccessToken,
		ChatID: ids[0],
		UserID: ids[1],
	}
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Println("Error Marshalling User", err.Error())
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("GHTOKENS"),
	}
	_, err = dynClient.PutItem(input)
	if err != nil {
		log.Println("Error adding user to DB")
	}
	log.Println("Added user to DB")
	// result, err := sqsClient.SendMessage(&sqs.SendMessageInput{
	// 	DelaySeconds: aws.Int64(10),
	// 	MessageAttributes: map[string]*sqs.MessageAttributeValue{
	// 		"UserID": &sqs.MessageAttributeValue{
	// 			DataType:    aws.String("String"),
	// 			StringValue: aws.String(ids[1]),
	// 		},
	// 	},
	// 	MessageBody: aws.String("Thanks for register to KFR CI"),
	// 	QueueUrl:    &qURL,
	// })
	// if err != nil {
	// 	log.Fatal("Error", err)
	// }
	// log.Println("Sended Message to Queue", *result.MessageId)
	http.Redirect(w, r, "https://t.me/kfr_cibot", 302)
}

//HandleMain func for testind propourse
func HandleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(map[string]string{"data": "Hello World!"})
	w.Write(res)
}
