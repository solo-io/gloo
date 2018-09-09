package engine

import (
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"github.com/vektah/gqlgen/neelance/schema"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/engine/resolvers"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/engine/exec"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/engine/router"
)

type Engine struct {
	sidecarAddr string
}

func NewEngine(sidecarAddr string) *Engine {
	return &Engine{sidecarAddr: sidecarAddr}
}

func (en *Engine) CreateGraphqlEndpoint(schema *v1.Schema, resolverMap *v1.ResolverMap) (*router.Endpoint, error, error) {
	resolverFactory := resolvers.NewResolverFactory(en.sidecarAddr, resolverMap)
	parsedSchema, err := parseSchemaString(schema)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse schema"), nil
	}
	executableResolvers, err := exec.NewExecutableResolvers(parsedSchema, resolverFactory.CreateResolver)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate executable resolvers from map")
	}
	executableSchema := exec.NewExecutableSchema(parsedSchema, executableResolvers)
	return &router.Endpoint{
		SchemaName: schema.Metadata.Name,
		RootPath:   "/" + schema.Metadata.Name,
		QueryPath:  "/" + schema.Metadata.Name + "/query",
		ExecSchema: executableSchema,
	}, nil, nil
}

func parseSchemaString(sch *v1.Schema) (*schema.Schema, error) {
	parsedSchema := schema.New()
	return parsedSchema, parsedSchema.Parse(sch.InlineSchema)
}
