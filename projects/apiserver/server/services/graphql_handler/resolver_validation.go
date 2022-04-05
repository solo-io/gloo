package graphql_handler

import (
	"github.com/ghodss/yaml"
	"github.com/rotisserie/eris"
	graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	InvalidResolverTypeErr = func(resolverType rpc_edge_v1.ValidateResolverYamlRequest_ResolverType) error {
		return eris.Errorf("invalid resolver type: %v", resolverType)
	}
	InvalidResolverYamlErr = func(err error) error {
		return eris.Wrapf(err, "invalid resolver yaml")
	}
)

func ValidateResolverYaml(resolverYaml string, resolverType rpc_edge_v1.ValidateResolverYamlRequest_ResolverType) error {
	// don't fail on empty yaml
	if resolverYaml == "" {
		return nil
	}

	switch resolverType {
	case rpc_edge_v1.ValidateResolverYamlRequest_REST_RESOLVER:
		return validateRestResolver(resolverYaml)
	case rpc_edge_v1.ValidateResolverYamlRequest_GRPC_RESOLVER:
		return validateGrpcResolver(resolverYaml)
	default:
		return InvalidResolverTypeErr(resolverType)
	}
}

func validateRestResolver(resolverYaml string) error {
	resolverJson, err := yamlToJson(resolverYaml)
	if err != nil {
		return err
	}

	resolver := &graphql_v1beta1.RESTResolver{}
	if err := protojson.Unmarshal(resolverJson, resolver); err != nil {
		return InvalidResolverYamlErr(err)
	}
	return nil
}

func validateGrpcResolver(resolverYaml string) error {
	resolverJson, err := yamlToJson(resolverYaml)
	if err != nil {
		return err
	}

	resolver := &graphql_v1beta1.GrpcResolver{}
	if err := protojson.Unmarshal(resolverJson, resolver); err != nil {
		return InvalidResolverYamlErr(err)
	}
	return nil
}

func yamlToJson(resolverYaml string) ([]byte, error) {
	resolverJson, err := yaml.YAMLToJSON([]byte(resolverYaml))
	if err != nil {
		return nil, eris.Wrapf(err, "failed to convert resolver yaml to json")
	}
	return resolverJson, nil
}
