package main

import (
	"context"
	"errors"
	"log"
	"os"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/go-github/github"
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

type githubClient struct {
	client *github.Client
	ctx    context.Context
	token string
}

func createClient(ctx context.Context, token string) *githubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &githubClient{github.NewClient(tc), ctx,token}

}

func (gc *githubClient) getRepos() [][]tb.InlineButton {
	replyKeys := [][]tb.InlineButton{}
	repos, _, err := gc.client.Repositories.List(gc.ctx, "", nil)
	if err != nil {
		log.Println(err)
	}
	for _, repo := range repos {
		log.Println(*repo.Name, *repo.SSHURL)
		replyKeys = append(replyKeys, []tb.InlineButton{{
			Unique: "repos",
			Text: *repo.Name,
			Data:  fmt.Sprintf("%s %s",gc.token,*repo.SSHURL),
		},
		})
	}
	return replyKeys
}

type userDB struct {
	Token  string `json:"token,omitempty"`
	ChatID string `json:"chat_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
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

func (kfr *kfrBot) handleRepos() {
	kfr.bot.Handle("/repos", func(m *tb.Message) {
		token, err := getToken(m.Sender)
		if err != nil {
			kfr.bot.Send(m.Sender, "First call /auth")
		} else {
			gc := createClient(context.Background(), token)
			kfr.bot.Send(m.Sender, "Repositories", &tb.ReplyMarkup{
				InlineKeyboard: gc.getRepos(),
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

func (kfr *kfrBot) handleRepoResponse() {
	buttons	:= tb.InlineButton{Unique:"repos"}
	kfr.bot.Handle(&buttons, func(c*tb.Callback){
		data := strings.Split(c.Data," ")
		kfr.bot.Respond(c, &tb.CallbackResponse{})
	})
}
