package types

import openapi "github.com/getkin/kin-openapi/openapi3"

// Copyright IBM Corp. 2018. All Rights Reserved.
// Node module openapi-to-graphql
// This file is licensed under the MIT License.
// License text available at https//opensource.org/licenses/MIT

/**
 * Type definitions for the data created during preprocessing.
 */

type ProcessedSecurityScheme struct {
	RawName string
	Def     *openapi.SecurityScheme

	/**
	 * Stores the names of the authentication credentials
	 * NOTE Structure depends on the type of the protocol (basic, API key...)
	 * NOTE Mainly used for the AnyAuth viewers
	 */
	Parameters map[string]string

	/**
	 * JSON schema to create the viewer for this security scheme from.
	 */
	Schema *openapi.Schema

	/**
	 * The OAS which this operation originated from
	 */
	Oas *openapi.T
}

type PreprocessingData struct {
	/**
	 * List of operation objects
	 */
	Operations map[string]*Operation

	/**
	 * List of Operation objects
	 */
	CallbackOperations map[string]Operation

	/**
	 * List of all the used object names to avoid collision
	 */
	UsedTypeNames []string

	/**
	 * List of data definitions for JSON schemas already used.
	 */
	Defs []*DataDefinition

	/**
	 * The security definitions contained in the OAS. References are resolved.
	 *
	 * NOTE Keys are sanitized
	 * NOTE Does not contain OAuth 2.0-related security schemes
	 */
	Security map[string]ProcessedSecurityScheme

	/**
	 * Mapping between sanitized strings and their original ones
	 */
	SaneMap map[string]string

	/**
	 * All of the provided OASs
	 */
	Oass []*openapi.T
}
