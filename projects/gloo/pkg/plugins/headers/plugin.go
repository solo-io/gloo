package headers

import (
	"os"
	"strconv"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"

	envoy_config_mutation_rules_v3 "github.com/envoyproxy/go-control-plane/envoy/config/common/mutation_rules/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_ehm_header_mutation_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/http/early_header_mutation/header_mutation/v3"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/api_conversion"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/headers"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var (
	_ plugins.RoutePlugin                 = new(plugin)
	_ plugins.VirtualHostPlugin           = new(plugin)
	_ plugins.WeightedDestinationPlugin   = new(plugin)
	_ plugins.HttpConnectionManagerPlugin = new(plugin)
)

const (
	ExtensionName = "headers"
)

var (
	MissingHeaderValueError = eris.New("header section of header value option cannot be nil")
	CantSetHostHeaderError  = eris.New("cannot set Host header in response headers")

	CantSetPseudoHeaderError = func(header string) error {
		return eris.Errorf(":-prefixed headers cannot be set: '%s'", header)
	}
)

// Puts Header Manipulation config on Routes, VirtualHosts, and Weighted Clusters
type plugin struct{}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(_ plugins.InitParams) {}

func (p *plugin) ProcessWeightedDestination(
	params plugins.RouteActionParams,
	in *v1.WeightedDestination,
	out *envoy_config_route_v3.WeightedCluster_ClusterWeight,
) error {
	headerManipulation := in.GetOptions().GetHeaderManipulation()
	if headerManipulation == nil {
		return nil
	}
	enforceMatchingNamespaces, err := getEnforceMatch()
	if err != nil {
		return err
	}
	upstreamNamespace := ""
	// Avoid the performance impact of looking up the upstream namespace if we don't need it
	// This is more important on routes and virtual hosts.
	if enforceMatchingNamespaces {
		us, err := upstreams.DestinationToUpstreamRef(in.GetDestination())
		if err == nil {
			upstreamNamespace = us.GetNamespace()
		}
	}
	headerSecretOptions := api_conversion.HeaderSecretOptions{EnforceNamespaceMatch: enforceMatchingNamespaces, UpstreamNamespace: upstreamNamespace}
	envoyHeader, err := convertHeaderConfig(headerManipulation, getSecretsFromSnapshot(params.Snapshot), headerSecretOptions)
	if err != nil {
		return err
	}

	out.RequestHeadersToAdd = envoyHeader.RequestHeadersToAdd
	out.RequestHeadersToRemove = envoyHeader.RequestHeadersToRemove
	out.ResponseHeadersToAdd = envoyHeader.ResponseHeadersToAdd
	out.ResponseHeadersToRemove = envoyHeader.ResponseHeadersToRemove

	return nil
}

func (p *plugin) ProcessVirtualHost(
	params plugins.VirtualHostParams,
	in *v1.VirtualHost,
	out *envoy_config_route_v3.VirtualHost,
) error {
	headerManipulation := in.GetOptions().GetHeaderManipulation()
	if headerManipulation == nil {
		return nil
	}
	enforceMatchingNamespaces, err := getEnforceMatch()
	if err != nil {
		return err
	}
	headerSecretOptions := api_conversion.HeaderSecretOptions{
		EnforceNamespaceMatch: enforceMatchingNamespaces,
	}
	// Avoid the performance impact of looking up the upstream namespace if we don't need it
	if enforceMatchingNamespaces && len(in.GetRoutes()) > 0 {
		usNamespace := getUpstreamNamespaceForRouteAction(params.Snapshot, in.GetRoutes()[0].GetRouteAction())
		if len(in.GetRoutes()) > 1 {
			for _, r := range in.GetRoutes()[1:] {
				// in order for the namespace match check to make sense, all the upstreams on the virtual host need to be the same
				if getUpstreamNamespaceForRouteAction(params.Snapshot, r.GetRouteAction()) != usNamespace {
					usNamespace = ""
					break
				}
			}
		}
		headerSecretOptions.UpstreamNamespace = usNamespace
	}
	envoyHeader, err := convertHeaderConfig(headerManipulation, getSecretsFromSnapshot(params.Snapshot), headerSecretOptions)
	if err != nil {
		return err
	}

	out.RequestHeadersToAdd = envoyHeader.RequestHeadersToAdd
	out.RequestHeadersToRemove = envoyHeader.RequestHeadersToRemove
	out.ResponseHeadersToAdd = envoyHeader.ResponseHeadersToAdd
	out.ResponseHeadersToRemove = envoyHeader.ResponseHeadersToRemove

	return nil
}

func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	headerManipulation := in.GetOptions().GetHeaderManipulation()
	if headerManipulation == nil {
		return nil
	}
	enforceMatchingNamespaces, err := getEnforceMatch()
	if err != nil {
		return err
	}
	headerSecretOptions := api_conversion.HeaderSecretOptions{EnforceNamespaceMatch: enforceMatchingNamespaces}
	// Avoid the performance impact of looking up the upstream namespace if we don't need it
	if enforceMatchingNamespaces {
		headerSecretOptions.UpstreamNamespace = getUpstreamNamespaceForRouteAction(params.Snapshot, in.GetRouteAction())
	}
	envoyHeader, err := convertHeaderConfig(headerManipulation, getSecretsFromSnapshot(params.Snapshot), headerSecretOptions)
	if err != nil {
		return err
	}

	out.RequestHeadersToAdd = envoyHeader.RequestHeadersToAdd
	out.RequestHeadersToRemove = envoyHeader.RequestHeadersToRemove
	out.ResponseHeadersToAdd = envoyHeader.ResponseHeadersToAdd
	out.ResponseHeadersToRemove = envoyHeader.ResponseHeadersToRemove

	return nil
}

func (p *plugin) ProcessHcmNetworkFilter(params plugins.Params, parentListener *v1.Listener,
	listener *v1.HttpListener, out *envoyhttp.HttpConnectionManager) error {

	in := listener.GetOptions().GetHttpConnectionManagerSettings()
	if in == nil {
		return nil
	}

	inManipulations := in.GetEarlyHeaderManipulation()
	if inManipulations == nil {
		return nil
	}

	requestAdd, err := api_conversion.ToEnvoyHeaderValueOptionList(inManipulations.GetHeadersToAdd(),
		getSecretsFromSnapshot(params.Snapshot), api_conversion.HeaderSecretOptions{})
	if err != nil {
		return err
	}

	outMutations := []*envoy_config_mutation_rules_v3.HeaderMutation{}

	for _, header := range requestAdd {
		outMutations = append(outMutations, &envoy_config_mutation_rules_v3.HeaderMutation{
			Action: &envoy_config_mutation_rules_v3.HeaderMutation_Append{
				Append: header,
			},
		})
	}

	for _, header := range inManipulations.GetHeadersToRemove() {
		outMutations = append(outMutations, &envoy_config_mutation_rules_v3.HeaderMutation{
			Action: &envoy_config_mutation_rules_v3.HeaderMutation_Remove{
				Remove: header,
			},
		})
	}

	typedConfig, err := utils.MessageToAny(&envoy_ehm_header_mutation_v3.HeaderMutation{
		Mutations: outMutations,
	})
	if err != nil {
		return err
	}

	out.EarlyHeaderMutationExtensions = []*corev3.TypedExtensionConfig{
		{
			Name:        "http.early_header_mutation.header_mutation",
			TypedConfig: typedConfig,
		},
	}

	return nil
}

type envoyHeaderManipulation struct {
	RequestHeadersToAdd     []*envoy_config_core_v3.HeaderValueOption
	RequestHeadersToRemove  []string
	ResponseHeadersToAdd    []*envoy_config_core_v3.HeaderValueOption
	ResponseHeadersToRemove []string
}

func getEnforceMatch() (bool, error) {

	enforceMatchStr := os.Getenv(api_conversion.MatchingNamespaceEnv)
	// ParseBool errors on empty string but we want to treat that as false
	if enforceMatchStr == "" {
		return false, nil
	}
	return strconv.ParseBool(enforceMatchStr)
}

// getUpstreamNamespaceForRouteAction finds the destination upstreams for a route action and if there's only one namespace
// between them, returns that namespace, otherwise returns an empty string.
func getUpstreamNamespaceForRouteAction(snapshot *v1snap.ApiSnapshot, action *v1.RouteAction) string {
	usRefs, err := pluginutils.DestinationUpstreams(snapshot, action)
	if err != nil || len(usRefs) == 0 {
		return ""
	}
	ns := usRefs[0].GetNamespace()
	// verify that all the upstreams in the list are in the same namespace
	// if not, we can't do a meaningful check for matching namespaces, so we will fail if headerSecretRef is set
	// otherwise that's weird but fine
	if len(usRefs) > 1 {
		for _, u := range usRefs[1:] {
			if u.GetNamespace() != ns {
				return ""
			}
		}
	}
	return ns
}
func getSecretsFromSnapshot(snapshot *v1snap.ApiSnapshot) *v1.SecretList {
	var secrets *v1.SecretList
	if snapshot == nil {
		secrets = &v1.SecretList{}
	} else {
		secrets = &snapshot.Secrets
	}
	return secrets
}

func convertHeaderConfig(in *headers.HeaderManipulation, secrets *v1.SecretList, secretOptions api_conversion.HeaderSecretOptions) (*envoyHeaderManipulation, error) {
	// request headers can either be made from a normal key/value pair, or.
	// they can be constructed from a supplied secret. To accomplish this, we use
	// a utility function that was originally created to accomplish this for health check headers.
	requestAdd, err := api_conversion.ToEnvoyHeaderValueOptionList(in.GetRequestHeadersToAdd(), secrets, secretOptions)
	if err != nil {
		return nil, err
	}
	// response headers have no reason to include secrets.
	responseAdd, err := convertResponseHeaderValueOption(in.GetResponseHeadersToAdd())
	if err != nil {
		return nil, err
	}

	return &envoyHeaderManipulation{
		RequestHeadersToAdd:     requestAdd,
		RequestHeadersToRemove:  in.GetRequestHeadersToRemove(),
		ResponseHeadersToAdd:    responseAdd,
		ResponseHeadersToRemove: in.GetResponseHeadersToRemove(),
	}, nil
}

func convertResponseHeaderValueOption(
	in []*headers.HeaderValueOption,
) ([]*envoy_config_core_v3.HeaderValueOption, error) {
	var out []*envoy_config_core_v3.HeaderValueOption
	for _, h := range in {
		header := h.GetHeader()
		if header == nil {
			return nil, MissingHeaderValueError
		}

		if strings.HasPrefix(header.GetKey(), ":") {
			return nil, CantSetPseudoHeaderError(header.GetKey())
		}

		if strings.EqualFold(header.GetKey(), "Host") {
			return nil, CantSetHostHeaderError
		}

		appendAction := envoy_config_core_v3.HeaderValueOption_APPEND_IF_EXISTS_OR_ADD
		if appendOption := h.GetAppend(); appendOption != nil {
			if appendOption.GetValue() == false {
				appendAction = envoy_config_core_v3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD
			}
		}

		out = append(out, &envoy_config_core_v3.HeaderValueOption{
			Header: &envoy_config_core_v3.HeaderValue{
				Key:   header.GetKey(),
				Value: header.GetValue(),
			},
			AppendAction: appendAction,
		})
	}
	return out, nil
}
