package types

/**
 * Custom type definitions for GraphQL.
 */

type GraphQlOperationType int64

const (
	GraphQlOperationType_Query GraphQlOperationType = iota
	GraphQlOperationType_Mutation
	GraphQlOperationType_Subscription
)

func (g GraphQlOperationType) String() string {
	switch g {
	case GraphQlOperationType_Query:
		return "Query"
	case GraphQlOperationType_Mutation:
		return "Mutation"
	case GraphQlOperationType_Subscription:
		return "Subscription"
	default:
		return ""
	}
}
