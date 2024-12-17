package translator

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/log"
	"google.golang.org/protobuf/types/known/wrapperspb"

	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/duration"

	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/golang/protobuf/proto"
	validationapi "github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
)

type FilterChainTranslator interface {
	ComputeFilterChains(params plugins.Params) []*plugins.ExtendedFilterChain
}

var (
	_ FilterChainTranslator = new(tcpFilterChainTranslator)
	_ FilterChainTranslator = new(httpFilterChainTranslator)
)

type tcpFilterChainTranslator struct {
	// List of TcpFilterChainPlugins to process
	tcpPlugins []plugins.TcpFilterChainPlugin
	// List of TcpFilterChainPlugins to process
	networkPlugins []plugins.NetworkFilterPlugin
	// The parent Listener, this is only used to associate errors with the parent resource
	parentListener *v1.Listener
	// The TcpListener used to generate the list of FilterChains
	listener *v1.TcpListener
	// The report used to store processing errors
	report *validationapi.TcpListenerReport

	// These values are optional (currently only available for HybridGateways)
	defaultSslConfig        *ssl.SslConfig
	prefixRanges            []*v3.CidrRange
	sourcePrefixRanges      []*v3.CidrRange
	passthroughCipherSuites []string
	destinationPort         uint32
}

// if the error is of ErrorWithKnownLevel type, extract the level; if the level
// is a warning, extract the TcpHost number and report the warning on that
// TcpHost. if the level is an error, report it on the TcpListener object
// itself. nothing currently reports errors on a TcpHost instance or warnings
// on a TcpListener instance, so these two branches are not currently supported
func (t *tcpFilterChainTranslator) reportCreateTcpFilterChainsError(err error) {
	getWarningType := func(errType error) validationapi.TcpHostReport_Warning_Type {
		switch errType.(type) {
		case *pluginutils.DestinationNotFoundError:
			return validationapi.TcpHostReport_Warning_InvalidDestinationWarning
		default:
			return validationapi.TcpHostReport_Warning_UnknownWarning
		}
	}

	reportTcpListenerError := func(errType error) {
		validation.AppendTCPListenerError(
			t.report,
			// currently only processing errors are reported
			validationapi.TcpListenerReport_Error_ProcessingError,
			fmt.Sprintf("listener %s: %s", t.parentListener.GetName(), errType.Error()))
	}

	reportError := func(errReport error) {
		var errType validation.ErrorWithKnownLevel
		switch {
		case errors.As(errReport, &errType):
			switch errType.ErrorLevel() {
			case validation.ErrorLevels_WARNING:
				if tcpHostNum := errType.GetContext().HostNum; tcpHostNum != nil {
					validation.AppendTcpHostWarning(
						t.report.GetTcpHostReports()[*tcpHostNum],
						getWarningType(errType.GetError()),
						fmt.Sprintf("listener %s: %s", t.parentListener.GetName(), errType.Error()))
					break
				}
				// hostNum was not provided, so just report as an error
				fallthrough
			case validation.ErrorLevels_ERROR:
				reportTcpListenerError(errType)
			}
			// if the error is not of ErrorWithKnownLevel type, report it
		// as an error on the TcpListener
		default:
			reportTcpListenerError(errType)
		}
	}

	if merr, ok := err.(*multierror.Error); ok {
		for _, unwrappedErr := range merr.WrappedErrors() {
			reportError(unwrappedErr)
		}
	} else {
		reportError(err)
	}
}

func (t *tcpFilterChainTranslator) ComputeFilterChains(params plugins.Params) []*plugins.ExtendedFilterChain {
	var filterChains []*envoy_config_listener_v3.FilterChain

	// 1. Generate the network filters
	networkFilters := sortNetworkFilters(t.computeNetworkFilters(params))

	// 2. Run the tcp filter chain plugins
	for _, plug := range t.tcpPlugins {
		pluginFilterChains, err := plug.CreateTcpFilterChains(params, t.parentListener, t.listener)
		if err != nil {
			t.reportCreateTcpFilterChainsError(err)
			continue
		}

		for _, pfc := range pluginFilterChains {
			pfc := pfc

			// 3. Add the network filters to each filter chain
			pfc.Filters = append(networkFilters, pfc.GetFilters()...)
			if t.defaultSslConfig != nil {
				if pfc.GetFilterChainMatch() == nil {
					pfc.FilterChainMatch = &envoy_config_listener_v3.FilterChainMatch{}
				}

				if pfc.GetFilterChainMatch().GetServerNames() == nil {
					pfc.GetFilterChainMatch().ServerNames = []string{}
				}

				if len(pfc.GetFilterChainMatch().GetServerNames()) == 0 {
					pfc.GetFilterChainMatch().ServerNames = t.defaultSslConfig.GetSniDomains()
				}
			}
			filterChains = append(filterChains, pfc)
		}
	}

	extFilterChains := make([]*plugins.ExtendedFilterChain, 0, len(filterChains))
	for _, fc := range filterChains {
		fc := fc
		extFilterChains = append(extFilterChains, &plugins.ExtendedFilterChain{FilterChain: fc, PassthroughCipherSuites: t.passthroughCipherSuites})
	}

	// 4. Apply SourcePrefixRange to FilterChainMatch, if defined
	if len(t.sourcePrefixRanges) > 0 {
		for _, fc := range extFilterChains {
			applySourcePrefixRangesToFilterChain(fc, t.sourcePrefixRanges)
		}
	}

	// 5. Apply destination PrefixRange to filter chain match,  if defined
	if len(t.prefixRanges) > 0 {
		for _, fc := range extFilterChains {
			applyDestinationPrefixRangesToFilterChain(fc, t.prefixRanges)
		}
	}

	// 6. Apply destination port to FilterChainMatch, if defined
	if t.destinationPort > 0 {
		for _, fc := range extFilterChains {
			applyDestinationPortFilterChain(fc, t.destinationPort)
		}
	}

	return extFilterChains
}

func (t *tcpFilterChainTranslator) computeNetworkFilters(params plugins.Params) []plugins.StagedNetworkFilter {
	var networkFilters []plugins.StagedNetworkFilter
	// Process the network filters.
	for _, plug := range t.networkPlugins {
		stagedFilters, err := plug.NetworkFiltersTCP(params, t.listener)
		if err != nil {
			reportTCPListenerProcessingError(params, t.report, err)
		}

		for _, nf := range stagedFilters {
			if nf.Filter == nil {
				log.Warnf("plugin %v implements NetworkFilters() but returned nil", plug.Name())
				continue
			}
			networkFilters = append(networkFilters, nf)
		}
	}
	return networkFilters
}

// An httpFilterChainTranslator configures a single set of NetworkFilters
// and then creates duplicate filter chains for each provided SslConfig.
type httpFilterChainTranslator struct {
	settings                *v1.Settings
	parentReport            *validationapi.ListenerReport
	networkFilterTranslator NetworkFilterTranslator
	sslConfigurations       []*ssl.SslConfig
	sslConfigTranslator     utils.SslConfigTranslator

	// These values are optional (currently only available for HybridListeners or AggregateListeners)
	defaultSslConfig   *ssl.SslConfig
	prefixRanges       []*v3.CidrRange
	sourcePrefixRanges []*v3.CidrRange
	destinationPort    uint32
}

func (h *httpFilterChainTranslator) ComputeFilterChains(params plugins.Params) []*plugins.ExtendedFilterChain {
	// 1. Generate all the network filters (including the HttpConnectionManager)
	networkFilters, err := h.networkFilterTranslator.ComputeNetworkFilters(params)
	if err != nil {
		contextutils.LoggerFrom(context.Background()).DPanic(err)
		return nil
	}
	if len(networkFilters) == 0 {
		return nil
	}

	// 2. Determine which, if any, SslConfigs are defined for this Listener
	sslConfigWithDefaults := h.getSslConfigurationWithDefaults()

	// 3. Create duplicate FilterChains for each unique SslConfig
	extFilters := h.createFilterChainsFromSslConfiguration(params.Snapshot, networkFilters, sslConfigWithDefaults)

	// 4. Apply SourcePrefixRange to FilterChainMatch, if defined
	if len(h.sourcePrefixRanges) > 0 {
		for _, fc := range extFilters {
			applySourcePrefixRangesToFilterChain(fc, h.sourcePrefixRanges)
		}
	}

	// 5. Apply destination PrefixRange to FilterChainMatch, if defined
	if len(h.prefixRanges) > 0 {
		for _, fc := range extFilters {
			applyDestinationPrefixRangesToFilterChain(fc, h.prefixRanges)
		}
	}

	// 6. Apply destination port to FilterChainMatch, if defined
	if h.destinationPort > 0 {
		for _, fc := range extFilters {
			applyDestinationPortFilterChain(fc, h.destinationPort)
		}
	}

	return extFilters
}

func (h *httpFilterChainTranslator) getSslConfigurationWithDefaults() []*ssl.SslConfig {
	mergedSslConfigurations := ConsolidateSslConfigurations(h.sslConfigurations)

	if h.defaultSslConfig == nil {
		return mergedSslConfigurations
	}

	// Merge each sslConfig with the default values
	var sslConfigWithDefaults []*ssl.SslConfig
	for _, ssl := range mergedSslConfigurations {
		sslConfigWithDefaults = append(sslConfigWithDefaults, MergeSslConfig(ssl, h.defaultSslConfig))
	}
	return sslConfigWithDefaults
}

func (h *httpFilterChainTranslator) createFilterChainsFromSslConfiguration(
	snap *v1snap.ApiSnapshot,
	networkFilters []*envoy_config_listener_v3.Filter,
	sslConfigurations []*ssl.SslConfig,
) []*plugins.ExtendedFilterChain {
	// if no ssl config is provided, return a single insecure filter chain
	if len(sslConfigurations) == 0 {
		return []*plugins.ExtendedFilterChain{{
			FilterChain: &envoy_config_listener_v3.FilterChain{
				Filters: networkFilters,
			},
		}}
	}

	// create a duplicate of the listener filter chain for each ssl cert we want to serve
	var secureFilterChains []*plugins.ExtendedFilterChain
	for _, sslConfig := range sslConfigurations {
		// get secrets
		downstreamTlsContext, err := h.sslConfigTranslator.ResolveDownstreamSslConfig(snap.Secrets, sslConfig)
		if err != nil {
			// if we are configured to warn on missing tls secret and we match that error, add a
			// warning instead of error to the report.
			if h.settings.GetGateway().GetValidation().GetWarnMissingTlsSecret().GetValue() &&
				errors.Is(err, utils.SslSecretNotFoundError) {
				// We add this as a warning to support eventual consistency with TLS Secret resources. In this way,
				// the Proxy producing this will not be considered Rejected, and the HTTPS Listener will still operate
				// as expected with a VirtualService in error.
				validation.AppendListenerWarning(h.parentReport, validationapi.ListenerReport_Warning_SSLConfigWarning, err.Error())
			} else {
				// If our error is any other than SslSecretNotFoundError, we assume it is due to a malformed secret or otherwise
				// irreconcilable issue.
				validation.AppendListenerError(h.parentReport, validationapi.ListenerReport_Error_SSLConfigError, err.Error())
			}
			continue
		}

		filterChain, err := newSslFilterChain(
			downstreamTlsContext,
			sslConfig.GetSniDomains(),
			networkFilters,
			sslConfig.GetTransportSocketConnectTimeout())
		if err != nil {
			validation.AppendListenerError(h.parentReport, validationapi.ListenerReport_Error_SSLConfigError, err.Error())
			continue
		}
		secureFilterChains = append(secureFilterChains, &plugins.ExtendedFilterChain{
			FilterChain:             filterChain,
			TerminatingCipherSuites: sslConfig.GetParameters().GetCipherSuites(),
		})
	}
	return secureFilterChains
}

func applySourcePrefixRangesToFilterChain(
	filterChain *plugins.ExtendedFilterChain,
	sourcePrefixRanges []*v3.CidrRange,
) {
	if filterChain == nil || len(sourcePrefixRanges) == 0 {
		// nothing to do
		return
	}

	ensureMatch(filterChain)
	filterChain.GetFilterChainMatch().SourcePrefixRanges = buildEnvoyRanges(sourcePrefixRanges)
}

func applyDestinationPrefixRangesToFilterChain(
	filterChain *plugins.ExtendedFilterChain,
	prefixRanges []*v3.CidrRange,
) {
	if filterChain == nil || len(prefixRanges) == 0 {
		// nothing to do
		return
	}

	ensureMatch(filterChain)
	filterChain.GetFilterChainMatch().PrefixRanges = buildEnvoyRanges(prefixRanges)
}

func applyDestinationPortFilterChain(filterChain *plugins.ExtendedFilterChain, destPort uint32) {
	if destPort == 0 || filterChain == nil {
		return
	}

	ensureMatch(filterChain)
	filterChain.GetFilterChainMatch().DestinationPort = &wrapperspb.UInt32Value{Value: destPort}
}

func ensureMatch(filterChain *plugins.ExtendedFilterChain) {
	if filterChain != nil && filterChain.GetFilterChainMatch() == nil {
		// create a FilterChainMatch, if necessary
		filterChain.FilterChainMatch = &envoy_config_listener_v3.FilterChainMatch{}
	}
}

func buildEnvoyRanges(prefixRanges []*v3.CidrRange) []*corev3.CidrRange {
	envoyPrefixRanges := make([]*envoy_config_core_v3.CidrRange, 0, len(prefixRanges))
	for _, pr := range prefixRanges {
		envoyPrefixRanges = append(envoyPrefixRanges, &envoy_config_core_v3.CidrRange{
			AddressPrefix: pr.GetAddressPrefix(),
			PrefixLen:     pr.GetPrefixLen(),
		})
	}
	return envoyPrefixRanges
}

func newSslFilterChain(
	downstreamTlsContext *envoyauth.DownstreamTlsContext,
	sniDomains []string,
	listenerFilters []*envoy_config_listener_v3.Filter,
	timeout *duration.Duration,
) (*envoy_config_listener_v3.FilterChain, error) {
	// copy listenerFilter so we can modify filter chain later without changing the filters on all of them!
	clonedListenerFilters := cloneListenerFilters(listenerFilters)
	typedConfig, err := utils.MessageToAny(downstreamTlsContext)
	if err != nil {
		return nil, err
	}
	return &envoy_config_listener_v3.FilterChain{
		FilterChainMatch: &envoy_config_listener_v3.FilterChainMatch{
			ServerNames: sniDomains,
		},
		Filters: clonedListenerFilters,
		TransportSocket: &envoy_config_core_v3.TransportSocket{
			Name:       wellknown.TransportSocketTls,
			ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{TypedConfig: typedConfig},
		},
		TransportSocketConnectTimeout: timeout,
	}, nil
}

func cloneListenerFilters(originalListenerFilters []*envoy_config_listener_v3.Filter) []*envoy_config_listener_v3.Filter {
	clonedListenerFilters := make([]*envoy_config_listener_v3.Filter, len(originalListenerFilters))
	for i, lf := range originalListenerFilters {
		clonedListenerFilters[i] = proto.Clone(lf).(*envoy_config_listener_v3.Filter)
	}

	return clonedListenerFilters
}

type multiFilterChainTranslator struct {
	translators []FilterChainTranslator
}

func (m *multiFilterChainTranslator) ComputeFilterChains(params plugins.Params) []*plugins.ExtendedFilterChain {
	var outFilterChains []*plugins.ExtendedFilterChain

	for _, translator := range m.translators {
		newFilterChains := translator.ComputeFilterChains(params)
		if newFilterChains != nil {
			outFilterChains = append(outFilterChains, newFilterChains...)
		}
	}

	return outFilterChains
}
