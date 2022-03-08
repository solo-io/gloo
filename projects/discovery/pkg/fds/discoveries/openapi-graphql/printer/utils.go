package printer

import (
	"fmt"
	"math"
	"reflect"

	. "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// Produces a GraphQL Value AST given a Golang value.
//
// Optionally, a GraphQL type may be provided, which will be used to
// disambiguate between value primitives.
//
// | JSON Value    | GraphQL Value        |
// | ------------- | -------------------- |
// | Object        | Input Object         |
// | Array         | List                 |
// | Boolean       | Boolean              |
// | String        | String / Enum Value  |
// | Number        | Int / Float          |

func AstFromValue(value interface{}, ttype Type) ast.Value {

	if ttype, ok := ttype.(*NonNull); ok {
		// Note: we're not checking that the result is non-null.
		// This function is not responsible for validating the input value.
		val := AstFromValue(value, ttype.OfType)
		return val
	}
	if isNullish(value) {
		return nil
	}
	valueVal := reflect.ValueOf(value)
	if !valueVal.IsValid() {
		return nil
	}
	if valueVal.Type().Kind() == reflect.Ptr {
		valueVal = valueVal.Elem()
	}
	if !valueVal.IsValid() {
		return nil
	}

	// Convert Golang slice to GraphQL list. If the Type is a list, but
	// the value is not an array, convert the value using the list's item type.
	if ttype, ok := ttype.(*List); ok {
		if valueVal.Type().Kind() == reflect.Slice {
			itemType := ttype.OfType
			values := []ast.Value{}
			for i := 0; i < valueVal.Len(); i++ {
				item := valueVal.Index(i).Interface()
				itemAST := AstFromValue(item, itemType)
				if itemAST != nil {
					values = append(values, itemAST)
				}
			}
			return ast.NewListValue(&ast.ListValue{
				Values: values,
			})
		}
		// Because GraphQL will accept single values as a "list of one" when
		// expecting a list, if there's a non-array value and an expected list type,
		// create an AST using the list's item type.
		val := AstFromValue(value, ttype.OfType)
		return val
	}

	if valueVal.Type().Kind() == reflect.Map {
		// TODO: implement AstFromValue from Map to Value
	}

	if value, ok := value.(bool); ok {
		return ast.NewBooleanValue(&ast.BooleanValue{
			Value: value,
		})
	}
	if value, ok := value.(int); ok {
		if ttype == Float {
			return ast.NewIntValue(&ast.IntValue{
				Value: fmt.Sprintf("%v.0", value),
			})
		}
		return ast.NewIntValue(&ast.IntValue{
			Value: fmt.Sprintf("%v", value),
		})
	}
	if value, ok := value.(float32); ok {
		return ast.NewFloatValue(&ast.FloatValue{
			Value: fmt.Sprintf("%v", value),
		})
	}
	if value, ok := value.(float64); ok {
		return ast.NewFloatValue(&ast.FloatValue{
			Value: fmt.Sprintf("%v", value),
		})
	}

	if value, ok := value.(string); ok {
		if _, ok := ttype.(*Enum); ok {
			return ast.NewEnumValue(&ast.EnumValue{
				Value: fmt.Sprintf("%v", value),
			})
		}
		return ast.NewStringValue(&ast.StringValue{
			Value: fmt.Sprintf("%v", value),
		})
	}

	// fallback, treat as string
	return ast.NewStringValue(&ast.StringValue{
		Value: fmt.Sprintf("%v", value),
	})
}

// Returns true if a value is null, undefined, or NaN.
func isNullish(value interface{}) bool {
	if value, ok := value.(string); ok {
		return value == ""
	}
	if value, ok := value.(int); ok {
		return math.IsNaN(float64(value))
	}
	if value, ok := value.(float32); ok {
		return math.IsNaN(float64(value))
	}
	if value, ok := value.(float64); ok {
		return math.IsNaN(value)
	}
	return value == nil
}
