package kube

import (
	"crypto/md5"
	"fmt"
	"reflect"
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

	GlueIngressClass = "glue"
)

type ingressConverter struct {
	errors chan error

	ingressLister v1beta1listers.IngressLister
	glueClient    clientset.Interface
}

func NewIngressConverter(cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}) (*ingressConverter, error) {
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
		ctrl.syncGlueResourcesWithIngresses,
		ingressInformer.Informer())

	go kubeInformerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()

	return ctrl, nil
}

func (c *ingressConverter) syncGlueResourcesWithIngresses(namespace, name string, v interface{}) {
	ingress, ok := v.(*v1beta1.Ingress)
	if !ok {
		return
	}
	// only react if it's our ingress class
	if !glueIngressClass(ingress) {
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
	desiredUpstreams, desiredRoutes, err := c.generateDesiredCrds(namespace)
	if err != nil {
		return fmt.Errorf("failed to generate desired crds: %v", err)
	}
	actualUpstreams, actualRoutes, err := c.getActualCrds(namespace)
	if err != nil {
		return fmt.Errorf("failed to list actual crds: %v", err)
	}
	if err := c.syncUpstreams(namespace, desiredUpstreams, actualUpstreams); err != nil {
		return fmt.Errorf("failed to sync actual with desired upstreams: %v", err)
	}
	if err := c.syncRoutes(namespace, desiredRoutes, actualRoutes); err != nil {
		return fmt.Errorf("failed to sync actual with desired routes: %v", err)
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

func (c *ingressConverter) syncRoutes(namespace string, desiredRoutes, actualRoutes []crdv1.Route) error {
	var (
		routesToCreate []crdv1.Route
		routesToUpdate []crdv1.Route
	)
	for _, desiredRoute := range desiredRoutes {
		var update bool
		for i, actualRoute := range actualRoutes {
			if desiredRoute.Name == actualRoute.Name {
				// modify existing route
				desiredRoute.ResourceVersion = actualRoute.ResourceVersion
				update = true
				if !reflect.DeepEqual(desiredRoute.Spec, actualRoute.Spec) {
					// only actually update if the spec has changed
					routesToUpdate = append(routesToUpdate, desiredRoute)
				}
				// remove it from the list we match against
				actualRoutes = append(actualRoutes[:i], actualRoutes[i+1:]...)
				break
			}
		}
		if !update {
			// desired was not found, mark for creation
			routesToCreate = append(routesToCreate, desiredRoute)
		}
	}
	for _, route := range routesToCreate {
		if _, err := c.glueClient.GlueV1().Routes(namespace).Create(&route); err != nil {
			return fmt.Errorf("failed to create route crd %s: %v", route.Name, err)
		}
	}
	for _, route := range routesToUpdate {
		if _, err := c.glueClient.GlueV1().Routes(namespace).Update(&route); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", route.Name, err)
		}
	}
	// only remaining are no longer desired, delete em!
	for _, route := range actualRoutes {
		if err := c.glueClient.GlueV1().Routes(namespace).Delete(route.Name, nil); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", route.Name, err)
		}
	}
	return nil
}

func (c *ingressConverter) getActualCrds(namespace string) ([]crdv1.Upstream, []crdv1.Route, error) {
	upstreams, err := c.glueClient.GlueV1().Upstreams(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get upstream crd list: %v", err)
	}
	routes, err := c.glueClient.GlueV1().Routes(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get route crd list: %v", err)
	}
	return upstreams.Items, routes.Items, nil
}

func (c *ingressConverter) generateDesiredCrds(namespace string) ([]crdv1.Upstream, []crdv1.Route, error) {
	ingressList, err := c.ingressLister.List(labels.Everything())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list ingresses: %v", err)
	}
	var (
		upstreams []crdv1.Upstream
		routes    []crdv1.Route
	)
	for _, ingress := range ingressList {
		// we only care about ingresses in the specific namespace
		if ingress.Namespace != namespace {
			continue
		}
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

func createDefaultResources(ingressName string, namespace string, backend *v1beta1.IngressBackend) (crdv1.Upstream, crdv1.Route) {
	us := v1.Upstream{
		Name: upstreamName(ingressName, *backend),
		Type: upstream.Kubernetes,
		Spec: upstream.ToMap(upstream.Spec{
			ServiceName:      backend.ServiceName,
			ServiceNamespace: namespace,
			ServicePortName:  portName(backend.ServicePort),
		}),
	}
	usCrd := crdv1.Upstream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      us.Name,
			Namespace: namespace,
		},
		Spec: crdv1.DeepCopyUpstream(us),
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
	routeCrd := crdv1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      routeName(route),
			Namespace: namespace,
		},
		Spec: crdv1.DeepCopyRoute(route),
	}

	return usCrd, routeCrd
}

func createResourcesForRule(ingressName string, namespace string, rule v1beta1.IngressRule) ([]crdv1.Upstream, []crdv1.Route) {
	var (
		upstreams []crdv1.Upstream
		routes    []crdv1.Route
	)
	host := rule.Host

	uniqueUpstreams := make(map[string]crdv1.Upstream)
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
				ServiceNamespace: namespace,
				ServicePortName:  portName(path.Backend.ServicePort),
			}),
		}
		usCrd := crdv1.Upstream{
			ObjectMeta: metav1.ObjectMeta{
				Name:      us.Name,
				Namespace: namespace,
			},
			Spec: crdv1.DeepCopyUpstream(us),
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
		routeCrd := crdv1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Name:      routeName(route),
				Namespace: namespace,
			},
			Spec: crdv1.DeepCopyRoute(route),
		}
		uniqueUpstreams[usCrd.Name] = usCrd
		routes = append(routes, routeCrd)
	}
	for _, usCrd := range uniqueUpstreams {
		upstreams = append(upstreams, usCrd)
	}
	return upstreams, routes
}

func portName(portVal intstr.IntOrString) string {
	if portVal.Type == intstr.String {
		return portVal.StrVal
	}
	return fmt.Sprintf("%d", portVal.IntVal)
}

func upstreamName(ingressName string, backend v1beta1.IngressBackend) string {
	return fmt.Sprintf("%s-%s-%s-%s", upstreamPrefix, ingressName, backend.ServiceName, portName(backend.ServicePort))
}

func routeName(route v1.Route) string {
	var pathName string
	if regex := route.Matcher.Path.Regex; regex != "" {
		pathName = regex
	}
	if prefix := route.Matcher.Path.Prefix; prefix != "" {
		pathName = prefix
	}
	if exact := route.Matcher.Path.Exact; exact != "" {
		pathName = exact
	}
	return fmt.Sprintf("%s-%s%s", routePrefix, route.Matcher.VirtualHost, pathToName(pathName))
}

func glueIngressClass(ingress *v1beta1.Ingress) bool {
	return ingress.Annotations["kubernetes.io/ingress.class"] == GlueIngressClass
}

func pathToName(path string) string {
	hash := md5.New()
	hash.Write([]byte(path))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
