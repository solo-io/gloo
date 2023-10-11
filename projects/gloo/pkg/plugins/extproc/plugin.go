package extproc

import (
	"reflect"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/rotisserie/eris"
	gloo_type_matcher_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extproc"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	_ plugins.Plugin            = new(plugin)
	_ plugins.HttpFilterPlugin  = new(plugin)
	_ plugins.VirtualHostPlugin = new(plugin)
	_ plugins.RoutePlugin       = new(plugin)

	NoServerRefErr    = eris.New("no external processing server configured")
	ServerNotFoundErr = func(usRef *core.ResourceRef) error {
		return eris.Errorf("external processing server upstream not found: %s", usRef.String())
	}
	NoFilterStageErr            = eris.New("no filter stage configured")
	MessageTimeoutOutOfRangeErr = func(seconds int64) error {
		return eris.Errorf("timeout duration (%vs) is out of range: must be within range [0s, 3600s]", seconds)
	}
	MaxMessageTimeoutErr = func(seconds int64, maxSeconds int64) error {
		return eris.Errorf("message timeout (%vs) cannot be greater than max message timeout (%vs)", seconds, maxSeconds)
	}
	UnsupportedMatchPatternErr = func(matcher *gloo_type_matcher_v3.StringMatcher) error {
		return eris.Errorf("unsupported string match pattern: %T", matcher.GetMatchPattern())
	}
	DisabledErr = eris.New("the disabled flag can only be set to true")
)

const (
	ExtensionName = "extproc"

	FilterName = "envoy.filters.http.ext_proc"
)

type plugin struct {
	globalSettings *extproc.Settings
}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {
	p.globalSettings = params.Settings.GetExtProc()
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter

	// get the listener-level settings
	disable := listener.GetOptions().GetDisableExtProc()
	listenerSettings := listener.GetOptions().GetExtProc()

	// if global and listener extproc settings are both nil, don't add the filter.
	// the merge function below could have handled this too, but doing the check
	// here to be explicit.
	if p.globalSettings == nil && listenerSettings == nil {
		return filters, nil
	}

	// if the listener explicitly disables extproc, don't add the filter
	if disable.GetValue() {
		return filters, nil
	}

	// do a shallow merge, with the listener-level settings taking precedence
	extProcSettings := mergeExtProcSettings(p.globalSettings, listenerSettings)
	if extProcSettings == nil {
		return filters, nil
	}

	// make sure filter stage is specified
	stage := extProcSettings.GetFilterStage()
	if stage == nil {
		return nil, NoFilterStageErr
	}

	// convert to envoy extproc filter
	extProcFilter, err := toEnvoyExternalProcessor(extProcSettings, params.Snapshot.Upstreams)
	if err != nil {
		return nil, eris.Wrap(err, "generating extproc filter config")
	}

	// create the staged filter
	convertedStage := plugins.ConvertFilterStage(stage)
	stagedFilter, err := plugins.NewStagedFilter(FilterName, extProcFilter, *convertedStage)
	if err != nil {
		return nil, eris.Wrap(err, "generating extproc staged filter")
	}
	filters = append(filters, stagedFilter)

	return filters, nil
}

func (p *plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoy_config_route_v3.VirtualHost) error {
	// get virtualhost-level extproc settings
	virtualHostSettings := in.GetOptions().GetExtProc()

	extProcPerRoute, err := toEnvoyExtProcPerRoute(virtualHostSettings, params.Snapshot.Upstreams)
	if err != nil {
		return eris.Wrap(err, "generating extproc vhost config")
	}
	// no extproc configured
	if extProcPerRoute == nil {
		return nil
	}
	return pluginutils.SetVhostPerFilterConfig(out, FilterName, extProcPerRoute)
}

func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	// get route-level extproc settings
	routeSettings := in.GetOptions().GetExtProc()

	extProcPerRoute, err := toEnvoyExtProcPerRoute(routeSettings, params.Snapshot.Upstreams)
	if err != nil {
		return eris.Wrap(err, "generating extproc route config")
	}
	// no extproc configured
	if extProcPerRoute == nil {
		return nil
	}
	return pluginutils.SetRoutePerFilterConfig(out, FilterName, extProcPerRoute)
}

// Copied/adapted from https://github.com/solo-io/gloo/blob/64804be15298c11663cfb5206f0ac959bfeb4686/projects/gateway/pkg/translator/merge.go
// Merge parent ExtProc settings into child ExtProc settings. Only zero-valued child fields will be
// overwritten by the fields from the parent.
func mergeExtProcSettings(parent *extproc.Settings, child *extproc.Settings) *extproc.Settings {
	if parent == nil {
		return child
	}

	if child == nil {
		// ok to return parent directly because we don't mutate it
		return parent
	}

	childValue, parentValue := reflect.ValueOf(child).Elem(), reflect.ValueOf(parent).Elem()

	for i := 0; i < childValue.NumField(); i++ {
		childField, parentField := childValue.Field(i), parentValue.Field(i)
		utils.ShallowMerge(childField, parentField, false)
	}

	return child
}
