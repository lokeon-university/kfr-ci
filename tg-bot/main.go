package main

import (
	"log"
)

func main() {
	b, err := newBot()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	b.newHandler("/start", b.handleStart)
	b.newHandler("/auth", b.handleOAuth)
	b.newHandler("/help", b.handleHelp)
	b.newHandler("/repos", b.handleRepositories)
	b.start()
}
