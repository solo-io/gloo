// Package translationretry schedules another translation after a transient failure.
package translationretry

import (
	"context"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/util/workqueue"
)

// Scheduler calls retry for a key after a delay that starts at one second and doubles, up to thirty
// seconds, for each retry of that key until Reset. A key has at most one pending retry. Retries stop
// when ctx ends.
type Scheduler struct {
	queue workqueue.TypedRateLimitingInterface[string]

	mu      sync.Mutex
	pending sets.Set[string]
}

// New returns a Scheduler that calls retry, one key at a time, for each key whose retry is due.
func New(ctx context.Context, retry func(key string)) *Scheduler {
	s := &Scheduler{
		queue: workqueue.NewTypedRateLimitingQueue(
			workqueue.NewTypedItemExponentialFailureRateLimiter[string](time.Second, 30*time.Second)),
		pending: sets.New[string](),
	}
	context.AfterFunc(ctx, s.queue.ShutDown)
	go func() {
		for {
			key, shutdown := s.queue.Get()
			if shutdown {
				return
			}
			// Cleared before calling retry, which may fail again and call Schedule.
			s.mu.Lock()
			s.pending.Delete(key)
			s.mu.Unlock()
			if ctx.Err() == nil {
				retry(key)
			}
			s.queue.Done(key)
		}
	}()
	return s
}

// Schedule calls retry for key after the key's current delay. While a retry for key is pending,
// Schedule does nothing, and the delay does not grow.
func (s *Scheduler) Schedule(key string) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pending.Has(key) {
		return
	}
	s.pending.Insert(key)
	s.queue.AddRateLimited(key)
}

// Reset restores key's initial delay. A retry already pending for key still runs.
func (s *Scheduler) Reset(key string) {
	if s != nil {
		s.queue.Forget(key)
	}
}
