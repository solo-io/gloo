package directives

import (
	"github.com/graphql-go/graphql/language/ast"
)

func NewGraphqlASTVisitor() *GraphqlASTVisitor {
	return &GraphqlASTVisitor{
		directiveVisitors: map[string]DirectiveVisitor{},
	}
}

type GraphqlASTVisitor struct {
	directiveVisitors map[string]DirectiveVisitor
}

type DirectiveVisitorParams struct {
	Directive       *ast.Directive
	DirectiveFields []*ast.FieldDefinition // non-nil on type directives (cannot combine into one type since we could have list of length one)
	DirectiveField  *ast.FieldDefinition   // non-nil on field directives
	Type            *ast.ObjectDefinition
}

// DirectiveVisitor returns an error, and a bool
// If bool is true, deletes the directive. If it is false, it leaves the directive.
// This is useful for directives such as @resolve or @cacheControl that we use only for control plane
// but we don't want to propogate to data plane.
type DirectiveVisitor func(params DirectiveVisitorParams) (bool, error)

func (g *GraphqlASTVisitor) AddDirectiveVisitor(directiveName string, visitor DirectiveVisitor) {
	g.directiveVisitors[directiveName] = visitor
}

func (g *GraphqlASTVisitor) Visit(root *ast.Document) error {
	if directiveVisitors := g.directiveVisitors; len(directiveVisitors) != 0 {
		for _, def := range root.Definitions {
			if objDef, ok := def.(*ast.ObjectDefinition); ok {
				var typeDirectiveIndicesToKeep []int
				// check type directives
				for i, directive := range objDef.Directives {
					if directive.Name == nil {
						continue
					}
					if visitFunc, ok := directiveVisitors[directive.Name.Value]; ok {
						deleteDirective, err := visitFunc(DirectiveVisitorParams{Directive: directive, DirectiveFields: objDef.Fields, Type: objDef})
						if err != nil {
							return err
						}
						if !deleteDirective {
							typeDirectiveIndicesToKeep = append(typeDirectiveIndicesToKeep, i)
						}
					}
				}
				var newTypeDirectives []*ast.Directive
				for _, idx := range typeDirectiveIndicesToKeep {
					newTypeDirectives = append(newTypeDirectives, objDef.Directives[idx])
				}
				objDef.Directives = newTypeDirectives
				// check field directives
				var fieldDirectiveIndicesToKeep []int
				for _, field := range def.(*ast.ObjectDefinition).Fields {
					for idx, directive := range field.Directives {
						if directive.Name == nil {
							continue
						}
						if visitFunc, ok := directiveVisitors[directive.Name.Value]; ok {
							deleteDirective, err := visitFunc(DirectiveVisitorParams{Directive: directive, DirectiveField: field, Type: objDef})
							if err != nil {
								return err
							}
							if !deleteDirective {
								fieldDirectiveIndicesToKeep = append(fieldDirectiveIndicesToKeep, idx)
							}
						}
					}
					var newFieldDirectives []*ast.Directive
					for _, idx := range typeDirectiveIndicesToKeep {
						newFieldDirectives = append(newFieldDirectives, objDef.Directives[idx])
					}
					field.Directives = newFieldDirectives
				}
			}
		}
	}
	return nil
}
