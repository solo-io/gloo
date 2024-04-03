package translator

import (
	"fmt"
	accesslogv3 "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/route/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"
	"reflect"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	validationapi "github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
	"github.com/solo-io/go-utils/contextutils"
)

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

	extFilterChains := l.filterChainTranslator.ComputeFilterChains(params)

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

	accessLogs := []*accesslogv3.AccessLog{}
	for _, al := range l.listener.Options.AccessLoggingService.GetAccessLog() {
		accessLog := &accesslogv3.AccessLog{
			Filter: convertAccessLogFilter(al.GetFilter()),
		}

		accessLogs = append(accessLogs, accessLog)
	}

	// This is upstream envoy definition we cannot mutate this struct
	out := &envoy_config_listener_v3.Listener{
		Name:         l.listener.GetName(),
		Address:      l.computeListenerAddress(),
		FilterChains: filterChains,
		AccessLog:    accessLogs,
	}

	for _, plug := range l.plugins {
		filterConverterPlug, ok := plug.(plugins.FilterChainMutatorPlugin)
		if !ok {
			continue
		}
		if err := filterConverterPlug.ProcessFilterChain(params, l.listener, cleanedExtendedChains, out); err != nil {
			validation.AppendListenerError(l.report,
				validationapi.ListenerReport_Error_ProcessingError,
				err.Error())
		}
	}

	CheckForFilterChainConsistency(out.GetFilterChains(), l.report, out)

	// run the Listener Plugins
	for _, listenerPlugin := range l.plugins {
		// Need to have the deprecated cipher information still available at this point in time
		if err := listenerPlugin.ProcessListener(params, l.listener, out); err != nil {
			validation.AppendListenerError(l.report,
				validationapi.ListenerReport_Error_ProcessingError,
				err.Error())
		}
	}

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

func convertAccessLogFilter(alf *als.AccessLogFilter) *accesslogv3.AccessLogFilter {
	filter := &accesslogv3.AccessLogFilter{}

	switch filterSpecifier := alf.GetFilterSpecifier().(type) {
	case *als.AccessLogFilter_StatusCodeFilter:
		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_StatusCodeFilter{
			StatusCodeFilter: &accesslogv3.StatusCodeFilter{
				Comparison: &accesslogv3.ComparisonFilter{
					Op: accesslogv3.ComparisonFilter_Op(filterSpecifier.StatusCodeFilter.GetComparison().GetOp()),
					Value: &envoy_config_core_v3.RuntimeUInt32{
						DefaultValue: filterSpecifier.StatusCodeFilter.GetComparison().GetValue().GetDefaultValue(),
						RuntimeKey:   filterSpecifier.StatusCodeFilter.GetComparison().GetValue().GetRuntimeKey(),
					},
				},
			},
		}
	case *als.AccessLogFilter_DurationFilter:
		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_DurationFilter{
			DurationFilter: &accesslogv3.DurationFilter{
				Comparison: &accesslogv3.ComparisonFilter{
					Op: accesslogv3.ComparisonFilter_Op(filterSpecifier.DurationFilter.GetComparison().GetOp()),
					Value: &envoy_config_core_v3.RuntimeUInt32{
						DefaultValue: filterSpecifier.DurationFilter.GetComparison().GetValue().GetDefaultValue(),
						RuntimeKey:   filterSpecifier.DurationFilter.GetComparison().GetValue().GetRuntimeKey(),
					},
				},
			},
		}
	case *als.AccessLogFilter_NotHealthCheckFilter:
		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_NotHealthCheckFilter{
			NotHealthCheckFilter: &accesslogv3.NotHealthCheckFilter{},
		}
	case *als.AccessLogFilter_TraceableFilter:
		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_TraceableFilter{
			TraceableFilter: &accesslogv3.TraceableFilter{},
		}
	case *als.AccessLogFilter_RuntimeFilter:
		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_RuntimeFilter{
			RuntimeFilter: &accesslogv3.RuntimeFilter{
				RuntimeKey: filterSpecifier.RuntimeFilter.GetRuntimeKey(),
				PercentSampled: &typev3.FractionalPercent{
					Numerator:   filterSpecifier.RuntimeFilter.GetPercentSampled().GetNumerator(),
					Denominator: typev3.FractionalPercent_DenominatorType(filterSpecifier.RuntimeFilter.GetPercentSampled().GetDenominator()),
				},
				UseIndependentRandomness: filterSpecifier.RuntimeFilter.GetUseIndependentRandomness(),
			},
		}
	case *als.AccessLogFilter_AndFilter:
		var andAlfs []*accesslogv3.AccessLogFilter

		for _, andAlf := range filterSpecifier.AndFilter.GetFilters() {
			andAlfs = append(andAlfs, convertAccessLogFilter(andAlf))
		}

		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_AndFilter{
			AndFilter: &accesslogv3.AndFilter{
				Filters: andAlfs,
			},
		}
	case *als.AccessLogFilter_OrFilter:
		var orAlfs []*accesslogv3.AccessLogFilter

		for _, orAlf := range filterSpecifier.OrFilter.GetFilters() {
			orAlfs = append(orAlfs, convertAccessLogFilter(orAlf))
		}

		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_OrFilter{
			OrFilter: &accesslogv3.OrFilter{
				Filters: orAlfs,
			},
		}
	case *als.AccessLogFilter_HeaderFilter:
		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_HeaderFilter{
			HeaderFilter: &accesslogv3.HeaderFilter{
				Header: &routev3.HeaderMatcher{
					Name:        filterSpecifier.HeaderFilter.Header.GetName(),
					InvertMatch: filterSpecifier.HeaderFilter.Header.GetInvertMatch(),
				},
			},
		}

		switch hms := filterSpecifier.HeaderFilter.Header.GetHeaderMatchSpecifier().(type) {
		case *v3.HeaderMatcher_ExactMatch:
			filter.GetFilterSpecifier().(*accesslogv3.AccessLogFilter_HeaderFilter).HeaderFilter.Header.HeaderMatchSpecifier = &routev3.HeaderMatcher_ExactMatch{
				ExactMatch: hms.ExactMatch,
			}
		case *v3.HeaderMatcher_SafeRegexMatch:
			filter.GetFilterSpecifier().(*accesslogv3.AccessLogFilter_HeaderFilter).HeaderFilter.Header.HeaderMatchSpecifier = &routev3.HeaderMatcher_SafeRegexMatch{
				SafeRegexMatch: &matcherv3.RegexMatcher{
					EngineType: &matcherv3.RegexMatcher_GoogleRe2{
						GoogleRe2: &matcherv3.RegexMatcher_GoogleRE2{
							MaxProgramSize: hms.SafeRegexMatch.GetGoogleRe2().GetMaxProgramSize(),
						},
					},
					Regex: hms.SafeRegexMatch.GetRegex(),
				},
			}
		case *v3.HeaderMatcher_RangeMatch:
			filter.GetFilterSpecifier().(*accesslogv3.AccessLogFilter_HeaderFilter).HeaderFilter.Header.HeaderMatchSpecifier = &routev3.HeaderMatcher_RangeMatch{
				RangeMatch: &typev3.Int64Range{
					Start: hms.RangeMatch.GetStart(),
					End:   hms.RangeMatch.GetEnd(),
				},
			}
		case *v3.HeaderMatcher_PresentMatch:
			filter.GetFilterSpecifier().(*accesslogv3.AccessLogFilter_HeaderFilter).HeaderFilter.Header.HeaderMatchSpecifier = &routev3.HeaderMatcher_PresentMatch{
				PresentMatch: hms.PresentMatch,
			}
		case *v3.HeaderMatcher_PrefixMatch:
			filter.GetFilterSpecifier().(*accesslogv3.AccessLogFilter_HeaderFilter).HeaderFilter.Header.HeaderMatchSpecifier = &routev3.HeaderMatcher_PrefixMatch{
				PrefixMatch: hms.PrefixMatch,
			}
		case *v3.HeaderMatcher_SuffixMatch:
			filter.GetFilterSpecifier().(*accesslogv3.AccessLogFilter_HeaderFilter).HeaderFilter.Header.HeaderMatchSpecifier = &routev3.HeaderMatcher_SuffixMatch{
				SuffixMatch: hms.SuffixMatch,
			}
		}
	case *als.AccessLogFilter_ResponseFlagFilter:
		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_ResponseFlagFilter{
			ResponseFlagFilter: &accesslogv3.ResponseFlagFilter{
				Flags: filterSpecifier.ResponseFlagFilter.GetFlags(),
			},
		}
	case *als.AccessLogFilter_GrpcStatusFilter:
		var statuses []accesslogv3.GrpcStatusFilter_Status
		for _, inStatus := range filterSpecifier.GrpcStatusFilter.GetStatuses() {
			statuses = append(statuses, accesslogv3.GrpcStatusFilter_Status(inStatus))
		}

		filter.FilterSpecifier = &accesslogv3.AccessLogFilter_GrpcStatusFilter{
			GrpcStatusFilter: &accesslogv3.GrpcStatusFilter{
				Statuses: statuses,
				Exclude:  filterSpecifier.GrpcStatusFilter.GetExclude(),
			},
		}
	}

	return filter
}