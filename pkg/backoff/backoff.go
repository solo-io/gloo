package backoff

import "time"

// Default values for ExponentialBackOff.
const (
	defaultInitialInterval = 500 * time.Millisecond
	defaultMaxElapsedTime  = 60 * time.Second
)

func WithBackoff(fn func() error, stop chan struct{}) error {
	// first try
	if err := fn(); err == nil {
		return nil
	}
	tilNextRetry := defaultInitialInterval
	var elapsed time.Duration
	for {
		select {
		// stopped by another goroutine
		case <-stop:
			return nil
		case <-time.After(tilNextRetry):
			elapsed += tilNextRetry
			tilNextRetry *= 2
			err := fn()
			if err == nil {
				return nil
			}
			if elapsed >= defaultMaxElapsedTime {
				return err
			}
		}
	}
}
