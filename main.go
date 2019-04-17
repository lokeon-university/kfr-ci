package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "TOKEN"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, repo := range repos {
		fmt.Println(*repo.Name, *repo.SSHURL)
	}
}
