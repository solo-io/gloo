package transformer

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	v32 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	envoytransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	osTransformation "github.com/solo-io/gloo/projects/gloo/pkg/plugins/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

// The enterprise transformer plugin is an extension of the open source transformation plugin.
// Supports building "transformers" such as the XSLT transformer using the same user-facing trasnformation api
// as the open source HeaderBody or TransformationTemplate transform.
type Plugin struct {
	plugin *osTransformation.Plugin
}

const (
	XsltTransformerFactoryName = "XsltTransformerFactory"
)

func (p *Plugin) PluginName() string {
	return osTransformation.PluginName
}

func (p *Plugin) IsUpgrade() bool {
	return true
}

var _ plugins.Plugin = new(Plugin)
var _ plugins.VirtualHostPlugin = new(Plugin)
var _ plugins.WeightedDestinationPlugin = new(Plugin)
var _ plugins.RoutePlugin = new(Plugin)
var _ plugins.HttpFilterPlugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	if p.plugin == nil {
		p.plugin = osTransformation.NewPlugin()
	}
	err := p.plugin.Init(params)
	if err != nil {
		return err
	}
	p.plugin.TranslateTransformation = translateTransformation
	return nil
}

func (p *Plugin) ProcessVirtualHost(
	params plugins.VirtualHostParams,
	in *v1.VirtualHost,
	out *envoy_config_route_v3.VirtualHost,
) error {
	return p.plugin.ProcessVirtualHost(params, in, out)
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	return p.plugin.ProcessRoute(params, in, out)
}

func (p *Plugin) ProcessWeightedDestination(
	params plugins.RouteParams,
	in *v1.WeightedDestination,
	out *envoy_config_route_v3.WeightedCluster_ClusterWeight,
) error {
	return p.plugin.ProcessWeightedDestination(params, in, out)
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	return p.plugin.HttpFilters(params, listener)
}

// Translates transformations that can include enterprise transformers, such as the XSLT transformer.
func translateTransformation(glooTransform *transformation.Transformation) (*envoytransformation.Transformation, error) {
	if glooTransform == nil {
		return nil, nil
	}
	out := &envoytransformation.Transformation{}
	switch typedTransformation := glooTransform.GetTransformationType().(type) {
	case *transformation.Transformation_XsltTransformation:
		{
			any := utils.MustMessageToAny(typedTransformation.XsltTransformation)
			out.TransformationType = &envoytransformation.Transformation_TransformerConfig{
				TransformerConfig: &v32.TypedExtensionConfig{
					// Arbitrary name for TypedExtension, will error if left empty
					Name:        XsltTransformerFactoryName,
					TypedConfig: any,
				},
			}
		}
	default:
		return osTransformation.TranslateTransformation(glooTransform)
	}
	return out, nil
}
