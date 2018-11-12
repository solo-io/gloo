package sqoop

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/sqoop"
	sqoopv1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
)

// ilackarms:: this is a very dumb plugin that looks for the sqoop service annotation and injects the service spec
// with a list of all the existing schemas.

// TODO ilackarms: move to a shared package
const GlooServiceAnnotation = "gloo.solo.io/service"

const GlooServiceAnnotation_Sqoop = "sqoop"

func hasSqoopAnnotation(upstream *v1.Upstream) bool {
	if upstream.Metadata.Annotations == nil {
		return false
	}
	return upstream.Metadata.Annotations[GlooServiceAnnotation] == GlooServiceAnnotation_Sqoop
}

type FunctionDiscoveryFactory struct {
	SchemaClient sqoopv1.SchemaClient
}

func (f *FunctionDiscoveryFactory) NewFunctionDiscovery(u *v1.Upstream) fds.UpstreamFunctionDiscovery {
	return &FunctionDiscovery{
		upstream:     u,
		schemaClient: f.SchemaClient,
	}
}

type FunctionDiscovery struct {
	upstream     *v1.Upstream
	schemaClient sqoopv1.SchemaClient
}

func (d *FunctionDiscovery) DetectFunctions(ctx context.Context, _ *url.URL, dependencies func() fds.Dependencies, updateFn func(fds.UpstreamMutator) error) error {
	// ilackarms: should this ever happen?
	if !hasSqoopAnnotation(d.upstream) {
		return nil
	}

	schemaList, err := d.schemaClient.List(d.upstream.Metadata.Namespace, clients.ListOpts{
		Ctx: ctx,
	})
	if err != nil {
		return err
	}

	var refs []core.ResourceRef
	for _, schema := range schemaList {
		refs = append(refs, schema.Metadata.Ref())
	}
	return updateFn(func(u *v1.Upstream) error {
		upstreamSpec, ok := u.UpstreamSpec.UpstreamType.(v1.ServiceSpecMutator)
		if !ok {
			return errors.New("upstream type does not contain service spec")
		}
		spec := upstreamSpec.GetServiceSpec()
		if spec == nil {
			spec = &plugins.ServiceSpec{}
		}
		sqoopSpec, ok := spec.PluginType.(*plugins.ServiceSpec_Sqoop)
		if !ok {
			sqoopSpec = &plugins.ServiceSpec_Sqoop{
				Sqoop: &sqoop.ServiceSpec{},
			}
		}

		sqoopSpec.Sqoop.Schemas = refs
		spec.PluginType = sqoopSpec

		upstreamSpec.SetServiceSpec(spec)
		return nil
	})
}

func (d *FunctionDiscovery) IsFunctional() bool {
	return hasSqoopAnnotation(d.upstream)
}

func (d *FunctionDiscovery) DetectType(_ context.Context, _ *url.URL) (*plugins.ServiceSpec, error) {
	return nil, nil
}
