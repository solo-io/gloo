package azure

import (
	"fmt"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/coreplugins/common"
	"github.com/solo-io/gloo/pkg/plugin"
)

func init() {
	plugin.Register(&Plugin{apiKeys: make(map[string]string)}, nil)
}

type Plugin struct {
	isNeeded bool
	hostname string
	apiKeys  map[string]string
}

const (
	// define Upstream type name
	UpstreamTypeAzure = "azure"

	// generic plugin info
	filterName  = "io.solo.azure_functions"
	pluginStage = plugin.OutAuth

	masterKeyName = "_master"

	// function-specific metadata
	functionHost = "host"
	functionPath = "path"
)

func (p *Plugin) GetDependencies(cfg *v1.Config) *plugin.Dependencies {
	deps := new(plugin.Dependencies)
	for _, upstream := range cfg.Upstreams {
		if upstream.Type != UpstreamTypeAzure {
			continue
		}
		azureUpstream, err := DecodeUpstreamSpec(upstream.Spec)
		if err != nil {
			// errors will be handled during validation
			// TODO: consider logging error here
			continue
		}
		deps.SecretRefs = append(deps.SecretRefs, azureUpstream.SecretRef)
	}
	return deps
}

func (p *Plugin) HttpFilters(_ *plugin.FilterPluginParams) []plugin.StagedFilter {
	defer func() {
		p.isNeeded = false
		p.hostname = ""
		p.apiKeys = make(map[string]string)
	}()

	if p.isNeeded {
		return []plugin.StagedFilter{{HttpFilter: &envoyhttp.HttpFilter{Name: filterName}, Stage: pluginStage}}
	}
	return nil
}

func (p *Plugin) ProcessRoute(_ *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	return nil
}

func (p *Plugin) ProcessUpstream(params *plugin.UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error {
	if in.Type != UpstreamTypeAzure {
		return nil
	}
	p.isNeeded = true

	out.Type = envoyapi.Cluster_LOGICAL_DNS
	out.DnsLookupFamily = envoyapi.Cluster_V4_ONLY

	azureUpstream, err := DecodeUpstreamSpec(in.Spec)
	if err != nil {
		return errors.Wrap(err, "invalid Azure upstream spec")
	}

	p.hostname = azureUpstream.GetHostname()
	out.Hosts = append(out.Hosts, &envoycore.Address{Address: &envoycore.Address_SocketAddress{SocketAddress: &envoycore.SocketAddress{
		Address:       p.hostname,
		PortSpecifier: &envoycore.SocketAddress_PortValue{PortValue: 443},
	}}})
	out.TlsContext = &envoyauth.UpstreamTlsContext{
		Sni: p.hostname,
	}

	if azureUpstream.SecretRef != "" {
		azureSecrets, ok := params.Secrets[azureUpstream.SecretRef]
		if !ok {
			return errors.Errorf("azure secrets for ref %v not found", azureUpstream.SecretRef)
		}

		p.apiKeys = azureSecrets
	}

	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	common.InitFilterMetadata(filterName, out.Metadata)
	out.Metadata.FilterMetadata[filterName] = &types.Struct{
		Fields: map[string]*types.Value{},
	}

	return nil
}

func (p *Plugin) ParseFunctionSpec(params *plugin.FunctionPluginParams, in v1.FunctionSpec) (*types.Struct, error) {
	if params.UpstreamType != UpstreamTypeAzure {
		return nil, nil
	}

	functionSpec, err := DecodeFunctionSpec(in)
	if err != nil {
		return nil, errors.Wrap(err, "invalid Azure Functions spec")
	}

	path, err := getPath(functionSpec, p.apiKeys)
	if err != nil {
		return nil, err
	}

	return getFunctionStruct(p.hostname, path), nil
}

func getApiKey(apiKeys map[string]string, keyNames []string) (string, error) {
	for _, keyName := range keyNames {
		apiKey, ok := apiKeys[keyName]
		if ok && apiKey != "" {
			return apiKey, nil
		}
	}

	return "", errors.New("secret not found")
}

func getPathParameters(functionSpec *FunctionSpec, apiKeys map[string]string) (string, error) {
	var keyNames []string
	switch functionSpec.AuthLevel {
	case "anonymous":
		return "", nil
	case "function":
		keyNames = []string{functionSpec.FunctionName, masterKeyName}
	case "admin":
		keyNames = []string{masterKeyName}
	}

	apiKey, err := getApiKey(apiKeys, keyNames)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("?code=%s", apiKey), nil
}

func getPath(functionSpec *FunctionSpec, apiKeys map[string]string) (string, error) {
	functionName := functionSpec.FunctionName

	pathParameters, err := getPathParameters(functionSpec, apiKeys)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/api/%s%s", functionName, pathParameters), nil
}

func getFunctionStruct(host string, path string) *types.Struct {
	return &types.Struct{
		Fields: map[string]*types.Value{
			functionHost: {Kind: &types.Value_StringValue{StringValue: host}},
			functionPath: {Kind: &types.Value_StringValue{StringValue: path}},
		},
	}
}
