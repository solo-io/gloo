package base

import (
	"encoding/json"

	"github.com/ilackarms/go-github-webhook-server/github"
	"github.com/nats-io/go-nats-streaming"
	"github.com/solo-io/gloo/pkg/log"
)

type Handler func(watch *github.WatchEvent) error

type Opts struct {
	ClientID  string
	ClusterID string
	NatsURL   string
	Subject   string
	Handler   Handler
}

func Run(opts Opts) error {
	conn, err := stan.Connect(opts.ClusterID, opts.ClientID, stan.NatsURL(opts.NatsURL))
	if err != nil {
		return err
	}

	log.Printf("Connected to nats-streaming")

	events := make(chan []byte)
	sub, err := conn.Subscribe(opts.Subject, func(msg *stan.Msg) {
		events <- []byte(msg.Data)
	})
	log.Printf("Subscribed to subject %s", opts.Subject)
	defer sub.Close()
	for {
		select {
		case data := <-events:
			var watch github.WatchEvent
			err := json.Unmarshal(data, &watch)
			if err != nil {
				log.Warnf("could not parse message as github watch event %s: %v", string(data), err)
				continue
			}
			if err := opts.Handler(&watch); err != nil {
				log.Warnf("error handling watch event: %v", err)
			}
		}
	}
}
