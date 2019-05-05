package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	p := &tb.Webhook{
		Endpoint: &tb.WebhookEndpoint{
			PublicURL: os.Getenv("TG_WEBHOOK"),
		},
	}
	b, err := newBot(p)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	b.newHandler("/start", b.handleStart)
	b.newHandler("/auth", b.handleOAuth)
	b.newHandler("/help", b.handleHelp)
	b.newHandler("/repos", b.handleRepositories)
	go b.start()
	r := mux.NewRouter()
	r.HandleFunc("/", statusOK).Methods("GET")
	r.HandleFunc("/notifications", b.updateStatus).Methods("POST")
	r.HandleFunc("/telegram", p.ServeHTTP).Methods("POST")
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
