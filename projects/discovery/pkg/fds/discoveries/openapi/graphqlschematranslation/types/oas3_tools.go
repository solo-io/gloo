package types

import openapi "github.com/getkin/kin-openapi/openapi3"

type SchemaNames struct {
	// Sorted in the following priority order
	FromRef    string
	FromSchema string
	FromPath   string

	/**
	Used when the preferred name is known, i.e. a new data def does not need to be created
	*/
	Preferred string
}

type Names struct {
	FromRef    string
	FromSchema string
	FromPath   string
}

type RequestSchemaAndNames struct {
	PayloadContentType string
	PayloadSchema      *openapi.Schema
	PayloadSchemaName  Names
	PayloadRequired    bool
}

type ResponseSchemaAndNames struct {
	ResponseContentType string
	ResponseSchema      *openapi.Schema
	ResponseSchemaName  Names
	StatusCode          string
}
