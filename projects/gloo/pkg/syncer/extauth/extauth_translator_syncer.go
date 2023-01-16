package extauth

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	"go.uber.org/zap"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/hashutils"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Compile-time assertion
var (
	_ syncer.TranslatorSyncerExtension = new(translatorSyncerExtension)
)

var (
	extauthConnectedStateDescription = "zero indicates gloo detected an error with the auth config and did not update its XDS snapshot, check the gloo logs for errors"
	extauthConnectedState            = stats.Int64("glooe.extauth/connected_state", extauthConnectedStateDescription, "1")

	extauthConnectedStateView = &view.View{
		Name:        "glooe.extauth/connected_state",
		Measure:     extauthConnectedState,
		Description: extauthConnectedStateDescription,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{},
	}
)

var (
	emptyTypedResources = map[string]envoycache.Resources{
		extauth.ExtAuthConfigType: {
			Version: "empty",
			Items:   map[string]envoycache.Resource{},
		},
	}
)

func init() {
	_ = view.Register(extauthConnectedStateView)
}

type translatorSyncerExtension struct {
	hasher func(resources []envoycache.Resource) (uint64, error)
}

func NewTranslatorSyncerExtension(_ context.Context, params syncer.TranslatorSyncerExtensionParams) syncer.TranslatorSyncerExtension {
	return &translatorSyncerExtension{
		hasher: params.Hasher,
	}
}

// ID returns the unique identifier for this TranslatorSyncerExtension
// This represents the Key in the SnapshotCache where Sync() will store results
func (s *translatorSyncerExtension) ID() string {
	return runner.ServerRole
}

// Sync processes an ApiSnapshot and updates reports with Errors/Warnings that it encounters
// and updates the SnapshotCache entry if possible
func (s *translatorSyncerExtension) Sync(
	ctx context.Context,
	snap *gloov1snap.ApiSnapshot,
	settings *gloov1.Settings,
	snapshotSetter syncer.SnapshotSetter,
	reports reporter.ResourceReports,
) {
	ctx = contextutils.WithLogger(ctx, "extAuthTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)

	if contextutils.GetLogLevel() == zap.DebugLevel {
		// only hash during debug because of the performance issues surrounding hashing
		snapHash := hashutils.MustHash(snap)
		logger.Debugf("begin auth sync %v (%v proxies, %v upstreams, %v endpoints, %v secrets, %v artifacts, %v auth configs)", snapHash,
			len(snap.Proxies), len(snap.Upstreams), len(snap.Endpoints), len(snap.Secrets), len(snap.Artifacts), len(snap.AuthConfigs))
		defer logger.Debugf("end auth sync %v", snapHash)
	}

	reports.Accept(snap.AuthConfigs.AsInputResources()...)
	reports.Accept(snap.Proxies.AsInputResources()...)

	// To avoid concurrency challenges, we produce an independent xDS Snapshot per translation run
	// This extension is invoked during our translation engine, as well as during validation of Gloo
	// Edge resources. It would be better to have some state that is persisted across runs because some
	// simple caching will save us CPU, but in the meantime we don't do that
	xdsSnapshotProducer := NewProxySourcedXdsSnapshotProducer()
	xdsSnapshot := xdsSnapshotProducer.ProduceXdsSnapshot(ctx, settings, snap, reports)

	var snapshotResources []envoycache.Resource
	for _, cfg := range xdsSnapshot {
		resource := extauth.NewExtAuthConfigXdsResourceWrapper(cfg)
		snapshotResources = append(snapshotResources, resource)
	}

	var extAuthSnapshot envoycache.Snapshot
	snapshotHash, err := s.hasher(snapshotResources)
	if err != nil {
		contextutils.LoggerFrom(ctx).DPanicw("error trying to hash snapshot resources for extauth translation", err)
	}

	if snapshotResources == nil || err != nil {
		// If there are no correctly formatted auth configs, use an empty configuration
		//
		// The SnapshotCache can now differentiate between nil and empty snapshotResources in a snapshot.
		// This was introduced with: https://github.com/solo-io/solo-kit/pull/410
		// A nil resource is not updated, whereas an empty resource is intended to be modified.
		//
		// The extauth service only becomes healthy after it has received auth configuration
		// from Gloo via xDS. Therefore, we must set the auth config resource to empty in the snapshot
		// so that extauth picks up the empty config, and becomes healthy
		extAuthSnapshot = envoycache.NewGenericSnapshot(emptyTypedResources)
	} else {
		snapshotVersion := fmt.Sprintf("%d", snapshotHash)
		extAuthSnapshot = envoycache.NewEasyGenericSnapshot(snapshotVersion, snapshotResources)
	}

	snapshotSetter.SetSnapshot(s.ID(), extAuthSnapshot)
	stats.Record(ctx, extauthConnectedState.M(int64(1)))
}
