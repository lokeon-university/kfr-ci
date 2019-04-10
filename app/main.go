package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lokeon-university/kfr-ci/app/webhooks"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/webhooks", webhooks.GitHubWebHookHandler).Methods("POST")
	r.HandleFunc("/auth", GitHubOAuthHandler).Methods("GET")
	http.Handle("/", r)
}
