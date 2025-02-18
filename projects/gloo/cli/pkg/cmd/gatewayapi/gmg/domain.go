package gmg

import (
	"bytes"
	sologatewayv1alpha1 "github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	admin_v3 "github.com/solo-io/solo-apis/client-go/admin.gloo.solo.io/v2"
	apimanagement_v2 "github.com/solo-io/solo-apis/client-go/apimanagement.gloo.solo.io/v2"
	infrastructure_v3 "github.com/solo-io/solo-apis/client-go/infrastructure.gloo.solo.io/v2"
	networking_v2 "github.com/solo-io/solo-apis/client-go/networking.gloo.solo.io/v2"
	observability_v2 "github.com/solo-io/solo-apis/client-go/observability.policy.gloo.solo.io/v2"
	resilience_v2 "github.com/solo-io/solo-apis/client-go/resilience.policy.gloo.solo.io/v2"
	security_v2 "github.com/solo-io/solo-apis/client-go/security.policy.gloo.solo.io/v2"
	traffic_v2 "github.com/solo-io/solo-apis/client-go/trafficcontrol.policy.gloo.solo.io/v2"
	appsv1 "k8s.io/api/apps/v1"
	apiv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"strings"

	"github.com/itchyny/json2yaml"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type Options struct {
	*options.Options

	InputFile                string
	Directory                string
	Overwrite                bool
	OverwriteSuffix          string
	Stats                    bool
	GCPRegex                 string
	RemoveGCPAUthConfig      bool
	RouteOptionStategicMerge bool
}

func (o *Options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.InputFile, "input-file", "", "File to convert")
	flags.BoolVar(&o.Overwrite, "overwrite", false, "Overwrite the existing files with the changes")
	flags.StringVar(&o.OverwriteSuffix, "suffix", "", "When writing to files add a suffix (to do side by side)")
	flags.BoolVar(&o.Stats, "stats", false, "Print stats about the conversion")
	flags.StringVar(&o.Directory, "dir", "", "Directory to read yaml/yml files")
}

type GlooMeshInput struct {
	FileName                   string
	YamlObjects                []string
	RouteTables                []*networking_v2.RouteTable
	VirtualDestinations        []*networking_v2.VirtualDestination
	VirtualGateways            []*networking_v2.VirtualGateway
	ExternalServices           []*networking_v2.ExternalService
	ExternalEndpoints          []*networking_v2.ExternalEndpoint
	RateLimitPolicies          []*traffic_v2.RateLimitPolicy
	ExtAuthPolicies            []*security_v2.ExtAuthPolicy
	CORSPolicies               []*security_v2.CORSPolicy
	JWTPolicies                []*security_v2.JWTPolicy
	HeaderManipulationPolicies []*traffic_v2.HeaderManipulationPolicy
	ConnectionPolicies         []*resilience_v2.ConnectionPolicy
	RateLimitServerConfigs     []*admin_v3.RateLimitServerConfig
	RateLimitClientConfigs     []*traffic_v2.RateLimitClientConfig
	Portals                    []*apimanagement_v2.Portal
	PortalGroups               []*apimanagement_v2.PortalGroup
	CloudProviders             []*infrastructure_v3.CloudProvider
	CloudResources             []*infrastructure_v3.CloudResources
	AccessLogPolicies          []*observability_v2.AccessLogPolicy
	ActiveHealthCheckPolicies  []*resilience_v2.ActiveHealthCheckPolicy
	FailoverPolicies           []*resilience_v2.FailoverPolicy
	FaultInjectionPolicies     []*resilience_v2.FaultInjectionPolicy
	ListenerConnectionPolicies []*resilience_v2.ListenerConnectionPolicy
	OutlierDetectionPolicies   []*resilience_v2.OutlierDetectionPolicy
	RetryTimeoutPolicies       []*resilience_v2.RetryTimeoutPolicy
	CSRFPolicies               []*security_v2.CSRFPolicy
	DLPPolicies                []*security_v2.DLPPolicy
	WAFPolicies                []*security_v2.WAFPolicy
	HTTPBufferPolicies         []*traffic_v2.HTTPBufferPolicy
	LoadBalancerPolicies       []*traffic_v2.LoadBalancerPolicy
	MirrorPolicies             []*traffic_v2.MirrorPolicy
	ProxyProtocolPolicies      []*traffic_v2.ProxyProtocolPolicy
	TransformationPolicies     []*traffic_v2.TransformationPolicy
}

type DelegateParentReference struct {
	Labels          map[string]string
	ParentName      string
	ParentNamespace string
}

type GatewayAPIOutput struct {
	FileName             string
	YamlObjects          []string
	DelegationReferences []*DelegateParentReference
	HTTPRoutes           []*gwv1.HTTPRoute
	RouteOptions         []*gatewaykube.RouteOption
	VirtualHostOptions   []*gatewaykube.VirtualHostOption
	Upstreams            []*glookube.Upstream
	AuthConfigs          []*v1.AuthConfig
	// Gateways             []*gwv1.Gateway
}

func (g *GatewayAPIOutput) HasItems() bool {

	if len(g.HTTPRoutes) > 0 {
		return true
	}
	if len(g.RouteOptions) > 0 {
		return true
	}
	if len(g.VirtualHostOptions) > 0 {
		return true
	}
	if len(g.Upstreams) > 0 {
		return true
	}
	if len(g.AuthConfigs) > 0 {
		return true
	}
	// if there are only yaml objects then skip because we didnt change anything in the file

	return false
}

func (g *GatewayAPIOutput) ToString() (string, error) {
	output := ""

	for _, y := range g.YamlObjects {
		output += "\n---\n" + y + "\n"
	}

	for _, obj := range g.Upstreams {
		o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), obj)
		if err != nil {
			return "", err
		}

		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, bytes.NewReader(o)); err != nil {
			return "", err
		}

		output += "\n---\n" + yaml.String()
	}

	for _, obj := range g.HTTPRoutes {
		o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), obj)
		if err != nil {
			return "", err
		}

		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, bytes.NewReader(o)); err != nil {
			return "", err
		}

		output += "\n---\n" + yaml.String()
	}
	for _, op := range g.RouteOptions {
		marshaller := YamlMarshaller{}
		yaml, err := marshaller.ToYaml(&op)
		if err != nil {
			return "", err
		}

		output += "\n---\n" + string(yaml)
	}
	for _, op := range g.AuthConfigs {
		marshaller := YamlMarshaller{}
		yaml, err := marshaller.ToYaml(&op)
		if err != nil {
			return "", err
		}

		output += "\n---\n" + string(yaml)
	}

	for _, op := range g.VirtualHostOptions {
		marshaller := YamlMarshaller{}
		yaml, err := marshaller.ToYaml(&op)
		if err != nil {
			return "", err
		}

		output += "\n---\n" + string(yaml)
	}

	// need to remove a few values
	//  creationTimestamp: null
	// status: {}
	// status:
	// parents: null
	output = strings.ReplaceAll(output, "  creationTimestamp: null\n", "")
	output = strings.ReplaceAll(output, "status:\n", "")
	output = strings.ReplaceAll(output, "parents: null\n", "")
	output = strings.ReplaceAll(output, "status: {}\n", "")
	output = strings.ReplaceAll(output, "\n\n\n", "\n")
	output = strings.ReplaceAll(output, "\n\n", "\n")
	output = strings.ReplaceAll(output, "spec: {}\n", "")

	// TODO remove leading and trailing ---
	// log.Printf("%s", output)
	return output, nil
}

var SchemeBuilder = runtime.SchemeBuilder{
	// K8s Gateway API resources
	gwv1.Install,
	apiv1beta1.Install,

	// Kubernetes Core resources
	corev1.AddToScheme,
	appsv1.AddToScheme,

	// Solo Kubernetes Gateway API resources
	sologatewayv1alpha1.AddToScheme,

	// Solo Gloo Mesh API resources
	networking_v2.AddToScheme,
	security_v2.AddToScheme,
	admin_v3.AddToScheme,
	resilience_v2.AddToScheme,
	traffic_v2.AddToScheme,
	observability_v2.AddToScheme,
	infrastructure_v3.AddToScheme,
	apimanagement_v2.AddToScheme,
}
