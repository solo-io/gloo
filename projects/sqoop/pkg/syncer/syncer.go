package syncer

import (
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"context"
	"github.com/solo-io/qloo/pkg/graphql"
	"github.com/solo-io/qloo/pkg/resolvers"
	"github.com/solo-io/qloo/pkg/exec"
)

type Syncer struct {}

func (s *Syncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {
	/*
	1. create new graphql endpoints
	2. update router with graphql endpoints
	3. resourceErrs := configErrs / writeReports(resourceErrs)
	4. configure gloo
	# tip: use multierror instead of early returns (when possible)
	 */
	 for _, schema := range snap.Schemas.List() {

	 }
}

// for supporting multiple schemas on the same port, essentially
func createRouterEndpoints() {}
func createEndpointForSchema() {}
func generateSkeletonResolvermap() {}
func generateEmptyResolvermap() {}
func configureGloo() {}

func (el *EventLoop) createGraphqlEndpoint(schema *v1.Schema, resolverMap *v1.ResolverMap) (*graphql.Endpoint, error, error) {
	resolverFactory := resolvers.NewResolverFactory(el.proxyAddr, resolverMap)
	parsedSchema, err := parseSchemaString(schema)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse schema"), nil
	}
	executableResolvers, err := exec.NewExecutableResolvers(parsedSchema, resolverFactory.CreateResolver)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate resolvers from map")
	}
	el.operator.ApplyResolvers(resolverMap)
	executableSchema := exec.NewExecutableSchema(parsedSchema, executableResolvers)
	return &graphql.Endpoint{
		SchemaName: schema.Name,
		RootPath:   "/" + schema.Name,
		QueryPath:  "/" + schema.Name + "/query",
		ExecSchema: executableSchema,
	}, nil, nil
}

