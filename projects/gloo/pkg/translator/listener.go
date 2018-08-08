package translator

import (
	"sort"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoylistener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	gogo_types "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
)

func (t *translator) computeListener(proxy *v1.Proxy, listener *v1.Listener, snap *v1.Snapshot, report reportErr) *envoyapi.Listener {
	validateListenerPorts(proxy, report)


	listenerFilters := t.computeListenerFilters(listener, snap, report)

	filterChains := computeListenerFilterChains(snap, listener, listenerFilters)

	out := &envoyapi.Listener{
		Name: listener.Name,
		Address: envoycore.Address{
			Address: &envoycore.Address_SocketAddress{
				SocketAddress: &envoycore.SocketAddress{
					Protocol: envoycore.TCP,
					Address:  listener.BindAddress,
					PortSpecifier: &envoycore.SocketAddress_PortValue{
						PortValue: listener.BindPort,
					},
					Ipv4Compat: true,
				},
			},
		},
		FilterChains: filterChains,
	}

	// run the Listener Plugins
	params := plugins.Params{
		Snapshot: snap,
	}
	for _, plug := range t.plugins {
		routePlugin, ok := plug.(plugins.ListenerPlugin)
		if !ok {
			continue
		}
		if err := routePlugin.ProcessListener(params, listener, out); err != nil {
			report(err, "plugin error on listener")
		}
	}

	return out
}

// create a duplicate of the listener filter chain for each ssl cert we want to serve
// plus one if there's an insecure one
func computeListenerFilterChains(snap *v1.Snapshot, listener *v1.Listener, listenerFilters []envoylistener.Filter) []envoylistener.FilterChain {
	// no filters = no filter chains
	// TODO(ilackarms): find another way to prevent the xds server from serving listeners with 0 filters
	// currently the translator does not add listeners with 0 filters to the xds snapshot
	if len(listenerFilters) == 0 {
		return nil
	}

	// if there are any insecure virtual services, need an insecure filter chain
	// that will match for them
	var addInsecureFilterChain bool
	var filterChains []envoylistener.FilterChain
	for _, virtualService := range snap.Cfg.VirtualServices {
		if virtualService.SslConfig == nil || virtualService.SslConfig.SslSecrets == nil {
			addInsecureFilterChain = true
			continue
		}
		var filterChain envoylistener.FilterChain
		// TODO(ilackarms): un-copypaste this bit
		switch sslSecrets := virtualService.SslConfig.SslSecrets.(type) {
		case *v1.SSLConfig_SecretRef:
			ref := sslSecrets.SecretRef
			certChain, privateKey, rootCa, err := getSslSecrets(ref, snap.Secrets)
			if err != nil {
				log.Warnf("skipping ssl vService with invalid secrets: %v", virtualService.Name)
				continue
			}
			domains := virtualService.SslConfig.SniDomains
			if len(domains) == 0 {
				domains = virtualService.Domains
			}
			filterChain = newSslFilterChain(certChain, privateKey, rootCa, true, domains, listenerFilters)
		case *v1.SSLConfig_SslFiles:
			certChain, privateKey, rootCa := sslSecrets.SslFiles.TlsCert, sslSecrets.SslFiles.TlsKey, sslSecrets.SslFiles.RootCa
			domains := virtualService.SslConfig.SniDomains
			if len(domains) == 0 {
				domains = virtualService.Domains
			}
			filterChain = newSslFilterChain(certChain, privateKey, rootCa, false, domains, listenerFilters)
		}
		filterChains = append(filterChains, filterChain)
	}

	if listener.SslConfig != nil {
		var filterChain envoylistener.FilterChain
		switch sslSecrets := listener.SslConfig.SslSecrets.(type) {
		case *v1.SSLConfig_SecretRef:
			ref := sslSecrets.SecretRef
			certChain, privateKey, rootCa, err := getSslSecrets(ref, snap.Secrets)
			if err != nil {
				log.Warnf("skipping ssl listener with invalid secrets: %v", listener.Name)
			} else {
				domains := listener.SslConfig.SniDomains
				filterChain = newSslFilterChain(certChain, privateKey, rootCa, true, domains, listenerFilters)
			}
		case *v1.SSLConfig_SslFiles:
			certChain, privateKey, rootCa := sslSecrets.SslFiles.TlsCert, sslSecrets.SslFiles.TlsKey, sslSecrets.SslFiles.RootCa
			domains := listener.SslConfig.SniDomains
			filterChain = newSslFilterChain(certChain, privateKey, rootCa, false, domains, listenerFilters)
		}
		filterChains = append(filterChains, filterChain)
	}

	// if 0 virtualservices are defined and no ssl config is provided for the listener
	// create a filter chain with no tls
	if addInsecureFilterChain || len(filterChains) == 0 {
		filterChains = append(filterChains, envoylistener.FilterChain{
			Filters: listenerFilters,
		})
	}

	return filterChains
}

func validateListenerPorts(proxy *v1.Proxy, report reportErr) {
	listenersByPort := make(map[uint32][]string)
	for _, listener := range proxy.Listeners {
		listenersByPort[listener.BindPort] = append(listenersByPort[listener.BindPort], listener.Name)
	}
	for port, listeners := range listenersByPort {
		if len(listeners) == 1 {
			continue
		}
		report(errors.Errorf("port %v is shared by listeners %v", port, listeners), "invalid listener config")
	}
}

func newSslFilterChain(certChain, privateKey, rootCa string, inline bool, sniDomains []string, listenerFilters []envoylistener.Filter) envoylistener.FilterChain {
	var certChainData, privateKeyData, rootCaData *envoycore.DataSource
	if !inline {
		certChainData = &envoycore.DataSource{
			Specifier: &envoycore.DataSource_Filename{
				Filename: certChain,
			},
		}
		privateKeyData = &envoycore.DataSource{
			Specifier: &envoycore.DataSource_Filename{
				Filename: privateKey,
			},
		}
		rootCaData = &envoycore.DataSource{
			Specifier: &envoycore.DataSource_Filename{
				Filename: rootCa,
			},
		}
	} else {
		certChainData = &envoycore.DataSource{
			Specifier: &envoycore.DataSource_InlineString{
				InlineString: certChain,
			},
		}
		privateKeyData = &envoycore.DataSource{
			Specifier: &envoycore.DataSource_InlineString{
				InlineString: privateKey,
			},
		}
		rootCaData = &envoycore.DataSource{
			Specifier: &envoycore.DataSource_InlineString{
				InlineString: rootCa,
			},
		}
	}
	var validationContext *envoyauth.CertificateValidationContext
	var requireClientCert *gogo_types.BoolValue
	if rootCa != "" {
		requireClientCert = &gogo_types.BoolValue{Value: true}
		validationContext = &envoyauth.CertificateValidationContext{
			TrustedCa: rootCaData,
		}
	}

	return envoylistener.FilterChain{
		Filters: listenerFilters,
		TlsContext: &envoyauth.DownstreamTlsContext{
			RequireClientCertificate: requireClientCert,
			CommonTlsContext: &envoyauth.CommonTlsContext{
				// default params
				TlsParams: &envoyauth.TlsParameters{},
				// TODO: configure client certificates
				TlsCertificates: []*envoyauth.TlsCertificate{
					{
						CertificateChain: certChainData,
						PrivateKey:       privateKeyData,
					},
				},
				ValidationContextType: &envoyauth.CommonTlsContext_ValidationContext{
					ValidationContext: validationContext,
				},
			},
		},
	}
}

func (t *translator) computeListenerFilters(listener *v1.Listener, snap *v1.Snapshot, report reportErr) []envoylistener.Filter {

var listenerFilters []plugins.StagedListenerFilter
	for _, plug := range t.plugins {
		// run the Listener Filter Plugins
		params := plugins.Params{
			Snapshot: snap,
		}
		for _, plug := range t.plugins {
			filterPlugin, ok := plug.(plugins.ListenerFilterPlugin)
			if !ok {
				continue
			}
			if err := filterPlugin.ProcessListenerFilter(params, listener, out); err != nil {
				report(err, "plugin error on listener filter")
			}
		}
		filterPlugin, ok := plug.(plugins.ListenerFilterPlugin)
		if !ok {
			continue
		}
		stagedFilters, err := filterPlugin.ListenerFilters(params, listener)
		if err != nil {
			report( err, "listener plugin error")
		}
		for _, listenerFilter := range stagedFilters {
			listenerFilters = append(listenerFilters, listenerFilter)
		}
	}

	// only add the http connection manager if listener has any virtual services
	if len(listener.VirtualServices) > 0 {
		// add the http connection manager filter after all the InAuth Listener Filters
		rdsName := routeConfigName(listener)
		httpConnMgr, err := t.computeHttpConnectionManager(listener, rdsName)
		if err != nil {
			cfgErrs.addError(proxy, err)
		}
		listenerFilters = append(listenerFilters, plugins.StagedListenerFilter{
			ListenerFilter: httpConnMgr,
			Stage:          plugins.PostInAuth,
		})
	}

	// sort filters by stage
	return sortListenerFilters(listenerFilters)
}

func sortListenerFilters(filters []plugins.StagedListenerFilter) []envoylistener.Filter {
	// sort them first by stage, then by name.
	less := func(i, j int) bool {
		filteri := filters[i]
		filterj := filters[j]
		if filteri.Stage != filterj.Stage {
			return filteri.Stage < filterj.Stage
		}
		return filteri.ListenerFilter.Name < filterj.ListenerFilter.Name
	}
	sort.SliceStable(filters, less)

	var sortedFilters []envoylistener.Filter
	for _, filter := range filters {
		sortedFilters = append(sortedFilters, filter.ListenerFilter)
	}

	return sortedFilters
}
