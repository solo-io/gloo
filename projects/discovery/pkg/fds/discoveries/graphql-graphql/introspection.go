package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rotisserie/eris"
	"github.com/wundergraph/graphql-go-tools/pkg/ast"
	"github.com/wundergraph/graphql-go-tools/pkg/astprinter"
	"github.com/wundergraph/graphql-go-tools/pkg/introspection"
)

/*
	Introspection is GraphQls way to identify the schema, types, directions, ect. on a GraphQL server. It will enable
	you to discover all of a server's schema. Gloo Uses Introspection to do just that, so that we can discover the schema
	of a GraphQL server, and serve it on the GraphQLAPI CRs (custom resources).
*/

type JSONResult struct {
	Data interface{} `json:"data"`
}

// postRequestTimeout is the default timeout for making the post request to get the introspectionData.
const postRequestTimeout = 5 * time.Second

// GetGraphQLSchema will return the schema of the GraphQL server using introspection.
func GetGraphQLSchema(u *url.URL) (string, error) {
	intro, err := GetIntrospectionResultsFromHost(u)
	if err != nil {
		return "", err
	}
	return GetSchemaFromIntrospectionJSON(intro)
}

// GetIntrospectionResultsFromHost will return the introspection results given a host. Note that this is done over POST Method.
func GetIntrospectionResultsFromHost(u *url.URL) ([]byte, error) {
	c := http.Client{
		Timeout: postRequestTimeout,
	}
	resp, err := c.Post(u.String(), "application/graphql", strings.NewReader(IntrospectionQuery))
	if err != nil {
		return nil, eris.Wrapf(err, "could not make graphql query at [%s]", u.Hostname())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonResult := JSONResult{}
	json.Unmarshal(body, &jsonResult)
	if data, ok := jsonResult.Data.(map[string]interface{}); ok {
		if _, hit := data["__schema"]; hit {
			data, err := json.Marshal(jsonResult.Data)
			if err != nil {
				return nil, eris.Wrapf(err, "error with introspection result: error with the data result [%s]", jsonResult.Data)
			}
			return data, nil
		}
	}
	return nil, eris.New(fmt.Sprintf("the request, to [%s], to retrieve the introspection: expected [{\"data\":\"result\"}]: instead returned the following [%s]", u.String(), string(body)))
}

// GetSchemaFromIntrospectionJSON will return the schema given an introspection JSON.
func GetSchemaFromIntrospectionJSON(introspectionJSON []byte) (string, error) {
	// given an introspection JSON create the schema from it
	converter := introspection.JsonConverter{}
	buf := bytes.NewBuffer(introspectionJSON)
	doc, err := converter.GraphQLDocument(buf)
	if err != nil {
		return "", eris.Wrap(err, "error converting introspection JSON to graphql schema")
	}
	filterDocument(doc)
	outWriter := &bytes.Buffer{}
	err = astprinter.PrintIndent(doc, nil, []byte("  "), outWriter)
	if err != nil {
		return "", eris.Wrap(err, "error formatting schema")
	}
	schemaOutputPretty := outWriter.Bytes()
	return string(schemaOutputPretty), nil
}

// these are the directives to remove
var removeDirectives = map[string]struct{}{
	"include":    {},
	"skip":       {},
	"deprecated": {},
}

// these are the types to remove
var removeObjectTypes = map[string]struct{}{
	"__InputValue": {},
	"__Directive":  {},
	"__Field":      {},
	"__Schema":     {},
	"__Type":       {},
	"__EnumValue":  {},
}

// these are the enums to remove
var removeEnumTypeDefinitions = map[string]struct{}{
	"__TypeKind":          {},
	"__DirectiveLocation": {},
}

// these are the scalars to remove
var removeScalarDefinitions = map[string]struct{}{
	"String":  {},
	"Boolean": {},
	"Int":     {},
	"ID":      {},
	"Float":   {},
}

// filterDocument will filter out all the pre-defined resources, types, and definitions on a
// graphql server. This is important because this can cause issues with validation in Envoy.
// Also this should filter out all the definitions that the client has not defined.
func filterDocument(doc *ast.Document) {
	nodesToRemove := []ast.Node{}

	for _, dir := range doc.Directives {
		for _, r := range dir.Arguments.Refs {
			bName := doc.DirectiveDefinitionNameBytes(r)
			if _, hit := removeDirectives[string(bName)]; hit {
				node, hit := doc.NodeByNameStr(bName.String())
				if hit {
					nodesToRemove = append(nodesToRemove, node)
				}
			}
		}
	}

	for i, _ := range doc.ObjectTypeDefinitions {
		name := doc.ObjectTypeDefinitionNameString(i)
		if _, hit := removeObjectTypes[name]; hit {
			node, hit := doc.NodeByNameStr(name)
			if hit {
				nodesToRemove = append(nodesToRemove, node)
			}
		}
	}

	for i, _ := range doc.EnumTypeDefinitions {
		name := doc.EnumTypeDefinitionNameString(i)
		if _, hit := removeEnumTypeDefinitions[name]; hit {
			node, hit := doc.NodeByNameStr(name)
			if hit {
				nodesToRemove = append(nodesToRemove, node)
			}
		}
	}

	for i, _ := range doc.ScalarTypeDefinitions {
		name := doc.ScalarTypeDefinitionNameString(i)
		if _, hit := removeScalarDefinitions[name]; hit {
			node, hit := doc.NodeByNameStr(name)
			if hit {
				nodesToRemove = append(nodesToRemove, node)
			}
		}
	}
	for _, n := range nodesToRemove {
		doc.RemoveRootNode(n)
	}
}
