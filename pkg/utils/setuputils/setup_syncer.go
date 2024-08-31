package setuputils

import (
	"context"
	"strings"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/gloo/pkg/utils"
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
	mSetupsRun = utils.MakeSumCounter("gloo.solo.io/setups_run", "The number of times the main setup loop has run")

	mNamespacesWatched      = statsutils.MakeLastValueCounter("gloo.solo.io/namespaces_watched", "The number of namespaces watched by the gloo controller")
	namespacesWatchedKey, _ = tag.NewKey("namespaces_watched")
	namespacedWatchedInView = &view.View{
		Name:        "gloo.solo.io/namespaces_watched",
		Description: "The number of namespaces watched by the gloo controller",
		Measure:     mNamespacesWatched,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{
			namespacesWatchedKey,
		},
	}
)

func init() {
	view.Register(namespacedWatchedInView)
}

// tell us how to setup
type SetupFunc func(ctx context.Context,
	kubeCache kube.SharedCache,
	inMemoryCache memory.InMemoryResourceCache,
	settings *v1.Settings,
	identity leaderelector.Identity) error

type SetupSyncer struct {
	settingsRef          *core.ResourceRef
	setupFunc            SetupFunc
	inMemoryCache        memory.InMemoryResourceCache
	identity             leaderelector.Identity
	previousSettingsHash uint64
}

func NewSetupSyncer(settingsRef *core.ResourceRef, setupFunc SetupFunc, identity leaderelector.Identity) *SetupSyncer {
	return &SetupSyncer{
		settingsRef:   settingsRef,
		setupFunc:     setupFunc,
		inMemoryCache: memory.NewInMemoryResourceCache(),
		identity:      identity,
	}
}

func (s *SetupSyncer) Sync(ctx context.Context, snap *v1.SetupSnapshot) error {
	settings, err := snap.Settings.Find(s.settingsRef.Strings())
	if err != nil {
		return errors.Wrapf(err, "finding bootstrap configuration")
	}
	ctx = settingsutil.WithSettings(ctx, settings)

	currentSettingsHash := settings.MustHash()

	// A sync is triggered if either :
	// - A namespace is added / deleted / modified
	// - The settings CR is changed

	// Check if the namespaces we watch has changed first since this will need to be computed if the settings CR has changed
	// If a namespace was modified, check if it affects the namespaces we should watch
	resyncRequired, err := settingsutil.UpdateNamespacesToWatch(settings, snap.Kubenamespaces)
	if err != nil {
		return err
	}

	// If the namespaces we should watch has not changed, check if the settings CR has changed
	if !resyncRequired {
		if s.previousSettingsHash == currentSettingsHash {
			return nil
		} else {
			contextutils.LoggerFrom(ctx).Debugw("received settings snapshot", zap.Any("settings", settings))
		}
	} else {
		watchedNamespaces := settingsutil.GetNamespacesToWatch(settings)
		contextutils.LoggerFrom(ctx).Debugw("received updated list of namespaces to watch", zap.Any("namespaces", watchedNamespaces))

		watchedNamespacesStr := strings.Join(watchedNamespaces, ",")
		statsutils.Measure(
			ctx,
			mNamespacesWatched,
			int64(len(watchedNamespaces)),
			tag.Insert(namespacesWatchedKey, watchedNamespacesStr),
		)
	}

	defer func() {
		// Run this after the function returns to ensure the hash is changed after we sync
		s.previousSettingsHash = currentSettingsHash
	}()

	utils.MeasureOne(
		ctx,
		mSetupsRun,
	)

	return s.setupFunc(ctx, kube.NewKubeCache(ctx), s.inMemoryCache, settings, s.identity)
}
