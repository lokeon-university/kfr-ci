package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
)

var ghOAuth *oauth2.Config

func setupGitHubOAuth() {
	ghOAuth = &oauth2.Config{
		ClientID:     os.Getenv("GH_OCID"),
		ClientSecret: os.Getenv("GH_OCIDS"),
		Endpoint:     gh.Endpoint,
		Scopes: []string{"repo","org:admin"},
	}
}


//GitHubOAuthHandler handle GitHub OAuth event.
func GitHubOAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	token, err := ghOAuth.Exchange(oauth2.NoContext, url.Get("code"))
	if err != nil {
		fmt.Println(token, err)
	}
	// TODO save user into database
}


