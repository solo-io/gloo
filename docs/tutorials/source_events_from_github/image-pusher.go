package main

import (
	"github.com/nats-io/go-nats-streaming"
	"github.com/solo-io/gloo/pkg/log"
)

func main() {
	log.Fatalf("err: %v", run())
}

func run() error {
	conn, err := stan.Connect("test-cluster", "image-pusher",
		stan.NatsURL("nats://nats-streaming.default.svc.cluster.local:4222"))
	if err != nil {
		return err
	}

	log.Printf("Connected to nats-streaming")

	events := make(chan []byte)

	sub, err := conn.Subscribe("github-webhooks", func(msg *stan.Msg) {
		events <- []byte(msg.Data)
	})
	log.Printf("Subscribed to github-webhooks topic")
	defer sub.Close()
	for {
		select {
		case event := <-events:
			if err := handleEvent(event); err != nil {
				log.Warnf("error handling event: %v", err)
			}
		}
	}
}

func handleEvent(rawEventData []byte) error {

}
