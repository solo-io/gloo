package test

import (
	"fmt"

	. "github.com/onsi/gomega"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
)

// Gets the type of the argument used in the query
// and verifies that the type is used in the correct
// input type
func getInputFieldType(document *ast.Document) string {
	var inputFieldType string
	for _, node := range document.Definitions {
		if t, ok := node.(*ast.ObjectDefinition); ok {
			if t.GetName().Value == "Query" {
				// There should only be one query we are testing at a time
				ExpectWithOffset(1, t.Fields).To(HaveLen(1))
				queryField := t.Fields[0]
				// there should only be one argument in this test that we are checking the type of
				ExpectWithOffset(1, queryField.Arguments).To(HaveLen(1))
				arg := queryField.Arguments[0]
				inputFieldType = printer.Print(arg.Type).(string)
			}
		}
	}
	for _, node := range document.Definitions {
		if inputType, ok := node.(*ast.InputObjectDefinition); ok {
			if inputType.GetName().Value == inputFieldType {
				for _, inputValueDef := range inputType.Fields {
					return printer.Print(inputValueDef.Type).(string)
				}
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
