package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/lokeon-university/kfr-ci/utils"
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

//GitHubOAuthHandler handle GitHub OAuth event.
func GitHubOAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	token, err := ghOAuth.Exchange(oauth2.NoContext, url.Get("code"))
	if err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.Atoi(url.Get("state"))
	key := datastore.IncompleteKey("users", nil)
	_, err = dbClient.Put(ctx, key, &utils.User{
		ID:    id,
		Token: token.AccessToken,
	})
	if err != nil {
		log.Printf("Unable to add user into db")
	}
	http.Redirect(w, r, "https://t.me/kfr_cibot", 302)
}

//HandleMain func for testind propourse
func HandleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(map[string]string{"data": "Hello World!"})
	w.Write(res)
}
