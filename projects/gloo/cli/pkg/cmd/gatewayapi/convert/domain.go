package convert

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/types"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/pflag"
)

type ErrorType string

const (
	ERROR_TYPE_UPDATE_OBJECT             ErrorType = "UPDATE_OBJECT"
	ERROR_TYPE_NOT_SUPPORTED                       = "NOT_SUPPORTED"
	ERROR_TYPE_IGNORED                             = "IGNORED"
	ERROR_TYPE_UNKNOWN_REFERENCE                   = "UNKNOWN_REFERENCE"
	ERROR_TYPE_NO_REFERENCES                       = "NO_REFERENCES"
	ERROR_TYPE_CEL_VALIDATION_CORRECTION           = "CEL_VALIDATION_CORRECTION"
)

type Options struct {
	*options.Options
	InputFile               string
	InputDir                string
	GlooSnapshotFile        string
	OutputDir               string
	Stats                   bool
	CombineRouteOptions     bool
	RetainFolderStructure   bool
	IncludeUnknownResources bool
	DeleteOutputDir         bool
	CreateNamespaces        bool
	ControlPlaneName        string
	ControlPlaneNamespace   string
}

func (opts *Options) validate() error {

	count := 0
	if opts.InputDir != "" {
		count++
	}
	if opts.InputFile != "" {
		count++
	}
	if opts.GlooSnapshotFile != "" {
		count++
	}
	if opts.ControlPlaneName != "" {
		count++
		if opts.ControlPlaneNamespace == "" {
			return fmt.Errorf("pod namespace must be specified")
		}
	}

	if count > 1 {
		return fmt.Errorf("only one of 'input-file' or 'directory' or 'input-snapshot' or `gloo-pod-name` can be specified")
	}
	if !opts.DeleteOutputDir && folderExists(opts.OutputDir) {
		return fmt.Errorf("output-dir already %s exists. It can be deleted with --delete-output-dir", opts.OutputDir)
	}
	return nil
}
func folderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
func (o *Options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.ControlPlaneName, "gloo-control-plane", "", "Name of the Gloo control plane pod")
	flags.StringVarP(&o.ControlPlaneNamespace, "gloo-control-plane-namespace", "n", "gloo-system", "Namespace of the Gloo control plane pod")
	flags.StringVar(&o.InputFile, "input-file", "", "Convert a single YAML file to the Gateway API")
	flags.StringVar(&o.InputDir, "input-dir", "", "InputDir to read yaml/yml files recursively")
	flags.StringVar(&o.GlooSnapshotFile, "input-snapshot", "", "Gloo input snapshot file location")
	flags.BoolVar(&o.Stats, "print-stats", false, "Print stats about the conversion")
	flags.BoolVar(&o.CombineRouteOptions, "combine-route-options", false, "Combine RouteOptions that are exactly the same and share them among the HTTPRoutes")
	flags.StringVar(&o.OutputDir, "output-dir", "./_output", "Output directory to write Gateway API configurations to. The directory must not exist before starting the migration. To delete and recreate the output directory, use the --recreate-output-dir option")
	flags.BoolVar(&o.RetainFolderStructure, "retain-input-folder-structure", false, "Arrange the generated Gateway API files in the same folder structure they were read from (input-dir only).")
	flags.BoolVar(&o.IncludeUnknownResources, "include-unknown", false, "Copy non-Gloo Gateway resources to the output directory without changing them. ")
	flags.BoolVar(&o.DeleteOutputDir, "delete-output-dir", false, "Delete the output directory if it already exists.")
	flags.BoolVar(&o.CreateNamespaces, "create-namespaces", false, "Create namespaces for the objects in a file.")
}

type GatewayAPIOutput struct {
	gatewayAPICache *GatewayAPICache
	edgeCache       *snapshot.Instance
	errors          map[ErrorType][]GlooError
}

func (o *GatewayAPIOutput) GetGatewayAPICache() *GatewayAPICache {
	if o.gatewayAPICache == nil {
		o.gatewayAPICache = &GatewayAPICache{}
	}
	return o.gatewayAPICache
}

func (o *GatewayAPIOutput) GetEdgeCache() *snapshot.Instance {
	if o.edgeCache == nil {
		o.edgeCache = &snapshot.Instance{}
	}
	return o.edgeCache
}
func (o *GatewayAPIOutput) EdgeCache(instance *snapshot.Instance) {
	o.edgeCache = instance
}

type GlooError struct {
	err       error
	errorType ErrorType
	name      string
	namespace string
	crdType   string
}

func (o *GatewayAPIOutput) AddError(errType ErrorType, msg string, args ...interface{}) {
	if o.errors == nil {
		o.errors = make(map[ErrorType][]GlooError)
	}
	if o.errors[errType] == nil {
		o.errors[errType] = make([]GlooError, 0)
	}

	o.errors[errType] = append(o.errors[errType], GlooError{
		err:       fmt.Errorf(msg, args...),
		errorType: errType,
		name:      "none",
		namespace: "none",
		crdType:   "none",
	})
}
func (o *GatewayAPIOutput) AddErrorFromWrapper(errType ErrorType, wrapper snapshot.Wrapper, msg string, args ...interface{}) {
	if o.errors == nil {
		o.errors = make(map[ErrorType][]GlooError)
	}
	if o.errors[errType] == nil {
		o.errors[errType] = make([]GlooError, 0)
	}
	o.errors[errType] = append(o.errors[errType], GlooError{
		err:       fmt.Errorf(msg, args...),
		errorType: errType,
		name:      wrapper.GetName(),
		namespace: wrapper.GetNamespace(),
		crdType:   fmt.Sprintf("%s/%s", wrapper.GetObjectKind().GroupVersionKind().Group, wrapper.GetObjectKind().GroupVersionKind().Kind),
	})
}

func NewGatewayAPIOutput() *GatewayAPIOutput {
	return &GatewayAPIOutput{
		gatewayAPICache: &GatewayAPICache{},
		edgeCache:       &snapshot.Instance{},
	}
}

type GatewayAPICache struct {
	YamlObjects         []*snapshot.YAMLWrapper
	HTTPRoutes          map[types.NamespacedName]*snapshot.HTTPRouteWrapper
	RouteOptions        map[types.NamespacedName]*snapshot.RouteOptionWrapper
	VirtualHostOptions  map[types.NamespacedName]*snapshot.VirtualHostOptionWrapper
	ListenerOptions     map[types.NamespacedName]*snapshot.ListenerOptionWrapper
	HTTPListenerOptions map[types.NamespacedName]*snapshot.HTTPListenerOptionWrapper
	DirectResponses     map[types.NamespacedName]*snapshot.DirectResponseWrapper

	Upstreams    map[types.NamespacedName]*snapshot.UpstreamWrapper
	AuthConfigs  map[types.NamespacedName]*snapshot.AuthConfigWrapper
	Gateways     map[types.NamespacedName]*snapshot.GatewayWrapper
	ListenerSets map[types.NamespacedName]*snapshot.ListenerSetWrapper
	Settings     map[types.NamespacedName]*snapshot.SettingsWrapper
}

func (g *GatewayAPICache) AddSettings(s *snapshot.SettingsWrapper) {
	if g.Settings == nil {
		g.Settings = make(map[types.NamespacedName]*snapshot.SettingsWrapper)
	}
	g.Settings[s.Index()] = s
}
func (g *GatewayAPICache) GetGateway(namespacedName types.NamespacedName) *snapshot.GatewayWrapper {
	if g.Gateways == nil {
		return nil
	}
	return g.Gateways[namespacedName]
}
func (g *GatewayAPICache) AddGateway(gw *snapshot.GatewayWrapper) {
	if g.Gateways == nil {
		g.Gateways = make(map[types.NamespacedName]*snapshot.GatewayWrapper)
	}
	g.Gateways[gw.Index()] = gw
}
func (g *GatewayAPICache) AddDirectResponse(d *snapshot.DirectResponseWrapper) {
	if g.DirectResponses == nil {
		g.DirectResponses = make(map[types.NamespacedName]*snapshot.DirectResponseWrapper)
	}
	g.DirectResponses[d.Index()] = d
}
func (g *GatewayAPICache) AddYAML(y *snapshot.YAMLWrapper) {
	if g.YamlObjects == nil {
		g.YamlObjects = []*snapshot.YAMLWrapper{}
	}
	g.YamlObjects = append(g.YamlObjects, y)
}

func (g *GatewayAPICache) AddHTTPRoute(route *snapshot.HTTPRouteWrapper) {
	if g.HTTPRoutes == nil {
		g.HTTPRoutes = make(map[types.NamespacedName]*snapshot.HTTPRouteWrapper)
	}
	g.HTTPRoutes[route.Index()] = route
}
func (g *GatewayAPICache) AddRouteOption(r *snapshot.RouteOptionWrapper) {
	if g.RouteOptions == nil {
		g.RouteOptions = make(map[types.NamespacedName]*snapshot.RouteOptionWrapper)
	}
	g.RouteOptions[r.Index()] = r
}
func (g *GatewayAPICache) AddVirtualHostOption(v *snapshot.VirtualHostOptionWrapper) {
	if g.VirtualHostOptions == nil {
		g.VirtualHostOptions = make(map[types.NamespacedName]*snapshot.VirtualHostOptionWrapper)
	}
	g.VirtualHostOptions[v.Index()] = v
}
func (g *GatewayAPICache) AddListenerOption(l *snapshot.ListenerOptionWrapper) {
	if g.ListenerOptions == nil {
		g.ListenerOptions = make(map[types.NamespacedName]*snapshot.ListenerOptionWrapper)
	}
	g.ListenerOptions[l.Index()] = l
}
func (g *GatewayAPICache) AddHTTPListenerOption(h *snapshot.HTTPListenerOptionWrapper) {
	if g.HTTPListenerOptions == nil {
		g.HTTPListenerOptions = make(map[types.NamespacedName]*snapshot.HTTPListenerOptionWrapper)
	}
	g.HTTPListenerOptions[h.Index()] = h
}
func (g *GatewayAPICache) AddUpstream(u *snapshot.UpstreamWrapper) {
	if g.Upstreams == nil {
		g.Upstreams = make(map[types.NamespacedName]*snapshot.UpstreamWrapper)
	}
	g.Upstreams[u.Index()] = u
}
func (g *GatewayAPICache) AddAuthConfig(a *snapshot.AuthConfigWrapper) {
	if g.AuthConfigs == nil {
		g.AuthConfigs = make(map[types.NamespacedName]*snapshot.AuthConfigWrapper)
	}
	g.AuthConfigs[a.Index()] = a
}
func (g *GatewayAPICache) AddListenerSet(l *snapshot.ListenerSetWrapper) {
	if g.ListenerSets == nil {
		g.ListenerSets = make(map[types.NamespacedName]*snapshot.ListenerSetWrapper)
	}
	g.ListenerSets[l.Index()] = l
}

func (o *GatewayAPIOutput) PreProcess(splitMatchers bool) error {

	if splitMatchers {
		if err := o.splitRouteMatchers(); err != nil {
			return err
		}
	}
	return nil
}

// we need to split the route matchers because prefix and exact matchers cause problems with rewrites
func (o *GatewayAPIOutput) splitRouteMatchers() error {
	for _, rt := range o.edgeCache.RouteTables() {
		var newRoutes []*gatewayv1.Route
		for _, route := range rt.Spec.GetRoutes() {
			editedRoute := generateRoutesForMethodMatchers(route)
			newRoutes = append(newRoutes, editedRoute)
		}
		rt.Spec.Routes = newRoutes

		o.edgeCache.AddRouteTable(rt)
	}
	return nil
}

func generateRoutesForMethodMatchers(route *gatewayv1.Route) *gatewayv1.Route {
	var newMatchers []*matchers.Matcher
	for _, m := range route.GetMatchers() {
		if len(m.GetMethods()) > 1 {
			// for each method we need to split out the matchers
			for _, method := range m.GetMethods() {
				newMatcher := &matchers.Matcher{
					PathSpecifier:   m.GetPathSpecifier(),
					CaseSensitive:   m.GetCaseSensitive(),
					Headers:         m.GetHeaders(),
					QueryParameters: m.GetQueryParameters(),
					Methods:         []string{method},
				}
				newMatchers = append(newMatchers, newMatcher)
			}
		} else {
			//it only has one so we just add it
			newMatchers = append(newMatchers, m)
		}
	}
	route.Matchers = newMatchers

	return route
}
