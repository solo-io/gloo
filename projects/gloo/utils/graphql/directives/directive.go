package directives

import (
	"fmt"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
)

// Directive is an interface for defining a type of directive that we support in graphql schemas.
type Directive interface {
	// Validate validates that the directive usage in the schema is syntactically correct.
	Validate(directiveVisitorParams DirectiveVisitorParams) error
}

// Allows us to print the specific character location of an error in a graphql schema
type locatable interface {
	GetLoc() *ast.Location
}

func NewGraphqlSchemaError(l locatable, description string, args ...interface{}) error {
	desc := fmt.Sprintf(description, args...)
	return gqlerrors.NewSyntaxError(l.GetLoc().Source, l.GetLoc().Start, desc)
}
