package translation

import (
	"os"
	"path"
	"strings"

	"github.com/rotisserie/eris"
)

var (
	// Env var that has the path to the graphql .js files
	GraphqlJsRootEnvVar = "GRAPHQL_JS_ROOT"
	// Env var that has the path to the proto dependencies that the .js files require
	GraphqlProtoRootEnvVar = "GRAPHQL_PROTO_ROOT"

	DefaultGraphqlJsRoot    = "/usr/local/bin/js/"
	DefaultGraphqlProtoRoot = "/usr/local/bin/js/proto/"
)

// Gets the path to the js root directory. May be a relative or absolute path.
// This path will always end with a "/".
func GetGraphqlJsRoot() string {
	if jsRoot := os.Getenv(GraphqlJsRootEnvVar); jsRoot != "" {
		if !strings.HasSuffix(jsRoot, "/") {
			jsRoot += "/"
		}
		return jsRoot
	}
	return DefaultGraphqlJsRoot
}

// Gets the absolute path to the proto root directory.
// This path will always end with a "/".
func GetGraphqlProtoRoot() (string, error) {
	if protoRoot := os.Getenv(GraphqlProtoRootEnvVar); protoRoot != "" {
		// JS needs the absolute path to the proto dir, so we join path to current dir + path from repository root
		currentDir, err := os.Getwd()
		if err != nil {
			return "", eris.Wrap(err, "unable to get current directory path for running graphql script")
		}
		absPath := path.Join(currentDir, protoRoot)
		if !strings.HasSuffix(absPath, "/") {
			absPath += "/"
		}
		return absPath, nil
	}
	return DefaultGraphqlProtoRoot, nil
}
