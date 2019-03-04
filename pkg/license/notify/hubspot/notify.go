package hubspot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/solo-io/solo-projects/pkg/license/model"
	"github.com/solo-io/solo-projects/pkg/license/notify"
)

const LicenseKey = "license_key"

var _ notify.Notifier = new(HubspotNotifier)

type HubspotNotifier struct {
	Hubspot *Hubspot
}

func (h *HubspotNotifier) Notify(ctx context.Context, ua model.UserInfo, key string) error {

	c := ContactUpdate{
		Properties: []PropertyUpdate{{
			Property: LicenseKey,
			Value:    key,
		}},
	}
	updateContact := func() error {
		_, err := h.Hubspot.UpsertContact(ctx, ua.Email, c)
		return err
	}
	err := updateContact()

	if err == RatelimitError {
		if err := h.retryRateLimit(updateContact); err != nil {
			// TODO: log log log log log!!!
		}
	}

	return err
}

func (h *HubspotNotifier) retryRateLimit(f func() error) error {
	ratelimiterRetrierOnce.Do(ratelimiterRetrierInstance.start)
	return ratelimiterRetrierInstance.add(f)
}

type ratelimiterRetrier struct {
	queue chan func() error
}

func (r *ratelimiterRetrier) add(f func() error) error {
	select {
	case r.queue <- f:
		return nil
	default:
		return fmt.Errorf("retry queue full")
	}
}

func (r *ratelimiterRetrier) start() {
	go func() {
		attempts := 1

		for f := range r.queue {
			err := f()
			if err == RatelimitError {
				time.Sleep(time.Second * time.Duration(attempts))
				attempts *= 2
			} else if err == nil {
				attempts = 1
			} else {
				// TODO: log error
			}
		}

	}()
}

var ratelimiterRetrierInstance = ratelimiterRetrier{
	queue: make(chan func() error, 10000),
}
var ratelimiterRetrierOnce sync.Once
