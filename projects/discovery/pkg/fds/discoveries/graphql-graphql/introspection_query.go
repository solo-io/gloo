package graphql

import (
	"bytes"
	"fmt"
	"strings"
)

// Introspection allows us to retrieve the definitions and schemas on a graphql server.
// for reference please review the latest version of [the graphql spec](https://spec.graphql.org/October2021/#sec-Schema-Introspection).
// This query was retrieved from [the GraphQL Go Repo](https://github.com/graphql-go/graphql/blob/master/testutil/introspection_query.go#L3).

func generateDepthOfType(depth int, b *bytes.Buffer) error {
	if b == nil {
		b = &bytes.Buffer{}
	}
	_, err := b.WriteString("ofType {kind,name \n")
	if err != nil {
		return err
	}
	if depth > 0 {
		err := generateDepthOfType(depth-1, b)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	return err
}

func formatString(s string) string {
	intro := strings.Join(strings.Split(s, "\n"), ",")
	intro = strings.ReplaceAll(intro, "{,", "{")
	intro = strings.Join(strings.Split(intro, " "), "")
	intro = strings.Join(strings.Split(intro, "\t"), "")
	intro = strings.Join(strings.Split(intro, ",}"), "}")
	return fmt.Sprintf(`{"query":"{%s}"}`, intro)
}

func init() {
	// The depth should be a reasonable depth for most types. This is used to describe types
	// that are based off other types.
	depth := 10
	b := &bytes.Buffer{}
	err := generateDepthOfType(depth, b)
	if err != nil {
		panic(err)
	}
	ofType := b.String()
	args := fmt.Sprintf(`args {
		name
		description
		defaultValue
		type {
			kind
			name
			%s
		}
	}`, ofType)
	query := fmt.Sprintf(`__schema {
		queryType { name }
		mutationType { name }
		subscriptionType { name }
		types {
			kind
			name
			description
			fields(includeDeprecated: true) {
				name
				description
				%s
				type {
					kind
					name	
					%s
				}
				isDeprecated
				deprecationReason
			}
			inputFields {
				defaultValue
				name
				description
				type {
					kind
					name
					%s
				}
			}
			interfaces {
				kind
				name
				%s
			}
			enumValues(includeDeprecated:true){
				name
				description
				isDeprecated
				deprecationReason
			}
			possibleTypes {
				kind
				name
				%s
			}
		}
		directives {
			name
			description
			locations
			onOperation
			onFragment
			onField
			%s
		}
	}`, args, ofType, ofType, ofType, ofType, args)
	// has to encapsulate in brackets
	IntrospectionQuery = formatString(query)
}

// IntrospectionQuery is the query used to gather the types and schema of a GraphQL server.
// In the future we may need to deepen, if this is the case we may want to dynamically build this query.
var IntrospectionQuery = ""
