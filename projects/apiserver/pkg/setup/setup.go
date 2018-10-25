package setup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"

	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/rs/cors"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/auth"
	apiserver "github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/graph"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/syncer"
)

// Setup initializes the apiserver
func Setup(ctx context.Context, port int, dev bool, debugMode bool, settings v1.SettingsClient, glooOpts bootstrap.Opts, gatewayOpts gatewaysetup.Opts, sqoopOpts sqoopsetup.Opts) error {
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
		Debug:            debugMode,
	})

	perTokenClientsets := NewPerTokenClientsets(settings, glooOpts, gatewayOpts, sqoopOpts)

	http.Handle("/playground", handler.Playground("Solo-ApiServer", "/query"))
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		var resolvers graph.ResolverRoot
		token := auth.GetToken(w, r)
		if token == "" {
			resolvers = apiserver.NewUnregisteredResolver()
		} else {
			clientset, err := perTokenClientsets.ClientsetForToken(token)
			if err != nil {
				contextutils.LoggerFrom(ctx).Errorf("failed to create clientset: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			resolvers = clientset.NewResolvers()
		}
		corsSettings.Handler(handler.GraphQL(
			graph.NewExecutableSchema(graph.Config{
				Resolvers: resolvers,
			}),
			handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
				rc := graphql.GetResolverContext(ctx)
				if debugMode {
					// note: many of our queries will produce cyclic objects
					// printing truncated string representations is an easy alternative to cycle detection
					fmt.Println("Entered", rc.Object, rc.Field.Name)
				}
				res, err = next(ctx)
				if debugMode {
					fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
				}
				return res, err
			}),
				handler.WebsocketUpgrader(websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
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

// PerTokenClientsets contains global settings and user-specific resource clients
// clients is a map from user token to resource clients
// the token is also used for authorizing actions on the resource clients
type PerTokenClientsets struct {
	lock        sync.RWMutex
	clients     map[string]*Clientset
	settings    v1.SettingsClient
	glooOpts    bootstrap.Opts
	gatewayOpts gatewaysetup.Opts
	sqoopOpts   sqoopsetup.Opts
}

// NewPerTokenClientsets - helper function
func NewPerTokenClientsets(settings v1.SettingsClient, glooOpts bootstrap.Opts, gatewayOpts gatewaysetup.Opts, sqoopOpts sqoopsetup.Opts) PerTokenClientsets {
	return PerTokenClientsets{
		clients:     make(map[string]*Clientset),
		settings:    settings,
		glooOpts:    glooOpts,
		gatewayOpts: gatewayOpts,
		sqoopOpts:   sqoopOpts,
	}
}

func (ptc PerTokenClientsets) ClientsetForToken(token string) (*Clientset, error) {
	ptc.lock.Lock()
	defer ptc.lock.Unlock()
	clientsetForToken, ok := ptc.clients[token]
	if ok {
		return clientsetForToken, nil
	}
	// the new clientset has a new cache
	clientset, err := NewClientSet(token, ptc.settings, ptc.glooOpts, ptc.gatewayOpts, ptc.sqoopOpts)
	if err != nil {
		return nil, err
	}
	ptc.clients[token] = clientset
	return clientset, nil
}

// Clientset is a collection of all the exposed resource clients
type Clientset struct {
	v1.UpstreamClient
	gatewayv1.VirtualServiceClient
	v1.SettingsClient
	v1.SecretClient
	v1.ArtifactClient
	sqoopv1.ResolverMapClient
	sqoopv1.SchemaClient
}

func setKubeFactoryCache(fact factory.ResourceClientFactory, cache *kube.KubeCache) {
	if kubeFactory, ok := fact.(*factory.KubeResourceClientFactory); ok {
		kubeFactory.SharedCache = cache
	}
}

// NewClientSet - helper function
// Warning! this will write to opts
// Todo: ilackarms: refactor this so opts is copied
func NewClientSet(token string, settings v1.SettingsClient, glooOpts bootstrap.Opts, gatewayOpts gatewaysetup.Opts, sqoopOpts sqoopsetup.Opts) (*Clientset, error) {
	cache := kube.NewKubeCache()
	// todo: be sure to add new resource clients here
	setKubeFactoryCache(glooOpts.Proxies, cache)
	setKubeFactoryCache(glooOpts.Upstreams, cache)
	setKubeFactoryCache(glooOpts.Artifacts, cache)
	setKubeFactoryCache(glooOpts.Secrets, cache)
	setKubeFactoryCache(gatewayOpts.VirtualServices, cache)
	setKubeFactoryCache(gatewayOpts.Gateways, cache)
	setKubeFactoryCache(sqoopOpts.Schemas, cache)
	setKubeFactoryCache(sqoopOpts.ResolverMaps, cache)
	upstreams, err := v1.NewUpstreamClientWithToken(glooOpts.Upstreams, token)
	if err != nil {
		return nil, err
	}
	secrets, err := v1.NewSecretClientWithToken(glooOpts.Secrets, token)
	if err != nil {
		return nil, err
	}
	artifacts, err := v1.NewArtifactClientWithToken(glooOpts.Artifacts, token)
	if err != nil {
		return nil, err
	}
	virtualServices, err := gatewayv1.NewVirtualServiceClientWithToken(gatewayOpts.VirtualServices, token)
	if err != nil {
		return nil, err
	}
	resolverMaps, err := sqoopv1.NewResolverMapClientWithToken(sqoopOpts.ResolverMaps, token)
	if err != nil {
		return nil, err
	}
	schemas, err := sqoopv1.NewSchemaClientWithToken(sqoopOpts.Schemas, token)
	if err != nil {
		return nil, err
	}
	if err := registerAll(upstreams, secrets, artifacts, virtualServices, resolverMaps, schemas); err != nil {
		return nil, err
	}
	return &Clientset{
		UpstreamClient:       upstreams,
		ArtifactClient:       artifacts,
		SecretClient:         secrets,
		ResolverMapClient:    resolverMaps,
		SchemaClient:         schemas,
		VirtualServiceClient: virtualServices,
		SettingsClient:       settings,
	}, nil
}

// NewResolvers - helper function
func (c Clientset) NewResolvers() *apiserver.ApiResolver {
	return apiserver.NewResolvers(c.UpstreamClient,
		c.SchemaClient,
		c.ArtifactClient,
		c.SettingsClient,
		c.SecretClient,
		c.VirtualServiceClient,
		c.ResolverMapClient)
}
