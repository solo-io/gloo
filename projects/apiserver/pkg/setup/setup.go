package setup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	apiserver "github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/graph"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/samples"
)

func Setup(port int) error {
	// TODO (ilackarms): pass in the factory with cli flags
	inputFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
		Cache: memory.NewInMemoryResourceCache(),
	})
	upstreams, err := v1.NewUpstreamClient(inputFactory)
	if err != nil {
		return err
	}
	virtualServices, err := gatewayv1.NewVirtualServiceClient(inputFactory)
	if err != nil {
		return err
	}
	resolverMaps, err := sqoopv1.NewResolverMapClient(inputFactory)
	if err != nil {
		return err
	}

	err = addSampleData(upstreams, virtualServices, resolverMaps)
	if err != nil {
		return err
	}

	http.Handle("/", handler.Playground("Solo-ApiServer", "/query"))
	http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{
		Resolvers: apiserver.NewResolvers(upstreams, virtualServices, resolverMaps),
	}),
		handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			rc := graphql.GetResolverContext(ctx)
			fmt.Println("Entered", rc.Object, rc.Field.Name)
			res, err = next(ctx)
			fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
			return res, err
		}),
	))

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func addSampleData(usClient v1.UpstreamClient,
	vsClient gatewayv1.VirtualServiceClient,
	rmClient sqoopv1.ResolverMapClient) error {
	upstreams, virtualServices, resolverMaps := sampleData()
	for _, us := range upstreams {
		_, err := usClient.Write(us, clients.WriteOpts{})
		if err != nil {
			return err
		}
	}
	for _, vs := range virtualServices {
		_, err := vsClient.Write(vs, clients.WriteOpts{})
		if err != nil {
			return err
		}
	}
	for _, rm := range resolverMaps {
		_, err := rmClient.Write(rm, clients.WriteOpts{})
		if err != nil {
			return err
		}
	}
	return nil
}

func sampleData() (v1.UpstreamList, gatewayv1.VirtualServiceList, sqoopv1.ResolverMapList) {
	return samples.Upstreams(), samples.VirtualServices(), sampleResolverMaps()
}

func sampleResolverMaps() sqoopv1.ResolverMapList {
	return sqoopv1.ResolverMapList{
		{
			Metadata: samples.MakeMetadata("resolvermap", "some-namespace", 1),
			Types: map[string]*sqoopv1.TypeResolver{
				"Foo": {
					Fields: map[string]*sqoopv1.FieldResolver{
						"field1": {Resolver: &sqoopv1.FieldResolver_GlooResolver{}},
						"field2": {Resolver: &sqoopv1.FieldResolver_TemplateResolver{}},
					},
				},
				"Bar": {
					Fields: map[string]*sqoopv1.FieldResolver{
						"field1": {Resolver: &sqoopv1.FieldResolver_GlooResolver{}},
						"field2": {Resolver: &sqoopv1.FieldResolver_TemplateResolver{}},
					},
				},
			},
		},
		{
			Metadata: samples.MakeMetadata("resolvermap", "some-namespace", 2),
			Types: map[string]*sqoopv1.TypeResolver{
				"Baz": {
					Fields: map[string]*sqoopv1.FieldResolver{
						"field1": {Resolver: &sqoopv1.FieldResolver_GlooResolver{}},
						"field2": {Resolver: &sqoopv1.FieldResolver_TemplateResolver{}},
					},
				},
				"Qux": {
					Fields: map[string]*sqoopv1.FieldResolver{
						"field1": {Resolver: &sqoopv1.FieldResolver_GlooResolver{}},
						"field2": {Resolver: &sqoopv1.FieldResolver_TemplateResolver{}},
					},
				},
			},
		},
	}
}
