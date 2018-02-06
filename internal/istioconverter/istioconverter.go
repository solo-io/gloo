package istioconverter

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	routing "istio.io/api/routing/v1alpha1"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/rest"

	"github.com/solo-io/glue/internal/pkg/kube/controller"
	"github.com/solo-io/glue/pkg/log"
	clientset "github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
)

const (
	IstioIngressClass = "istio"
	istioIngressName  = "istio-ingress"
)

var (
	configDescriptor = model.ConfigDescriptor{
		model.RouteRule,
		model.DestinationPolicy,
	}
)

type istioConverter struct {
	errors chan error

	glueClient    clientset.Interface
	istioCrdStore model.ConfigStoreCache
	ingressLister v1beta1listers.IngressLister
}

func NewIstioConverter(kubeconfig, domainSuffix string, cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}) (*istioConverter, error) {
	glueClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create glue clientset: %v", err)
	}

	configClient, err := crd.NewClient(kubeconfig, configDescriptor, domainSuffix)
	if err != nil {
		return nil, multierror.Prefix(err, "failed to open a config client.")
	}
	configController := crd.NewController(configClient, kube.ControllerOptions{
		WatchedNamespace: "",
		ResyncPeriod:     resyncDuration,
		DomainSuffix:     domainSuffix,
	})

	converter := &istioConverter{
		errors:        make(chan error),
		glueClient:    glueClient,
		istioCrdStore: configController,
	}

	for _, schema := range configDescriptor {
		configController.RegisterEventHandler(schema.Type, converter.configHandler)
	}

	if err := converter.initializeIngressWatcher(cfg, resyncDuration, stopCh); err != nil {
		return nil, fmt.Errorf("failed to start ingress controller %v", err)
	}
	go configController.Run(stopCh)

	return converter, nil
}

func (c *istioConverter) Error() <-chan error {
	return c.errors
}

func (c *istioConverter) configHandler(cfg model.Config, e model.Event) {
	log.Printf("Istio Event: %v: %v", e, cfg)
}

func (c *istioConverter) ingressHandler(namespace, name string, v interface{}) {
	ingress, ok := v.(*v1beta1.Ingress)
	if !ok {
		return
	}
	// only react if it's our ingress class
	if !istioIngressClass(ingress) {
		return
	}
	log.Printf("Ingress Event: %v", v)
}

func (c *istioConverter) sync() error {
	routeRules, err := c.listRouteRules()
	if err != nil {
		return fmt.Errorf("getting route rules: %v", err)
	}
	destinationPolicies, err := c.listDestinationPolicies()
	if err != nil {
		return fmt.Errorf("getting destination policies: %v", err)
	}
	ingresses, err := c.listIngresses()
	if err != nil {
		return fmt.Errorf("getting ingresses: %v", err)
	}

}

func (c *istioConverter) generateDesiredCrds(routeRules []*routing.RouteRule,
	destinationPolicies []*routing.DestinationPolicy,
	ingresses []*v1beta1.Ingress) ([]crdv1.Upstream, []crdv1.Route, error) {

	routeRules = ingressRouteRules(routeRules)
	destinationPolicies = ingressDestPolicies(destinationPolicies)
	ingresses = istioIngresses(ingresses)

	var (
		upstreams []crdv1.Upstream
		routes    []crdv1.Route
	)
	for _, ingress := range ingresses {
		if ingress.Spec.Backend != nil {
			us, route := createDefaultResources(ingress.Name, namespace, ingress.Spec.Backend)
			upstreams = append(upstreams, us)
			routes = append(routes, route)
		}
		for _, rule := range ingress.Spec.Rules {
			ruleUpstreams, ruleRoutes := createResourcesForRule(ingress.Name, ingress.Namespace, rule)
			upstreams = append(upstreams, ruleUpstreams...)
			routes = append(routes, ruleRoutes...)
		}
	}
	return upstreams, routes, nil
}

func ingressRouteRules(rules []*routing.RouteRule) []*routing.RouteRule {
	var wanted []*routing.RouteRule
	for _, rule := range rules {
		// if source is nil, we'll apply it. if source is the ingress, we'll apply it
		if rule.Match.Source == nil || rule.Match.Source.Name == istioIngressName {
			wanted = append(wanted, rule)
		}
	}
	return wanted
}

func (c *istioConverter) listRouteRules() ([]*routing.RouteRule, error) {
	configObjects, err := c.istioCrdStore.List(model.RouteRule.Type, model.NamespaceAll)
	if err != nil {
		return nil, fmt.Errorf("istio crd list failed: %v", err)
	}
	var routeRules []*routing.RouteRule
	for _, cfg := range configObjects {
		routeRule, ok := cfg.Spec.(*routing.RouteRule)
		if !ok {
			return nil, fmt.Errorf("%v was not a route rule. what!!!???", cfg)
		}
		routeRules = append(routeRules, routeRule)
	}
	return routeRules, nil
}

func (c *istioConverter) listDestinationPolicies() ([]*routing.DestinationPolicy, error) {
	configObjects, err := c.istioCrdStore.List(model.RouteRule.Type, model.NamespaceAll)
	if err != nil {
		return nil, fmt.Errorf("istio crd list failed: %v", err)
	}
	var destinationPolicies []*routing.DestinationPolicy
	for _, cfg := range configObjects {
		destinationPolicy, ok := cfg.Spec.(*routing.DestinationPolicy)
		if !ok {
			return nil, fmt.Errorf("%v was not a destinationpolicy. what!!!???", cfg)
		}
		destinationPolicies = append(destinationPolicies, destinationPolicy)
	}
	return destinationPolicies, nil
}

func (c *istioConverter) listIngresses() ([]*v1beta1.Ingress, error) {
	ingressList, err := c.ingressLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %v", err)
	}
	return ingressList, nil
}

func istioIngressClass(ingress *v1beta1.Ingress) bool {
	return ingress.Annotations["kubernetes.io/ingress.class"] == IstioIngressClass
}

func (c *istioConverter) initializeIngressWatcher(cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}) error {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to create kube clientset: %v", err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, resyncDuration)
	ingressInformer := kubeInformerFactory.Extensions().V1beta1().Ingresses()

	kubeController := controller.NewController("istio-ingress-controller", kubeClient,
		c.ingressHandler,
		ingressInformer.Informer())

	go kubeInformerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()
	c.ingressLister = ingressInformer.Lister()
	return nil
}
