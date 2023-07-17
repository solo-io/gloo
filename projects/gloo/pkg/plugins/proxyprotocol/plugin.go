package proxyprotocol

import (
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_proxy_protocol "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/proxy_protocol/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	"github.com/golang/protobuf/proto"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/proxyprotocol"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

var (
	_ plugins.Plugin         = new(plugin)
	_ plugins.ListenerPlugin = new(plugin)
)

const (
	// FilterName to pass to envoy. This was originally a solo specific version
	// now that our change are accepted up stream use the upstream name
	FilterName = wellknown.ProxyProtocol
)

type plugin struct{}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return proxyprotocol.ExtensionName
}

func (p *plugin) Init(_ plugins.InitParams) {
}

func (p *plugin) ProcessListener(params plugins.Params, in *v1.Listener, out *envoy_config_listener_v3.Listener) error {
	if !in.GetUseProxyProto().GetValue() && in.GetOptions().GetProxyProtocol() == nil {
		// If UseProxyProto is not defined on the listener, or it is false, do not append the filter
		return nil
	}

	var listenerFilter *envoy_config_listener_v3.ListenerFilter
	var err error
	// only use our custom filter if we use enterprise-only features to limit blast radius if there's an issue
	// with our proxy protocol filter
	if proxyprotocol.UsesEnterpriseOnlyFeatures(in) {
		listenerFilter, err = GenerateProxyProtocolListenerFilter(in)
		if err != nil {
			return err
		}
	} else {
		listenerFilter, err = proxyprotocol.GenerateProxyProtocolListenerFilter(in)
		if err != nil {
			return err
		}
	}

	// consider refactor on main to recreate a ListenerFilterPlugin and use that instead
	// to add filter stages for listener filters to capture this relationship there.
	// not doing right now since this would cause merge conflicts with ongoing work to support
	// a hybrid gateway in gloo edge.
	//
	// For the pr with context around the hybrid gateway and potential merge conflicts, see
	// https://github.com/solo-io/gloo/pull/5585/files#diff-3acbfb86193ba2ede91a589285a4dcaa27cf59279a271ce188afc08c50b55d2dL108-L111
	addedFilter := false
	for idx, lf := range out.ListenerFilters {
		if lf.GetName() == wellknown.TlsInspector {
			// we should put the filter here, it must come before the tls inspector
			tlsInspector := proto.Clone(out.ListenerFilters[idx]).(*envoy_config_listener_v3.ListenerFilter)
			out.ListenerFilters[idx] = listenerFilter
			out.ListenerFilters = append(out.GetListenerFilters(), tlsInspector)
			addedFilter = true
			break
		}
	}
	if !addedFilter {
		out.ListenerFilters = append(out.GetListenerFilters(), listenerFilter)
	}
	return nil
}

// GenerateProxyProtocolListenerFilter generates a proxy protocol listener filter
// This is not to be confused with the upstream setting as this unwraps proxy protocol
// rather than setting it on the request upstream which is done in another location.
func GenerateProxyProtocolListenerFilter(in *v1.Listener) (*envoy_config_listener_v3.ListenerFilter, error) {

	proxyProtocol := envoy_proxy_protocol.ProxyProtocol{}
	if pp := in.GetOptions().GetProxyProtocol(); pp != nil {
		proxyProtocol.AllowRequestsWithoutProxyProtocol = pp.AllowRequestsWithoutProxyProtocol
		var rules []*envoy_proxy_protocol.ProxyProtocol_Rule
		for _, rule := range pp.GetRules() {

			var tlvPresent *envoy_proxy_protocol.ProxyProtocol_KeyValuePair
			if rule.OnTlvPresent != nil {
				tlvPresent = &envoy_proxy_protocol.ProxyProtocol_KeyValuePair{
					MetadataNamespace: rule.GetOnTlvPresent().GetMetadataNamespace(),
					Key:               rule.GetOnTlvPresent().GetKey(),
				}
			}

			rules = append(rules, &envoy_proxy_protocol.ProxyProtocol_Rule{
				TlvType:      rule.TlvType,
				OnTlvPresent: tlvPresent,
			})
		}
		proxyProtocol.Rules = rules
	}

	msg, err := utils.MessageToAny(&proxyProtocol)
	if err != nil {
		return nil, err
	}

	return &envoy_config_listener_v3.ListenerFilter{
		Name: FilterName,
		ConfigType: &envoy_config_listener_v3.ListenerFilter_TypedConfig{
			TypedConfig: msg,
		},
	}, nil
}
