package setuputils

import (
	"context"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	"github.com/solo-io/gloo/pkg/utils/statsutils"
	"go.opencensus.io/tag"

	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"go.uber.org/zap"
)

var (
	mSetupsRun = statsutils.MakeSumCounter("gloo.solo.io/setups_run", "The number of times the main setup loop has run")

	mNamespacesWatched      = statsutils.MakeGauge("gloo.solo.io/namespaces_watched", "The number of namespaces watched by the gloo controller")
	namespacesWatchedKey, _ = tag.NewKey("namespaces_watched")
)

// tell us how to setup
type SetupFunc func(ctx context.Context,
	kubeCache kube.SharedCache,
	inMemoryCache memory.InMemoryResourceCache,
	settings *v1.Settings,
	identity leaderelector.Identity) error

const (
	// FirstSyncTimeoutEnv is a Go duration ("10m") bounding how long a held restart waits for the first sync.
	FirstSyncTimeoutEnv = "GLOO_SETUP_FIRST_SYNC_TIMEOUT"

	defaultFirstSyncTimeout = 10 * time.Minute
)

func firstSyncTimeout() time.Duration {
	if d, err := time.ParseDuration(os.Getenv(FirstSyncTimeoutEnv)); err == nil && d > 0 {
		return d
	}
	return defaultFirstSyncTimeout
}

type firstSyncGateKey struct{}

// firstSyncGate is placed on each setup run's ctx; a component arms it via HoldRestartsUntilFirstSync.
type firstSyncGate struct {
	held    atomic.Bool
	ready   chan struct{}
	release func()
}

func newFirstSyncGate() *firstSyncGate {
	gate := &firstSyncGate{ready: make(chan struct{})}
	gate.release = sync.OnceFunc(func() { close(gate.ready) })
	return gate
}

// HoldRestartsUntilFirstSync holds settings and namespace restarts of the setup run owning ctx until the
// returned func is called. Callers in the same run share one release, so the first call frees the hold.
// Outside a setup run it returns a no-op.
func HoldRestartsUntilFirstSync(ctx context.Context) func() {
	gate, ok := ctx.Value(firstSyncGateKey{}).(*firstSyncGate)
	if !ok {
		return func() {}
	}
	gate.held.Store(true)
	return gate.release
}

type SetupSyncer struct {
	settingsRef         *core.ResourceRef
	setupFunc           SetupFunc
	inMemoryCache       memory.InMemoryResourceCache
	identity            leaderelector.Identity
	currentSettingsHash uint64

	// forceEmit re-sends the current snapshot through the setup emitter; nil disables deferral.
	forceEmit chan<- struct{}
	// gate belongs to the current setup run; restarts are held while it is armed and not yet released.
	gate *firstSyncGate
	// restartDeferred is set once a change has been deferred for the current gate, so only one emitWhenReady runs per setup run.
	restartDeferred bool
}

// NewSetupSyncer builds the syncer; forceEmit is the setup emitter's emit channel and may be nil.
func NewSetupSyncer(settingsRef *core.ResourceRef, setupFunc SetupFunc, identity leaderelector.Identity, forceEmit chan<- struct{}) *SetupSyncer {
	return &SetupSyncer{
		settingsRef:   settingsRef,
		setupFunc:     setupFunc,
		inMemoryCache: memory.NewInMemoryResourceCache(),
		identity:      identity,
		forceEmit:     forceEmit,
	}
}

func (s *SetupSyncer) Sync(ctx context.Context, snap *v1.SetupSnapshot) error {
	settings, err := snap.Settings.Find(s.settingsRef.Strings())
	if err != nil {
		return errors.Wrapf(err, "finding bootstrap configuration")
	}
	ctx = settingsutil.WithSettings(ctx, settings)

	contextutils.LoggerFrom(ctx).Debugw("received settings snapshot", zap.Any("settings", settings))

	_, err = settingsutil.UpdateNamespacesToWatch(settings, snap.Kubenamespaces)
	if err != nil {
		return err
	}

	watchedNamespaces := settingsutil.GetNamespacesToWatch(settings)
	contextutils.LoggerFrom(ctx).Debugw("received updated list of namespaces to watch", zap.Any("namespaces", strings.Join(watchedNamespaces, ",")))

	statsutils.Measure(
		ctx,
		mNamespacesWatched,
		int64(len(watchedNamespaces)),
	)

	statsutils.MeasureOne(
		ctx,
		mSetupsRun,
	)

	s.gate, s.restartDeferred = newFirstSyncGate(), false
	ctx = context.WithValue(ctx, firstSyncGateKey{}, s.gate)

	err = s.setupFunc(ctx, kube.NewKubeCache(ctx), s.inMemoryCache, settings, s.identity)
	if err != nil {
		s.gate = nil // a failed run holds nothing; the next change restarts it at once
		return err
	}

	s.currentSettingsHash = settings.MustHash()
	return nil
}

// ShouldSync compares two snapshots and determines whether a sync is needed based on the following conditions
// 1. The settings CR has changed
// 2. A namespace is added / deleted / modified that changes the namespaces to watch
// If the running setup holds restarts (HoldRestartsUntilFirstSync), a needed sync is deferred until it
// releases them, so that setup publishes xDS at least once. The snapshot is then re-delivered via forceEmit.
func (s *SetupSyncer) ShouldSync(ctx context.Context, oldSnapshot, newSnapshot *v1.SetupSnapshot) bool {
	if !s.snapshotChanged(ctx, oldSnapshot, newSnapshot) {
		return false
	}
	// Sync now before the first run, when deferral is disabled, after a failed run, or when nothing holds restarts.
	if s.forceEmit == nil || s.gate == nil || !s.gate.held.Load() {
		return true
	}
	// The hold has been released, so the restart can go ahead.
	select {
	case <-s.gate.ready:
		return true
	default:
	}
	// Defer the change. Later changes during the same hold are covered by the one re-send, which delivers the latest snapshot.
	if !s.restartDeferred {
		s.restartDeferred = true
		contextutils.LoggerFrom(ctx).Infow("deferring setup restart until the current setup has run its first sync")
		go s.emitWhenReady(ctx, s.gate)
	}
	return false
}

// emitWhenReady waits until the setup releases its hold (or the timeout passes), then re-sends the
// snapshot so the event loop calls ShouldSync again for the deferred change.
func (s *SetupSyncer) emitWhenReady(ctx context.Context, gate *firstSyncGate) {
	timeout := firstSyncTimeout()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-gate.ready:
	case <-timer.C:
		contextutils.LoggerFrom(ctx).Warnw("setup did not finish its first sync before the timeout; restarting setup anyway",
			zap.Duration("timeout", timeout))
		// Drop the hold so ShouldSync lets the re-sent snapshot restart setup.
		gate.held.Store(false)
	case <-ctx.Done():
		return
	}
	// The emitter may be busy, so the send also gives up on shutdown.
	select {
	case s.forceEmit <- struct{}{}:
	case <-ctx.Done():
	}
}

func (s *SetupSyncer) snapshotChanged(ctx context.Context, oldSnapshot, newSnapshot *v1.SetupSnapshot) bool {
	// Basic sanity checks. Return a true if there is an error to ensure a sync to get into a good state
	if oldSnapshot == nil || newSnapshot == nil {
		return true
	}

	newSettings, err := newSnapshot.Settings.Find(s.settingsRef.Strings())
	if err != nil {
		return true
	}

	oldSettings, err := oldSnapshot.Settings.Find(s.settingsRef.Strings())
	if err != nil {
		return true
	}

	// A sync is triggered if either :
	// 1. The settings CR is changed
	// 2. A namespace is added / deleted / modified

	// 1. Check whether the settings CR is changed
	if s.currentSettingsHash != newSettings.MustHash() {
		contextutils.LoggerFrom(ctx).Debugw("syncing since settings have been updated")
		return true
	}

	// 2. Check if a namespace is added / deleted / modified
	// If a namespace was modified, check if it changes the namespaces to watch
	contextutils.LoggerFrom(ctx).Debugw("received updated list of namespaces", zap.Any("namespaces", newSnapshot.Kubenamespaces))
	newNamespacesToWatch, err := settingsutil.GenerateNamespacesToWatch(newSettings, newSnapshot.Kubenamespaces)
	if err != nil {
		return true
	}
	oldNamespacesToWatch := settingsutil.GetNamespacesToWatch(oldSettings)
	namespacesChanged := !slices.Equal(newNamespacesToWatch, oldNamespacesToWatch)

	contextutils.LoggerFrom(ctx).Debugw("list of namespaces to watch", zap.Any("oldNamespacesToWatch", oldNamespacesToWatch), zap.Any("newNamespacesToWatch", newNamespacesToWatch), zap.Any("namespacesChanged", namespacesChanged))

	return namespacesChanged
}
