package graphql

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/goutils"
	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/solo-io/go-utils/cliutils"
	. "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi-graphql/graphqlschematranslation/types"
)

func (t *OasToGqlTranslator) PreprocessOas(oass []*openapi.T, options SchemaOptions) *PreprocessingData {
	var data = &PreprocessingData{
		Operations:         map[string]*Operation{},
		CallbackOperations: map[string]Operation{},
		UsedTypeNames:      []string{GraphQlOperationType_Query.String(), GraphQlOperationType_Mutation.String(), GraphQlOperationType_Subscription.String()},
		Defs:               []*DataDefinition{},
		Security:           map[string]ProcessedSecurityScheme{},
		SaneMap:            map[string]string{},
		Oass:               oass,
	}

	for index, oas := range oass {
		SetOASInfoAndTitle(index, oas)
		currentSecurity := GetProcessedSecuritySchemes(oas, data)
		commonSecurityPropertyName := getCommonPropertyNames(data.Security, currentSecurity)
		for _, name := range commonSecurityPropertyName {
			t.handleWarningf(DUPLICATE_SECURITY_SCHEME,
				"The security scheme from OAS "+currentSecurity[name].Oas.Info.Title+" will be ignored",
				Location{Oas: oas},
				"Multiple OpenApiSpecs share security schemes with the same name %s", name,
			)
		}

		// Don't overwrite prexisting security schemes from other oas specs
		for key, secScheme := range currentSecurity {
			if _, ok := data.Security[key]; !ok {
				data.Security[key] = secScheme
			}
		}

		for path, pathItem := range oas.Paths {
			for httpMethodName, op := range pathItem.Operations() {
				if httpMethodName == http.MethodGet {
					t.TranslateQueryOperation(data, oas, path, pathItem, op)
				} else {
					// all other http method types get translated to mutation operation
					t.TranslateMutationOperation(data, oas, path, pathItem, httpMethodName, op)
				}
			}
		}
	}
	return data
}

// Set the OpenApiSpec title to something if it is empty
func SetOASInfoAndTitle(index int, oas *openapi.T) {
	if oas.Info == nil {
		oas.Info = &openapi.Info{}
	}
	if oas.Info.Title == "" {
		oas.Info.Title = fmt.Sprintf("OpenApiSpec #%d", index)
	} else {
		oas.Info.Title = fmt.Sprintf("OpenApiSpec '%s'", oas.Info.Title)
	}
}

func getCommonPropertyNames(security map[string]ProcessedSecurityScheme, security2 map[string]ProcessedSecurityScheme) []string {
	var result []string
	for key := range security {
		if _, ok := security2[key]; ok {
			result = append(result, key)
		}
	}
	return result
}

/**
 * NOTE: THIS IS CURRENTLY UNUSED as we don't support security schemes right now
 * Extracts the security schemes from given OAS and organizes the information in
 * a data structure that is easier for OpenAPI-to-GraphQL to use
 *
 * Here is the structure of the data:
 * {
 *   {string} [sanitized name] { Contains information about the security protocol
 *     {string} rawName           Stores the raw security protocol name
 *     {object} def               Definition provided by OAS
 *     {object} parameters        Stores the names of the authentication credentials
 *                                  NOTE: Structure will depend on the type of the protocol
 *                                    (e.g. basic authentication, API key, etc.)
 *                                  NOTE: Mainly used for the AnyAuth viewers
 *     {object} schema            Stores the GraphQL schema to create the viewers
 *   }
 * }
 *
 * Here is an example:
 * {
 *   MyApiKey: {
 *     rawName: "My_api_key",
 *     def: { ... },
 *     parameters: {
 *       apiKey: MyKeyApiKey
 *     },
 *     schema: { ... }
 *   }
 *   MyBasicAuth: {
 *     rawName: "My_basic_auth",
 *     def: { ... },
 *     parameters: {
 *       username: MyBasicAuthUsername,
 *       password: MyBasicAuthPassword,
 *     },
 *     schema: { ... }
 *   }
 * }
 */
func GetProcessedSecuritySchemes(oas *openapi.T, data *PreprocessingData) map[string]ProcessedSecurityScheme {
	var result map[string]ProcessedSecurityScheme
	security := oas.Components.SecuritySchemes
	var (
		description string
		parameters  map[string]string
		schema      *openapi.Schema
	)
	for key, protocol := range security {
		switch protocol.Value.Type {
		case "apiKey":
			{
				description = fmt.Sprintf("API Key credentials for the security protocol %s", key)
				if len(data.Oass) > 1 && oas.Info != nil {
					description += fmt.Sprintf("in OpenApiSpec %s", oas.Info.Title)
				}

				parameters = map[string]string{
					"apiKey": sanitizeString(key+"_apiKey", CaseStyle_camelCase),
				}

				schema = &openapi.Schema{
					Type:        "object",
					Description: description,
					Properties: map[string]*openapi.SchemaRef{
						"apiKey": {
							Value: &openapi.Schema{
								Type: "string",
							},
						},
					},
				}
			}
		case "http":
			{
				switch protocol.Value.Scheme {
				case "basic":
					{
						description = "Basic auth credentials for security protocol '" + key + "'"
						parameters = map[string]string{
							"username": sanitizeString(key+"_username", CaseStyle_camelCase),
							"password": sanitizeString(key+"_password", CaseStyle_camelCase),
						}
						schema = &openapi.Schema{
							Type:        "object",
							Description: description,
							Properties: map[string]*openapi.SchemaRef{
								"username": {
									Value: &openapi.Schema{
										Type: "string",
									},
								},
								"password": {
									Value: &openapi.Schema{
										Type: "string",
									},
								},
							},
						}
					}
				default:
					{
						//handleWarning
						fmt.Println("Currently unsupported HTTP authentication protocol type 'http' and scheme '" + protocol.Value.Scheme + "' in OAS")
					}
				}
			}
		case "openIdConnect":
			fmt.Println("OIDC is not currently supported security scheme")
		case "oauth2":
			fmt.Println("Oauth2 OAS support is provided using the `tokenJSONpath` option")
		}

		result = map[string]ProcessedSecurityScheme{
			key: {
				RawName:    "key",
				Def:        protocol.Value,
				Parameters: parameters,
				Schema:     schema,
				Oas:        oas,
			},
		}
	}
	return result
}

func (t *OasToGqlTranslator) TranslateQueryOperation(data *PreprocessingData, oas *openapi.T, path string, pathItem *openapi.PathItem, operation *openapi.Operation) {
	graphqlOperationType := GraphQlOperationType_Query
	operationString := FormatOperationString(http.MethodGet, path, oas)
	operationData := t.ProcessOperation(
		path,
		http.MethodGet,
		operationString,
		graphqlOperationType,
		operation,
		pathItem,
		oas,
		data)
	if _, ok := data.Operations[operationData.OperationId]; !ok {
		data.Operations[operationData.OperationId] = operationData
	} else {
		t.handleWarningf(DUPLICATE_OPERATIONID, "", Location{Oas: oas},
			"OperationId %s already exists in the openapischema. The duplicate operation will be ignored.", operationData.OperationId)
	}
}

func (t *OasToGqlTranslator) TranslateMutationOperation(data *PreprocessingData, oas *openapi.T, path string, pathItem *openapi.PathItem, httpMethod string, operation *openapi.Operation) {
	// check Put, Delete, etc.
	graphqlOperationType := GraphQlOperationType_Mutation
	operationString := FormatOperationString(httpMethod, path, oas)
	operationData := t.ProcessOperation(
		path,
		httpMethod,
		operationString,
		graphqlOperationType,
		operation,
		pathItem,
		oas,
		data)
	if _, ok := data.Operations[operationData.OperationId]; !ok {
		data.Operations[operationData.OperationId] = operationData
	}

}

func (t *OasToGqlTranslator) ProcessOperation(path,
	method,
	operationString string,
	operationType GraphQlOperationType,
	operation *openapi.Operation,
	pathItem *openapi.PathItem,
	oas *openapi.T,
	data *PreprocessingData,
) *Operation {

	// Generation operation Id
	operationId := operation.OperationID
	// operationId may not exist
	if len(operationId) == 0 {
		operationId = GenerateOperationId(method, path)
	}

	// Generate operation description
	description := operation.Description
	if len(description) == 0 {
		description = operation.Summary
	}
	description += "\n\nEquivalent to " + operationString

	// Generate request schema for operation
	requestSchema := GetRequestSchemaAndNames(path, operation, oas)
	var payloadDefinition *DataDefinition
	if requestSchema.PayloadSchema != nil {
		payloadDefinition = t.CreateDataDef(requestSchema.PayloadSchemaName, requestSchema.PayloadSchema, true, data, nil, oas)
	}

	resSchemaAndName := t.GetResponseSchemaAndNames(path, method, operation, oas, data)
	links := t.GetLinks(path, method, operation, oas, data)
	responseDefinition := t.CreateDataDef(
		resSchemaAndName.ResponseSchemaName,
		resSchemaAndName.ResponseSchema,
		false,
		data,
		links,
		oas)

	parameters := GetParameters(path, method, operation, pathItem, oas)
	servers := GetServers(operation, pathItem, oas)

	securityRequirements := GetSecurityRequirements(operation, data.Security, oas)

	return &Operation{
		OperationId:          operationId,
		OperationString:      operationString,
		Description:          description,
		Path:                 path,
		Method:               method,
		PayloadContentType:   requestSchema.PayloadContentType,
		PayloadDefinition:    payloadDefinition,
		PayloadRequired:      requestSchema.PayloadRequired,
		ResponseContentType:  resSchemaAndName.ResponseContentType,
		ResponseDefinition:   responseDefinition,
		Parameters:           parameters,
		SecurityRequirements: securityRequirements,
		Servers:              servers,
		InViewer:             false,
		OperationType:        operationType,
		StatusCode:           resSchemaAndName.StatusCode,
		Oas:                  oas,
	}
}

func GetServers(operation *openapi.Operation, item *openapi.PathItem, oas *openapi.T) openapi.Servers {
	var servers openapi.Servers
	if operation.Servers != nil && len(*operation.Servers) > 0 {
		servers = *operation.Servers
	} else if len(item.Servers) > 0 {
		servers = oas.Servers
	} else if len(oas.Servers) > 0 {
		servers = oas.Servers
	}

	if len(servers) == 0 {
		servers = append(servers, &openapi.Server{
			URL: "/",
		})
	}
	return servers
}

// https://swagger.io/docs/specification/authentication/
func GetSecurityRequirements(operation *openapi.Operation, securitySchemes map[string]ProcessedSecurityScheme, oas *openapi.T) []string {
	var results []string

	// First, consider global requirements
	globalSecurity := oas.Security
	if len(globalSecurity) > 0 {
		for _, secReq := range globalSecurity {
			for schemaKey, _ := range secReq {
				if obj, ok := securitySchemes[schemaKey]; ok && obj.Def.Type != "oauth2" {
					results = append(results, schemaKey)
				}
			}
		}
	}

	// Second, consider operation requirements
	localSecurity := operation.Security
	if localSecurity != nil {
		for _, secReq := range *localSecurity {
			for schemaKey, _ := range secReq {
				if obj, ok := securitySchemes[schemaKey]; ok && obj.Def.Type != "oauth2" {
					if !cliutils.Contains(results, schemaKey) {
						results = append(results, schemaKey)
					}
				}
			}
		}
	}
	return results

}

func (t *OasToGqlTranslator) CreateDataDef(names Names, schema *openapi.Schema, isInputObjectType bool, data *PreprocessingData, links map[string]*openapi.Link, oas *openapi.T) *DataDefinition {
	saneLinks := map[string]*openapi.Link{}
	for linkKey, link := range links {
		saneLinks[sanitizeString(linkKey, CaseStyle_camelCase)] = link
	}

	preferredName := GetPreferredName(names)

	if schema == nil {
		t.handleWarningf(MISSING_SCHEMA, "",
			Location{Oas: oas}, "Could not create data definition for schema with name %s", preferredName)
		return &DataDefinition{
			PreferredName: preferredName,
			Schema: &openapi.Schema{
				ExtensionProps: openapi.ExtensionProps{},
				Type:           "object",
				Description:    "Placeholder schema because this does not have a response schema",
			},
			TargetGraphQLType: "json",
		}
	}

	existingDataDef := GetExistingDataDef(preferredName, schema, data.Defs)
	if existingDataDef != nil {
		// found existing data definition

		if len(saneLinks) != 0 {
			if len(existingDataDef.Links) != 0 {
				for link, linkRef := range saneLinks {
					if _, ok := existingDataDef.Links[link]; !ok {
						existingDataDef.Links[link] = linkRef
					}
				}
			} else {
				existingDataDef.Links = saneLinks
			}
		}
		return existingDataDef
	}

	usedNames := map[string]bool{}
	for _, name := range data.UsedTypeNames {
		usedNames[name] = true
	}
	for name, _ := range data.SaneMap {
		usedNames[name] = true
	}
	name := GetSchemaName(names, usedNames, CaseStyle_PascalCase)
	saneName := sanitizeString(name, CaseStyle_PascalCase)
	saneInputName := goutils.Capitalize(saneName + "Input")

	targetGraphQlType := GetSchemaTargetGraphQlType(schema, data)

	dataDef := &DataDefinition{
		PreferredName:              preferredName,
		Schema:                     schema,
		TargetGraphQLType:          targetGraphQlType,
		SubDefinitions:             nil,
		Links:                      saneLinks,
		GraphQLTypeName:            saneName,
		GraphQLInputObjectTypeName: saneInputName,
	}

	if cliutils.Contains([]string{"object", "list", "enum"}, targetGraphQlType) {
		data.UsedTypeNames = append(data.UsedTypeNames, saneName, saneInputName)
		data.Defs = append(data.Defs, dataDef)
	}
	//todo-  support AnyOf and AllOf
	if len(targetGraphQlType) > 0 {
		if targetGraphQlType == "list" {
			if itemsSchema := schema.Items.Value; itemsSchema != nil {
				itemName := saneName + "ListItem"
				if refName := GetRefName(schema.Items.Ref); refName != "" {
					itemName = refName
				}
				itemNames := Names{
					FromRef: itemName,
				}
				subDefinition := t.CreateDataDef(itemNames, itemsSchema, isInputObjectType, data, nil, oas)
				dataDef.SubDefinitions = &SubDefinitions{
					ListType: subDefinition,
				}
			}
		} else if targetGraphQlType == "object" {
			dataDef.SubDefinitions = &SubDefinitions{
				InputObjectType: map[string]*DataDefinition{},
			}
			if len(schema.Properties) > 0 {
				t.addObjectPropertiesToDataDef(dataDef, schema, isInputObjectType, data, oas)
			} else {
				t.handleWarningf(OBJECT_MISSING_PROPERTIES, "", Location{Oas: oas},
					"Schema %s does not have any properties. Using JSON type.", schema.Title)
				dataDef.TargetGraphQLType = "json"
			}
		}
	} else {
		t.handleWarningf(UNKNOWN_TARGET_TYPE, "", Location{Oas: oas},
			"No target graphqltype could be identified for schema %s. Using JSON type.", schema.Title)
		dataDef.TargetGraphQLType = "json"
	}
	return dataDef
}
