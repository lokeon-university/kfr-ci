package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	p := &tb.Webhook{
		//Listen: ":" + os.Getenv("PORT"),
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
	r.HandleFunc("/", HandleMain).Methods("GET")
	r.HandleFunc("/status", b.updateStatus).Methods("POST")
	r.HandleFunc("/tg", p.ServeHTTP).Methods("POST")
	srv := &http.Server{
		Addr:         "127.0.0.1:" + os.Getenv("PORT"),
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

//HandleMain func for testind propourse
func HandleMain(w http.ResponseWriter, r *http.Request) {
	log.Println("readed")
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(map[string]string{"data": "Hello World!"})
	w.Write(res)
}
