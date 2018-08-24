package azure

import (
	"context"
	"fmt"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/pluginutils"
)

const (
	pluginStage   = plugins.OutAuth
	masterKeyName = "_master"
)

func init() {
	plugins.RegisterFunc(NewAzurePlugin)
}

type plugin struct {
	recordedUpstreams map[string]*azure.UpstreamSpec
	apiKeys           map[string]string
	ctx               context.Context
}

func NewAzurePlugin() plugins.Plugin {
	return &plugin{recordedUpstreams: make(map[string]*azure.UpstreamSpec)}
}

func (p *plugin) Init(params plugins.InitParams) error {
	p.ctx = params.Ctx
	return nil
}

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, out *envoyapi.Cluster) error {
	upstreamSpec, ok := in.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Azure)
	if !ok {
		// not ours
		return nil
	}
	azureUpstream := upstreamSpec.Azure

	// configure Envoy cluster routing info
	out.Type = envoyapi.Cluster_LOGICAL_DNS
	// TODO(yuval-k): why do we need to make sure we use ipv4 only dns?
	out.DnsLookupFamily = envoyapi.Cluster_V4_ONLY
	hostname := GetHostname(upstreamSpec.Azure)

	pluginutils.EnvoySingleEndpointLoadAssignment(out, hostname, 443)

	out.TlsContext = &envoyauth.UpstreamTlsContext{
		// TODO(yuval-k): Add verification context
		Sni: hostname,
	}

	if azureUpstream.SecretRef != "" {
		// TODO: namespace
		azureSecrets, err := params.Snapshot.SecretList.Find("", azureUpstream.SecretRef)
		if err != nil {
			return errors.Wrapf(err, "azure secrets for ref %v not found", azureUpstream.SecretRef)
		}
		p.apiKeys = azureSecrets.Data
	}

	// TODO(yuval-k): What about namespace?!
	p.recordedUpstreams[in.Metadata.Name] = upstreamSpec.Azure

	return nil
}

func (p *plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	return pluginutils.MarkHeaders(p.ctx, in, out, func(spec *v1.Destination) ([]*envoycore.HeaderValueOption, error) {
		// check if it's aws destination
		if spec.DestinationSpec == nil {
			return nil, nil
		}
		azureDestinationSpec, ok := spec.DestinationSpec.DestinationType.(*v1.DestinationSpec_Azure)
		if !ok {
			return nil, nil
		}
		// get upstream
		upstreamSpec, ok := p.recordedUpstreams[spec.UpstreamName]
		if !ok {
			// TODO(yuval-k): panic in debug
			return nil, errors.Errorf("%v is not an Azure upstream", spec.UpstreamName)
		}
		// should be aws upstream

		// get function
		functionName := azureDestinationSpec.Azure.FunctionName
		for _, functionSpec := range upstreamSpec.Functions {
			if functionSpec.FunctionName == functionName {
				path, err := getPath(functionSpec, p.apiKeys)
				if err != nil {
					return nil, err
				}
				hostname := GetHostname(upstreamSpec)
				// TODO: this is removed by: https://github.com/envoyproxy/envoy/pull/4220
				// TODO: use transformation filter?
				ret := []*envoycore.HeaderValueOption{
					header(":path", path),
					header(":authority", hostname),
				}

				return ret, nil
			}
		}
		return nil, errors.Errorf("unknown function %v", functionName)
	})
}

func header(k, v string) *envoycore.HeaderValueOption {
	return &envoycore.HeaderValueOption{
		Header: &envoycore.HeaderValue{
			Key:   k,
			Value: v,
		},
		Append: &types.BoolValue{Value: false},
	}
}

func getPath(functionSpec *azure.UpstreamSpec_FunctionSpec, apiKeys map[string]string) (string, error) {
	functionName := functionSpec.FunctionName

	pathParameters, err := getPathParameters(functionSpec, apiKeys)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/api/%s%s", functionName, pathParameters), nil
}

func getPathParameters(functionSpec *azure.UpstreamSpec_FunctionSpec, apiKeys map[string]string) (string, error) {
	var keyNames []string
	switch functionSpec.AuthLevel {
	case "anonymous":
		return "", nil
	case "function":
		// TODO(talnordan): Consider whether using the "function" authentication level should require
		// using a function key and not a master key. This is a product decision. From the technical
		// point of view, a master key does satisfy the "function" authentication level.
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

func getApiKey(apiKeys map[string]string, keyNames []string) (string, error) {
	for _, keyName := range keyNames {
		apiKey, ok := apiKeys[keyName]
		if ok && apiKey != "" {
			return apiKey, nil
		}
	}

	return "", fmt.Errorf("secret not found for key names %v", keyNames)
}
