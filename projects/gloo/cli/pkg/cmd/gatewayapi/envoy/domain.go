package envoy

import (
	"bytes"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	sologatewayv1alpha1 "github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/pflag"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"math/rand"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apiv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
	RandomSuffix = 4
	RandomSeed   = 1
)

var (
	runtimeScheme *runtime.Scheme
	codecs        serializer.CodecFactory
	letterRunes   = []rune("abcdefghijklmnopqrstuvwxyz")
	r             = rand.New(rand.NewSource(RandomSeed))
	SchemeBuilder = runtime.SchemeBuilder{
		// K8s Gateway API resources
		gwv1.Install,
		apiv1beta1.Install,

		// Kubernetes Core resources
		corev1.AddToScheme,
		appsv1.AddToScheme,

		// Solo Kubernetes Gateway API resources
		sologatewayv1alpha1.AddToScheme,
	}
)

type Options struct {
	*options.Options
	InputFile          string
	OutputDir          string
	FolderPerNamespace bool
	Stats              bool
}

func (o *Options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.InputFile, "input-file", "", "File to convert")
	flags.BoolVar(&o.Stats, "stats", false, "Print stats about the conversion")
	flags.StringVar(&o.OutputDir, "_output", "./_output",
		"Where to write files")
	flags.BoolVar(&o.FolderPerNamespace, "folder-per-namespace", false,
		"Organize files by namespace")
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

func RandStringRunes(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}
