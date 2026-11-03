// Package translationretry schedules another translation after a transient failure.
package translationretry

import (
	"context"
	"sync"
	"time"
)

// Scheduler coalesces retries per key and backs off repeated failures. Call Reset
// after success or deletion. Callbacks must return promptly; they request fresh
// work rather than translating a captured snapshot.
type Scheduler struct {
	ctx              context.Context
	mu               sync.Mutex
	entries          map[string]*entry
	initial, maximum time.Duration
	retry            func(context.Context, string)
}

type entry struct {
	delay time.Duration
	timer *time.Timer
}

// New creates a scheduler with retries from one second up to thirty seconds.
// Cancelling ctx stops pending retries.
func New(ctx context.Context, retry func(context.Context, string)) *Scheduler {
	return newScheduler(ctx, time.Second, 30*time.Second, retry)
}

func newScheduler(ctx context.Context, initial, maximum time.Duration, retry func(context.Context, string)) *Scheduler {
	s := &Scheduler{ctx: ctx, entries: make(map[string]*entry), initial: initial, maximum: maximum, retry: retry}
	context.AfterFunc(ctx, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		for key, e := range s.entries {
			if e.timer != nil {
				e.timer.Stop()
			}
			delete(s.entries, key)
		}
	})
	return s
}

func (s *Scheduler) Schedule(key string) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx.Err() != nil {
		return
	}
	e := s.entries[key]
	if e == nil {
		e = &entry{delay: s.initial}
		s.entries[key] = e
	}
	if e.timer != nil {
		return
	}
	e.timer = time.AfterFunc(e.delay, func() {
		s.mu.Lock()
		if s.ctx.Err() != nil || s.entries[key] != e {
			s.mu.Unlock()
			return
		}
		// Clear pending before waking the consumer: it may immediately fail
		// again and schedule the next attempt.
		e.timer = nil
		s.mu.Unlock()
		s.retry(s.ctx, key)
	})
	e.delay = min(2*e.delay, s.maximum)
}

func (s *Scheduler) Reset(key string) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if e := s.entries[key]; e != nil {
		if e.timer != nil {
			e.timer.Stop()
		}
		delete(s.entries, key)
	}
}
