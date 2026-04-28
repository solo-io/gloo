package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	"github.com/solo-io/gloo/projects/sds/pkg/server"
	"github.com/solo-io/go-utils/contextutils"
)

// sdsUpdateDebounce is the quiet period after the last fsnotify event before reloading
// certs from disk, so writers (e.g. Istio) can finish updating key and cert files.
const sdsUpdateDebounce = 500 * time.Millisecond

func Run(ctx context.Context, secrets []server.Secret, sdsClient, sdsServerAddress string) error {
	ctx, cancel := context.WithCancel(ctx)

	// Set up the gRPC server
	sdsServer := server.SetupEnvoySDS(secrets, sdsClient, sdsServerAddress)
	// Run the gRPC Server
	serverStopped, err := sdsServer.Run(ctx) // runs the grpc server in internal goroutines
	if err != nil {
		cancel()
		return err
	}

	// Initialize the SDS config
	err = sdsServer.UpdateSDSConfig(ctx)
	if err != nil {
		cancel()
		return err
	}

	// create a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		cancel()
		return err
	}
	defer watcher.Close()

	// Wire in signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigs)

	// call watchFiles here before calling it again in the
	// goroutine, otherwise the two calls may race when adding
	// watches to `watcher`
	watchFiles(ctx, watcher, secrets)

	go func() {
		runWatcherLoop(ctx, watcher, func(ctx context.Context) {
			if err := sdsServer.UpdateSDSConfig(ctx); err != nil {
				contextutils.LoggerFrom(ctx).Warnw("failed to update SDS config after cert file change", zap.Error(err))
			}
			if ctx.Err() != nil {
				return
			}
			watchFiles(ctx, watcher, secrets)
		})
	}()

	select {
	case <-sigs:
	case <-ctx.Done():
	}
	cancel()
	select {
	case <-serverStopped:
		return nil
	case <-time.After(3 * time.Second):
		return nil
	}
}

func runWatcherLoop(ctx context.Context, watcher *fsnotify.Watcher, onDebouncedUpdate func(context.Context)) {
	debounceTimer := time.NewTimer(sdsUpdateDebounce)
	stopAndDrainTimer(debounceTimer)
	defer debounceTimer.Stop()

	var pendingUpdate bool
	for {
		select {
		case event := <-watcher.Events:
			contextutils.LoggerFrom(ctx).Infow("received event", zap.Any("event", event))
			pendingUpdate = true
			stopAndDrainTimer(debounceTimer)
			debounceTimer.Reset(sdsUpdateDebounce)
		case err := <-watcher.Errors:
			contextutils.LoggerFrom(ctx).Warnw("Received error from file watcher", zap.Error(err))
		case <-debounceTimer.C:
			if !pendingUpdate {
				continue
			}
			pendingUpdate = false
			if ctx.Err() != nil {
				return
			}
			onDebouncedUpdate(ctx)
		case <-ctx.Done():
			stopAndDrainTimer(debounceTimer)
			return
		}
	}
}

func stopAndDrainTimer(timer *time.Timer) {
	if timer == nil {
		return
	}
	if timer.Stop() {
		return
	}
	select {
	case <-timer.C:
	default:
	}
}

func watchFiles(ctx context.Context, watcher *fsnotify.Watcher, secrets []server.Secret) {
	for _, s := range secrets {
		contextutils.LoggerFrom(ctx).Infow("watcher started", zap.String("sslKeyFile", s.SslKeyFile), zap.String("sshCertFile", s.SslCertFile), zap.String("sslCaFile", s.SslCaFile))
		if err := watcher.Add(s.SslKeyFile); err != nil {
			contextutils.LoggerFrom(ctx).Warn(zap.Error(err))
		}
		if err := watcher.Add(s.SslCertFile); err != nil {
			contextutils.LoggerFrom(ctx).Warn(zap.Error(err))
		}
		if err := watcher.Add(s.SslCaFile); err != nil {
			contextutils.LoggerFrom(ctx).Warn(zap.Error(err))
		}
	}
}
