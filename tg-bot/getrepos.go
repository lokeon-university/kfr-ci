package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	tb "gopkg.in/tucnak/telebot.v2"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	Config: aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ID"), os.Getenv("AWS_SECRET"), ""),
	},
}))

var dynClient = dynamodb.New(sess)
var tableName = "GHTOKENS"

type repositories struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	SSHURL        string `json:"sshurl,omitempty"`
	NameWithOwner string `json:"name_with_owner,omitempty"`
}

type userDB struct {
	Token  string `json:"token,omitempty"`
	ChatID string `json:"chat_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type query struct {
	Viewer struct {
		ID           string
		Login        string
		Repositories struct {
			Nodes []repositories
		} `graphql:"repositories(last: 20)"`
	}
}

//	CAMBIAR
func runQuery(token string) (data query) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)
	data = query{}
	err := client.Query(context.Background(), &data, nil)
	if err != nil {
		log.Println(err)
	}
	return
}

func getToken(user *tb.User) (string, error) {
	result, err := dynClient.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(strconv.Itoa(user.ID)),
			},
		},
	})
	if err != nil {
		log.Println(err)
		log.Panicln("Error getting user from DB")
	}
	item := userDB{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Panicln("Error Unmarshaling user from DB")
	}
	if item.Token == "" {
		return "", errors.New("User is Not register")
	}
	return item.Token, nil
}

func (q *query) getReposList() [][]tb.InlineButton {
	replyKeys := [][]tb.InlineButton{}
	for _, repo := range q.Viewer.Repositories.Nodes {
		replyKeys = append(replyKeys, []tb.InlineButton{{
			Text: repo.NameWithOwner,
			URL:  repo.SSHURL,
			//TODO change URL for a Callback
		},
		})
	}
	return replyKeys
}

func (kfr *kfrBot) handleRepos() {
	kfr.bot.Handle("/repos", func(m *tb.Message) {
		token, err := getToken(m.Sender)
		if err != nil {
			kfr.bot.Send(m.Sender, "First call /auth")
		} else {
			repos := runQuery(token)
			kfr.bot.Send(m.Sender, "Repositories", &tb.ReplyMarkup{
				InlineKeyboard: repos.getReposList(),
			})
		}
	})
	log.Println("Handled Repos")
}

func (kfr *kfrBot) handleHelp() {
	kfr.bot.Handle("/help", func(m *tb.Message) {
		kfr.bot.Send(m.Sender, `/repos -> Devuelve una lista con los repositorios de una cuenta previamente registrada.
		/auth -> Registra a un usuario mediante su cuenta de Github.`)
	})
	log.Println("Handled Help")
}
