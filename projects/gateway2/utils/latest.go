package utils

import (
	"context"
	"sync"
)

// Latest retains a complete snapshot for independent, replaceable consumers.
// Unlike a queue, reading does not remove the value. Each consumer tracks the
// version it has successfully applied and starts at zero to replay current state.
// Values must be immutable after publication; consumers must copy before mutation.
// Use NewLatest to initialize it.
type Latest[T any] struct {
	mu      sync.Mutex
	value   T
	version uint64
	changed chan struct{}
}

func NewLatest[T any]() *Latest[T] {
	return &Latest[T]{changed: make(chan struct{})}
}

// Publish replaces the current snapshot, including an explicitly empty value.
func (l *Latest[T]) Publish(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.value = value
	l.version++
	close(l.changed)
	l.changed = make(chan struct{})
}

// Wait returns the latest snapshot newer than after, or waits for publication.
// Cancellation takes precedence over an already available snapshot.
func (l *Latest[T]) Wait(ctx context.Context, after uint64) (T, uint64, error) {
	for {
		l.mu.Lock()
		if err := ctx.Err(); err != nil {
			l.mu.Unlock()
			var zero T
			return zero, 0, err
		}
		if l.version > after {
			value, version := l.value, l.version
			l.mu.Unlock()
			return value, version, nil
		}
		changed := l.changed
		l.mu.Unlock()
		select {
		case <-ctx.Done():
			var zero T
			return zero, 0, ctx.Err()
		case <-changed:
		}
	}
}
