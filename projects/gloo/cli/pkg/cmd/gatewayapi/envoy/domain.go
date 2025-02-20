package envoy

import (
	"bytes"
	"encoding/json"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	sologatewayv1alpha1 "github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/spf13/pflag"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"math/rand"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apiv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/yaml"
	"strings"

	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
)

var (
	r = rand.New(rand.NewSource(RandomSeed))
)

type Options struct {
	*options.Options
	InputFile string
	Stats     bool
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

func (o *Options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.InputFile, "input-file", "", "File to convert")
	flags.BoolVar(&o.Stats, "stats", false, "Print stats about the conversion")
}

type EnvoySnapshot struct {
	Listeners adminv3.ListenersConfigDump
	Clusters  adminv3.ClustersConfigDump
	Routes    adminv3.RoutesConfigDump
}

func (e *EnvoySnapshot) GetClusterByName(clusterName string) (*envoy_config_cluster_v3.Cluster, error) {

	for _, cluster := range e.Clusters.DynamicActiveClusters {
		var cl envoy_config_cluster_v3.Cluster
		if err := cluster.Cluster.UnmarshalTo(&cl); err != nil {
			return nil, err
		}
		if cl.Name == clusterName {
			return &cl, nil
		}
	}
	return nil, nil
}
func (e *EnvoySnapshot) GetRouteByName(routeName string) (*route.RouteConfiguration, error) {

	for _, rt := range e.Routes.DynamicRouteConfigs {
		var rc route.RouteConfiguration
		if err := rt.RouteConfig.UnmarshalTo(&rc); err != nil {
			return nil, err
		}
		if rc.Name == routeName {
			return &rc, nil
		}
	}
	return nil, nil
}

type GatewayAPIOutput struct {
	HTTPRoutes         []*gwv1.HTTPRoute
	RouteOptions       []*gatewaykube.RouteOption
	VirtualHostOptions []*gatewaykube.VirtualHostOption
	Upstreams          []*glookube.Upstream
	AuthConfigs        []*v1.AuthConfig
	Gateways           []*gwv1.Gateway
}

type YamlMarshaller struct{}

func (YamlMarshaller) ToYaml(resource interface{}) ([]byte, error) {
	switch typedResource := resource.(type) {
	case nil:
		return []byte{}, nil
	case proto.Message:
		buf := &bytes.Buffer{}
		if err := (&jsonpb.Marshaler{OrigName: true}).Marshal(buf, typedResource); err != nil {
			return nil, err
		}
		return yaml.JSONToYAML(buf.Bytes())
	default:
		return yaml.Marshal(resource)
	}
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

// removeNullFields recursively removes fields with null values from a map
func removeNullFields(m map[string]interface{}) {
	for k, v := range m {
		if v == nil {
			delete(m, k)
		} else if nestedMap, ok := v.(map[string]interface{}); ok {
			removeNullFields(nestedMap)
		} else if nestedSlice, ok := v.([]interface{}); ok {
			for _, item := range nestedSlice {
				if nestedItemMap, ok := item.(map[string]interface{}); ok {
					removeNullFields(nestedItemMap)
				}
			}
		}
	}
}
func (g *GatewayAPIOutput) ToString() (string, error) {
	output := ""

	for _, gw := range g.Gateways {
		o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), gw)
		if err != nil {
			return "", err
		}

		var data map[string]interface{}
		if err := json.Unmarshal(o, &data); err != nil {
			return "", err
		}
		removeNullFields(data)

		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return "", err
		}
		output += "\n---\n" + string(yamlData)
	}

	for _, obj := range g.Upstreams {
		o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), obj)
		if err != nil {
			return "", err
		}

		var data map[string]interface{}
		if err := json.Unmarshal(o, &data); err != nil {
			return "", err
		}
		removeNullFields(data)

		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return "", err
		}
		output += "\n---\n" + string(yamlData)
	}

	for _, obj := range g.HTTPRoutes {
		o, err := runtime.Encode(codecs.LegacyCodec(corev1.SchemeGroupVersion, gwv1.SchemeGroupVersion, gatewaykube.SchemeGroupVersion, glookube.SchemeGroupVersion), obj)
		if err != nil {
			return "", err
		}

		var data map[string]interface{}
		if err := json.Unmarshal(o, &data); err != nil {
			return "", err
		}
		removeNullFields(data)

		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return "", err
		}
		output += "\n---\n" + string(yamlData)
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
}
