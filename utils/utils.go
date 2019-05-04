package utils

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//String return address of string
func String(str string) *string {
	return &str
}

// User struct data
type User struct {
	ID    int
	Token string
}

// GitHubClient struct to get access to user account
type GitHubClient struct {
	client *github.Client
	ctx    context.Context
}

// NewGitHubClient create new client to access gihub api
func NewGitHubClient(ctx context.Context, token string) *GitHubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &GitHubClient{github.NewClient(tc), ctx}
}

// SetWebhook create hook on given repo
func (gc *GitHubClient) SetWebhook(name, owner, ID string) (int, error) {
	hook := &github.Hook{
		Active: github.Bool(true),
		Events: []string{"push"},
		Config: map[string]interface{}{
			"content_type": "json",
			"url":          os.Getenv("GH_WEBHOOK") + ID,
			"secret":       os.Getenv("GH_APPSECRET"),
		}}
	_, res, err := gc.client.Repositories.CreateHook(gc.ctx, owner, name, hook)
	return res.StatusCode, err
}

// GetRespositories of user
func (gc *GitHubClient) GetRespositories() ([]*github.Repository, error) {
	repos, _, err := gc.client.Repositories.List(gc.ctx, "", nil)
	if err != nil {
		return nil, err
	}
	return repos, nil
}
