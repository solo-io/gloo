package kube

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/ingress/core/pkg/ingress/status"
	"k8s.io/ingress/core/pkg/ingress/store"

	"github.com/solo-io/glue/internal/pkg/kube/controller"
	"github.com/solo-io/glue/internal/pkg/kube/upstream"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
	clientset "github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
)

const (
	resourcePrefix    = "glue-generated"
	upstreamPrefix    = resourcePrefix + "-upstream"
	virtualHostPrefix = resourcePrefix + "-virtualHost"

	defaultVirtualHost = "default"

	GlueIngressClass = "glue"
)

type ingressConverter struct {
	errors             chan error
	useAsGlobalIngress bool
	// name of the kubernetes service for the ingress (envoy)
	ingressService string

	ingressLister v1beta1listers.IngressLister
	glueClient    clientset.Interface
}

func NewIngressConverter(cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}, useAsGlobalIngress bool, ingressNamespace, ingressService string) (*ingressConverter, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	glueClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create glue clientset: %v", err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, resyncDuration)
	ingressInformer := kubeInformerFactory.Extensions().V1beta1().Ingresses()

	ctrl := &ingressConverter{
		errors:             make(chan error),
		useAsGlobalIngress: useAsGlobalIngress,
		ingressService:     ingressService,

		ingressLister: ingressInformer.Lister(),
		glueClient:    glueClient,
	}

	kubeController := controller.NewController("glue-ingress-controller", kubeClient,
		ctrl.syncGlueResourcesWithIngresses,
		ingressInformer.Informer())

	go kubeInformerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()
	go func() {
		ctrl.syncIngressStatuses(kubeClient, ingressInformer.Informer().GetStore(), ingressNamespace, ingressService, stopCh)
	}()

	return ctrl, nil
}

const (
	ingressElectionID = "kube-ingress-importer-leader"
)

func (c *ingressConverter) syncIngressStatuses(client kubernetes.Interface,
	ingressStore cache.Store,
	ingressNamespace, ingressService string,
	stopCh <-chan struct{}) {
	sync := status.NewStatusSyncer(status.Config{
		Client:              client,
		IngressLister:       store.IngressLister{Store: ingressStore},
		ElectionID:          ingressElectionID, // TODO: configurable?
		PublishService:      ingressNamespace + "/" + ingressService,
		DefaultIngressClass: GlueIngressClass,
		IngressClass:        GlueIngressClass,
		CustomIngressStatus: func(ingress *v1beta1.Ingress) []corev1.LoadBalancerIngress { return nil },
	})
	defer sync.Shutdown()
	sync.Run(stopCh)
}

func (c *ingressConverter) syncGlueResourcesWithIngresses(namespace, name string, v interface{}) {
	ingress, ok := v.(*v1beta1.Ingress)
	if !ok {
		return
	}
	// only react if it's an ingress we care about
	if !c.isOurIngress(ingress) {
		return
	}
	log.Debugf("syncing glue config items after ingress %v/%v changed", namespace, name)
	if err := c.syncGlueResources(namespace); err != nil {
		c.errors <- err
	}
}

func (c *ingressConverter) Error() <-chan error {
	return c.errors
}

func (c *ingressConverter) syncGlueResources(namespace string) error {
	desiredUpstreams, desiredVirtualHosts, err := c.generateDesiredCrds(namespace)
	if err != nil {
		return fmt.Errorf("failed to generate desired crds: %v", err)
	}
	actualUpstreams, actualVirtualHosts, err := c.getActualCrds(namespace)
	if err != nil {
		return fmt.Errorf("failed to list actual crds: %v", err)
	}
	if err := c.syncUpstreams(namespace, desiredUpstreams, actualUpstreams); err != nil {
		return fmt.Errorf("failed to sync actual with desired upstreams: %v", err)
	}
	if err := c.syncVirtualHosts(namespace, desiredVirtualHosts, actualVirtualHosts); err != nil {
		return fmt.Errorf("failed to sync actual with desired virtualHosts: %v", err)
	}
	return nil
}

func (c *ingressConverter) syncUpstreams(namespace string, desiredUpstreams, actualUpstreams []crdv1.Upstream) error {
	var (
		upstreamsToCreate []crdv1.Upstream
		upstreamsToUpdate []crdv1.Upstream
	)
	for _, desiredUpstream := range desiredUpstreams {
		var update bool
		for i, actualUpstream := range actualUpstreams {
			if desiredUpstream.Name == actualUpstream.Name {
				// modify existing upstream
				desiredUpstream.ResourceVersion = actualUpstream.ResourceVersion
				update = true
				if !reflect.DeepEqual(desiredUpstream.Spec, actualUpstream.Spec) {
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
		if _, err := c.glueClient.GlueV1().Upstreams(namespace).Create(&us); err != nil {
			return fmt.Errorf("failed to create upstream crd %s: %v", us.Name, err)
		}
	}
	for _, us := range upstreamsToUpdate {
		if _, err := c.glueClient.GlueV1().Upstreams(namespace).Update(&us); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", us.Name, err)
		}
	}
	// only remaining are no longer desired, delete em!
	for _, us := range actualUpstreams {
		if err := c.glueClient.GlueV1().Upstreams(namespace).Delete(us.Name, nil); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", us.Name, err)
		}
	}
	return nil
}

func (c *ingressConverter) syncVirtualHosts(namespace string, desiredVirtualHosts, actualVirtualHosts []crdv1.VirtualHost) error {
	var (
		virtualHostsToCreate []crdv1.VirtualHost
		virtualHostsToUpdate []crdv1.VirtualHost
	)
	for _, desiredVirtualHost := range desiredVirtualHosts {
		var update bool
		for i, actualVirtualHost := range actualVirtualHosts {
			if desiredVirtualHost.Name == actualVirtualHost.Name {
				// modify existing virtualHost
				desiredVirtualHost.ResourceVersion = actualVirtualHost.ResourceVersion
				update = true
				if !reflect.DeepEqual(desiredVirtualHost.Spec, actualVirtualHost.Spec) {
					// only actually update if the spec has changed
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
		if _, err := c.glueClient.GlueV1().VirtualHosts(namespace).Create(&virtualHost); err != nil {
			return fmt.Errorf("failed to create virtualHost crd %s: %v", virtualHost.Name, err)
		}
	}
	for _, virtualHost := range virtualHostsToUpdate {
		if _, err := c.glueClient.GlueV1().VirtualHosts(namespace).Update(&virtualHost); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", virtualHost.Name, err)
		}
	}
	// only remaining are no longer desired, delete em!
	for _, virtualHost := range actualVirtualHosts {
		if err := c.glueClient.GlueV1().VirtualHosts(namespace).Delete(virtualHost.Name, nil); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", virtualHost.Name, err)
		}
	}
	return nil
}

func (c *ingressConverter) getActualCrds(namespace string) ([]crdv1.Upstream, []crdv1.VirtualHost, error) {
	upstreams, err := c.glueClient.GlueV1().Upstreams(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get upstream crd list: %v", err)
	}
	virtualHosts, err := c.glueClient.GlueV1().VirtualHosts(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get virtual host crd list: %v", err)
	}
	return upstreams.Items, virtualHosts.Items, nil
}

func (c *ingressConverter) generateDesiredCrds(namespace string) ([]crdv1.Upstream, []crdv1.VirtualHost, error) {
	ingressList, err := c.ingressLister.List(labels.Everything())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list ingresses: %v", err)
	}
	upstreamsByName := make(map[string]v1.Upstream)
	routesByVirtualHostName := make(map[string][]v1.Route)
	sslsByVirtualHostName := make(map[string]v1.SSLConfig)
	for _, ingress := range ingressList {
		// we only care about ingresses in the specific namespace, if namespace is given
		if namespace != "" && ingress.Namespace != namespace {
			continue
		}
		// only care if it's our ingress class, or we're the global default
		if !c.isOurIngress(ingress) {
			continue
		}
		// configure ssl for each host
		for _, tls := range ingress.Spec.TLS {
			if len(tls.Hosts) == 0 {
				sslsByVirtualHostName[defaultVirtualHost] = v1.SSLConfig{SecretRef: tls.SecretName}
			}
			for _, host := range tls.Hosts {
				sslsByVirtualHostName[host] = v1.SSLConfig{SecretRef: tls.SecretName}
			}
		}
		// default virtualhost
		if ingress.Spec.Backend != nil {
			us := newUpstreamFromBackend(ingress.Namespace, *ingress.Spec.Backend)
			if _, ok := routesByVirtualHostName[defaultVirtualHost]; ok {
				runtime.HandleError(errors.Errorf("default backend was redefined in ingress %v, ignoring", ingress.Name))
			} else {
				routesByVirtualHostName[defaultVirtualHost] = []v1.Route{
					{
						Matcher: v1.Matcher{
							Path: v1.Path{
								Prefix: "/",
							},
						},
						Destination: v1.Destination{
							SingleDestination: v1.SingleDestination{
								UpstreamDestination: &v1.UpstreamDestination{
									UpstreamName: us.Name,
								},
							},
						},
					},
				}
			}
		}
		for _, rule := range ingress.Spec.Rules {
			addRoutesAndUpstreams(ingress.Namespace, rule, upstreamsByName, routesByVirtualHostName)
		}
	}
	uniqueVirtualHosts := make(map[string]v1.VirtualHost)
	for host, routes := range routesByVirtualHostName {
		uniqueVirtualHosts[host] = v1.VirtualHost{
			Name: host,
			// kubernetes only supports a single domain per virtualhost
			Domains:   []string{host},
			Routes:    routes,
			SSLConfig: sslsByVirtualHostName[host],
		}
	}
	return upstreams, virtualHosts, nil
}

func addRoutesAndUpstreams(namespace string, rule v1beta1.IngressRule, upstreams map[string]v1.Upstream, routes map[string][]v1.Route) {
	if rule.HTTP == nil {
		return
	}
	for _, path := range rule.HTTP.Paths {
		generatedUpstream := newUpstreamFromBackend(namespace, path.Backend)
		upstreams[generatedUpstream.Name] = generatedUpstream
		host := rule.Host
		if host == "" {
			host = defaultVirtualHost
		}
		routes[rule.Host] = append(routes[rule.Host], v1.Route{
			Matcher: v1.Matcher{
				Path: v1.Path{
					Regex: path.Path,
				},
			},
			Destination: v1.Destination{
				SingleDestination: v1.SingleDestination{
					UpstreamDestination: &v1.UpstreamDestination{
						UpstreamName: generatedUpstream.Name,
					},
				},
			},
		})
	}
}

func newUpstreamFromBackend(namespace string, backend v1beta1.IngressBackend) v1.Upstream {
	return v1.Upstream{
		Name: upstreamName(namespace, backend),
		Type: upstream.Kubernetes,
		Spec: upstream.ToMap(upstream.Spec{
			ServiceName:      backend.ServiceName,
			ServiceNamespace: namespace,
			ServicePortName:  backend.ServicePort.String(),
		}),
	}
}

func upstreamName(namespace string, backend v1beta1.IngressBackend) string {
	return fmt.Sprintf("%s-%s-%s-%s", upstreamPrefix, namespace, backend.ServiceName, backend.ServicePort.String())
}

func (c *ingressConverter) isOurIngress(ingress *v1beta1.Ingress) bool {
	return c.useAsGlobalIngress || ingress.Annotations["kubernetes.io/ingress.class"] == GlueIngressClass
}

func pathToName(path string) string {
	hash := md5.New()
	hash.Write([]byte(path))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
