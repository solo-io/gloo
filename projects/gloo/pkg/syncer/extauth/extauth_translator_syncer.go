package extauth

import (
	"context"
	"fmt"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/runner"

	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"go.uber.org/zap"

	"github.com/mitchellh/hashstructure"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/go-utils/contextutils"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	extAuthPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

type ExtAuthTranslatorSyncerExtension struct {
	reporter reporter.Reporter
}

func NewTranslatorSyncerExtension(params syncer.TranslatorSyncerExtensionParams) *ExtAuthTranslatorSyncerExtension {
	return &ExtAuthTranslatorSyncerExtension{reporter: params.Reporter}
}

// TODO(marco): report errors on auth config resources once we have the strongly typed API. Currently it is not possible
//  to do this consistently, since we need to parse the raw extension to get to the auth config, an operation that might itself fail.
func (s *ExtAuthTranslatorSyncerExtension) Sync(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache envoycache.SnapshotCache) error {
	ctx = contextutils.WithLogger(ctx, "extAuthTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	snapHash := hashutils.MustHash(snap)
	logger.Infof("begin auth sync %v (%v proxies, %v upstreams, %v endpoints, %v secrets, %v artifacts, %v auth configs)", snapHash,
		len(snap.Proxies), len(snap.Upstreams), len(snap.Endpoints), len(snap.Secrets), len(snap.Artifacts), len(snap.AuthConfigs))
	defer logger.Infof("end auth sync %v", snapHash)

	return s.SyncAndSet(ctx, snap, xdsCache)
}

type SnapshotSetter interface {
	SetSnapshot(node string, snapshot envoycache.Snapshot) error
}

func (s *ExtAuthTranslatorSyncerExtension) SyncAndSet(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache SnapshotSetter) error {
	helper := newHelper()
	reports := make(reporter.ResourceReports)
	reports.Accept(snap.AuthConfigs.AsInputResources()...)

	for _, cfg := range snap.AuthConfigs {
		// Validate auth config
		extAuthPlugin.ValidateAuthConfig(cfg, reports)
		configRef := cfg.GetMetadata().Ref()
		if err := helper.processAuthExtension(ctx, snap, &configRef); err != nil {
			reports.AddError(cfg, err)
			return err
		}
	}

	var resources []envoycache.Resource
	for _, cfg := range helper.translatedConfigs {
		resource := extauth.NewExtAuthConfigXdsResourceWrapper(cfg)
		resources = append(resources, resource)
	}
	h, err := hashstructure.Hash(resources, nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).With(zap.Error(err)).DPanic("error hashing ext auth")
		return err
	}
	extAuthSnapshot := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", h), resources)
	_ = xdsCache.SetSnapshot(runner.ExtAuthServerRole, extAuthSnapshot)
	if err := s.reporter.WriteReports(ctx, reports, nil); err != nil {
		contextutils.LoggerFrom(ctx).Warnf("Failed writing report for auth configs: %v", err)
		return errors.Wrapf(err, "writing reports")
	}
	return nil
}

// This translation helper contains a map where each key is the unique identifier of an AuthConfig and the corresponding
// value is the translated config. We use it avoid translating the same configuration multiple times.
type helper struct {
	translatedConfigs map[string]*extauth.ExtAuthConfig
}

func newHelper() *helper {
	return &helper{
		translatedConfigs: make(map[string]*extauth.ExtAuthConfig),
	}
}

func (h *helper) processAuthExtension(ctx context.Context, snap *gloov1.ApiSnapshot, configRef *core.ResourceRef) error {
	if configRef == nil {
		// Just return if there is nothing to translate
		return nil
	}

	// Don't perform duplicate work if we already have translated this resource
	if _, ok := h.translatedConfigs[configRef.Key()]; ok {
		return nil
	}

	translatedConfig, err := extAuthPlugin.TranslateExtAuthConfig(ctx, snap, configRef)
	if err != nil {
		return err
	} else if translatedConfig == nil {
		// Do nothing if config is empty or consists only of custom auth configs
		return nil
	}

	h.translatedConfigs[configRef.Key()] = translatedConfig
	return nil
}
