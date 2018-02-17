package controller

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	kubeplugin "github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/kubecontroller"
)

const (
	resourcePrefix = "gloo-generated"
	upstreamPrefix = resourcePrefix + "-upstream"
)

type ServiceController struct {
	errors             chan error
	useAsGlobalIngress bool

	serviceLister kubelisters.ServiceLister
	upstreams     storage.Interface
	runFunc       func(stop <-chan struct{})
}

func NewServiceController(cfg *rest.Config,
	configStore storage.Interface,
	resyncDuration time.Duration,
	useAsGlobalIngress bool) (*ServiceController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	// attempt to register upstreams if they don't exist
	if err := configStore.V1().Register(); err != nil {
		return nil, errors.Wrap(err, "failed to register upstreams")
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, resyncDuration)
	serviceInformer := kubeInformerFactory.Core().V1().Services()

	c := &ServiceController{
		errors:             make(chan error),
		useAsGlobalIngress: useAsGlobalIngress,

		serviceLister: serviceInformer.Lister(),
		upstreams:     configStore,
	}

	kubeController := kubecontroller.NewController("gloo-ingress-controller", kubeClient,
		kubecontroller.NewSyncHandler(c.syncGlooUpstreamsWithKubeServices),
		serviceInformer.Informer())

	c.runFunc = func(stop <-chan struct{}) {
		go kubeInformerFactory.Start(stop)
		go kubeController.Run(2, stop)
		<-stop
		log.Printf("ingress controller stopped")
	}

	return c, nil
}

func (c *ServiceController) Run(stop <-chan struct{}) {
	c.runFunc(stop)
}

func (c *ServiceController) Error() <-chan error {
	return c.errors
}

func (c *ServiceController) syncGlooUpstreamsWithKubeServices() {
	if err := c.syncGlooUpstreams(); err != nil {
		c.errors <- err
	}
}

func (c *ServiceController) syncGlooUpstreams() error {
	desiredUpstreams, err := c.generateDesiredUpstreams()
	if err != nil {
		return fmt.Errorf("failed to generate desired upstreams: %v", err)
	}
	actualUpstreams, err := c.getActualUpstreams()
	if err != nil {
		return fmt.Errorf("failed to list actual upstreams: %v", err)
	}
	if err := c.syncUpstreams(desiredUpstreams, actualUpstreams); err != nil {
		return fmt.Errorf("failed to sync actual with desired upstreams: %v", err)
	}
	return nil
}

func (c *ServiceController) getActualUpstreams() ([]*v1.Upstream, error) {
	upstreams, err := c.upstreams.V1().Upstreams().List()
	if err != nil {
		return nil, fmt.Errorf("failed to get upstream crd list: %v", err)
	}
	return upstreams, nil
}

func (c *ServiceController) generateDesiredUpstreams() ([]*v1.Upstream, error) {
	serviceList, err := c.serviceLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %v", err)
	}
	var upstreams []*v1.Upstream
	for _, svc := range serviceList {
		for _, port := range svc.Spec.Ports {
			spec, err := kubeplugin.EncodeUpstreamSpec(kubeplugin.UpstreamSpec{
				ServiceNamespace: svc.Namespace,
				ServiceName:      svc.Name,
				ServicePort:      port.String(),
			})
			if err != nil {
				runtime.HandleError(err)
				continue
			}
			upstream := &v1.Upstream{
				Name: upstreamName(svc.Namespace, svc.Name, port.String()),
				Type: kubeplugin.UpstreamTypeKube,
				Spec: spec,
			}
			upstreams = append(upstreams, upstream)
		}
	}
	return upstreams, nil
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

func (c *ServiceController) syncUpstreams(desiredUpstreams, actualUpstreams []*v1.Upstream) error {
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
		if _, err := c.upstreams.V1().Upstreams().Create(us); err != nil {
			return fmt.Errorf("failed to create upstream crd %s: %v", us.Name, err)
		}
	}
	for _, us := range upstreamsToUpdate {
		if _, err := c.upstreams.V1().Upstreams().Update(us); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", us.Name, err)
		}
	}
	// only remaining are no longer desired, delete em!
	for _, us := range actualUpstreams {
		if err := c.upstreams.V1().Upstreams().Delete(us.Name); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", us.Name, err)
		}
	}
	return nil
}

func upstreamName(serviceNamespace, serviceName, servicePort string) string {
	return fmt.Sprintf("%s-%s-%s-%s", upstreamPrefix, serviceNamespace, serviceName, servicePort)
}
