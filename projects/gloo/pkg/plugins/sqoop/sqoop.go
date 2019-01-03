package sqoop

/*
if this destination spec has rest service spec
this will grab the parameters from the route extention
*/
import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	glooplugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins"
	transformapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/transformation"
	"github.com/solo-io/solo-projects/projects/sqoop/pkg/engine"
)

type UpstreamWithServiceSpec interface {
	GetServiceSpec() *glooplugins.ServiceSpec
}

type plugin struct {
	transformsAdded   *bool
	recordedUpstreams map[core.ResourceRef]*glooplugins.ServiceSpec_Sqoop
	ctx               context.Context
}

func NewPlugin(transformsAdded *bool) plugins.Plugin {
	return &plugin{transformsAdded: transformsAdded}
}

func (p *plugin) Init(params plugins.InitParams) error {
	p.ctx = params.Ctx
	p.recordedUpstreams = make(map[core.ResourceRef]*glooplugins.ServiceSpec_Sqoop)
	return nil
}

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, _ *envoyapi.Cluster) error {
	if withServiceSpec, ok := in.UpstreamSpec.UpstreamType.(UpstreamWithServiceSpec); ok {
		serviceSpec := withServiceSpec.GetServiceSpec()
		if serviceSpec == nil {
			return nil
		}

		if serviceSpec.PluginType == nil {
			return nil
		}

		sqoopServiceSpec, ok := serviceSpec.PluginType.(*glooplugins.ServiceSpec_Sqoop)
		if !ok {
			return nil
		}

		if sqoopServiceSpec.Sqoop == nil {
			return errors.Errorf("%v has an empty sqoop service spec", in.Metadata.Ref())
		}
		p.recordedUpstreams[in.Metadata.Ref()] = sqoopServiceSpec
	}
	return nil
}

func (p *plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	return pluginutils.MarkPerFilterConfig(p.ctx, in, out, transformation.FilterName, func(spec *v1.Destination) (proto.Message, error) {
		// check if it's rest destination
		if spec.DestinationSpec == nil {
			return nil, nil
		}
		sqoopDestinationSpec, ok := spec.DestinationSpec.DestinationType.(*v1.DestinationSpec_Sqoop)
		if !ok {
			return nil, nil
		}
		if sqoopDestinationSpec.Sqoop == nil {
			return nil, errors.Errorf("route %+v has an empty sqoop destination spec", in.Matcher)
		}
		// get upstream
		sqoopServiceSpec, ok := p.recordedUpstreams[spec.Upstream]
		if !ok {
			return nil, errors.Errorf("%v does not have a sqoop service spec", spec.Upstream)
		}

		var schemaNameValid bool
		for _, n := range sqoopServiceSpec.Sqoop.Schemas {
			if n.Equal(sqoopDestinationSpec.Sqoop.Schema) {
				schemaNameValid = true
				break
			}
		}
		if !schemaNameValid {
			return nil, errors.Errorf("schema %v not found, available: %v",
				sqoopDestinationSpec.Sqoop.Schema, sqoopServiceSpec.Sqoop.Schemas)
		}
		*p.transformsAdded = true

		// transform all routes to be the correct path
		var path string
		if sqoopDestinationSpec.Sqoop.Playground {
			path = engine.SqoopPlaygroundPath(sqoopDestinationSpec.Sqoop.Schema)
		} else {
			path = engine.SqoopQueryPath(sqoopDestinationSpec.Sqoop.Schema)
		}

		return &transformapi.RouteTransformations{
			RequestTransformation: &transformapi.Transformation{
				TransformationType: &transformapi.Transformation_TransformationTemplate{
					TransformationTemplate: &transformapi.TransformationTemplate{
						Headers: map[string]*transformapi.InjaTemplate{
							":path": {Text: path},
						},
					},
				},
			},
		}, nil
	})
}
