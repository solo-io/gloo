package setup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	apiserver "github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/graph"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
)

func Setup() error {
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

	http.Handle("/", handler.Playground("Starwars", "/query"))
	http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{
		Resolvers: &apiserver.ApiResolver{
			Upstreams:       upstreams,
			VirtualServices: virtualServices,
			ResolverMaps:    resolverMaps,
			// TODO(ilackarms): just make these private functions, remove converter
			Converter: &apiserver.Converter{},
		},
	}),
		handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			rc := graphql.GetResolverContext(ctx)
			fmt.Println("Entered", rc.Object, rc.Field.Name)
			res, err = next(ctx)
			fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
			return res, err
		}),
	))

	return http.ListenAndServe(":8080", nil)
}
