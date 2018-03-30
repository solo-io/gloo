package main

import (
	"fmt"

	"os"

	"github.com/ilackarms/go-github-webhook-server/github"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/docs/tutorials/source_events_from_github/base"
	"github.com/solo-io/gloo/pkg/log"
)

func main() {
	log.Fatalf("err:", run())
}

func run() error {
	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")

	// Slack client
	client := slack.New(slackToken)

	log.Printf("starting slack client")

	opts := base.Opts{
		ClientID:  os.Getenv("HOSTNAME"),
		ClusterID: "test-cluster",
		NatsURL:   "nats://nats-streaming.default.svc.cluster.local:4222",
		Subject:   "github-webhooks",
		Handler:   handleWatch(client, slackChannel),
	}
	base.Run(opts)
	log.Printf("terminated")
	return nil
}

func handleWatch(client *slack.Client, channel string) func(watch *github.WatchEvent) error {
	return func(watch *github.WatchEvent) error {
		sender := watch.Sender.Login
		text := fmt.Sprintf("Another star for gloo! Thanks, %v !", sender)

		// Send a Message
		_, _, err := client.PostMessage(channel, text, slack.PostMessageParameters{})
		if err != nil {
			return errors.Wrapf(err, "sending message")
		}
		log.Printf("successful post: ", text)

		return nil
	}
}
