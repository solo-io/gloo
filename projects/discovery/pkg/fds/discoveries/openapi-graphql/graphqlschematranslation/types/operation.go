package types

import (
	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/graphql-go/graphql"
)

/**
 * Type definitions for the objects created during preprocessing for every
 * operation in the OAS.
 */

type SubDefinitions struct {
	ListType        *DataDefinition            // For GraphQL list type
	InputObjectType map[string]*DataDefinition // For GraphQL (input) object type
	UnionType       []*DataDefinition          // For GraphQL union type
}

type DataDefinition_GraphQlType struct {
	graphql.Object
	graphql.List
	graphql.Union
	graphql.Enum
	graphql.Scalar
}

type DataDefinition_GraphQlInputObjectType struct {
	graphql.List
	graphql.InputObject
}

type DataDefinition struct {
	// OAS-related

	// Ideal name for the GraphQL type and is used with the schema to identify a specific GraphQL type
	PreferredName string

	// The schema of the data type, why may have gone through some resolution, and is used with preferredName to identify a specific GraphQL type
	Schema *openapi.Schema

	/**
	 * Similar to the required property in object schemas but because of certain
	 * keywords to combine schemas, e.g. "allOf", this resolves the required
	 * property in all member schemas
	 */
	Required []string

	// The type GraphQL type this dataDefintion will be created into
	TargetGraphQLType string

	// Collapsed link objects from all operations returning the same response data
	Links map[string]*openapi.Link

	/**
	 * Data definitions of subschemas in the schema
	 *
	 * I.e. If the dataDef is a list type, the subDefinition is a reference to the
	 * list item type
	 *
	 * Or if the dataDef is an object type, the subDefinitions are references to
	 * the field types
	 *
	 * Or if the dataDef is a union type, the subDefinitions are references to
	 * the member types
	 */
	SubDefinitions *SubDefinitions

	// GraphQL-related

	// The potential name of the GraphQL type if it is created
	GraphQLTypeName string

	// The potential name of the GraphQL input object type if it is created
	GraphQLInputObjectTypeName string

	// The GraphQL type if it is created
	GraphQlType graphql.Type

	// The GraphQL input object type if it is created
	GraphQLInputObjectType graphql.Type
}

type Operation struct {
	/**
	 * Identifier of the operation - may be created by concatenating method & path
	 */
	OperationId string

	/**
	 * A combination of the operation method and path (and the title of the OAS
	 * where the operation originates from if multiple OASs are provided) in the
	 * form of
	 *
	 * {title of OAS (if applicable)} {method in ALL_CAPS} {path}
	 *
	 * Used for documentation and logging
	 */
	OperationString string

	/**
	 * Human-readable description of the operation
	 */
	Description string

	/**
	 * URL path of this operation
	 */
	Path string

	/**
	 * HTTP method for this operation
	 */
	Method string

	/**
	 * Content-type of the request payload
	 */
	PayloadContentType string

	/**
	 * Information about the request payload (if any)
	 */
	PayloadDefinition *DataDefinition

	/**
	 * Determines wheter request payload is required for the request
	 */
	PayloadRequired bool

	/**
	 * Content-type of the request payload
	 */
	ResponseContentType string

	/**
	 * Information about the response payload
	 */
	ResponseDefinition *DataDefinition

	/**
	 * List of parameters of the operation
	 */
	Parameters openapi.Parameters

	/**
	 * List of keys of security schemes required by this operation
	 *
	 * NOTE Keys are sanitized
	 * NOTE Does not contain OAuth 2.0-related security schemes
	 */
	//todo - Unused in first iteration of implementation
	SecurityRequirements []string

	/**
	 * (Local) server definitions of the operation.
	 */
	Servers openapi.Servers

	/**
	 * Whether this operation should be placed in an authentication viewer
	 * (cannot be true if "viewer" option passed to OpenAPI-to-GraphQL is false).
	 */
	InViewer bool

	/**
	 * Type of root operation type, i.e. whether the generated field should be
	 * added to the Query, Mutation, or Subscription root operation
	 */
	OperationType GraphQlOperationType

	/**
	 * The success HTTP code, 200-299, destined to become a GraphQL object type
	 */
	StatusCode string

	/**
	 * The OAS which this operation originated from
	 */
	Oas *openapi.T
}
