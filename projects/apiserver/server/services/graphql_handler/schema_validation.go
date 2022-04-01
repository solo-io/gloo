package graphql_handler

import (
	"github.com/graphql-go/graphql/language/parser"
	"github.com/rotisserie/eris"
	graphql_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1alpha1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/directives"
)

func ValidateSchemaDefinition(req *rpc_edge_v1.ValidateSchemaDefinitionRequest) error {
	switch req.GetInput().(type) {
	case *rpc_edge_v1.ValidateSchemaDefinitionRequest_SchemaDefinition:
		return validateInternal(req.GetSchemaDefinition(), map[string]*graphql_v1alpha1.Resolution{})
	case *rpc_edge_v1.ValidateSchemaDefinitionRequest_Spec:
		return validateInternal(req.GetSpec().GetExecutableSchema().GetSchemaDefinition(),
			req.GetSpec().GetExecutableSchema().GetExecutor().GetLocal().GetResolutions())
	default:
		return eris.Errorf("request must specify either schema definition or spec: %v", req)
	}
}

// Validates the following:
// 1. the schema definition string can be parsed
// 2. all usages of supported directives are syntactically correct
// 3. resolver names referenced via `@resolve` directives in the schema definition have a corresponding
//    entry in the resolutions map.
func validateInternal(schema string, resolutions map[string]*graphql_v1alpha1.Resolution) error {
	doc, err := parser.Parse(parser.ParseParams{Source: schema})
	if err != nil {
		return eris.Wrap(err, "unable to parse graphql schema")
	}

	visitor := directives.NewGraphqlASTVisitor()
	// resolve directive
	visitor.AddDirectiveVisitor(directives.RESOLVER_DIRECTIVE, func(directiveVisitorParams directives.DirectiveVisitorParams) (bool, error) {
		// validate correct usage of the resolve directive
		resolveDirective := directives.NewResolveDirective()
		err := resolveDirective.Validate(directiveVisitorParams)
		if err != nil {
			return false, err
		}

		// check if referenced resolver exists in the resolutions map
		resolution := resolutions[resolveDirective.ResolverName]
		if resolution == nil {
			return false, directives.NewGraphqlSchemaError(resolveDirective.ResolverNameAstValue,
				"resolver %s is not defined", resolveDirective.ResolverName)
		}

		return false, nil
	})

	// cacheControl directive
	visitor.AddDirectiveVisitor(directives.CACHE_CONTROL_DIRECTIVE, func(directiveVisitorParams directives.DirectiveVisitorParams) (bool, error) {
		// validate correct usage of the cacheControl directive
		cacheControlDirective := directives.NewCacheControlDirective()
		return cacheControlDirective.Validate(directiveVisitorParams)
	})

	return visitor.Visit(doc)
}
