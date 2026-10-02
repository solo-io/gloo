package setuputils

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSetupErrorsStopOnCancellationWithoutChannelClose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- waitForSetupErrors(ctx, make(chan error), true) }()
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("setup ignored cancellation")
	}
}
func TestSetupErrorsPropagateFailure(t *testing.T) {
	errs := make(chan error, 1)
	want := errors.New("setup failed")
	errs <- want
	if err := waitForSetupErrors(context.Background(), errs, true); err != want {
		t.Fatal(err)
	}
}

func TestSetupErrorsContinueWhenNotFatal(t *testing.T) {
	errs := make(chan error, 1)
	errs <- errors.New("recoverable setup error")
	close(errs)
	if err := waitForSetupErrors(context.Background(), errs, false); err != nil {
		t.Fatal(err)
	}
}
