package grpc

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/common/annotations"
	"github.com/solo-io/gloo/pkg/coreplugins/common"
	"github.com/solo-io/gloo/pkg/plugin"
)

type Plugin struct{}

const (
	filterName  = "envoy.grpc_json_transcoder"
	pluginStage = plugin.PreOutAuth

	ServiceTypeGRPC = "HTTP-Functions"
)

/*
 - need to create envoy routes for the grpc service. this is how we map the functions to the cluster
   like so:
              - match: { prefix: "/bookstore.Bookstore" }
                route: { cluster: service_google }

          http_filters:
          - name: envoy.grpc_json_transcoder
            config:
              proto_descriptor: ./proto.pb
              services: [bookstore.Bookstore]
          - name: envoy.router

for every service in the file, create a d

*/

func (p *Plugin) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func isOurs(in *v1.Upstream) bool {
	if in.Metadata == nil || in.Metadata.Annotations == nil {
		return false
	}
	return in.Metadata.Annotations[annotations.ServiceType] == ServiceTypeGRPC
}

func (p *Plugin) ProcessUpstream(params *plugin.UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error {
	if !isOurs(in) {
		return nil
	}

	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	common.InitFilterMetadata(filterName, out.Metadata)
	out.Metadata.FilterMetadata[filterName] = &types.Struct{
		Fields: make(map[string]*types.Value),
	}

	return nil
}
