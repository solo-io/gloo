package directives

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/rotisserie/eris"
)

const (
	RESOLVER_DIRECTIVE     = "resolve"
	RESOLVER_NAME_ARGUMENT = "name"
)

type resolveDirective struct {
	ResolverNameAstValue ast.Value
	ResolverName         string
}

func NewResolveDirective() *resolveDirective {
	return &resolveDirective{}
}

// Validate validates that the resolve directive usage is syntactically correct.
// At the end of successful validation, `ResolverNameAstValue` and `ResolverName` will be
// populated.
func (d *resolveDirective) Validate(directiveVisitorParams DirectiveVisitorParams) error {
	if directiveVisitorParams.DirectiveField == nil {
		return eris.Errorf(`"%s" directive must only be used on fields`, RESOLVER_DIRECTIVE)
	}
	arguments := map[string]ast.Value{}
	directive := directiveVisitorParams.Directive
	for _, argument := range directive.Arguments {
		arguments[argument.Name.Value] = argument.Value
	}
	resolverName, ok := arguments[RESOLVER_NAME_ARGUMENT]
	if !ok {
		return NewGraphqlSchemaError(directive, `the "%s" directive must have a "%s" argument to reference a resolver`,
			RESOLVER_DIRECTIVE, RESOLVER_NAME_ARGUMENT)
	}
	if resolverName.GetKind() != kinds.StringValue {
		return NewGraphqlSchemaError(resolverName, `"%s" argument must be a string value`, RESOLVER_NAME_ARGUMENT)
	}
	name := resolverName.GetValue().(string)

	d.ResolverNameAstValue = resolverName
	d.ResolverName = name
	return nil
}
