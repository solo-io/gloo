package translationretry

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestRetryBackoffResetAndCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		calls := make(chan string, 10)
		s := New(ctx, func(key string) { calls <- key })
		expectCall := func(key string, delay time.Duration) {
			t.Helper()
			time.Sleep(delay - time.Nanosecond)
			synctest.Wait()
			select {
			case got := <-calls:
				t.Fatalf("retry for %s fired before %s", got, delay)
			default:
			}
			time.Sleep(time.Nanosecond)
			synctest.Wait()
			select {
			case got := <-calls:
				if got != key {
					t.Fatalf("retry key = %s, want %s", got, key)
				}
			default:
				t.Fatalf("retry for %s did not fire after %s", key, delay)
			}
		}
		expectNoCall := func() {
			t.Helper()
			time.Sleep(time.Minute)
			synctest.Wait()
			select {
			case key := <-calls:
				t.Fatalf("unexpected retry for %s", key)
			default:
			}
		}

		// Backoff: each retry doubles the delay, up to 30s.
		for _, delay := range []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second, 30 * time.Second, 30 * time.Second} {
			s.Schedule("proxy")
			expectCall("proxy", delay)
		}

		// Reset: the next retry is back to 1s.
		s.Reset("proxy")
		s.Schedule("proxy")
		expectCall("proxy", time.Second)

		// A Schedule while a retry is pending adds no second retry and does not grow the delay.
		s.Reset("proxy")
		s.Schedule("proxy")
		s.Schedule("proxy")
		expectCall("proxy", time.Second)
		expectNoCall()
		s.Schedule("proxy")
		expectCall("proxy", 2*time.Second)

		// Cancellation: a retry pending when ctx ends never runs.
		s.Schedule("shutdown")
		cancel()
		expectNoCall()
	})
}

// A retry that fails again calls Schedule from inside the retry callback, as the syncers do.
func TestRetryCanBeScheduledByConsumerImmediately(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		attempts := 0
		var s *Scheduler
		// The first two attempts fail and reschedule. Three attempts take 1s + 2s + 4s = 7s.
		s = New(ctx, func(key string) {
			attempts++
			if attempts < 3 {
				s.Schedule(key)
			}
		})
		s.Schedule("proxy")
		time.Sleep(7 * time.Second)
		synctest.Wait()
		if attempts != 3 {
			t.Fatalf("got %d attempts, want 3", attempts)
		}
	})
}
