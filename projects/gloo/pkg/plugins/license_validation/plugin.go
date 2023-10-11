package license_validation

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/pkg/license"
)

var (
	_ plugins.Plugin = new(plugin)
)

// The license validation plugin validates that a valid Enterprise license is present and logs a warning at the error
// level if not
// This is done in the Init so that it occurs on every translation loop
// This is intended to be sufficiently noisy that users can reasonably be expected to know they are using an invalid license
type plugin struct {
	licensedFeatureProvider *license.LicensedFeatureProvider
}

func NewPlugin(lfp *license.LicensedFeatureProvider) *plugin {
	return &plugin{
		licensedFeatureProvider: lfp,
	}
}

func (p *plugin) Name() string {
	return "license-validator"
}

func (p *plugin) Init(params plugins.InitParams) {
	state := p.licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
	if state.Reason != "" {
		contextutils.LoggerFrom(params.Ctx).Errorf("LICENSE ERROR: %v\n", state.Reason)
	}
}
