package ingress

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/glue-storage"
	kubeplugin "github.com/solo-io/glue/internal/plugins/kubernetes"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/kubecontroller"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

const (
	resourcePrefix    = "glue-generated"
	upstreamPrefix    = resourcePrefix + "-upstream"
	virtualHostPrefix = resourcePrefix + "-virtualhost"

	defaultVirtualHost = "default"

	GlueIngressClass = "glue"
)

type ingressController struct {
	errors             chan error
	useAsGlobalIngress bool

	ingressLister v1beta1listers.IngressLister
	configObjects storage.Interface
	runFunc       func(stop <-chan struct{})
}

func NewIngressController(cfg *rest.Config,
	configStore storage.Interface,
	resyncDuration time.Duration,
	useAsGlobalIngress bool) (*ingressController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	// attempt to register configObjects if they don't exist
	if err := configStore.V1().Register(); err != nil {
		return nil, errors.Wrap(err, "failed to register configObjects")
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, resyncDuration)
	ingressInformer := kubeInformerFactory.Extensions().V1beta1().Ingresses()

	c := &ingressController{
		errors:             make(chan error),
		useAsGlobalIngress: useAsGlobalIngress,

		ingressLister: ingressInformer.Lister(),
		configObjects: configStore,
	}

	kubeController := kubecontroller.NewController("glue-ingress-controller", kubeClient,
		kubecontroller.NewSyncHandler(c.syncGlueResourcesWithIngresses),
		ingressInformer.Informer())

	c.runFunc = func(stop <-chan struct{}) {
		wg := &sync.WaitGroup{}
		go func(stop <-chan struct{}) {
			wg.Add(1)
			kubeInformerFactory.Start(stop)
			wg.Done()
		}(stop)
		go func() {
			kubeController.Run(2, stop)
		}()
		wg.Wait()
	}

	return c, nil
}

func (c *ingressController) Run(stop <-chan struct{}) {
	c.runFunc(stop)
}

func (c *ingressController) Error() <-chan error {
	return c.errors
}

func (c *ingressController) syncGlueResourcesWithIngresses() {
	if err := c.syncGlueResources(); err != nil {
		c.errors <- err
	}
}

func (c *ingressController) syncGlueResources() error {
	desiredUpstreams, desiredVirtualHosts, err := c.generateDesiredResources()
	if err != nil {
		return fmt.Errorf("failed to generate desired configObjects: %v", err)
	}
	actualUpstreams, actualVirtualHosts, err := c.getActualResources()
	if err != nil {
		return fmt.Errorf("failed to list actual configObjects: %v", err)
	}
	if err := c.syncUpstreams(desiredUpstreams, actualUpstreams); err != nil {
		return fmt.Errorf("failed to sync actual with desired upstreams: %v", err)
	}
	if err := c.syncVirtualHosts(desiredVirtualHosts, actualVirtualHosts); err != nil {
		return fmt.Errorf("failed to sync actual with desired virtualHosts: %v", err)
	}
	return nil
}

func (c *ingressController) getActualResources() ([]*v1.Upstream, []*v1.VirtualHost, error) {
	upstreams, err := c.configObjects.V1().Upstreams().List()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get upstream crd list: %v", err)
	}
	virtualHosts, err := c.configObjects.V1().VirtualHosts().List()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get virtual host crd list: %v", err)
	}
	return upstreams, virtualHosts, nil
}

func (c *ingressController) generateDesiredResources() ([]*v1.Upstream, []*v1.VirtualHost, error) {
	ingressList, err := c.ingressLister.List(labels.Everything())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list ingresses: %v", err)
	}
	upstreamsByName := make(map[string]*v1.Upstream)
	routesByHostName := make(map[string][]*v1.Route)
	sslsByHostName := make(map[string]*v1.SSLConfig)
	// deterministic way to sort ingresses
	sort.SliceStable(ingressList, func(i, j int) bool {
		return strings.Compare(ingressList[i].Name, ingressList[j].Name) > 0
	})
	for _, ingress := range ingressList {
		// only care if it's our ingress class, or we're the global default
		if !isOurIngress(c.useAsGlobalIngress, ingress) {
			continue
		}
		// configure ssl for each host
		for _, tls := range ingress.Spec.TLS {
			if len(tls.Hosts) == 0 {
				sslsByHostName[defaultVirtualHost] = &v1.SSLConfig{SecretRef: tls.SecretName}
			}
			for _, host := range tls.Hosts {
				sslsByHostName[host] = &v1.SSLConfig{SecretRef: tls.SecretName}
			}
		}
		// default virtualhost
		if ingress.Spec.Backend != nil {
			us, err := newUpstreamFromBackend(ingress.Namespace, *ingress.Spec.Backend)
			if err != nil {
				return nil, nil, errors.Wrap(err, "internal err: failed to backend to upstream")
			}
			if _, ok := routesByHostName[defaultVirtualHost]; ok {
				runtime.HandleError(errors.Errorf("default backend was redefined in ingress %v, ignoring", ingress.Name))
			} else {
				routesByHostName[defaultVirtualHost] = []*v1.Route{
					{
						Matcher: &v1.Matcher{
							Path: &v1.Matcher_PathPrefix{
								PathPrefix: "/",
							},
						},
						SingleDestination: &v1.Destination{
							DestinationType: &v1.Destination_Upstream{
								Upstream: &v1.UpstreamDestination{
									Name: us.Name,
								},
							},
						},
					},
				}
			}
		}
		for _, rule := range ingress.Spec.Rules {
			addRoutesAndUpstreams(ingress.Namespace, rule, upstreamsByName, routesByHostName)
		}
	}
	uniqueVirtualHosts := make(map[string]*v1.VirtualHost)
	for host, routes := range routesByHostName {
		// sort routes by path length
		// equal length sorted by string compare
		// longest routes should come first
		sortRoutes(routes)
		// TODO: evaluate
		// set default virtualhost to match *
		domains := []string{host}
		if host == defaultVirtualHost {
			domains[0] = "*"
		}
		uniqueVirtualHosts[host] = &v1.VirtualHost{
			Name: host,
			// kubernetes only supports a single domain per virtualhost
			Domains:   domains,
			Routes:    routes,
			SslConfig: sslsByHostName[host],
		}
	}
	var (
		upstreams    []*v1.Upstream
		virtualHosts []*v1.VirtualHost
	)
	for _, us := range upstreamsByName {
		upstreams = append(upstreams, us)
	}
	for name, virtualHost := range uniqueVirtualHosts {
		if name != defaultVirtualHost {
			name = fmt.Sprintf("%s-%s", virtualHostPrefix, name)
		}
		virtualHosts = append(virtualHosts, virtualHost)
	}
	return upstreams, virtualHosts, nil
}

func getPathStr(route *v1.Route) string {
	switch path := route.Matcher.Path.(type) {
	case *v1.Matcher_PathPrefix:
		return path.PathPrefix
	case *v1.Matcher_PathRegex:
		return path.PathRegex
	case *v1.Matcher_PathExact:
		return path.PathExact
	}
	return ""
}

func sortRoutes(routes []*v1.Route) {
	sort.SliceStable(routes, func(i, j int) bool {
		p1 := getPathStr(routes[i])
		p2 := getPathStr(routes[j])
		l1 := len(p1)
		l2 := len(p2)
		if l1 == l2 {
			return strings.Compare(p1, p2) < 0
		}
		// longer = comes first
		return l1 > l2
	})
}

func (c *ingressController) syncUpstreams(desiredUpstreams, actualUpstreams []*v1.Upstream) error {
	var (
		upstreamsToCreate []*v1.Upstream
		upstreamsToUpdate []*v1.Upstream
	)
	for _, desiredUpstream := range desiredUpstreams {
		var update bool
		for i, actualUpstream := range actualUpstreams {
			if desiredUpstream.Name == actualUpstream.Name {
				// modify existing upstream
				desiredUpstream.Metadata = actualUpstream.GetMetadata()
				update = true
				if !desiredUpstream.Equal(actualUpstream) {
					// only actually update if the spec has changed
					upstreamsToUpdate = append(upstreamsToUpdate, desiredUpstream)
				}
				// remove it from the list we match against
				actualUpstreams = append(actualUpstreams[:i], actualUpstreams[i+1:]...)
				break
			}
		}
		if !update {
			// desired was not found, mark for creation
			upstreamsToCreate = append(upstreamsToCreate, desiredUpstream)
		}
	}
	for _, us := range upstreamsToCreate {
		if _, err := c.configObjects.V1().Upstreams().Create(us); err != nil {
			return fmt.Errorf("failed to create upstream crd %s: %v", us.Name, err)
		}
	}
	for _, us := range upstreamsToUpdate {
		if _, err := c.configObjects.V1().Upstreams().Update(us); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", us.Name, err)
		}
	}
	// only remaining are no longer desired, delete em!
	for _, us := range actualUpstreams {
		if err := c.configObjects.V1().Upstreams().Delete(us.Name); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", us.Name, err)
		}
	}
	return nil
}

func (c *ingressController) syncVirtualHosts(desiredVirtualHosts, actualVirtualHosts []*v1.VirtualHost) error {
	var (
		virtualHostsToCreate []*v1.VirtualHost
		virtualHostsToUpdate []*v1.VirtualHost
	)
	for _, desiredVirtualHost := range desiredVirtualHosts {
		var update bool
		for i, actualVirtualHost := range actualVirtualHosts {
			if desiredVirtualHost.Name == actualVirtualHost.Name {
				// modify existing virtualHost
				desiredVirtualHost.Metadata = actualVirtualHost.GetMetadata()
				update = true
				// only actually update if the spec has changed
				if !desiredVirtualHost.Equal(actualVirtualHost) {
					virtualHostsToUpdate = append(virtualHostsToUpdate, desiredVirtualHost)
				}
				// remove it from the list we match against
				actualVirtualHosts = append(actualVirtualHosts[:i], actualVirtualHosts[i+1:]...)
				break
			}
		}
		if !update {
			// desired was not found, mark for creation
			virtualHostsToCreate = append(virtualHostsToCreate, desiredVirtualHost)
		}
	}
	for _, virtualHost := range virtualHostsToCreate {
		if _, err := c.configObjects.V1().VirtualHosts().Create(virtualHost); err != nil {
			return fmt.Errorf("failed to create virtualHost crd %s: %v", virtualHost.Name, err)
		}
	}
	for _, virtualHost := range virtualHostsToUpdate {
		if _, err := c.configObjects.V1().VirtualHosts().Update(virtualHost); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", virtualHost.Name, err)
		}
	}
	// only remaining are no longer desired, delete em!
	for _, virtualHost := range actualVirtualHosts {
		if err := c.configObjects.V1().VirtualHosts().Delete(virtualHost.Name); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", virtualHost.Name, err)
		}
	}
	return nil
}

func addRoutesAndUpstreams(namespace string, rule v1beta1.IngressRule, upstreams map[string]*v1.Upstream, routes map[string][]*v1.Route) {
	if rule.HTTP == nil {
		return
	}
	for _, path := range rule.HTTP.Paths {
		generatedUpstream, err := newUpstreamFromBackend(namespace, path.Backend)
		if err != nil {
			continue
		}
		upstreams[generatedUpstream.Name] = generatedUpstream
		host := rule.Host
		if host == "" {
			host = defaultVirtualHost
		}
		routes[rule.Host] = append(routes[rule.Host], &v1.Route{
			Matcher: &v1.Matcher{
				Path: &v1.Matcher_PathRegex{
					PathRegex: path.Path,
				},
			},
			SingleDestination: &v1.Destination{
				DestinationType: &v1.Destination_Upstream{
					Upstream: &v1.UpstreamDestination{
						Name: generatedUpstream.Name,
					},
				},
			},
		})
	}
}

func newUpstreamFromBackend(namespace string, backend v1beta1.IngressBackend) (*v1.Upstream, error) {
	spec, err := kubeplugin.EncodeUpstreamSpec(kubeplugin.UpstreamSpec{
		ServiceName:      backend.ServiceName,
		ServiceNamespace: namespace,
		ServicePort:      backend.ServicePort.String(),
	})
	return &v1.Upstream{
		Name: upstreamName(namespace, backend),
		Type: kubeplugin.UpstreamTypeKube,
		Spec: spec,
	}, err
}

func upstreamName(namespace string, backend v1beta1.IngressBackend) string {
	return fmt.Sprintf("%s-%s-%s-%s", upstreamPrefix, namespace, backend.ServiceName, backend.ServicePort.String())
}

func isOurIngress(useAsGlobalIngress bool, ingress *v1beta1.Ingress) bool {
	return useAsGlobalIngress || ingress.Annotations["kubernetes.io/ingress.class"] == GlueIngressClass
}
