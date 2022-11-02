package validation

import (
	"context"

	ratelimit "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO-JAKE add in extension Supported GVKs, or get rid of supported GVKs all together...
var GvkToSupportedDeleteExtensionResources = map[schema.GroupVersionKind]bool{
	ratelimit.RateLimitConfigGVK: true,
}

var GvkToSupportedExtensionResources = map[schema.GroupVersionKind]bool{
	ratelimit.RateLimitConfigGVK: true,
}

// NewValidator will create a new validator for the purpose of validating extensions.
func NewValidator(extensions []syncer.TranslatorSyncerExtension, settings *v1.Settings) validator {
	return validator{
		extensions: extensions,
		settings:   settings,
	}
}

type Validator interface {
	// Validate will validate the proxies of a snapshot, and return the resource reports corresponding to the
	// resources in the snapshot.
	Validate(ctx context.Context, snapshot *gloosnapshot.ApiSnapshot) reporter.ResourceReports
}

type validator struct {
	extensions []syncer.TranslatorSyncerExtension
	settings   *v1.Settings
}

// Validate will sync the extensions with the snapshot. This uses a no-op snapshot setter, so that no changes to the
// snapshot occur.
func (v *validator) Validate(ctx context.Context, snapshot *gloosnapshot.ApiSnapshot) reporter.ResourceReports {
	// TODO-JAKE how does this get done with the translation?
	// TODO-JAKE make some tests for this?
	noOp := &syncer.NoOpSnapshotSetter{}
	reports := reporter.ResourceReports{}
	for _, ex := range v.extensions {
		// TODO-JAKE we can add in the settings for the ext-auth syncer
		intermediateReports := make(reporter.ResourceReports)
		ex.Sync(ctx, snapshot, v.settings, noOp, intermediateReports)
		reports.Merge(intermediateReports)
	}
	return reports
}
