package setup

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/rs/cors"
	gatewayV1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/auth"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/config"
	apiServer "github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/graph"
	sqoopV1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
)

// Setup initializes the api server
func Setup(ctx context.Context, port int, dev bool, debugMode bool, opts ApiServerOpts) error {

	// fail fast if the environment is not correctly configured
	config.ValidateEnvVars()

	// initial resource registration
	if err := initialResourceRegistration(opts); err != nil {
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

	perTokenClientsets := NewPerTokenClientsets(opts)

	http.Handle("/playground", handler.Playground("Solo-ApiServer", "/query"))
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		var resolvers graph.ResolverRoot
		token := auth.GetToken(w, r)
		if token == "" {
			resolvers = apiServer.NewUnregisteredResolver()
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

// PerTokenClientsets contains global settings and user-specific resource clients
// clients is a map from user token to resource clients
// the token is also used for authorizing actions on the resource clients
type PerTokenClientsets struct {
	lock    sync.RWMutex
	clients map[string]*ClientSet
	opts    ApiServerOpts
}

// NewPerTokenClientsets - helper function
func NewPerTokenClientsets(opts ApiServerOpts) PerTokenClientsets {
	return PerTokenClientsets{
		clients: make(map[string]*ClientSet),
		opts:    opts,
	}
}

func (ptc PerTokenClientsets) ClientsetForToken(token string) (*ClientSet, error) {
	ptc.lock.Lock()
	defer ptc.lock.Unlock()
	clientsetForToken, ok := ptc.clients[token]
	if ok {
		return clientsetForToken, nil
	}
	// the new clientset has a new cache
	clientset, err := NewClientSet(token, ptc.opts)
	if err != nil {
		return nil, err
	}
	ptc.clients[token] = clientset
	return clientset, nil
}

// ClientSet is a collection of all the exposed resource clients
type ClientSet struct {
	v1.UpstreamClient
	gatewayV1.VirtualServiceClient
	v1.SettingsClient
	v1.SecretClient
	v1.ArtifactClient
	sqoopV1.ResolverMapClient
	sqoopV1.SchemaClient
}

// NewClientSet - helper function
// Warning! this will write to opts
// Todo: ilackarms: refactor this so opts is copied
func NewClientSet(token string, opts ApiServerOpts) (*ClientSet, error) {

	// Create a new cache for the token
	cache := kube.NewKubeCache()

	// todo: be sure to add new resource clients here
	for _, rcFactory := range []factory.ResourceClientFactory{
		opts.UpstreamsRCF,
		opts.VirtualServicesRCF,
		opts.SchemasRCF,
		opts.ResolverMapsRCF,
		opts.SecretsRCF,
		opts.ArtifactsRCF,
	} {
		setKubeFactoryCache(rcFactory, cache)
	}

	upstreams, err := v1.NewUpstreamClientWithToken(opts.UpstreamsRCF, token)
	if err != nil {
		return nil, err
	}
	secrets, err := v1.NewSecretClientWithToken(opts.SecretsRCF, token)
	if err != nil {
		return nil, err
	}
	artifacts, err := v1.NewArtifactClientWithToken(opts.ArtifactsRCF, token)
	if err != nil {
		return nil, err
	}
	virtualServices, err := gatewayV1.NewVirtualServiceClientWithToken(opts.VirtualServicesRCF, token)
	if err != nil {
		return nil, err
	}
	resolverMaps, err := sqoopV1.NewResolverMapClientWithToken(opts.ResolverMapsRCF, token)
	if err != nil {
		return nil, err
	}
	schemas, err := sqoopV1.NewSchemaClientWithToken(opts.SchemasRCF, token)
	if err != nil {
		return nil, err
	}

	if err := registerAll(upstreams, secrets, artifacts, virtualServices, resolverMaps, schemas); err != nil {
		return nil, err
	}
	return &ClientSet{
		UpstreamClient:       upstreams,
		ArtifactClient:       artifacts,
		SecretClient:         secrets,
		ResolverMapClient:    resolverMaps,
		SchemaClient:         schemas,
		VirtualServiceClient: virtualServices,
		SettingsClient:       opts.SettingsClient,
	}, nil
}

// NewResolvers - helper function
func (c ClientSet) NewResolvers() *apiServer.ApiResolver {
	return apiServer.NewResolvers(c.UpstreamClient,
		c.SchemaClient,
		c.ArtifactClient,
		c.SettingsClient,
		c.SecretClient,
		c.VirtualServiceClient,
		c.ResolverMapClient)
}

func setKubeFactoryCache(fact factory.ResourceClientFactory, cache *kube.KubeCache) {
	if kubeFactory, ok := fact.(*factory.KubeResourceClientFactory); ok {
		kubeFactory.SharedCache = cache
	}
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

func initialResourceRegistration(opts ApiServerOpts) error {
	upstreams, err := v1.NewUpstreamClient(opts.UpstreamsRCF)
	if err != nil {
		return err
	}
	if err := upstreams.Register(); err != nil {
		return err
	}
	secrets, err := v1.NewSecretClient(opts.SecretsRCF)
	if err != nil {
		return err
	}
	if err := secrets.Register(); err != nil {
		return err
	}
	artifacts, err := v1.NewArtifactClient(opts.ArtifactsRCF)
	if err != nil {
		return err
	}
	if err := artifacts.Register(); err != nil {
		return err
	}
	virtualServices, err := gatewayV1.NewVirtualServiceClient(opts.VirtualServicesRCF)
	if err != nil {
		return err
	}
	if err := virtualServices.Register(); err != nil {
		return err
	}
	resolverMaps, err := sqoopV1.NewResolverMapClient(opts.ResolverMapsRCF)
	if err != nil {
		return err
	}
	if err := resolverMaps.Register(); err != nil {
		return err
	}
	schemas, err := sqoopV1.NewSchemaClient(opts.SchemasRCF)
	if err != nil {
		return err
	}
	if err := schemas.Register(); err != nil {
		return err
	}
	return nil
}
