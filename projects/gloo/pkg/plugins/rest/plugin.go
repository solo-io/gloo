package rest

import (
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/proto"
	transformapi "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	glooplugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/transformation"
	transformutils "github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/errors"
)

var (
	_ plugins.Plugin         = new(plugin)
	_ plugins.UpstreamPlugin = new(plugin)
	_ plugins.RoutePlugin    = new(plugin)
)

const (
	ExtensionName = "rest"
)

/*
if this destination spec has rest service spec
this will grab the parameters from the route extension
*/
type plugin struct {
	recordedUpstreams map[string]*glooplugins.ServiceSpec_Rest
}

func NewPlugin() plugins.Plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(_ plugins.InitParams) {
	p.recordedUpstreams = make(map[string]*glooplugins.ServiceSpec_Rest)
}

type UpstreamWithServiceSpec interface {
	GetServiceSpec() *glooplugins.ServiceSpec
}

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, _ *envoy_config_cluster_v3.Cluster) error {
	if withServiceSpec, ok := in.GetUpstreamType().(UpstreamWithServiceSpec); ok {
		serviceSpec := withServiceSpec.GetServiceSpec()
		if serviceSpec == nil {
			return nil
		}

		if serviceSpec.GetPluginType() == nil {
			return nil
		}

		restServiceSpec, ok := serviceSpec.GetPluginType().(*glooplugins.ServiceSpec_Rest)
		if !ok {
			return nil
		}
		if restServiceSpec.Rest == nil {
			return errors.Errorf("%v has an empty rest service spec", in.GetMetadata().Ref())
		}
		p.recordedUpstreams[in.GetMetadata().Ref().Key()] = restServiceSpec
	}
	return nil
}

func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	return pluginutils.MarkPerFilterConfig(params.Ctx, params.Snapshot, in, out, transformation.FilterName,
		func(spec *v1.Destination) (proto.Message, error) {
			// check if it's rest destination
			if spec.GetDestinationSpec() == nil {
				return nil, nil
			}
			restDestinationSpec, ok := spec.GetDestinationSpec().GetDestinationType().(*v1.DestinationSpec_Rest)
			if !ok {
				return nil, nil
			}

			// get upstream
			upstreamRef, err := upstreams.DestinationToUpstreamRef(spec)
			if err != nil {
				contextutils.LoggerFrom(params.Ctx).Error(err)
				return nil, err
			}
			restServiceSpec, ok := p.recordedUpstreams[upstreamRef.Key()]
			if !ok {
				return nil, errors.Errorf("%s does not have a rest service spec", upstreamRef)
			}
			funcname := restDestinationSpec.Rest.GetFunctionName()
			transformationorig := restServiceSpec.Rest.GetTransformations()[funcname]
			if transformationorig == nil {
				return nil, errors.Errorf("unknown function %v", funcname)
			}

			// copy to prevent changing the original in memory.
			transformation := *transformationorig

			// add extensions from the destination spec
			transformation.Extractors, err = transformutils.CreateRequestExtractors(params.Ctx, restDestinationSpec.Rest.GetParameters())
			if err != nil {
				return nil, err
			}

			// get function
			ret := &transformapi.RouteTransformations{
				RequestTransformation: &transformapi.Transformation{
					TransformationType: &transformapi.Transformation_TransformationTemplate{
						TransformationTemplate: &transformation,
					},
				},
			}

			if restDestinationSpec.Rest.GetResponseTransformation() != nil {
				// TODO(yuval-k): should we add \ support response parameters?
				ret.ResponseTransformation = &transformapi.Transformation{
					TransformationType: &transformapi.Transformation_TransformationTemplate{
						TransformationTemplate: restDestinationSpec.Rest.GetResponseTransformation(),
					},
				}
			}

			return ret, nil
		},
	)
}
