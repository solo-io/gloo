package snapshot

import (
	"fmt"

	gateway_solo_io "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	github_com_solo_io_gloo_projects_gloo_pkg_api_external_solo_ratelimit "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	gloo_solo_io "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	enterprise_gloo_solo_io "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	graphql_gloo_solo_io "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"istio.io/istio/pkg/kube/krt"
	"k8s.io/apimachinery/pkg/types"
)

type Named interface {
	GetName() string
	GetNamespace() string
}
type GlooCollection[T any] interface {
	Find(namespace string, name string) (T, error)
	List() []T
	IsFindEfficient() bool
}

type Md interface {
	GetMetadata() *core.Metadata
}

type ResourceWrapper[T Md] struct {
	Resource T
}

func (t ResourceWrapper[T]) GetName() string {
	return t.Resource.GetMetadata().GetName()
}
func (t ResourceWrapper[T]) GetNamespace() string {
	return t.Resource.GetMetadata().GetNamespace()
}

type SliceCollection[T Md] []T

func (t SliceCollection[T]) Find(namespace string, name string) (T, error) {
	for _, t := range t {
		if t.GetMetadata().GetName() == name && t.GetMetadata().GetNamespace() == namespace {
			return t, nil
		}
	}
	var zero T
	return zero, fmt.Errorf("list did not find %T %v.%v", zero, namespace, name)
}

func (t SliceCollection[T]) List() []T {
	return t
}

func (t SliceCollection[T]) IsFindEfficient() bool {
	return false
}

type KrtCollection[T any] struct {
	C    krt.Collection[T]
	Kctx krt.HandlerContext
}

func (t KrtCollection[T]) Find(namespace string, name string) (T, error) {
	ret := krt.Fetch(t.Kctx, t.C, krt.FilterObjectName(types.NamespacedName{Name: name, Namespace: namespace}))
	if len(ret) != 1 {
		var zero T
		return zero, fmt.Errorf("list did not find %T %v.%v", zero, namespace, name)
	}
	return ret[0], nil
}

func (t KrtCollection[T]) List() []T {
	return krt.Fetch(t.Kctx, t.C)
}

func (t KrtCollection[T]) IsFindEfficient() bool {
	return true
}

func FromApiSnapshot(snapshot *v1snap.ApiSnapshot) *Snapshot {
	return &Snapshot{
		Artifacts:          SliceCollection[*gloo_solo_io.Artifact](snapshot.Artifacts),
		Endpoints:          SliceCollection[*gloo_solo_io.Endpoint](snapshot.Endpoints),
		Proxies:            SliceCollection[*gloo_solo_io.Proxy](snapshot.Proxies),
		UpstreamGroups:     SliceCollection[*gloo_solo_io.UpstreamGroup](snapshot.UpstreamGroups),
		Secrets:            SliceCollection[*gloo_solo_io.Secret](snapshot.Secrets),
		Upstreams:          SliceCollection[*gloo_solo_io.Upstream](snapshot.Upstreams),
		AuthConfigs:        SliceCollection[*enterprise_gloo_solo_io.AuthConfig](snapshot.AuthConfigs),
		Ratelimitconfigs:   SliceCollection[*github_com_solo_io_gloo_projects_gloo_pkg_api_external_solo_ratelimit.RateLimitConfig](snapshot.Ratelimitconfigs),
		VirtualServices:    SliceCollection[*gateway_solo_io.VirtualService](snapshot.VirtualServices),
		RouteTables:        SliceCollection[*gateway_solo_io.RouteTable](snapshot.RouteTables),
		Gateways:           SliceCollection[*gateway_solo_io.Gateway](snapshot.Gateways),
		VirtualHostOptions: SliceCollection[*gateway_solo_io.VirtualHostOption](snapshot.VirtualHostOptions),
		RouteOptions:       SliceCollection[*gateway_solo_io.RouteOption](snapshot.RouteOptions),
		HttpGateways:       SliceCollection[*gateway_solo_io.MatchableHttpGateway](snapshot.HttpGateways),
		TcpGateways:        SliceCollection[*gateway_solo_io.MatchableTcpGateway](snapshot.TcpGateways),
		GraphqlApis:        SliceCollection[*graphql_gloo_solo_io.GraphQLApi](snapshot.GraphqlApis),
	}
}

type ArtifactList = GlooCollection[*gloo_solo_io.Artifact]
type EndpointList = GlooCollection[*gloo_solo_io.Endpoint]
type ProxieList = GlooCollection[*gloo_solo_io.Proxy]
type UpstreamGroupList = GlooCollection[*gloo_solo_io.UpstreamGroup]
type SecretList = GlooCollection[*gloo_solo_io.Secret]
type UpstreamList = GlooCollection[*gloo_solo_io.Upstream]
type AuthConfigList = GlooCollection[*enterprise_gloo_solo_io.AuthConfig]
type RatelimitconfigList = GlooCollection[*github_com_solo_io_gloo_projects_gloo_pkg_api_external_solo_ratelimit.RateLimitConfig]
type VirtualServiceList = GlooCollection[*gateway_solo_io.VirtualService]
type RouteTableList = GlooCollection[*gateway_solo_io.RouteTable]
type GatewayList = GlooCollection[*gateway_solo_io.Gateway]
type VirtualHostOptionList = GlooCollection[*gateway_solo_io.VirtualHostOption]
type RouteOptionList = GlooCollection[*gateway_solo_io.RouteOption]
type HttpGatewayList = GlooCollection[*gateway_solo_io.MatchableHttpGateway]
type TcpGatewayList = GlooCollection[*gateway_solo_io.MatchableTcpGateway]
type GraphqlApiList = GlooCollection[*graphql_gloo_solo_io.GraphQLApi]

type Snapshot struct {
	Artifacts          GlooCollection[*gloo_solo_io.Artifact]
	Endpoints          GlooCollection[*gloo_solo_io.Endpoint]
	Proxies            GlooCollection[*gloo_solo_io.Proxy]
	UpstreamGroups     GlooCollection[*gloo_solo_io.UpstreamGroup]
	Secrets            GlooCollection[*gloo_solo_io.Secret]
	Upstreams          GlooCollection[*gloo_solo_io.Upstream]
	AuthConfigs        GlooCollection[*enterprise_gloo_solo_io.AuthConfig]
	Ratelimitconfigs   GlooCollection[*github_com_solo_io_gloo_projects_gloo_pkg_api_external_solo_ratelimit.RateLimitConfig]
	VirtualServices    GlooCollection[*gateway_solo_io.VirtualService]
	RouteTables        GlooCollection[*gateway_solo_io.RouteTable]
	Gateways           GlooCollection[*gateway_solo_io.Gateway]
	VirtualHostOptions GlooCollection[*gateway_solo_io.VirtualHostOption]
	RouteOptions       GlooCollection[*gateway_solo_io.RouteOption]
	HttpGateways       GlooCollection[*gateway_solo_io.MatchableHttpGateway]
	TcpGateways        GlooCollection[*gateway_solo_io.MatchableTcpGateway]
	GraphqlApis        GlooCollection[*graphql_gloo_solo_io.GraphQLApi]
}
