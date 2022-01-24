package proxyprotocol

import (
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/proto"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/proxy_protocol"
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
	FilterName = "io.solo.envoy.filters.listener.proxy_protocol"
)

type plugin struct{}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return proxyprotocol.ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) error {
	return nil
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
		listenerFilter, err = GenerateCustomProxyProtocolListenerFilter(in)
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

func GenerateCustomProxyProtocolListenerFilter(in *v1.Listener) (*envoy_config_listener_v3.ListenerFilter, error) {
	customProxyProtocol := proxy_protocol.CustomProxyProtocol{}
	if pp := in.GetOptions().GetProxyProtocol(); pp != nil {
		customProxyProtocol.AllowRequestsWithoutProxyProtocol = pp.AllowRequestsWithoutProxyProtocol
		var rules []*proxy_protocol.CustomProxyProtocol_Rule
		for _, rule := range pp.GetRules() {

			var tlvPresent *proxy_protocol.CustomProxyProtocol_KeyValuePair
			if rule.OnTlvPresent != nil {
				tlvPresent = &proxy_protocol.CustomProxyProtocol_KeyValuePair{
					MetadataNamespace: rule.GetOnTlvPresent().GetMetadataNamespace(),
					Key:               rule.GetOnTlvPresent().GetKey(),
				}
			}

			rules = append(rules, &proxy_protocol.CustomProxyProtocol_Rule{
				TlvType:      rule.TlvType,
				OnTlvPresent: tlvPresent,
			})
		}
		customProxyProtocol.Rules = rules
	}

	msg, err := utils.MessageToAny(&customProxyProtocol)
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
