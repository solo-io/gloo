package printer

import (
	"fmt"
	"strings"

	. "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
	. "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi/printer/types"
)

func PrintFilteredSchema(schema *Schema) string {
	var toPrint []string
	// Print Schema Definition
	toPrint = append(toPrint, PrintSchemaDefinition(schema))
	// Print Non built in Directives
	filteredDirectives := FilterDirectives(schema.Directives())
	toPrint = append(toPrint, PrintDirectives(filteredDirectives)...)
	// Print Types
	filteredTypes := FilterTypes(schema.TypeMap())
	toPrint = append(toPrint, PrintTypes(filteredTypes)...)

	var printingArray []string
	for _, out := range toPrint {
		if out != "" {
			printingArray = append(printingArray, out)
		}
	}
	return strings.Join(printingArray, "\n\n")
}

func FilterTypes(types TypeMap) TypeMap {
	builtInTypes := map[string]bool{
		"String":              true,
		"Int":                 true,
		"Float":               true,
		"Boolean":             true,
		"ID":                  true,
		"__Schema":            true,
		"__Directive":         true,
		"__DirectiveLocation": true,
		"__Type":              true,
		"__Field":             true,
		"__InputValue":        true,
		"__EnumValue":         true,
		"__TypeKind":          true,
	}
	var ret = make(TypeMap)
	for typeName, t := range types {
		if _, ok := builtInTypes[typeName]; !ok {
			ret[typeName] = t
		}
	}
	return ret
}

func FilterDirectives(directives []*Directive) []*Directive {
	specifiedDirectives := map[*Directive]bool{
		IncludeDirective: true,
		SkipDirective:    true,
	}

	var ret []*Directive
	for _, d := range directives {
		if _, ok := specifiedDirectives[d]; !ok {
			ret = append(ret, d)
		}
	}
	return ret
}

func PrintTypes(types TypeMap) []string {
	var ret []string
	for typeName, t := range types {
		ret = append(ret, PrintType(typeName, t))
	}
	return ret
}

func PrintType(typeName string, t Type) string {
	if scalar := IsScalarType(t); scalar != nil {
		return PrintScalar(typeName, scalar)
	}
	if obj := IsObjectType(t); obj != nil {
		return PrintObject(typeName, obj)
	}
	if i := IsInterfaceType(t); i != nil {
		return PrintInterface(typeName, i)
	}
	if u := IsUnionType(t); u != nil {
		return PrintUnion(typeName, u)
	}
	if e := IsEnumType(t); e != nil {
		return PrintEnum(typeName, e)
	}
	// Istanbul ignore else (See: 'https://github.com/graphql/graphql-js/Issues/2618')
	if io := IsInputObjectType(t); io != nil {
		return PrintInputObject(typeName, io)
	}

	return ""
}

func PrintScalar(name string, t *Scalar) string {
	return PrintDescription(&PrintDescriptionParams{Description: t.Description()}) +
		"scalar " + name /* PrintSpecifiedByUrl -- not supported in graphql-go*/
}

func PrintObject(name string, t *Object) string {
	return PrintDescription(&PrintDescriptionParams{
		Description: t.Description(),
	}) + "type " + name + PrintImplementedInterfaces(t) + PrintFields(t)
}

type Fieldable interface {
	Fields() FieldDefinitionMap
}

func PrintFields(t Fieldable) string {
	var ret []string
	first := true
	for name, field := range t.Fields() {
		ret = append(ret, PrintDescription(&PrintDescriptionParams{
			Description:  field.Description,
			Indentation:  "  ",
			FirstInBlock: first,
		})+"  "+name+
			PrintArgs(&PrintArgsParams{
				Args:        field.Args,
				Indentation: "  ",
			})+": "+
			field.Type.String()+
			PrintDeprecated(field.DeprecationReason))
		first = false
	}
	return PrintBlock(ret)
}

func PrintBlock(block []string) string {
	if len(block) == 0 {
		return ""
	}
	return " {\n" + strings.Join(block, "\n") + "\n}"
}

func PrintImplementedInterfaces(t *Object) string {
	interfaces := t.Interfaces()
	if len(interfaces) == 0 {
		return ""
	}
	ret := " implements "
	var iFaces []string
	for _, i := range interfaces {
		iFaces = append(iFaces, i.Name())
	}
	return ret + strings.Join(iFaces, " & ")
}

func PrintInterface(name string, t *Interface) string {
	return PrintDescription(&PrintDescriptionParams{
		Description: t.Description(),
	}) + "interface " + name +
		" implements " + t.Name() +
		PrintFields(t)
}

func PrintUnion(name string, t *Union) string {
	types := t.Types()
	possibleTypes := ""
	var typeNames []string
	for _, t := range t.Types() {
		typeNames = append(typeNames, t.Name())
	}
	if len(types) > 0 {
		possibleTypes = " = " + strings.Join(typeNames, " | ")
	}
	return PrintDescription(&PrintDescriptionParams{Description: t.Description()}) +
		"union " + name + possibleTypes
}

func PrintEnum(name string, t *Enum) string {
	var values []string
	for i, val := range t.Values() {
		values = append(values, PrintDescription(&PrintDescriptionParams{
			Description:  val.Description,
			Indentation:  "  ",
			FirstInBlock: i == 0,
		})+"  "+val.Name) /* + PrintDeprecated(val)*/
	}

	return PrintDescription(&PrintDescriptionParams{Description: t.Description()}) +
		"enum " + name + PrintBlock(values)
}

func PrintInputObject(typeName string, t *InputObject) string {
	var fields []string
	first := true
	for _, field := range t.Fields() {
		fields = append(fields, PrintDescription(&PrintDescriptionParams{
			Description:  field.Description(),
			Indentation:  "  ",
			FirstInBlock: first,
		})+"  "+PrintInputObjectValue(field))
		first = false
	}

	return PrintDescription(&PrintDescriptionParams{
		Description: t.Description(),
	}) + "input " + typeName + PrintBlock(fields)
}

func PrintInputObjectValue(field *InputObjectField) string {
	if field == nil {
		return ""
	}
	defaultAst := AstFromValue(field.DefaultValue, field.Type)
	fieldDecl := field.Name() + ": " + field.Type.String()
	if defaultAst != nil && defaultAst.GetValue() != nil {
		fieldDecl += " = " + fmt.Sprintf("%s", printer.Print(defaultAst))
	}
	return fieldDecl /* + PrintDeprecated(arg.Reason)*/
}

func PrintDeprecated(reason string) string {
	if reason == "" {
		return reason
	}
	const DEFAULT_DEPRECATION_REASON = "No longer supported"
	if reason != DEFAULT_DEPRECATION_REASON {
		//const astValue = printer.Print({ kind: Kind.STRING, value: reason });
		astValue := printer.Print(ast.NewStringValue(&ast.StringValue{
			Value: reason,
		}))

		return ` @deprecated(reason: ` + fmt.Sprintf("%s", astValue) + ")"
	}
	return " @deprecated"
}

func IsScalarType(t Type) *Scalar {
	scalar, _ := t.(*Scalar)
	return scalar
}

func IsObjectType(t Type) *Object {
	obj, _ := t.(*Object)
	return obj
}

func IsInterfaceType(t Type) *Interface {
	i, _ := t.(*Interface)
	return i
}

func IsUnionType(t Type) *Union {
	u, _ := t.(*Union)
	return u
}

func IsEnumType(t Type) *Enum {
	e, _ := t.(*Enum)
	return e
}

func IsInputObjectType(t Type) *InputObject {
	i, _ := t.(*InputObject)
	return i
}

func PrintDirectives(directives []*Directive) []string {
	var ret []string
	for _, directive := range directives {

		directiveDesc := PrintDescription(&PrintDescriptionParams{
			Description: directive.Description,
		})

		repeatable := ""
		//todo - support printing repeatable directives when graphql-go supports rpt directives
		/*if directive.IsRepeatable {
		      repeatable := " repeatable"
		  }
		*/

		result := fmt.Sprintf("%sdirective @%s%s%s on %s",
			directiveDesc,
			directive.Name,
			PrintArgs(&PrintArgsParams{Args: directive.Args}),
			repeatable,
			strings.Join(directive.Locations, " | "))
		ret = append(ret, result)
	}
	return ret
}

func PrintArgs(params *PrintArgsParams) string {
	if params == nil || len(params.Args) == 0 {
		return ""
	}

	args := params.Args
	var someArgHasDescription bool
	for _, arg := range args {
		if arg.Description() != "" {
			someArgHasDescription = true
			break
		}
	}
	if !someArgHasDescription {
		var argsInputValue []string
		for _, arg := range args {
			argsInputValue = append(argsInputValue, PrintInputValue(arg))
		}
		return "(" + strings.Join(argsInputValue, ", ") + ")"
	}

	var toRet []string
	for i, arg := range args {
		toRet = append(toRet, PrintDescription(&PrintDescriptionParams{
			Description:  arg.Description(),
			Indentation:  "  " + params.Indentation,
			FirstInBlock: i == 0,
		})+"  "+params.Indentation+PrintInputValue(arg))
	}

	return "(\n" + strings.Join(toRet, "\n") + "\n" + params.Indentation + ")"
}

func PrintInputValue(arg *Argument) string {
	if arg == nil {
		return ""
	}
	defaultAst := AstFromValue(arg.DefaultValue, arg.Type)
	argDecl := arg.Name() + ": " + arg.Type.String()
	if defaultAst != nil && defaultAst.GetValue() != nil {
		argDecl += " = " + fmt.Sprintf("%s", printer.Print(defaultAst))
	}
	return argDecl /* + PrintDeprecated(arg.Reason)*/
}

func PrintSchemaDefinition(schema *Schema) string {
	if IsSchemaOfCommonNames(schema) {
		return ""
	}

	var operationTypes []string
	if q := schema.QueryType(); q != nil {
		str := fmt.Sprintf("  query: %s", q.Name())
		operationTypes = append(operationTypes, str)
	}

	if q := schema.MutationType(); q != nil {
		str := fmt.Sprintf("  mutation: %s", q.Name())
		operationTypes = append(operationTypes, str)
	}

	if q := schema.SubscriptionType(); q != nil {
		str := fmt.Sprintf("  subscription: %s", q.Name())
		operationTypes = append(operationTypes, str)
	}

	schemaStr := fmt.Sprintf("schema {\n%s\n}", strings.Join(operationTypes, "\n"))
	return PrintDescription(NewPrintDescriptionParams()) + schemaStr
}

func PrintDescription(args *PrintDescriptionParams) string {
	if args == nil || args.Description == "" {
		return ""
	}

	preferMultiLines := len(args.Description) > 70
	blockString := PrintBlockString(args.Description, preferMultiLines)
	prefix := args.Indentation
	if args.Indentation != "" && !args.FirstInBlock {
		prefix = "\n" + args.Indentation
	}

	return prefix + strings.ReplaceAll(blockString, "\n", "\n"+args.Indentation) + "\n"
}

func PrintBlockString(value string, wrapLines bool) string {
	isSingleLine := !strings.Contains(value, "\n")
	hasLeadingSpace := strings.HasPrefix(value, " ") || strings.HasPrefix(value, "\t")
	hasTrailingQuote := strings.HasSuffix(value, "\"")
	hasTrailingSlash := strings.HasSuffix(value, "\\")
	printAsMultipleLines := !isSingleLine || hasTrailingQuote || hasTrailingSlash || wrapLines

	result := ""
	if printAsMultipleLines && !(isSingleLine && hasLeadingSpace) {
		result += "\n"
	}
	result += value
	if printAsMultipleLines {
		result += "\n"
	}
	if result == "" {
		return ""
	}
	return `"""` + strings.ReplaceAll(result, `"""`, `\\"""`) + `"""`
}

func IsSchemaOfCommonNames(schema *Schema) bool {
	queryType := schema.QueryType()
	if queryType != nil && queryType.Name() != "Query" {
		return false
	}
	if mutationType := schema.MutationType(); mutationType != nil && mutationType.Name() != "Mutation" {
		return false
	}
	if subscriptionType := schema.MutationType(); subscriptionType != nil && subscriptionType.Name() != "Mutation" {
		return false
	}

	return true
}
