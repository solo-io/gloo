package test

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
)

func getInputFieldType(document *ast.Document) string {
	for _, node := range document.Definitions {
		if inputType, ok := node.(*ast.InputObjectDefinition); ok {
			for _, inputValueDef := range inputType.Fields {
				return fmt.Sprintf("%s", printer.Print(inputValueDef.Type))
			}
		}
	}
	return ""
}

func getFieldType(document *ast.Document) string {
	for _, node := range document.Definitions {
		if inputType, ok := node.(*ast.ObjectDefinition); ok {
			for _, valueDef := range inputType.Fields {
				return fmt.Sprintf("%s", printer.Print(valueDef.Type))
			}
		}
	}
	return ""
}

func getTypeDefinition(document *ast.Document, inputType bool, name string) ast.Definition {
	for _, def := range document.Definitions {
		if inputType {
			if inputDef, ok := def.(*ast.InputObjectDefinition); ok {
				if inputDef.GetName().Value == name {
					return inputDef
				}
			}
		} else {
			if regularDef, ok := def.(*ast.ObjectDefinition); ok {
				if regularDef.GetName().Value == name {
					return regularDef
				}
			}
		}
	}
	return nil
}

func getEnumDefinition(doc *ast.Document, enumTypeName string) *ast.EnumDefinition {
	for _, def := range doc.Definitions {
		if enumdef, ok := def.(*ast.EnumDefinition); ok {
			if enumdef.Name.Value == enumTypeName {
				return enumdef
			}
		}
	}
	return nil
}

func getFieldsWithType(definition ast.Definition, fieldName string, typeName string) bool {
	if inputObjectDef, ok := definition.(*ast.InputObjectDefinition); ok {
		for _, inputField := range inputObjectDef.Fields {
			if inputField.Name.Value == fieldName {
				if fmt.Sprintf("%s", printer.Print(inputField.Type)) == typeName {
					return true
				}
			}
		}
	} else if def, ok := definition.(*ast.ObjectDefinition); ok {
		for _, inputField := range def.Fields {
			if inputField.Name.Value == fieldName {
				if fmt.Sprintf("%s", printer.Print(inputField.Type)) == typeName {
					return true
				}
			}
		}
	}
	return false

}
