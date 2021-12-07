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
