package main

import (
	"flag"
	"log"
	"os"

	"github.com/thiagozs/go-slack/options"
	"github.com/thiagozs/go-slack/slackr"
)

var (
	msg     string
	channel string
)

func main() {
	flag.StringVar(&msg, "msg", "", "Message to send")
	flag.StringVar(&channel, "ch", "", "Channel id")

	flag.Parse()

	if msg == "" {
		log.Fatalf("error: msg is required")
	}

	if channel == "" {
		log.Fatalf("error: channel is required")
	}

	token := os.Getenv("SLACKBOT_TOKEN_GHACTIONS")

	lopts := []options.Options{
		options.CfgDebug(false),
		options.CfgToken(token),
	}

	slk, err := slackr.NewSlackClient(lopts)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("Channel: %s\n", channel)
	log.Printf("Message: %s\n", msg)

	if err := slk.SendMessageChannel(channel, msg); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("done\n")

}
