package graphql

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
	Directive      *ast.Directive
	DirectiveField *ast.FieldDefinition
	Type           *ast.ObjectDefinition
}

type DirectiveVisitor func(params DirectiveVisitorParams) error

func (g *GraphqlASTVisitor) AddDirectiveVisitor(directiveName string, visitor DirectiveVisitor) {
	g.directiveVisitors[directiveName] = visitor
}

func (g *GraphqlASTVisitor) Visit(root *ast.Document) error {
	if directiveVisitors := g.directiveVisitors; len(directiveVisitors) != 0 {
		for _, def := range root.Definitions {
			if d, ok := def.(*ast.ObjectDefinition); ok {
				for _, field := range def.(*ast.ObjectDefinition).Fields {
					for _, directive := range field.Directives {
						if directive.Name == nil {
							continue
						}
						if visitFunc, ok := directiveVisitors[directive.Name.Value]; ok {
							err := visitFunc(DirectiveVisitorParams{Directive: directive, DirectiveField: field, Type: d})
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}
