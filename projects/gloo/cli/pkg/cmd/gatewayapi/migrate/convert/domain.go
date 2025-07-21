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
	DisableListenerSets     bool
}

func (o *Options) Validate() error {

	count := 0
	if o.InputDir != "" {
		count++
	}
	if o.InputFile != "" {
		count++
	}
	if o.GlooSnapshotFile != "" {
		count++
	}
	if o.ControlPlaneName != "" {
		count++
		if o.ControlPlaneNamespace == "" {
			return fmt.Errorf("pod namespace must be specified")
		}
	}

	if count > 1 {
		return fmt.Errorf("only one of 'input-file' or 'directory' or 'input-snapshot' or `gloo-pod-name` can be specified")
	}
	if !o.DeleteOutputDir && FolderExists(o.OutputDir) {
		return fmt.Errorf("output-dir already %s exists. It can be deleted with --delete-output-dir", o.OutputDir)
	}
	return nil
}
func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func (o *Options) AddToFlags(flags *pflag.FlagSet) {
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
	flags.BoolVar(&o.DisableListenerSets, "disable-listenersets", false, "Do not create listenersets and bind hosts and ports directly to Gateway.")

}

type GatewayAPIOutput struct {
	gatewayAPICache *GatewayAPICache
	edgeCache       *snapshot.Instance
	errors          map[ErrorType][]GlooError
}

func (g *GatewayAPIOutput) SetGatewayAPICache(cache *GatewayAPICache) {
	g.gatewayAPICache = cache
}

func (g *GatewayAPIOutput) SetEdgeCache(cache *snapshot.Instance) {
	g.edgeCache = cache
}

func (g *GatewayAPIOutput) GetGatewayAPICache() *GatewayAPICache {
	if g.gatewayAPICache == nil {
		g.gatewayAPICache = &GatewayAPICache{}
	}
	return g.gatewayAPICache
}

func (g *GatewayAPIOutput) GetEdgeCache() *snapshot.Instance {
	if g.edgeCache == nil {
		g.edgeCache = &snapshot.Instance{}
	}
	return g.edgeCache
}
func (g *GatewayAPIOutput) EdgeCache(instance *snapshot.Instance) {
	g.edgeCache = instance
}

type GlooError struct {
	err       error
	errorType ErrorType
	name      string
	namespace string
	crdType   string
}

func (g *GatewayAPIOutput) AddError(errType ErrorType, msg string, args ...interface{}) {
	if g.errors == nil {
		g.errors = make(map[ErrorType][]GlooError)
	}
	if g.errors[errType] == nil {
		g.errors[errType] = make([]GlooError, 0)
	}

	g.errors[errType] = append(g.errors[errType], GlooError{
		err:       fmt.Errorf(msg, args...),
		errorType: errType,
		name:      "none",
		namespace: "none",
		crdType:   "none",
	})
}
func (g *GatewayAPIOutput) AddErrorFromWrapper(errType ErrorType, wrapper snapshot.Wrapper, msg string, args ...interface{}) {
	if g.errors == nil {
		g.errors = make(map[ErrorType][]GlooError)
	}
	if g.errors[errType] == nil {
		g.errors[errType] = make([]GlooError, 0)
	}
	g.errors[errType] = append(g.errors[errType], GlooError{
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
	YamlObjects  []*snapshot.YAMLWrapper
	Gateways     map[types.NamespacedName]*snapshot.GatewayWrapper
	ListenerSets map[types.NamespacedName]*snapshot.ListenerSetWrapper
	HTTPRoutes   map[types.NamespacedName]*snapshot.HTTPRouteWrapper

	AuthConfigs          map[types.NamespacedName]*snapshot.AuthConfigWrapper
	DirectResponses      map[types.NamespacedName]*snapshot.DirectResponseWrapper
	Backends             map[types.NamespacedName]*snapshot.BackendWrapper
	HTTPListenerPolicies map[types.NamespacedName]*snapshot.HTTPListenerPolicyWrapper
	GlooTrafficPolicies  map[types.NamespacedName]*snapshot.GlooTrafficPolicyWrapper
	GatewayExtensions    map[types.NamespacedName]*snapshot.GatewayExtensionWrapper
	KGatewayParameters   map[types.NamespacedName]*snapshot.GatewayParametersWrapper
	BackendConfigPolicy  map[types.NamespacedName]*snapshot.BackendConfigPolicyWrapper
}

func (g *GatewayAPICache) AddBackendConfigPolicy(policy *snapshot.BackendConfigPolicyWrapper) {
	if g.BackendConfigPolicy == nil {
		g.BackendConfigPolicy = make(map[types.NamespacedName]*snapshot.BackendConfigPolicyWrapper)
	}
	g.BackendConfigPolicy[policy.Index()] = policy
}

func (g *GatewayAPICache) AddHTTPListenerPolicy(policy *snapshot.HTTPListenerPolicyWrapper) {
	if g.HTTPListenerPolicies == nil {
		g.HTTPListenerPolicies = make(map[types.NamespacedName]*snapshot.HTTPListenerPolicyWrapper)
	}
	g.HTTPListenerPolicies[policy.Index()] = policy
}

func (g *GatewayAPICache) AddBackend(b *snapshot.BackendWrapper) {
	if g.Backends == nil {
		g.Backends = make(map[types.NamespacedName]*snapshot.BackendWrapper)
	}
	g.Backends[b.Index()] = b
}
func (g *GatewayAPICache) GetBackend(namespacedName types.NamespacedName) *snapshot.BackendWrapper {
	if g.Backends == nil {
		return nil
	}
	return g.Backends[namespacedName]
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
func (g *GatewayAPICache) AddGatewayExtension(gew *snapshot.GatewayExtensionWrapper) {
	if g.GatewayExtensions == nil {
		g.GatewayExtensions = make(map[types.NamespacedName]*snapshot.GatewayExtensionWrapper)
	}
	g.GatewayExtensions[gew.Index()] = gew
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
func (g *GatewayAPICache) AddListenerSet(l *snapshot.ListenerSetWrapper) {
	if g.ListenerSets == nil {
		g.ListenerSets = make(map[types.NamespacedName]*snapshot.ListenerSetWrapper)
	}
	g.ListenerSets[l.Index()] = l
}
func (g *GatewayAPICache) AddAuthConfig(a *snapshot.AuthConfigWrapper) {
	if g.AuthConfigs == nil {
		g.AuthConfigs = make(map[types.NamespacedName]*snapshot.AuthConfigWrapper)
	}
	g.AuthConfigs[a.Index()] = a
}
func (g *GatewayAPICache) AddGlooTrafficPolicy(gtp *snapshot.GlooTrafficPolicyWrapper) {
	if g.GlooTrafficPolicies == nil {
		g.GlooTrafficPolicies = make(map[types.NamespacedName]*snapshot.GlooTrafficPolicyWrapper)
	}
	g.GlooTrafficPolicies[gtp.Index()] = gtp
}

func (g *GatewayAPIOutput) PreProcess(splitMatchers bool) error {

	if splitMatchers {
		if err := g.splitRouteMatchers(); err != nil {
			return err
		}
	}
	return nil
}

// we need to split the route matchers because prefix and exact matchers cause problems with rewrites
func (g *GatewayAPIOutput) splitRouteMatchers() error {
	for _, rt := range g.edgeCache.RouteTables() {
		var newRoutes []*gatewayv1.Route
		for _, route := range rt.Spec.GetRoutes() {
			editedRoute := generateRoutesForMethodMatchers(route)
			newRoutes = append(newRoutes, editedRoute)
		}
		rt.Spec.Routes = newRoutes

		g.edgeCache.AddRouteTable(rt)
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
