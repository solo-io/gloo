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
		s := newScheduler(ctx, 10*time.Millisecond, 40*time.Millisecond, func(_ context.Context, key string) { calls <- key })
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
		for _, delay := range []time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 40 * time.Millisecond, 40 * time.Millisecond} {
			s.Schedule("proxy")
			s.Schedule("proxy") // Duplicate requests must not create duplicate retries.
			expectCall("proxy", delay)
		}
		s.Reset("proxy")
		s.Schedule("proxy")
		expectCall("proxy", 10*time.Millisecond)
		s.Schedule("deleted")
		s.Reset("deleted")
		s.Schedule("shutdown")
		cancel()
		time.Sleep(time.Second)
		synctest.Wait()
		select {
		case key := <-calls:
			t.Fatalf("retry for %s fired after reset or cancellation", key)
		default:
		}
	})
}

func TestRetryCanBeScheduledByConsumerImmediately(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		attempts := 0
		var s *Scheduler
		s = newScheduler(ctx, time.Second, time.Minute, func(_ context.Context, key string) {
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
