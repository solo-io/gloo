package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// fakeEnvoy writes an executable that ends itself with the given shell command.
func fakeEnvoy(t *testing.T, script string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "envoy")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunEnvoyValidateClassifiesExits(t *testing.T) {
	tests := []struct {
		name            string
		script          string
		wantErr         bool
		wantInterrupted bool
		wantCrashed     bool
	}{
		{name: "valid config", script: "exit 0"},
		{name: "invalid config", script: "echo bad config; exit 1", wantErr: true},
		{name: "killed from outside", script: "kill -KILL $$", wantErr: true, wantInterrupted: true},
		{name: "terminated from outside", script: "kill -TERM $$", wantErr: true, wantInterrupted: true},
		{name: "segfault", script: "kill -SEGV $$", wantErr: true, wantInterrupted: true, wantCrashed: true},
		{name: "abort", script: "kill -ABRT $$", wantErr: true, wantInterrupted: true, wantCrashed: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunEnvoyValidate(context.Background(), fakeEnvoy(t, tt.script), "{}")

			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got := errors.Is(err, ErrValidationInterrupted); got != tt.wantInterrupted {
				t.Fatalf("errors.Is(err, ErrValidationInterrupted) = %v, want %v (err: %v)", got, tt.wantInterrupted, err)
			}
			if got := errors.Is(err, ErrValidationCrashed); got != tt.wantCrashed {
				t.Fatalf("errors.Is(err, ErrValidationCrashed) = %v, want %v (err: %v)", got, tt.wantCrashed, err)
			}
		})
	}
}

func TestRunEnvoyValidateLogsCrashes(t *testing.T) {
	tests := []struct {
		name       string
		script     string
		wantSignal string
	}{
		{name: "segfault", script: "kill -SEGV $$", wantSignal: "segmentation fault"},
		{name: "killed from outside", script: "kill -KILL $$"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture the error logs RunEnvoyValidate writes.
			core, logs := observer.New(zap.ErrorLevel)
			ctx := contextutils.WithExistingLogger(context.Background(), zap.New(core).Sugar())

			// Run a fake envoy that ends itself with the signal.
			_ = RunEnvoyValidate(ctx, fakeEnvoy(t, tt.script), "{}")

			// Only a crash is logged.
			entries := logs.All()
			if tt.wantSignal == "" {
				if len(entries) != 0 {
					t.Fatalf("want no error logs, got %d: %q", len(entries), entries[0].Message)
				}
				return
			}
			if len(entries) != 1 {
				t.Fatalf("want 1 error log, got %d", len(entries))
			}
			if sig := entries[0].ContextMap()["signal"]; sig != tt.wantSignal {
				t.Fatalf("logged signal = %v, want %q", sig, tt.wantSignal)
			}
		})
	}
}

// The context ends before the call, so envoy never starts.
func TestRunEnvoyValidateWithEndedContextIsInterrupted(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := RunEnvoyValidate(ctx, fakeEnvoy(t, "exit 0"), "{}")

	if !errors.Is(err, ErrValidationInterrupted) {
		t.Fatalf("a validation whose context already ended must be an interruption, got %v", err)
	}
}
