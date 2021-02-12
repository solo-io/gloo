package leftmost_xff_address

import (
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/xff_offset"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/leftmost_xff_address"
)

const (
	SoloXffOffsetFilterName = "io.solo.filters.http.solo_xff_offset"
)

var (
	// want this to be the first filter to run
	XffFilterStage = plugins.BeforeStage(plugins.FaultStage)
)

// This plugin configured envoy to interpret the x-forwarded-for http header and set the downstream remote address differently
// The x-forwarded-for header contains a list of IP addresses.
// This filter takes the Offset'th IP address from the left of the header value, and sets that IP as the downstream remote address.
type Plugin struct {
}

var (
	_ plugins.Plugin           = new(Plugin)
	_ plugins.HttpFilterPlugin = new(Plugin)
	_ plugins.Upgradable       = new(Plugin)
)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *Plugin) PluginName() string {
	return leftmost_xff_address.ExtensionName
}

func (p *Plugin) IsUpgrade() bool {
	return true
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {

	var filters []plugins.StagedHttpFilter

	if enableLeftmostXff := listener.GetOptions().GetLeftmostXffAddress(); enableLeftmostXff.GetValue() == false {
		return filters, nil
	}

	leftmostXffConfig := &SoloXffOffset{Offset: 0}
	stagedFilter, err := plugins.NewStagedFilterWithConfig(SoloXffOffsetFilterName, leftmostXffConfig, XffFilterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}
