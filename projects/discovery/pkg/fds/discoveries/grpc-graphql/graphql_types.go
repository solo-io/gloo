package grpc

// Defined the GraphQL Scalar Types here as constants
const (
	GRAPHQL_STRING           = "String"
	GRAPHQL_BOOLEAN          = "Boolean"
	GRAPHQL_INT              = "Int"
	GRAPHQL_FLOAT            = "Float"
	GRAPHQL_NON_NULL_STRING  = "String!"
	GRAPHQL_NON_NULL_BOOLEAN = "Boolean!"
	GRAPHQL_NON_NULL_INT     = "Int!"
	GRAPHQL_NON_NULL_FLOAT   = "Float!"
)

// GetNonNullType will return the non null type else it will return the type back, if it is not a non null type.
func GetNonNullType(t string) string {
	switch t {
	case GRAPHQL_NON_NULL_BOOLEAN:
		return GRAPHQL_BOOLEAN
	case GRAPHQL_NON_NULL_FLOAT:
		return GRAPHQL_FLOAT
	case GRAPHQL_NON_NULL_INT:
		return GRAPHQL_INT
	case GRAPHQL_NON_NULL_STRING:
		return GRAPHQL_STRING
	default:
		return t
	}
}
