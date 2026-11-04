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

// firstSyncTimeout returns the duration in FirstSyncTimeoutEnv. An unset value is the default; a value
// that does not parse or is not positive logs a warning and is also the default.
func firstSyncTimeout(ctx context.Context) time.Duration {
	value := os.Getenv(FirstSyncTimeoutEnv)
	if value == "" {
		return defaultFirstSyncTimeout
	}
	d, err := time.ParseDuration(value)
	if err != nil || d <= 0 {
		contextutils.LoggerFrom(ctx).Warnw("ignoring first-sync timeout; expected a positive Go duration such as 10m",
			zap.String("env", FirstSyncTimeoutEnv), zap.String("value", value), zap.Duration("default", defaultFirstSyncTimeout))
		return defaultFirstSyncTimeout
	}
	return d
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

// wait blocks until the gate is released, the first-sync timeout passes, or ctx ends. It returns at once
// if nothing armed the gate.
func (g *firstSyncGate) wait(ctx context.Context) {
	if !g.held.Load() {
		return
	}
	contextutils.LoggerFrom(ctx).Infow("holding setup restarts until the current setup has run its first sync")
	timeout := firstSyncTimeout(ctx)
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-g.ready:
	case <-timer.C:
		contextutils.LoggerFrom(ctx).Warnw("setup did not finish its first sync before the timeout; handling setup changes anyway",
			zap.Duration("timeout", timeout))
	case <-ctx.Done():
	}
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
}

func NewSetupSyncer(settingsRef *core.ResourceRef, setupFunc SetupFunc, identity leaderelector.Identity) *SetupSyncer {
	return &SetupSyncer{
		settingsRef:   settingsRef,
		setupFunc:     setupFunc,
		inMemoryCache: memory.NewInMemoryResourceCache(),
		identity:      identity,
	}
}

// Sync runs the setup function. If a component in the run holds restarts (HoldRestartsUntilFirstSync),
// Sync then waits for it to release them. The setup event loop reads no snapshots while Sync runs, and
// the setup emitter re-sends a snapshot it could not deliver, so a settings or namespace change made
// during the wait restarts setup once Sync returns.
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

	gate := newFirstSyncGate()
	ctx = context.WithValue(ctx, firstSyncGateKey{}, gate)

	err = s.setupFunc(ctx, kube.NewKubeCache(ctx), s.inMemoryCache, settings, s.identity)
	if err != nil {
		return err
	}

	s.currentSettingsHash = settings.MustHash()
	gate.wait(ctx)
	return nil
}

// ShouldSync compares two snapshots and determines whether a sync is needed based on the following conditions
// 1. The settings CR has changed
// 2. A namespace is added / deleted / modified that changes the namespaces to watch
func (s *SetupSyncer) ShouldSync(ctx context.Context, oldSnapshot, newSnapshot *v1.SetupSnapshot) bool {
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
