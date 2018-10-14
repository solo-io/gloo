package setup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/rs/cors"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	apiserver "github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/graph"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/syncer"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"sync"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

func Setup(ctx context.Context, port int, dev bool, settings v1.SettingsClient, glooOpts bootstrap.Opts, gatewayOpts gatewaysetup.Opts, sqoopOpts sqoopsetup.Opts) error {
	// initial resource registration
	upstreams, err := v1.NewUpstreamClient(glooOpts.Upstreams)
	if err != nil {
		return err
	}
	if err := upstreams.Register(); err != nil {
		return err
	}
	secrets, err := v1.NewSecretClient(glooOpts.Secrets)
	if err != nil {
		return err
	}
	if err := secrets.Register(); err != nil {
		return err
	}
	artifacts, err := v1.NewArtifactClient(glooOpts.Artifacts)
	if err != nil {
		return err
	}
	if err := artifacts.Register(); err != nil {
		return err
	}
	virtualServices, err := gatewayv1.NewVirtualServiceClient(gatewayOpts.VirtualServices)
	if err != nil {
		return err
	}
	if err := virtualServices.Register(); err != nil {
		return err
	}
	resolverMaps, err := sqoopv1.NewResolverMapClient(sqoopOpts.ResolverMaps)
	if err != nil {
		return err
	}
	if err := resolverMaps.Register(); err != nil {
		return err
	}
	schemas, err := sqoopv1.NewSchemaClient(sqoopOpts.Schemas)
	if err != nil {
		return err
	}
	if err := schemas.Register(); err != nil {
		return err
	}

	// serve the query route such that it can be accessed from our UI during development
	corsSettings := cors.New(cors.Options{
		// the development server started by react-scripts defaults to ports 3000, 3001, etc. depending on what's available
		// TODO: Pass debug and CORS urls as flags
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002", "http://localhost:8000", "localhost/:1", "http://localhost:8082", "http://localhost"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// TODO ilackarms: refactor this to be less kube-specific, move to its own setup syncer or special ClientFactory
	lock := sync.Mutex{}
	insertCacheIntoRcFactory := func(perTokenCaches map[string]*kube.KubeCache, token string, fact factory.ResourceClientFactory) {
		lock.Lock()
		defer lock.Unlock()
		switch fact := fact.(type) {
		case *factory.KubeResourceClientFactory:
			cacheForToken, ok := perTokenCaches[token]
			if !ok {
				cacheForToken = kube.NewKubeCache()
				perTokenCaches[token] = cacheForToken
			}
			fact.SharedCache = cacheForToken
		default:
			contextutils.LoggerFrom(ctx).Warnf("not initializing a per-token cache for resource client %v", fact)
		}
	}
	perTokenCaches := make(map[string]*kube.KubeCache)

	http.Handle("/playground", handler.Playground("Solo-ApiServer", "/query"))
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		insertCacheIntoRcFactory(perTokenCaches, token, glooOpts.Upstreams)
		upstreams, err := v1.NewUpstreamClientWithToken(glooOpts.Upstreams, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		insertCacheIntoRcFactory(perTokenCaches, token, glooOpts.Secrets)
		secrets, err := v1.NewSecretClientWithToken(glooOpts.Secrets, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		insertCacheIntoRcFactory(perTokenCaches, token, glooOpts.Artifacts)
		artifacts, err := v1.NewArtifactClientWithToken(glooOpts.Artifacts, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		insertCacheIntoRcFactory(perTokenCaches, token, gatewayOpts.VirtualServices)
		virtualServices, err := gatewayv1.NewVirtualServiceClientWithToken(gatewayOpts.VirtualServices, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		insertCacheIntoRcFactory(perTokenCaches, token, sqoopOpts.ResolverMaps)
		resolverMaps, err := sqoopv1.NewResolverMapClientWithToken(sqoopOpts.ResolverMaps, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		insertCacheIntoRcFactory(perTokenCaches, token, sqoopOpts.Schemas)
		schemas, err := sqoopv1.NewSchemaClientWithToken(sqoopOpts.Schemas, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if dev {
			if err := registerAll(upstreams, secrets, artifacts, virtualServices, resolverMaps, schemas); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		corsSettings.Handler(handler.GraphQL(
			graph.NewExecutableSchema(graph.Config{
				Resolvers: apiserver.NewResolvers(upstreams, schemas, artifacts, settings, secrets, virtualServices, resolverMaps),
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
