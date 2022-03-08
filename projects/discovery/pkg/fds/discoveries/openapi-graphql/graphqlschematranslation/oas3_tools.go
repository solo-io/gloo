package graphql

import (
	"fmt"
	"net/http"
	url2 "net/url"
	"regexp"
	"strings"

	"github.com/Masterminds/goutils"
	openapi "github.com/getkin/kin-openapi/openapi3"
	. "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi-graphql/graphqlschematranslation/types"
)

var SupportedRequestContentTypes = []string{
	"application/json",
	"applicaiton/x-www-form-urlencoded",
	"*/*",
}

func GetRequestSchemaAndNames(
	path string,
	operation *openapi.Operation,
	oas *openapi.T,
) *RequestSchemaAndNames {
	payloadContentType, requestBodyObj := GetRequestBodyObject(operation, oas)
	if payloadContentType != "" {
		payloadSchemaObj := requestBodyObj.Content.Get(payloadContentType).Schema
		payloadSchema := payloadSchemaObj.Value
		payloadRequired := requestBodyObj.Required

		names := Names{
			FromRef:    GetRefName(payloadSchemaObj.Ref),
			FromSchema: payloadSchema.Title,
			FromPath:   InferResourceNameFromPath(path),
		}
		supportedRequestContentType := false
		for _, ct := range SupportedRequestContentTypes {
			if strings.Contains(payloadContentType, ct) {
				supportedRequestContentType = true
			}
		}
		if !supportedRequestContentType {
			// If there is no supported content type, then create an input argument that will just be
			// the payload as a string
			reducedName := ReduceStringArray(strings.Split(payloadContentType, "/"), func(name, term string, _ int) string {
				return name + goutils.Capitalize(term)
			}, "")
			saneContentTypeName := goutils.Uncapitalize(reducedName)
			names = Names{
				FromPath: saneContentTypeName,
			}
			description := fmt.Sprintf("String represents payload of content type '%s'", payloadContentType)
			if payloadSchema.Description != "" {
				description += "\n\nOriginal top level description: " + payloadSchema.Description
			}
			payloadSchema = &openapi.Schema{
				Type:        "string",
				Description: description,
			}
		}

		return &RequestSchemaAndNames{
			PayloadContentType: payloadContentType,
			PayloadSchema:      payloadSchema,
			PayloadSchemaName:  names,
			PayloadRequired:    payloadRequired,
		}

	}

	return &RequestSchemaAndNames{
		PayloadRequired: false,
	}
}

func (t *OasToGqlTranslator) GetResponseSchemaAndNames(path string, method string, operation *openapi.Operation, oas *openapi.T, data *PreprocessingData) *ResponseSchemaAndNames {
	statusCode := t.GetResponseStatusCode(path, method, operation, oas, data)
	if statusCode == "" {
		return &ResponseSchemaAndNames{}
	}
	responseContentType, responseObject := GetResponseObject(operation, statusCode, oas)
	if len(responseContentType) > 0 {
		responseSchema := responseObject.Content.Get(responseContentType).Schema
		var fromRefName = GetRefName(responseSchema.Ref)
		names := Names{
			FromRef:    fromRefName,
			FromSchema: responseSchema.Value.Title,
			FromPath:   InferResourceNameFromPath(path),
		}
		return &ResponseSchemaAndNames{
			ResponseContentType: responseContentType,
			ResponseSchema:      responseSchema.Value,
			ResponseSchemaName:  names,
			StatusCode:          statusCode,
		}
	} else {
		return &ResponseSchemaAndNames{
			ResponseSchemaName: Names{
				FromPath: InferResourceNameFromPath(path),
			},
			ResponseContentType: "application/json",
		}
	}
}

// Looks for BaseUrl in OpenApi Operation spec (See https://swagger.io/docs/specification/paths-and-operations/ "Overriding Global Servers")
// If operation server override is not available, uses global server.
// Strips everything but path from the server URL, as we use upstream refs here instead of full URLs for the request.
// Returns base url path and true if found, else returns empty string and false
func (t *OasToGqlTranslator) GetBaseUrlPath(operation *Operation) (string, error) {
	// Remove trailing slash from url
	if len(operation.Servers) > 0 {
		url := BuildServerUrl(operation.Servers[0])
		if len(operation.Servers) > 1 {
			t.handleWarning(MORE_THAN_ONE_SERVER,
				"More than one server provided in operation",
				Location{Operation: operation},
				"Randomly selecting first url: "+url)
		}
		URL, err := url2.Parse(url)
		if err != nil || URL == nil {
			t.handleWarningf(INVALID_SERVER_URL, "",
				Location{Operation: operation},
				"Invalid server URL provided in operation: %s, attempting to use top level server URL", url)
		} else {
			return strings.TrimSuffix(URL.Path, "/"), nil
		}
	}

	oas := operation.Oas
	// oas.Servers is guaranteed to have atleast one entry by the openapi-spec parser
	// https://swagger.io/specification/#openapi-object
	// > If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of /.
	url := BuildServerUrl(oas.Servers[0])
	if len(oas.Servers) > 1 {
		t.handleWarningf(MORE_THAN_ONE_SERVER,
			"Randomly selecting first url: "+url,
			Location{Oas: operation.Oas},
			"More than one server provided in open api spec servers")
	}
	URL, err := url2.Parse(url)
	if err != nil || URL == nil {
		return "", t.createErrorf(INVALID_SERVER_URL, "",
			Location{Operation: operation},
			"Invalid server URL provided in operation: %s", url)
	}
	return strings.TrimSuffix(URL.Path, "/"), nil
}

func BuildServerUrl(server *openapi.Server) string {
	url := server.URL
	for k, v := range server.Variables {
		url = strings.ReplaceAll(url, "{"+k+"}", v.Default)
	}
	return url
}

func GenerateOperationId(method, path string) string {
	str := fmt.Sprintf("%s %s", strings.ToLower(method), path)
	return sanitizeString(str, CaseStyle_camelCase)
}

func (t *OasToGqlTranslator) GetResponseStatusCode(path string, method string, operation *openapi.Operation, oas *openapi.T, data *PreprocessingData) string {
	successCodeRe := regexp.MustCompile(`2[0-9]{2}|2XX`)
	// Look for success codes first, then return default if none exist
	var successCodes []string
	for code, _ := range operation.Responses {
		if successCodeRe.Match([]byte(code)) {
			successCodes = append(successCodes, code)
		}
	}
	if len(successCodes) == 1 {
		return successCodes[0]
	} else if len(successCodes) > 1 {
		t.handleWarningf(MULTIPLE_RESPONSES,
			"The response object with the HTTP code "+successCodes[0]+" will be selected",
			Location{Path: []string{path, method}},
			"Operation %s contains multiple possible successful response object (HTTP Code 200-299 or 2XX). Only one can be chosen.",
			FormatOperationString(method, path, oas))
	}
	if operation.Responses.Default() != nil {
		return "default"
	}
	return ""
}

func (t *OasToGqlTranslator) GetLinks(path, method string, operation *openapi.Operation, oas *openapi.T, data *PreprocessingData) map[string]*openapi.Link {
	var links = map[string]*openapi.Link{}
	statusCode := t.GetResponseStatusCode(path, method, operation, oas, data)
	if statusCode == "" {
		return links
	}
	response := operation.Responses[statusCode].Value
	for linkKey, linkRef := range response.Links {
		links[linkKey] = linkRef.Value
	}
	return links
}

func stringToHttpMethod(method string) string {
	switch strings.ToLower(method) {
	case "get":
		return http.MethodGet

	case "put":
		return http.MethodPut

	case "post":
		return http.MethodPost

	case "patch":
		return http.MethodPatch

	case "delete":
		return http.MethodDelete

	case "options":
		return http.MethodOptions

	case "head":
		return http.MethodHead

	default:
		return ""
	}
}
