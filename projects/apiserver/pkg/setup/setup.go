package setup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/rs/cors"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	apiserver "github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/graph"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/setup"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/setup"
	"github.com/solo-io/solo-kit/samples"
)

func Setup(port int, dev bool, glooOpts bootstrap.Opts, gatewayOpts gatewaysetup.Opts, sqoopOpts sqoopsetup.Opts) error {
	// override with memory stuff
	// TODO(ilackarms): move this into a bootstrap package where it can be shared
	if dev {
		inMemory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}

		glooOpts.Secrets = inMemory
		glooOpts.Proxies = inMemory
		glooOpts.Upstreams = inMemory
		gatewayOpts.VirtualServices = inMemory
		sqoopOpts.Schemas = inMemory
		sqoopOpts.ResolverMaps = inMemory

		err := addSampleData(inMemory)
		if err != nil {
			return err
		}
	}

	// serve the query route such that it can be accessed from our UI during development
	corsSettings := cors.New(cors.Options{
		// the development server started by react-scripts defaults to ports 3000, 3001, etc. depending on what's available
		// TODO: Pass debug and CORS urls as flags
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002", "http://localhost:8000", "localhost/:1", "http://localhost:8082"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	http.Handle("/playground", handler.Playground("Solo-ApiServer", "/query"))
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		upstreams, err := v1.NewUpstreamClientWithToken(glooOpts.Upstreams, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		secrets, err := v1.NewSecretClientWithToken(glooOpts.Secrets, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		artifacts, err := v1.NewArtifactClientWithToken(glooOpts.Artifacts, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		virtualServices, err := gatewayv1.NewVirtualServiceClientWithToken(gatewayOpts.VirtualServices, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resolverMaps, err := sqoopv1.NewResolverMapClientWithToken(sqoopOpts.ResolverMaps, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		schemas, err := sqoopv1.NewSchemaClientWithToken(sqoopOpts.Schemas, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := registerAll(upstreams, secrets, artifacts, virtualServices, resolverMaps, schemas); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		corsSettings.Handler(handler.GraphQL(
			graph.NewExecutableSchema(graph.Config{
				Resolvers: apiserver.NewResolvers(upstreams, schemas, artifacts, secrets, virtualServices, resolverMaps),
			}),
			handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
				rc := graphql.GetResolverContext(ctx)
				fmt.Println("Entered", rc.Object, rc.Field.Name)
				res, err = next(ctx)
				fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
				return res, err
			}),
		)).ServeHTTP(w, r)
	})

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

// TODO(ilackarms): move to solo kit
type registrant interface {
	Register() error
}

func registerAll(clients ...registrant) error {
	for _, client := range clients {
		if err := client.Register(); err != nil {
			return err
		}
	}
	return nil
}

func addSampleData(inputFactory factory.ResourceClientFactory) error {
	usClient, err := v1.NewUpstreamClient(inputFactory)
	if err != nil {
		return err
	}
	vsClient, err := gatewayv1.NewVirtualServiceClient(inputFactory)
	if err != nil {
		return err
	}
	rmClient, err := sqoopv1.NewResolverMapClient(inputFactory)
	if err != nil {
		return err
	}
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
	return samples.Upstreams(), samples.VirtualServices(), samples.ResolverMaps()
}
