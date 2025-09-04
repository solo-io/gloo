package translator

import (
	"fmt"
	"reflect"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/projects/gloo/constants"
	validationapi "github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
)

// logComputeListener is a helper function that logs listener messages only when COMPUTE_LISTENER_LOGS is enabled
func logComputeListener(logger *zap.SugaredLogger, msg string, keysAndValues ...interface{}) {
	if envutils.IsEnvTruthy(constants.ComputeListenerLogsEnv) {
		// Add the issue label to all gated logs
		keysAndValues = append([]interface{}{"issue", "8539"}, keysAndValues...)
		logger.Infow(msg, keysAndValues...)
	}
}

type ListenerTranslator interface {
	// A single Gloo Listener produces a single Envoy listener
	// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/listeners/listeners#arch-overview-listeners
	ComputeListener(params plugins.Params) *envoy_config_listener_v3.Listener
}

var _ ListenerTranslator = new(emptyListenerTranslator)
var _ ListenerTranslator = new(listenerTranslatorInstance)

type emptyListenerTranslator struct {
}

func (e *emptyListenerTranslator) ComputeListener(params plugins.Params) *envoy_config_listener_v3.Listener {
	return nil
}

type listenerTranslatorInstance struct {
	plugins               []plugins.ListenerPlugin
	listener              *v1.Listener
	report                *validationapi.ListenerReport
	filterChainTranslator FilterChainTranslator
}

func (l *listenerTranslatorInstance) ComputeListener(params plugins.Params) *envoy_config_listener_v3.Listener {
	params.Ctx = contextutils.WithLogger(params.Ctx, "compute_listener."+l.listener.GetName())
	logger := contextutils.LoggerFrom(params.Ctx)

	logComputeListener(logger, "Starting listener computation",
		"listener_name", l.listener.GetName(),
		"bind_address", l.listener.GetBindAddress(),
		"bind_port", l.listener.GetBindPort(),
		"listener_type", fmt.Sprintf("%T", l.listener.GetListenerType()))

	extFilterChains := l.filterChainTranslator.ComputeFilterChains(params)

	logComputeListener(logger, "Computed filter chains",
		"listener_name", l.listener.GetName(),
		"extended_filter_chain_count", len(extFilterChains))

	// our format for setting up filterchains to be extended has a non-zero chance
	// of having nil filterchains which makes follow up  logic more complex
	// as we are already iterating over the list we add this as well.
	cleanedExtendedChains := make([]*plugins.ExtendedFilterChain, 0, len(extFilterChains))
	// unwrap all filterChains before putting into envoy listener
	filterChains := make([]*envoy_config_listener_v3.FilterChain, 0, len(extFilterChains))
	for _, extFilterChain := range extFilterChains {
		extFilterChain := extFilterChain
		if extFilterChain != nil && extFilterChain.FilterChain != nil {
			filterChains = append(filterChains, extFilterChain.FilterChain)
			cleanedExtendedChains = append(cleanedExtendedChains, extFilterChain)
		}
	}

	logComputeListener(logger, "Processed filter chains",
		"listener_name", l.listener.GetName(),
		"valid_filter_chain_count", len(filterChains),
		"cleaned_extended_chain_count", len(cleanedExtendedChains))

	// This is upstream envoy definition we cannot mutate this struct
	out := &envoy_config_listener_v3.Listener{
		Name:         l.listener.GetName(),
		Address:      l.computeListenerAddress(),
		FilterChains: filterChains,
	}

	logComputeListener(logger, "Processing filter chain mutator plugins",
		"listener_name", l.listener.GetName(),
		"plugin_count", len(l.plugins))

	for _, plug := range l.plugins {
		filterConverterPlug, ok := plug.(plugins.FilterChainMutatorPlugin)
		if !ok {
			continue
		}
		logComputeListener(logger, "Processing filter chain with plugin",
			"listener_name", l.listener.GetName(),
			"plugin_name", plug.Name())

		if err := filterConverterPlug.ProcessFilterChain(params, l.listener, cleanedExtendedChains, out); err != nil {
			logComputeListener(logger, "Error processing filter chain with plugin",
				"listener_name", l.listener.GetName(),
				"plugin_name", plug.Name(),
				"error", err.Error())
			validation.AppendListenerError(l.report,
				validationapi.ListenerReport_Error_ProcessingError,
				err.Error())
		}
	}

	CheckForFilterChainConsistency(out.GetFilterChains(), l.report, out)

	logComputeListener(logger, "Processing listener plugins",
		"listener_name", l.listener.GetName(),
		"plugin_count", len(l.plugins))

	// run the Listener Plugins
	for _, listenerPlugin := range l.plugins {
		logComputeListener(logger, "Processing listener with plugin",
			"listener_name", l.listener.GetName(),
			"plugin_name", listenerPlugin.Name())

		// Need to have the deprecated cipher information still available at this point in time
		if err := listenerPlugin.ProcessListener(params, l.listener, out); err != nil {
			logComputeListener(logger, "Error processing listener with plugin",
				"listener_name", l.listener.GetName(),
				"plugin_name", listenerPlugin.Name(),
				"error", err.Error())
			validation.AppendListenerError(l.report,
				validationapi.ListenerReport_Error_ProcessingError,
				err.Error())
		}
	}

	logComputeListener(logger, "Completed listener computation",
		"listener_name", l.listener.GetName(),
		"envoy_listener_name", out.GetName(),
		"final_filter_chain_count", len(out.GetFilterChains()))

	return out
}

// computeListenerAddress returns the Address that this listener will listen for traffic
func (l *listenerTranslatorInstance) computeListenerAddress() *envoy_config_core_v3.Address {
	_, isIpv4Address, err := IsIpv4Address(l.listener.GetBindAddress())
	if err != nil {
		validation.AppendListenerError(l.report,
			validationapi.ListenerReport_Error_ProcessingError,
			err.Error(),
		)
	}

	return &envoy_config_core_v3.Address{
		Address: &envoy_config_core_v3.Address_SocketAddress{
			SocketAddress: &envoy_config_core_v3.SocketAddress{
				Protocol: envoy_config_core_v3.SocketAddress_TCP,
				Address:  l.listener.GetBindAddress(),
				PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
					PortValue: l.listener.GetBindPort(),
				},
				// As of Envoy 1.22: https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.22/v1.22.0.html
				// the Ipv4Compat flag can only be set on Ipv6 address and Ipv4-mapped Ipv6 address.
				// Check if this is a non-padded pure ipv4 address and unset the compat flag if so.
				Ipv4Compat: !isIpv4Address,
			},
		},
	}
}

func validateListenerPorts(proxy *v1.Proxy, listenerReport *validationapi.ListenerReport) {
	listenersByPort := make(map[uint32][]int)
	for i, listener := range proxy.GetListeners() {
		listenersByPort[listener.GetBindPort()] = append(listenersByPort[listener.GetBindPort()], i)
	}
	for port, listeners := range listenersByPort {
		if len(listeners) == 1 {
			continue
		}
		var listenerNames []string
		for _, idx := range listeners {
			listenerNames = append(listenerNames, proxy.GetListeners()[idx].GetName())
		}
		validation.AppendListenerError(listenerReport,
			validationapi.ListenerReport_Error_BindPortNotUniqueError,
			fmt.Sprintf("port %v is shared by listeners %v", port, listeners),
		)
	}
}

// CheckForFilterChainConsistency to avoid the envoy error that occurs here:
// https://github.com/envoyproxy/envoy/blob/v1.15.0/source/server/filter_chain_manager_impl.cc#L162-L166
// Note: this is NOT address non-equal but overlapping FilterChainMatches, which is a separate check here:
// https://github.com/envoyproxy/envoy/blob/50ef0945fa2c5da4bff7627c3abf41fdd3b7cffd/source/server/filter_chain_manager_impl.cc#L218-L354
// Given the complexity of the overlap detection implementation, we don't want to duplicate that behavior here.
// We may want to consider invoking envoy from a library to detect overlapping and other issues, which would build
// off this discussion: https://github.com/solo-io/gloo/issues/2114
// This also checks that if we are using matchers we have the required names on all filterchains
// Visible for testing
func CheckForFilterChainConsistency(filterChains []*envoy_config_listener_v3.FilterChain, listenerReport *validationapi.ListenerReport, out *envoy_config_listener_v3.Listener) {
	usingListenerLevelMatcher := out.GetFilterChainMatcher() != nil
	for idx1, filterChain := range filterChains {
		if usingListenerLevelMatcher && filterChain.GetName() == "" {
			// only need to validate that the filterchain has a name

			validation.AppendListenerError(listenerReport,
				validationapi.ListenerReport_Error_ProcessingError, "Tried to make a filter chain without a name ")

		}

		for idx2, otherFilterChain := range filterChains {

			// only need to compare each pair once
			if idx2 <= idx1 {
				continue
			}
			if usingListenerLevelMatcher && filterChain.GetName() == otherFilterChain.GetName() {
				validation.AppendListenerError(listenerReport,
					validationapi.ListenerReport_Error_NameNotUniqueError, fmt.Sprintf("Tried to make a filter chain with the same name as another "+
						" FilterChain {%v}", otherFilterChain.GetName()))

			} else if !usingListenerLevelMatcher && reflect.DeepEqual(filterChain.GetFilterChainMatch(), otherFilterChain.GetFilterChainMatch()) {
				validation.AppendListenerError(listenerReport,
					validationapi.ListenerReport_Error_SSLConfigError, fmt.Sprintf("Tried to apply multiple filter chains "+
						"with the same FilterChainMatch {%v}. This is usually caused by overlapping sniDomains or multiple empty sniDomains in virtual services", filterChain.GetFilterChainMatch()))
			}
		}
	}
}
