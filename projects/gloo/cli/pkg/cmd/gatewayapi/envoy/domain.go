package envoy

import (
	"fmt"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"path/filepath"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	"strings"

	_ "github.com/cncf/xds/go/udpa/type/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	_ "istio.io/api/envoy/config/filter/http/alpn/v2alpha1"

	"sigs.k8s.io/yaml"
)

type Options struct {
	*options.Options
	InputFile          string
	OutputDir          string
	RouteTableFile     string
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
	flags.StringVar(&o.RouteTableFile, "route-table-file", "",
		"Organize HTTPRoutes with the same domains as the RouteTables")
}

type Outputs struct {
	OutputDir          string
	FolderPerNamespace bool
	Gateway            *gwv1.Gateway
	httpRoutes         map[string]*gwv1.HTTPRoute
	routeOptions       map[string]*gatewaykube.RouteOption
	virtualHostOptions map[string]*gatewaykube.VirtualHostOption
	upstreams          map[string]*glookube.Upstream
	Settings           v1.Settings
	envoyClusters      map[string]*envoy_config_cluster_v3.Cluster
	Errors             []error
}

func (o *Outputs) GetUpstream(name string) *glookube.Upstream {
	if o.upstreams == nil {
		return nil
	}
	return o.upstreams[name]

}

func (o *Outputs) AddUpstream(upstream *glookube.Upstream) {
	if o.upstreams == nil {
		o.upstreams = make(map[string]*glookube.Upstream)
	}
	name := fmt.Sprintf("%v/%v", upstream.Namespace, upstream.Name)
	if _, exists := o.upstreams[name]; exists {
		o.Errors = append(o.Errors, fmt.Errorf("overwriting Upstream %s", name))
	}
	o.upstreams[name] = upstream
}

func (o *Outputs) AddRouteOption(ro *gatewaykube.RouteOption) {
	if o.routeOptions == nil {
		o.routeOptions = make(map[string]*gatewaykube.RouteOption)
	}
	name := fmt.Sprintf("%v/%v", ro.Namespace, ro.Name)
	if _, exists := o.routeOptions[name]; exists {
		o.Errors = append(o.Errors, fmt.Errorf("overwriting RouteOption %s", name))
	}
	o.routeOptions[name] = ro
}
func (o *Outputs) AddVirtualHostOption(vho *gatewaykube.VirtualHostOption) {
	if o.virtualHostOptions == nil {
		o.virtualHostOptions = make(map[string]*gatewaykube.VirtualHostOption)
	}
	name := fmt.Sprintf("%v/%v", vho.Namespace, vho.Name)
	if _, exists := o.virtualHostOptions[name]; exists {
		o.Errors = append(o.Errors, fmt.Errorf("overwriting VirtualHostOption %s", name))
	}
	o.virtualHostOptions[name] = vho
}
func (o *Outputs) AddRoute(r *gwv1.HTTPRoute) {
	if o.httpRoutes == nil {
		o.httpRoutes = make(map[string]*gwv1.HTTPRoute)
	}
	name := fmt.Sprintf("%v/%v", r.Namespace, r.Name)
	if _, exists := o.httpRoutes[name]; exists {
		o.Errors = append(o.Errors, fmt.Errorf("overwriting HTTPRoute %s", name))
	}
	o.httpRoutes[name] = r
}

func (o *Outputs) AddClusters(dump *adminv3.ClustersConfigDump) error {
	if o.envoyClusters == nil {
		o.envoyClusters = make(map[string]*envoy_config_cluster_v3.Cluster)
	}
	for _, cluster := range dump.DynamicActiveClusters {
		var cl envoy_config_cluster_v3.Cluster
		if err := cluster.Cluster.UnmarshalTo(&cl); err != nil {
			return err
		}
		o.envoyClusters[cl.Name] = &cl
	}
	return nil
}

func (o *Outputs) GetClusterByName(clusterName string) *envoy_config_cluster_v3.Cluster {
	if o.envoyClusters == nil {
		return nil
	}

	return o.envoyClusters[clusterName]
}

func removeNullYamlFields(yamlData []byte) []byte {
	stringData := strings.ReplaceAll(string(yamlData), "  creationTimestamp: null\n", "")
	stringData = strings.ReplaceAll(stringData, "status:\n", "")
	stringData = strings.ReplaceAll(stringData, "parents: null\n", "")
	stringData = strings.ReplaceAll(stringData, "status: {}\n", "")
	stringData = strings.ReplaceAll(stringData, "\n\n\n", "\n")
	stringData = strings.ReplaceAll(stringData, "\n\n", "\n")
	stringData = strings.ReplaceAll(stringData, "spec: {}\n", "")
	return []byte(stringData)
}

func (o *Outputs) Write() error {
	fmt.Printf("Writing files\n")
	// Write Gateway
	gwYaml, err := yaml.Marshal(o.Gateway)
	if err != nil {
		return err
	}
	folder := o.OutputDir
	if o.FolderPerNamespace {
		folder = filepath.Join(o.OutputDir, o.Gateway.Namespace)
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(fmt.Sprintf("%s/gateway-%s.yaml", folder, o.Gateway.Name), removeNullYamlFields(gwYaml), 0644)
	if err != nil {
		return err
	}

	// Write Routes
	fmt.Printf("\tHTTPRoutes: %d\n", len(o.httpRoutes))
	for _, route := range o.httpRoutes {
		if o.FolderPerNamespace {
			folder = filepath.Join(o.OutputDir, route.Namespace)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		routeYaml, err := yaml.Marshal(route)
		if err != nil {
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/httproute-%s.yaml", folder, route.Name), removeNullYamlFields(routeYaml), 0644)
		if err != nil {
			return err
		}
	}
	fmt.Printf("\tRouteOptions: %d\n", len(o.routeOptions))
	for _, routeOption := range o.routeOptions {
		if o.FolderPerNamespace {
			folder = filepath.Join(o.OutputDir, routeOption.Namespace)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		routeOptionYaml, err := yaml.Marshal(routeOption)
		if err != nil {
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/routeoption-%s.yaml", folder, routeOption.Name), removeNullYamlFields(routeOptionYaml), 0644)
		if err != nil {
			return err
		}
	}
	fmt.Printf("\tUpstreams: %d\n", len(o.upstreams))
	for _, upstream := range o.upstreams {
		if o.FolderPerNamespace {
			folder = filepath.Join(o.OutputDir, upstream.Namespace)
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		upstreamYaml, err := yaml.Marshal(upstream)
		if err != nil {
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/upstream-%s.yaml", folder, upstream.Name), removeNullYamlFields(upstreamYaml), 0644)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(fmt.Sprintf("%s/errors.txt", o.OutputDir))
	if err != nil {
		return err
	}
	for _, err := range o.Errors {
		f.WriteString(err.Error() + "\n")
	}
	f.Close()

	// Write Settings
	// Marshal Settings to JSON first, then convert to YAML
	settingsJson, err := protojson.Marshal(&o.Settings)
	if err != nil {
		return err
	}
	settingsYaml, err := yaml.JSONToYAML(settingsJson)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/settings.yaml", o.OutputDir), settingsYaml, 0644)
	fmt.Printf("File writes complete\n")
	return err
}
