package translator

import (
	"fmt"
	"net"
	"strings"

	envoyal "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	errors "github.com/rotisserie/eris"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// UpstreamToClusterName converts an Upstream ref to a cluster name to be used in envoy.
func UpstreamToClusterName(upstream *core.ResourceRef) string {

	// For non-namespaced resources, return only name
	if upstream.GetNamespace() == "" {
		return upstream.GetName()
	}

	// Don't use dots in the name as it messes up prometheus stats
	return fmt.Sprintf("%s_%s", upstream.GetName(), upstream.GetNamespace())
}

// UpstreamToClusterStatsName gets a cluster stats name from which key info about the upstream can be parsed out.
// Currently only kube-type upstreams are handled, and the rest will fall back to using the cluster name.
func UpstreamToClusterStatsName(upstream *gloov1.Upstream) string {
	statsName := UpstreamToClusterName(upstream.GetMetadata().Ref())

	switch upstreamType := upstream.GetUpstreamType().(type) {
	case *v1.Upstream_Kube:
		upstreamName := upstream.GetMetadata().GetName()
		// Add an identifying prefix if it's a "real" upstream (fake upstreams already have such a prefix).
		if !kubernetes.IsFakeKubeUpstream(upstreamName) {
			upstreamName = fmt.Sprintf("%s%s", kubernetes.KubeUpstreamStatsPrefix, upstreamName)
		}
		statsName = fmt.Sprintf("%s_%s_%s_%s_%v",
			upstreamName,
			upstream.GetMetadata().GetNamespace(),
			upstreamType.Kube.GetServiceNamespace(),
			upstreamType.Kube.GetServiceName(),
			upstreamType.Kube.GetServicePort(),
		)
	}

	// although envoy will do this before emitting stats, we sanitize here because some of our tests call
	// this function to get expected stats names
	return strings.ReplaceAll(statsName, ":", "_")
}

// ClusterToUpstreamRef converts an envoy cluster name back to an upstream ref.
// (this is currently only used in the tunneling plugin and the old UI)
// This does the inverse of UpstreamToClusterName
func ClusterToUpstreamRef(cluster string) (*core.ResourceRef, error) {

	split := strings.Split(cluster, "_")
	if len(split) > 2 || len(split) < 1 {
		return nil, errors.Errorf("unable to convert cluster %s back to upstream ref", cluster)
	}

	ref := &core.ResourceRef{
		Name: split[0],
	}

	if len(split) == 2 {
		ref.Namespace = split[1]
	}
	return ref, nil
}

func NewFilterWithTypedConfig(name string, config proto.Message) (*envoy_config_listener_v3.Filter, error) {

	s := &envoy_config_listener_v3.Filter{
		Name: name,
	}

	if config != nil {
		marshalledConf, err := utils.MessageToAny(config)
		if err != nil {
			// this should NEVER HAPPEN!
			return &envoy_config_listener_v3.Filter{}, err
		}

		s.ConfigType = &envoy_config_listener_v3.Filter_TypedConfig{
			TypedConfig: marshalledConf,
		}
	}

	return s, nil
}

func NewAccessLogWithConfig(name string, config proto.Message) (envoyal.AccessLog, error) {
	s := envoyal.AccessLog{
		Name: name,
	}

	if config != nil {
		marshalledConf, err := utils.MessageToAny(config)
		if err != nil {
			// this should NEVER HAPPEN!
			return envoyal.AccessLog{}, err
		}

		s.ConfigType = &envoyal.AccessLog_TypedConfig{
			TypedConfig: marshalledConf,
		}
	}

	return s, nil
}

func ParseTypedConfig(c typedConfigObject, config proto.Message) error {
	anyPb := c.GetTypedConfig()
	if anyPb != nil {
		return ptypes.UnmarshalAny(anyPb, config)
	}
	return nil
}

type typedConfigObject interface {
	GetTypedConfig() *any.Any
}

// IsIpv4Address returns whether
// the provided address is valid IPv4, is pure(unmapped) IPv4, and if there was an error in the bindaddr
// This is used to distinguish between IPv4 and IPv6 addresses
func IsIpv4Address(bindAddress string) (validIpv4, strictIPv4 bool, err error) {
	bindIP := net.ParseIP(bindAddress)
	if bindIP == nil {
		// If bindAddress is not a valid textual representation of an IP address
		return false, false, errors.Errorf("bindAddress %s is not a valid IP address", bindAddress)

	} else if bindIP.To4() == nil {
		// If bindIP is not an IPv4 address, To4 returns nil.
		// so this is not an acceptable ipv4
		return false, false, nil
	}
	return true, isPureIPv4Address(bindAddress), nil

}

// isPureIPv4Address checks the string to see if it is
// ipv4 and not ipv4 mapped into ipv6 space and not ipv6.
// Used as the standard net.Parse smashes everything to ipv6.
// Basically false if ::ffff:0.0.0.0 and true if 0.0.0.0
func isPureIPv4Address(ipString string) bool {
	for i := range len(ipString) {
		switch ipString[i] {
		case '.':
			return true
		case ':':
			return false
		}
	}
	return false
}
