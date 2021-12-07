package graphql

type MitigationType uint32

const (
	/**
	 * Problems with the OAS
	 *
	 * Should be caught by the module oas-validator
	 */
	INVALID_OAS MitigationType = iota
	UNNAMED_PARAMETER

	// General problems
	AMBIGUOUS_UNION_MEMBERS
	CANNOT_GET_FIELD_TYPE
	COMBINE_SCHEMAS
	DUPLICATE_FIELD_NAME
	DUPLICATE_LINK_KEY
	INVALID_HTTP_METHOD
	INPUT_UNION
	MISSING_RESPONSE_SCHEMA
	MISSING_SCHEMA
	MULTIPLE_RESPONSES
	NON_APPLICATION_JSON_SCHEMA
	OBJECT_MISSING_PROPERTIES
	UNKNOWN_TARGET_TYPE
	UNRESOLVABLE_SCHEMA
	UNSUPPORTED_HTTP_SECURITY_SCHEME
	UNSUPPORTED_JSON_SCHEMA_KEYWORD
	CALLBACKS_MULTIPLE_OPERATION_OBJECTS
	EXCESSIVE_NESTING
	INVALID_LIST_TYPE
	GRAPHQL_SCHEMA_CREATION_ERR

	// Links
	AMBIGUOUS_LINK
	LINK_NAME_COLLISION
	UNRESOLVABLE_LINK
	LINK_LOCATION_UNSUPPORTED
	MORE_THAN_ONE_LINK_PARAM_TEMPLATE
	LINK_PARAM_TEMPLATE_PRESENT
	LINK_UNSUPPORTED_EXTRACTION

	// Multiple OAS
	DUPLICATE_OPERATIONID
	DUPLICATE_SECURITY_SCHEME
	MULTIPLE_OAS_SAME_TITLE

	// Miscellaneous
	OAUTH_SECURITY_SCHEME
	MORE_THAN_ONE_SERVER

	INVALID_SERVER_URL
)

var Mitigations = map[MitigationType]string{
	/**
	 * Problems with the OAS
	 *
	 * Should be caught by the module oas-validator
	 */
	INVALID_OAS:       "Ignore issue and continue.",
	UNNAMED_PARAMETER: "Ignore parameter.",

	// Server
	//NO_SERVER_URL:        "Skip processing oas and continue",
	MORE_THAN_ONE_SERVER: "Choose first server and continue.",
	INVALID_SERVER_URL:   "Skip processing operation and continue.",

	// General problems
	AMBIGUOUS_UNION_MEMBERS:              "Ignore issue and continue.",
	CANNOT_GET_FIELD_TYPE:                "Ignore field and continue.",
	COMBINE_SCHEMAS:                      "Ignore combine schema keyword and continue.",
	DUPLICATE_FIELD_NAME:                 "Ignore field and maintain preexisting field.",
	DUPLICATE_LINK_KEY:                   "Ignore link and maintain preexisting link.",
	INPUT_UNION:                          "The data will be stored in an arbitrary JSON type.",
	INVALID_HTTP_METHOD:                  "Ignore operation and continue.",
	MISSING_RESPONSE_SCHEMA:              "Ignore operation.",
	MISSING_SCHEMA:                       "Use arbitrary JSON type.",
	MULTIPLE_RESPONSES:                   "Select first response object with successful status code (200-299).",
	NON_APPLICATION_JSON_SCHEMA:          "Ignore schema",
	OBJECT_MISSING_PROPERTIES:            "The (sub-)object will be stored in an arbitray JSON type.",
	UNKNOWN_TARGET_TYPE:                  "The data will be stored in an arbitrary JSON type.",
	UNRESOLVABLE_SCHEMA:                  "Ignore and continue. May lead to unexpected behavior.",
	UNSUPPORTED_HTTP_SECURITY_SCHEME:     "Ignore security scheme.",
	UNSUPPORTED_JSON_SCHEMA_KEYWORD:      "Ignore keyword and continue.",
	CALLBACKS_MULTIPLE_OPERATION_OBJECTS: "Select arbitrary operation object",
	EXCESSIVE_NESTING:                    "Stopping nesting at 50.",
	INVALID_LIST_TYPE:                    "Not creating list type.",
	GRAPHQL_SCHEMA_CREATION_ERR:          "Not creating schema.",

	// Links
	AMBIGUOUS_LINK:              `Use first occurrence of "#/".`,
	LINK_NAME_COLLISION:         "Ignore link and maintain preexisting field.",
	UNRESOLVABLE_LINK:           "Ignore link.",
	LINK_PARAM_TEMPLATE_PRESENT: "Link param templates are not currently supported. Ignoring everything outside the extraction in the template.",
	LINK_UNSUPPORTED_EXTRACTION: "Link will be ignored.",

	// Multiple OAS
	DUPLICATE_OPERATIONID:     "Ignore operation and maintain preexisting operation.",
	DUPLICATE_SECURITY_SCHEME: "Ignore security scheme and maintain preexisting scheme.",
	MULTIPLE_OAS_SAME_TITLE:   "Ignore issue and continue.",

	// Miscellaneous
	OAUTH_SECURITY_SCHEME: `Ignore security scheme`,
}
