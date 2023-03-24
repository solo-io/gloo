package graphql

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries"

	graphql_plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/graphql"
)

func NewFunctionDiscoveryFactory(opts bootstrap.Opts) fds.FunctionDiscoveryFactory {
	// Allow disabling of fds for GraphQL purposes, default to enabled
	if gqlEnabled := opts.Settings.GetDiscovery().GetFdsOptions().GetGraphqlEnabled(); gqlEnabled != nil && gqlEnabled.GetValue() == false {
		return nil
	}
	return &FunctionDiscoveryFactory{
		DetectionTimeout: time.Second * 15,
	}
}

var _ fds.FunctionDiscoveryFactory = new(FunctionDiscoveryFactory)

func (f *FunctionDiscoveryFactory) FunctionDiscoveryFactoryName() string {
	return "GraphqlFunctionDiscoveryFactory"
}

// FunctionDiscoveryFactory returns a FunctionDiscovery that can be used to discover functions
type FunctionDiscoveryFactory struct {
	DetectionTimeout time.Duration
	Artifacts        v1.ArtifactClient
}

// NewFunctionDiscovery returns a FunctionDiscovery that can be used to discover functions
func (f *FunctionDiscoveryFactory) NewFunctionDiscovery(u *v1.Upstream, clients fds.AdditionalClients) fds.UpstreamFunctionDiscovery {
	return &GraphqlDiscovery{
		upstream:         u,
		client:           clients.GraphqlClient,
		detectionTimeout: f.DetectionTimeout,
	}
}

var _ fds.UpstreamFunctionDiscovery = new(GraphqlDiscovery)

type GraphqlDiscovery struct {
	upstream         *v1.Upstream
	client           v1beta1.GraphQLApiClient
	detectionTimeout time.Duration
}

// IsFunctional returns true if the upstream has already been discovered as a
// GraphQL upstream with introspection.
func (f *GraphqlDiscovery) IsFunctional() bool {
	// as long as the url of the GraphQL Service Spec has not been generated, then it is Functional
	spec := f.getGraphQLSpec(f.upstream)
	return spec != nil && spec.GetEndpoint().GetUrl() != ""
}

// DetectType will ensure that the GraphQL Service Spec of the Upstream is able to produce a schema when queried with a
// GraphQL introspection query.
func (f *GraphqlDiscovery) DetectType(ctx context.Context, url *url.URL) (*plugins.ServiceSpec, error) {
	log := contextutils.LoggerFrom(ctx)
	log.Debugf("attempting to detect GraphQL for the upstream [%s.%s] with url [%s]", f.upstream.GetMetadata().GetNamespace(), f.upstream.GetMetadata().GetName(), url)
	// by default discovery will only look up the path graphql, but users can create there own path, by suppling it in the upstream ServiceSpec manually
	// this will bypass `Detect Type()` and go automatically to the `Detect Functions()`
	url.Path = "graphql"
	switch url.Scheme {
	case "http":
		// do nothing
	case "https":
		// do nothing
	case "tcp":
		// if it is a tcp address, assume it is plain http
		url.Scheme = "http"
	default:
		return nil, fmt.Errorf("unsupported baseurl for graphql discovery %v", url)
	}
	_, err := GetGraphQLSchema(url)
	if err != nil {
		return nil, err
	}

	svcInfo := &plugins.ServiceSpec{
		PluginType: &plugins.ServiceSpec_Graphql{
			Graphql: &graphql_plugins.ServiceSpec{
				Endpoint: &graphql_plugins.ServiceSpec_Endpoint{
					Url: url.String(),
				},
			},
		},
	}
	return svcInfo, nil
}

// DetectFunctions will generate the GraphQL API for the respective upstream.
func (f *GraphqlDiscovery) DetectFunctions(ctx context.Context, url *url.URL, _ func() fds.Dependencies, updatecb func(fds.UpstreamMutator) error) error {
	spec := f.getGraphQLSpec(f.upstream)
	if spec == nil {
		return eris.New(fmt.Sprintf("the upstream %s.%s needs to have a GraphQL Service Spec", f.upstream.GetMetadata().GetNamespace(), f.upstream.GetMetadata().GetName()))
	}
	err := contextutils.NewExponentialBackoff(contextutils.ExponentialBackoff{
		MaxDuration: &f.detectionTimeout,
	}).Backoff(ctx, func(ctx context.Context) error {
		err := f.writeGraphQLAPI()
		return err
	})

	// no need to record any other error. ExpoentialBackoff will log a debug
	// log on all errors, that are not context issues.
	if err != nil {
		if ctx.Err() != nil {
			return multierror.Append(err, ctx.Err())
		}
		// ignore other errors as we would like to continue forever.
	}
	if err := contextutils.Sleep(ctx, 30*time.Second); err != nil {
		return err
	}
	return nil
}

// getGraphQLSpec returns the upstream's GraphQL Service Spec, if it exists
func (f *GraphqlDiscovery) getGraphQLSpec(u *v1.Upstream) *graphql_plugins.ServiceSpec {
	ssg, ok := u.GetUpstreamType().(v1.ServiceSpecGetter)
	if !ok {
		return nil
	}
	graphqlServiceSpec, ok := ssg.GetServiceSpec().GetPluginType().(*plugins.ServiceSpec_Graphql)
	if !ok {
		return nil
	}
	return graphqlServiceSpec.Graphql
}

func (f *GraphqlDiscovery) writeGraphQLAPI() error {
	nn := fmt.Sprintf("%s.%s", f.upstream.GetMetadata().GetNamespace(), f.upstream.GetMetadata().GetName())
	spec := f.getGraphQLSpec(f.upstream)
	// we already added the url to the upstream
	u := spec.Endpoint.GetUrl()
	urlObject, err := url.Parse(u)
	if err != nil {
		return eris.Wrapf(err, "the url did not parse for graphql [%s] %s", u, nn)
	}
	schema, err := GetGraphQLSchema(urlObject)
	if err != nil {
		return eris.Wrap(err, fmt.Sprintf("need the upstreams endpoint to be able to query for graphQL introspection %s", nn))
	}
	options := discoveries.GraphQLApiOptions{
		Remote: &discoveries.RemoteExecutor{
			Path: urlObject.Path,
		},
		Schema: schema,
	}
	resource, err := discoveries.NewGraphQLApi(f.upstream, options)
	if err != nil {
		return eris.Wrapf(err, "unable to instantiate the remote GraphQLApi, : [%s.%s]", f.upstream.GetMetadata().GetNamespace(), f.upstream.GetMetadata().GetName())
	}
	_, err = f.client.Write(resource, clients.WriteOpts{})
	if err != nil {
		return eris.Wrapf(err, fmt.Sprintf("could not write graphQL api for url [%s] %s", urlObject.Hostname(), nn))
	}
	return nil
}
