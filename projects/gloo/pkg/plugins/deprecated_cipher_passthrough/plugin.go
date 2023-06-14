package deprecated_cipher_passthrough

import (
	"context"
	"fmt"
	"strconv"

	core_v3 "github.com/cncf/xds/go/xds/core/v3"
	matcher_v3 "github.com/cncf/xds/go/xds/type/matcher/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_tls_inspector "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/tls_inspector/v3"
	network_inputs_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/matching/common_inputs/network/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-multierror"
	server_name_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/matching/custom_matchers/server_name/v3"
	cipher_inputs_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/matching/inputs/cipher_detection_input/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	gloo_dcp "github.com/solo-io/gloo/projects/gloo/pkg/plugins/deprecated_cipher_passthrough"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	_ plugins.Plugin                   = new(plugin)
	_ plugins.ListenerPlugin           = new(plugin)
	_ plugins.FilterChainMutatorPlugin = new(plugin)
)

const (
	ExtensionName = gloo_dcp.ExtensionName
	// This needs to match the constant in Envoy!!!!!
	PassthroughInputName = "PASSTHROUGH_FILTER_CHAIN"
)

// Temporary intermediate representation of the deprecated cipher passthrough matchers.
// these types are easy to work with and aggregate matchers into one by one.
// We then convert these to ugly envoy protos.

// This is similar to CidrRange proto, but comparable so it can be used as a key in a maps
type comparableCidrRange struct {
	// IPv4 or IPv6 address, e.g. ``192.0.0.0`` or ``2001:db8::``.
	AddressPrefix string
	// Length of prefix, e.g. 0, 32. Defaults to 0 when unset.
	PrefixLen uint32
}

// As we only support 2 matchings for deprecated ciphers, no map is needed here
// just a struct with 2 options default and passthrough.
// Additionally this needs to specify passthrough ciphers and optionally
// terminating ciphers to override the default native support checks
//
//	This is the leaf of the IR.
type deprecatedCipherMapping struct {
	PassthroughCipherSuites            []uint32
	TerminatingCipherSuites            []uint32
	PassthroughCipherSuitesFilterChain *envoy_config_listener_v3.FilterChain
	DefaultFilterChain                 *envoy_config_listener_v3.FilterChain
}

// A map with cidr ranges as the key
type sourceIpCidrMap map[comparableCidrRange]*deprecatedCipherMapping

// This is the root type of the IR: a map with server names as the key
type serverNameMap map[string]sourceIpCidrMap

// magic value when ip ranges not specified
var noIpRanges = comparableCidrRange{
	AddressPrefix: "",
	PrefixLen:     0,
}

// end if IR

type plugin struct {
	isUsed bool
}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Init(_ plugins.InitParams) {
	p.isUsed = false
}

func (p *plugin) Name() string {
	return ExtensionName
}

// ProcessListener implements the plugins.ListenerPlugin interface.
// Currently a no-op to allow for signature of listener translator to not change
// MUST run after tls inspector plugin.
func (p *plugin) ProcessListener(params plugins.Params, in *v1.Listener, out *envoy_config_listener_v3.Listener) error {
	if !p.isUsed {
		return nil
	}
	// mutate listener
	cipherSoloInspector := &envoy_config_listener_v3.ListenerFilter{
		Name: "envoy.filters.listener.tls_cipher_inspector",
		ConfigType: &envoy_config_listener_v3.ListenerFilter_TypedConfig{
			TypedConfig: &any.Any{TypeUrl: "type.googleapis.com/envoy.config.filter.listener.tls_cipher_inspector.v3.TlsCipherInspector"},
		},
	}

	hasTLS := false
	for _, lf := range out.ListenerFilters {
		if lf.GetName() == wellknown.TlsInspector {
			// we should put the filter here, it must come before the tls inspector
			hasTLS = true
			break
		}
	}
	if !hasTLS {
		// add both tls inspector and ours
		// add tls inspector first
		configEnvoy := &envoy_tls_inspector.TlsInspector{}
		msg, err := utils.MessageToAny(configEnvoy)
		if err != nil {
			return err
		}
		tlsInspector := &envoy_config_listener_v3.ListenerFilter{
			Name: wellknown.TlsInspector,
			ConfigType: &envoy_config_listener_v3.ListenerFilter_TypedConfig{
				TypedConfig: msg,
			},
		}
		out.ListenerFilters = append(out.GetListenerFilters(), tlsInspector)

	}
	out.ListenerFilters = append(out.GetListenerFilters(), cipherSoloInspector)
	return nil
}

// ProcessFilterChain implements the mutator type to allow for extendedfilterchain consumption
func (p *plugin) ProcessFilterChain(params plugins.Params, in *v1.Listener, inFilters []*plugins.ExtendedFilterChain, out *envoy_config_listener_v3.Listener) error {

	m, fc, err := ConvertFilterChain(params.Ctx, inFilters)
	if err != nil {
		return err
	}

	if m == nil {
		return nil
	}
	out.FilterChainMatcher = m
	out.FilterChains = fc

	p.isUsed = true

	return nil
}

func haveDeperecatedCiphers(fcm []*plugins.ExtendedFilterChain) bool {
	for _, fc := range fcm {
		if len(fc.PassthroughCipherSuites) != 0 {
			return true
		}
	}
	return false
}

/*
ConvertFilterChain is modeled after FilterChainManagerImpl::addFilterChains in envoy
and converts filter chains to the new matcher framework.
*/
func ConvertFilterChain(ctx context.Context, fcm []*plugins.ExtendedFilterChain) (*matcher_v3.Matcher, []*envoy_config_listener_v3.FilterChain, error) {
	var ret []*envoy_config_listener_v3.FilterChain
	if !haveDeperecatedCiphers(fcm) {
		//easy case, NOP
		for _, fc := range fcm {
			fc := fc // shouldnt need this but paranoia is good
			ret = append(ret, fc.FilterChain)
		}
		return nil, ret, nil
	}

	// Convert existing filter chains matchers to temporary intermediate representation that will be easy to convert to the new matcher framework
	serverNameMap, filterChains, err := filterChainsToMatcherIR(ctx, fcm)
	if err != nil {
		return nil, nil, err
	}
	// We have easy to work with intermediate representation, so now we can directly convert it to envoy config.
	m, err := convertIrToEnvoyMatcher(ctx, serverNameMap)
	return m, filterChains, err
}

func filterChainsToMatcherIR(ctx context.Context, fcm []*plugins.ExtendedFilterChain) (serverNameMap, []*envoy_config_listener_v3.FilterChain, error) {
	serverNameMap := make(serverNameMap)
	var filterChains []*envoy_config_listener_v3.FilterChain
	for i, fc := range fcm {
		fc := fc
		err := addFilterChainServerNamesToMap(ctx, serverNameMap, fc)
		if err != nil {
			return nil, nil, err
		}
		// according to docs, If matcher is specified, all filter_chains  must have a
		// non-empty and unique name field and not specify filter_chain_match
		fc.FilterChainMatch = nil
		fc.Name = fmt.Sprintf("filter_chain_%d", i)
		filterChains = append(filterChains, fc.FilterChain)
	}
	return serverNameMap, filterChains, nil
}

func toTypedExtensionConfig(ctx context.Context, name string, msg proto.Message) *core_v3.TypedExtensionConfig {
	any, err := utils.MessageToAny(msg)
	if err != nil {
		// this should never happen
		contextutils.LoggerFrom(ctx).DPanicf("unable to marshal message %v to any %w", msg, err)
	}

	// Create a new typed extension config with the given name and marshalled message.
	return &core_v3.TypedExtensionConfig{
		Name:        name,
		TypedConfig: any,
	}

}

// convert IR to envoy matcher. cannot fail, as the IR must be directly convertable to envoy api.
func convertIrToEnvoyMatcher(ctx context.Context, serverNameMap serverNameMap) (*matcher_v3.Matcher, error) {

	matcher := matcher_v3.Matcher{}
	snm := &server_name_v3.ServerNameMatcher{}
	for serverName, sourceIpCidrMap := range serverNameMap {
		// Create a server name matcher:
		sourceOnMatch := sourceCidrMapOnMatch(ctx, sourceIpCidrMap)
		if serverName == "" {
			matcher.OnNoMatch = sourceOnMatch
		} else {
			snm.ServerNameMatchers = append(snm.ServerNameMatchers, &server_name_v3.ServerNameMatcher_ServerNameSetMatcher{
				ServerNames: []string{serverName},
				OnMatch:     sourceOnMatch,
			})
		}
	}

	matcher.MatcherType = &matcher_v3.Matcher_MatcherTree_{
		MatcherTree: &matcher_v3.Matcher_MatcherTree{
			Input: toTypedExtensionConfig(ctx, "envoy.matching.inputs.server_name", &network_inputs_v3.ServerNameInput{}),
			TreeType: &matcher_v3.Matcher_MatcherTree_CustomMatch{
				CustomMatch: toTypedExtensionConfig(ctx, "envoy.matching.custom_matchers.server_name_matcher", snm),
			},
		},
	}
	return &matcher, nil
}

func sourceCidrMapOnMatch(ctx context.Context, sourceIpCidrMap sourceIpCidrMap) *matcher_v3.Matcher_OnMatch {
	matcher := &matcher_v3.Matcher{}

	onMatch := &matcher_v3.Matcher_OnMatch{
		OnMatch: &matcher_v3.Matcher_OnMatch_Matcher{
			Matcher: matcher,
		},
	}

	ipTrieMatcher := &matcher_v3.IPMatcher{}
	for sourceIpCidr, deprecatedCipherMap := range sourceIpCidrMap {
		if sourceIpCidr == noIpRanges {
			dcom := deprecatedCipherOnMatch(ctx, deprecatedCipherMap)
			if len(sourceIpCidrMap) == 1 {
				// simple case, no ip ranges, skip this table
				return dcom
			}
			matcher.OnNoMatch = dcom
		} else {
			dcom := deprecatedCipherOnMatch(ctx, deprecatedCipherMap)

			ipTrieMatcher.RangeMatchers = append(ipTrieMatcher.RangeMatchers, &matcher_v3.IPMatcher_IPRangeMatcher{
				Ranges: []*core_v3.CidrRange{{
					AddressPrefix: sourceIpCidr.AddressPrefix,
					PrefixLen:     wrapperspb.UInt32(sourceIpCidr.PrefixLen),
				}},
				OnMatch:   dcom,
				Exclusive: true,
			})
		}
	}

	matcher.MatcherType = &matcher_v3.Matcher_MatcherTree_{
		MatcherTree: &matcher_v3.Matcher_MatcherTree{
			Input: toTypedExtensionConfig(ctx, "envoy.matching.inputs.source_ip", &network_inputs_v3.SourceIPInput{}),
			TreeType: &matcher_v3.Matcher_MatcherTree_CustomMatch{
				CustomMatch: toTypedExtensionConfig(ctx, "envoy.matching.custom_matchers.trie_matcher", ipTrieMatcher),
			},
		},
	}
	return onMatch
}

// this is the terminal function that will map passthrough ciphers to the string action
// that is the filter chain name
func deprecatedCipherOnMatch(ctx context.Context, deprecatedCipherMap *deprecatedCipherMapping) *matcher_v3.Matcher_OnMatch {

	var defaultChainAction *matcher_v3.Matcher_OnMatch
	if deprecatedCipherMap.DefaultFilterChain != nil {
		defaultChainAction = &matcher_v3.Matcher_OnMatch{
			OnMatch: &matcher_v3.Matcher_OnMatch_Action{
				Action: toTypedExtensionConfig(ctx, "string", wrapperspb.String(deprecatedCipherMap.DefaultFilterChain.Name)),
			},
		}
	}

	// dont add the node if its not needed
	if deprecatedCipherMap.PassthroughCipherSuitesFilterChain == nil {
		return defaultChainAction
	}

	matcher := &matcher_v3.Matcher{}

	onMatch := &matcher_v3.Matcher_OnMatch{
		OnMatch: &matcher_v3.Matcher_OnMatch_Matcher{
			Matcher: matcher,
		},
	}

	matcherMap := map[string]*matcher_v3.Matcher_OnMatch{}

	if deprecatedCipherMap.PassthroughCipherSuitesFilterChain != nil {
		matcherMap[PassthroughInputName] = &matcher_v3.Matcher_OnMatch{
			OnMatch: &matcher_v3.Matcher_OnMatch_Action{
				Action: toTypedExtensionConfig(ctx, "string", wrapperspb.String(deprecatedCipherMap.PassthroughCipherSuitesFilterChain.Name)),
			},
		}
	}

	matcher.MatcherType = &matcher_v3.Matcher_MatcherTree_{
		MatcherTree: &matcher_v3.Matcher_MatcherTree{
			Input: toTypedExtensionConfig(ctx, "envoy.matching.inputs.cipher_detection_input", &cipher_inputs_v3.CipherDetectionInput{
				PassthroughCiphers: deprecatedCipherMap.PassthroughCipherSuites,
				TerminatingCiphers: deprecatedCipherMap.TerminatingCipherSuites,
			}),
			TreeType: &matcher_v3.Matcher_MatcherTree_ExactMatchMap{
				ExactMatchMap: &matcher_v3.Matcher_MatcherTree_MatchMap{
					Map: matcherMap,
				},
			},
		},
	}
	matcher.OnNoMatch = defaultChainAction

	return onMatch
}

// Command to generate this table:
// openssl ciphers -V | awk '{print "\x22"$3"\x22"": "$1","}' | sed -r 's/,0x//g'
var nameConversion = map[string]uint32{
	"TLS_AES_256_GCM_SHA384":        0x1302,
	"TLS_CHACHA20_POLY1305_SHA256":  0x1303,
	"TLS_AES_128_GCM_SHA256":        0x1301,
	"ECDHE-ECDSA-AES256-GCM-SHA384": 0xC02C,
	"ECDHE-RSA-AES256-GCM-SHA384":   0xC030,
	"DHE-RSA-AES256-GCM-SHA384":     0x009F,
	"ECDHE-ECDSA-CHACHA20-POLY1305": 0xCCA9,
	"ECDHE-RSA-CHACHA20-POLY1305":   0xCCA8,
	"DHE-RSA-CHACHA20-POLY1305":     0xCCAA,
	"ECDHE-ECDSA-AES128-GCM-SHA256": 0xC02B,
	"ECDHE-RSA-AES128-GCM-SHA256":   0xC02F,
	"DHE-RSA-AES128-GCM-SHA256":     0x009E,
	"ECDHE-ECDSA-AES256-SHA384":     0xC024,
	"ECDHE-RSA-AES256-SHA384":       0xC028,
	"DHE-RSA-AES256-SHA256":         0x006B,
	"ECDHE-ECDSA-AES128-SHA256":     0xC023,
	"ECDHE-RSA-AES128-SHA256":       0xC027,
	"DHE-RSA-AES128-SHA256":         0x0067,
	"ECDHE-ECDSA-AES256-SHA":        0xC00A,
	"ECDHE-RSA-AES256-SHA":          0xC014,
	"DHE-RSA-AES256-SHA":            0x0039,
	"ECDHE-ECDSA-AES128-SHA":        0xC009,
	"ECDHE-RSA-AES128-SHA":          0xC013,
	"DHE-RSA-AES128-SHA":            0x0033,
	"RSA-PSK-AES256-GCM-SHA384":     0x00AD,
	"DHE-PSK-AES256-GCM-SHA384":     0x00AB,
	"RSA-PSK-CHACHA20-POLY1305":     0xCCAE,
	"DHE-PSK-CHACHA20-POLY1305":     0xCCAD,
	"ECDHE-PSK-CHACHA20-POLY1305":   0xCCAC,
	"AES256-GCM-SHA384":             0x009D,
	"PSK-AES256-GCM-SHA384":         0x00A9,
	"PSK-CHACHA20-POLY1305":         0xCCAB,
	"RSA-PSK-AES128-GCM-SHA256":     0x00AC,
	"DHE-PSK-AES128-GCM-SHA256":     0x00AA,
	"AES128-GCM-SHA256":             0x009C,
	"PSK-AES128-GCM-SHA256":         0x00A8,
	"AES256-SHA256":                 0x003D,
	"AES128-SHA256":                 0x003C,
	"ECDHE-PSK-AES256-CBC-SHA384":   0xC038,
	"ECDHE-PSK-AES256-CBC-SHA":      0xC036,
	"SRP-RSA-AES-256-CBC-SHA":       0xC021,
	"SRP-AES-256-CBC-SHA":           0xC020,
	"RSA-PSK-AES256-CBC-SHA384":     0x00B7,
	"DHE-PSK-AES256-CBC-SHA384":     0x00B3,
	"RSA-PSK-AES256-CBC-SHA":        0x0095,
	"DHE-PSK-AES256-CBC-SHA":        0x0091,
	"AES256-SHA":                    0x0035,
	"PSK-AES256-CBC-SHA384":         0x00AF,
	"PSK-AES256-CBC-SHA":            0x008D,
	"ECDHE-PSK-AES128-CBC-SHA256":   0xC037,
	"ECDHE-PSK-AES128-CBC-SHA":      0xC035,
	"SRP-RSA-AES-128-CBC-SHA":       0xC01E,
	"SRP-AES-128-CBC-SHA":           0xC01D,
	"RSA-PSK-AES128-CBC-SHA256":     0x00B6,
	"DHE-PSK-AES128-CBC-SHA256":     0x00B2,
	"RSA-PSK-AES128-CBC-SHA":        0x0094,
	"DHE-PSK-AES128-CBC-SHA":        0x0090,
	"AES128-SHA":                    0x002F,
	"PSK-AES128-CBC-SHA256":         0x00AE,
	"PSK-AES128-CBC-SHA":            0x008C,
}

func convertCipherNameToU32(cipherNames []string) ([]uint32, error) {
	// preallocate space but in the terminating cipher instance we may have partial
	// sucess so we dont want any nils in the slice returned
	u32 := make([]uint32, 0, len(cipherNames))
	var multiErr *multierror.Error
	for _, name := range cipherNames {
		intRep, ok := nameConversion[name]
		if !ok {
			rep64, err := strconv.ParseUint(name, 0, 32)
			if err != nil {
				multiErr = multierror.Append(multiErr, fmt.Errorf("unsupported cipher name %s", name))
			} else {
				intRep = uint32(rep64)
			}
		}
		// ie if its valid
		if intRep != 0 {
			u32 = append(u32, intRep)
		}
	}
	return u32, multiErr.ErrorOrNil()
}

func addFilterChainServerNamesToMap(ctx context.Context, serverNameMap serverNameMap, fc *plugins.ExtendedFilterChain) error {
	serverNames := fc.GetFilterChainMatch().GetServerNames()
	if len(serverNames) == 0 {
		return serverNameMap.addServerNameToMap(ctx, "", fc)
	}
	for _, serverName := range serverNames {
		err := serverNameMap.addServerNameToMap(ctx, serverName, fc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m serverNameMap) addServerNameToMap(ctx context.Context, srvname string, fc *plugins.ExtendedFilterChain) error {
	if m[srvname] == nil {
		m[srvname] = make(sourceIpCidrMap)
	}
	sourceIpRanges := fc.GetFilterChainMatch().GetSourcePrefixRanges()
	if len(sourceIpRanges) == 0 {
		return m[srvname].addSourceIpToMap(ctx, noIpRanges, fc)
	}
	for _, ipRange := range sourceIpRanges {
		cirdRange := comparableCidrRange{
			AddressPrefix: ipRange.GetAddressPrefix(),
			PrefixLen:     ipRange.GetPrefixLen().GetValue(),
		}
		err := m[srvname].addSourceIpToMap(ctx, cirdRange, fc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m sourceIpCidrMap) addSourceIpToMap(ctx context.Context,
	prefix comparableCidrRange, fc *plugins.ExtendedFilterChain) error {
	if m[prefix] == nil {
		m[prefix] = &deprecatedCipherMapping{}
	}

	return m[prefix].addPassthroughCiphers(ctx, fc.PassthroughCipherSuites, fc)
}

func (m *deprecatedCipherMapping) addPassthroughCiphers(ctx context.Context,
	passthroughCipherSuites []string, fc *plugins.ExtendedFilterChain) error {
	if len(passthroughCipherSuites) == 0 {
		if m.DefaultFilterChain != nil {
			return fmt.Errorf("multiple filter chains with overlapping matching rules are defined")
		}
		m.DefaultFilterChain = fc.FilterChain
		// if the terminating ciphers are specified, we need to override the default
		terminatingCiphers, err := convertCipherNameToU32(fc.TerminatingCipherSuites)
		if err != nil {
			// for now dont log until we get logger plumbed in
			contextutils.LoggerFrom(ctx).Errorf(
				"unable to convert sslconfig's noted cipher based on lookup: %s", err.Error())
		}
		m.TerminatingCipherSuites = terminatingCiphers

	} else {
		if m.PassthroughCipherSuitesFilterChain != nil {
			return fmt.Errorf("multiple filter chains with overlapping matching rules are defined")
		}

		passthroughCiphers, err := convertCipherNameToU32(passthroughCipherSuites)
		if err != nil {
			return err
		}
		m.PassthroughCipherSuites = passthroughCiphers
		m.PassthroughCipherSuitesFilterChain = fc.FilterChain
	}

	// check that we dont have overlap
	if m.PassthroughCipherSuites != nil && m.TerminatingCipherSuites != nil {
		for _, cipher := range m.PassthroughCipherSuites {
			for _, termCipher := range m.TerminatingCipherSuites {
				if cipher == termCipher {
					return fmt.Errorf("multiple filter chains have shared ciphers between passthrough and non, cipher:%v", cipher)
				}
			}
		}
	}
	return nil
}
