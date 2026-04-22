package run

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestRunWatcherLoopSerializesDebouncedUpdates(t *testing.T) {
	watcher, watchedFile := newTestWatcher(t)
	defer watcher.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firstUpdateStarted := make(chan struct{}, 1)
	secondUpdateStarted := make(chan struct{}, 1)
	releaseFirstUpdate := make(chan struct{})
	updatesStarted := int32(0)
	activeUpdates := int32(0)
	maxActiveUpdates := int32(0)

	done := make(chan struct{})
	go func() {
		runWatcherLoop(ctx, watcher, func(context.Context) {
			active := atomic.AddInt32(&activeUpdates, 1)
			for {
				currentMax := atomic.LoadInt32(&maxActiveUpdates)
				if active <= currentMax || atomic.CompareAndSwapInt32(&maxActiveUpdates, currentMax, active) {
					break
				}
			}

			updateNumber := atomic.AddInt32(&updatesStarted, 1)
			switch updateNumber {
			case 1:
				firstUpdateStarted <- struct{}{}
				<-releaseFirstUpdate
			case 2:
				secondUpdateStarted <- struct{}{}
			}

			atomic.AddInt32(&activeUpdates, -1)
		})
		close(done)
	}()

	triggerWatcherEvent(t, watchedFile, "first")
	waitForSignal(t, firstUpdateStarted, "first debounced update")

	triggerWatcherEvent(t, watchedFile, "second")
	time.Sleep(sdsUpdateDebounce + 200*time.Millisecond)

	if got := atomic.LoadInt32(&updatesStarted); got != 1 {
		t.Fatalf("expected the second update to wait for the first to finish, started=%d", got)
	}

	close(releaseFirstUpdate)
	waitForSignal(t, secondUpdateStarted, "second debounced update")

	cancel()
	waitForSignal(t, done, "watcher loop shutdown")

	if got := atomic.LoadInt32(&updatesStarted); got != 2 {
		t.Fatalf("expected exactly two debounced updates, got %d", got)
	}
	if got := atomic.LoadInt32(&maxActiveUpdates); got != 1 {
		t.Fatalf("expected updates to run serially, max concurrency=%d", got)
	}
}

func TestRunWatcherLoopSkipsPendingUpdateAfterCancel(t *testing.T) {
	watcher, watchedFile := newTestWatcher(t)
	defer watcher.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firstUpdateStarted := make(chan struct{}, 1)
	releaseFirstUpdate := make(chan struct{})
	updatesStarted := int32(0)

	done := make(chan struct{})
	go func() {
		runWatcherLoop(ctx, watcher, func(context.Context) {
			updateNumber := atomic.AddInt32(&updatesStarted, 1)
			if updateNumber == 1 {
				firstUpdateStarted <- struct{}{}
				<-releaseFirstUpdate
			}
		})
		close(done)
	}()

	triggerWatcherEvent(t, watchedFile, "first")
	waitForSignal(t, firstUpdateStarted, "first debounced update")

	triggerWatcherEvent(t, watchedFile, "second")
	cancel()
	close(releaseFirstUpdate)

	waitForSignal(t, done, "watcher loop shutdown")

	if got := atomic.LoadInt32(&updatesStarted); got != 1 {
		t.Fatalf("expected pending debounced update to be dropped on cancel, got %d updates", got)
	}
}

func newTestWatcher(t *testing.T) (*fsnotify.Watcher, string) {
	t.Helper()

	dir := t.TempDir()
	watchedFile := filepath.Join(dir, "tls.crt")
	if err := os.WriteFile(watchedFile, []byte("initial"), 0o600); err != nil {
		t.Fatalf("writing watched file: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("creating watcher: %v", err)
	}
	if err := watcher.Add(watchedFile); err != nil {
		watcher.Close()
		t.Fatalf("adding file to watcher: %v", err)
	}

	return watcher, watchedFile
}

func triggerWatcherEvent(t *testing.T, watchedFile, contents string) {
	t.Helper()

	if err := os.WriteFile(watchedFile, []byte(contents), 0o600); err != nil {
		t.Fatalf("writing watched file: %v", err)
	}
}

func waitForSignal(t *testing.T, ch <-chan struct{}, description string) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", description)
	}
}
