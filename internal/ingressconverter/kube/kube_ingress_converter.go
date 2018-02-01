package kube

import (
	"fmt"
	"time"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/rest"

	clientset "github.com/solo-io/glue/internal/configwatcher/kube/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue/internal/configwatcher/kube/crd/solo.io/v1"
	"github.com/solo-io/glue/internal/pkg/kube/controller"
	"github.com/solo-io/glue/internal/pkg/kube/upstream"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
)

const (
	resourcePrefix = "glue-generated"
	upstreamPrefix = resourcePrefix + "-upstream"
	routePrefix    = resourcePrefix + "-route"

	defaultRouteWeight = 10000
)

type ingressConverter struct {
	errors chan error

	ingressLister v1beta1listers.IngressLister
	glueClient    clientset.Interface
}

func newIngressConverter(cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}) (*ingressConverter, error) {
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
		errors:        make(chan error),
		ingressLister: ingressInformer.Lister(),
		glueClient:    glueClient,
	}

	kubeController := controller.NewController("glue-ingress-controller", kubeClient,
		ctrl.syncGlueResourcesWithIngress,
		ingressInformer.Informer())

	go kubeInformerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()

	return ctrl, nil
}

func (c *ingressConverter) syncGlueResourcesWithIngress(namespace, name string, _ interface{}) {
	if err := func() error {
		log.Printf("syncing glue config items after ingress %v/%v changed", namespace, name)
		ingressList, err := c.ingressLister.List(labels.Everything())
		if err != nil {
			return fmt.Errorf("failed to list ingresses: %v", err)
		}
		var (
			upstreams []v1.Upstream
			routes    []v1.Route
		)
		for _, ingress := range ingressList {
			if ingress.Spec.Backend != nil {
				us, route := createDefaultRoute(ingress.Name, ingress.Namespace, ingress.Spec.Backend)
				upstreams = append(upstreams, us)
				routes = append(routes, route)
			}
			for _, rule := range ingress.Spec.Rules {
				ruleUpstreams, ruleRoutes := createResourcesForRule(ingress.Name, ingress.Namespace, rule)
				upstreams = append(upstreams, ruleUpstreams...)
				routes = append(routes, ruleRoutes...)
			}
		}
		return c.writeCrds(namespace, upstreams, routes)
	}(); err != nil {
		c.errors <- err
	}
}

func (c *ingressConverter) Error() <-chan error {
	return c.errors
}

func (c *ingressConverter) writeCrds(namespace string, upstreams []v1.Upstream, routes []v1.Route) error {
	for _, us := range upstreams {
		crdUpstream := &crdv1.Upstream{
			ObjectMeta: metav1.ObjectMeta{
				Name:      us.Name,
				Namespace: namespace,
			},
			Spec: crdv1.DeepCopyUpstream(us),
		}
		if _, err := c.glueClient.GlueV1().Upstreams(namespace).Create(crdUpstream); err != nil {
			return fmt.Errorf("failed to create upstream crd %s: %v", crdUpstream.Name, err)
		}
	}
	for _, route := range routes {
		crdRoute := &crdv1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Name:      routeName(route),
				Namespace: namespace,
			},
			Spec: crdv1.DeepCopyRoute(route),
		}
		if _, err := c.glueClient.GlueV1().Routes(namespace).Create(crdRoute); err != nil {
			return fmt.Errorf("failed to create route crd %s: %v", crdRoute.Name, err)
		}
	}
	return nil
}

func createDefaultRoute(ingressName string, ingressNamespace string, backend *v1beta1.IngressBackend) (v1.Upstream, v1.Route) {
	us := v1.Upstream{
		Name: upstreamName(ingressName, *backend),
		Type: upstream.Kubernetes,
		Spec: upstream.ToMap(upstream.Spec{
			ServiceName:      backend.ServiceName,
			ServiceNamespace: ingressNamespace,
			ServicePortName:  portName(backend.ServicePort),
		}),
	}

	route := v1.Route{
		Matcher: v1.Matcher{
			Path: v1.Path{
				Prefix: "/",
			},
		},
		Destination: v1.Destination{
			UpstreamDestination: &v1.UpstreamDestination{
				UpstreamName: us.Name,
			},
		},
		Weight: defaultRouteWeight,
	}
	return us, route
}

func createResourcesForRule(ingressName string, ingressNamespace string, rule v1beta1.IngressRule) ([]v1.Upstream, []v1.Route) {
	var (
		upstreams []v1.Upstream
		routes    []v1.Route
	)
	host := rule.Host

	for i, path := range rule.IngressRuleValue.HTTP.Paths {
		pathRegex := path.Path
		if pathRegex == "" {
			pathRegex = "/"
		}
		us := v1.Upstream{
			Name: upstreamName(ingressName, path.Backend),
			Type: upstream.Kubernetes,
			Spec: upstream.ToMap(upstream.Spec{
				ServiceName:      path.Backend.ServiceName,
				ServiceNamespace: ingressNamespace,
				ServicePortName:  portName(path.Backend.ServicePort),
			}),
		}
		route := v1.Route{
			Matcher: v1.Matcher{
				Path: v1.Path{
					Regex: pathRegex,
				},
				VirtualHost: host,
			},
			Destination: v1.Destination{
				UpstreamDestination: &v1.UpstreamDestination{
					UpstreamName: us.Name,
				},
			},
			Weight: len(rule.IngressRuleValue.HTTP.Paths) - i,
		}
		upstreams = append(upstreams, us)
		routes = append(routes, route)
	}
	return upstreams, routes
}

func portName(portVal intstr.IntOrString) string {
	if portVal.Type == intstr.String {
		return portVal.StrVal
	}
	return fmt.Sprintf("%s", portVal.IntVal)
}

func upstreamName(ingressName string, backend v1beta1.IngressBackend) string {
	return fmt.Sprintf("%s-%s-%s-%s", upstreamPrefix, ingressName, backend.ServiceName, portName(backend.ServicePort))
}

func routeName(route v1.Route) string {
	return fmt.Sprintf("%s%s%s", routePrefix, route.Matcher.VirtualHost, route.Destination.UpstreamDestination.UpstreamName)
}
