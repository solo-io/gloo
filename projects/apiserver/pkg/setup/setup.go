package setup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/config"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/rs/cors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/auth"
	apiServer "github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/graph"
)

// Setup initializes the api server
func Setup(ctx context.Context, port int, debugMode bool) error {

	// Serve the query route such that it can be accessed from our UI during development
	corsSettings := cors.New(cors.Options{
		AllowedOrigins:   config.CorsAllowedOrigins,
		AllowedHeaders:   config.CorsAllowedHeaders,
		AllowCredentials: true,
		Debug:            debugMode,
	})

	// Clientset registry
	perTokenClientsets := NewPerTokenClientsets()

	http.Handle("/playground", handler.Playground("Solo-ApiServer", "/query"))
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		var resolvers graph.ResolverRoot
		token := auth.GetToken(w, r)
		if token == "" {
			resolvers = apiServer.NewUnregisteredResolver()
		} else {

			// get from cache or create anew
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
			getResolverMiddleware(debugMode),
			handler.WebsocketUpgrader(websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}),
		)).ServeHTTP(w, r)
	})

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func getResolverMiddleware(debugMode bool) handler.Option {
	return handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
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
	})
}
