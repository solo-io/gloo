package graphql

import (
	"net/http"
	"log"
	"runtime/debug"
	"context"
	"errors"

	"github.com/99designs/gqlgen/handler"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/graph"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
)

func NewHandler(upstreams v1.UpstreamClient, virtualServices gatewayv1.VirtualServiceClient, resolverMaps sqoopv1.ResolverMapClient) http.Handler {
	return handler.GraphQL(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &ApiResolver{
				Upstreams:       upstreams,
				VirtualServices: virtualServices,
				ResolverMaps:    resolverMaps,
				Converter:       &Converter{},
			},
		}),
		// TODO(ilackarms): handle recovering panics
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			// send this panic somewhere
			log.Print(err)
			debug.PrintStack()
			return errors.New("user message on panic")
		}),
	)
}
