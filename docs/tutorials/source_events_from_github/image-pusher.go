package main

import (
	"net/http"

	"github.com/ilackarms/go-github-webhook-server/github"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/docs/tutorials/source_events_from_github/base"
)

func main() {
	opts := base.Opts{
		ClientID:  "image-pusher",
		ClusterID: "test-cluster",
		NatsURL:   "nats://nats-streaming.default.svc.cluster.local:4222",
		Subject:   "github-webhooks",
		Handler:   handleWatch,
	}
	base.Run(opts)
}

func handleWatch(watch *github.WatchEvent) error {
	imgUrl := watch.Sender.AvatarURL
	resp, err := http.Get(imgUrl)
	if err != nil {
		return errors.Wrap(err, "downloading image from url "+imgUrl)
	}
}
