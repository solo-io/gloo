package kubernetes

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/kubecontroller"
	kubev1 "k8s.io/api/core/v1"
	kubev1resources "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/endpointdiscovery"
)

type endpointController struct {
	endpoints       chan endpointdiscovery.EndpointGroups
	errors          chan error
	endpointsLister kubelisters.EndpointsLister
	servicesLister  kubelisters.ServiceLister
	podsLister      kubelisters.PodLister
	upstreamSpecs   map[string]*UpstreamSpec
	runFunc         func(stop <-chan struct{})
}

func newEndpointController(cfg *rest.Config, resyncDuration time.Duration) (*endpointController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, resyncDuration)
	endpointInformer := informerFactory.Core().V1().Endpoints()
	serviceInformer := informerFactory.Core().V1().Services()
	podInformer := informerFactory.Core().V1().Pods()

	c := &endpointController{
		endpoints:       make(chan endpointdiscovery.EndpointGroups),
		errors:          make(chan error),
		endpointsLister: endpointInformer.Lister(),
		servicesLister:  serviceInformer.Lister(),
		podsLister:      podInformer.Lister(),
	}

	kubeController := kubecontroller.NewController("glue-endpoints-controller",
		kubeClient,
		kubecontroller.NewSyncHandler(c.syncEndpoints),
		endpointInformer.Informer(),
		serviceInformer.Informer(),
		podInformer.Informer())

	c.runFunc = func(stop <-chan struct{}) {
		go informerFactory.Start(stop)
		go kubeController.Run(2, stop)
		<-stop
		log.Printf("kube endpoint discovery stopped")
	}

	return c, nil
}

func (c *endpointController) Run(stop <-chan struct{}) {
	c.runFunc(stop)
}

// triggers an update
func (c *endpointController) TrackUpstreams(upstreams []*v1.Upstream) {
	if c.upstreamSpecs == nil {
		c.upstreamSpecs = make(map[string]*UpstreamSpec)
	}
	for _, us := range upstreams {
		spec, err := DecodeUpstreamSpec(us.Spec)
		if err != nil {
			runtime.HandleError(err)
			continue
		}
		c.upstreamSpecs[us.Name] = spec
	}
	c.syncEndpoints()
}

func (c *endpointController) Endpoints() <-chan endpointdiscovery.EndpointGroups {
	return c.endpoints
}

func (c *endpointController) Error() <-chan error {
	return c.errors
}

// pushes EndpointGroups or error to channel
func (c *endpointController) syncEndpoints() {
	endpointGroups, err := c.getUpdatedEndpoints()
	if err != nil {
		c.errors <- err
		return
	}
	// ignore empty configs / no secrets to watch
	if len(endpointGroups) == 0 {
		return
	}
	c.endpoints <- endpointGroups
}

// retrieves secrets from kubernetes
func (c *endpointController) getUpdatedEndpoints() (endpointdiscovery.EndpointGroups, error) {
	serviceList, err := c.servicesLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving endpoints: %v", err)
	}
	endpointList, err := c.endpointsLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving endpoints: %v", err)
	}
	podList, err := c.podsLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving pods: %v", err)
	}

	endpointGroups := make(endpointdiscovery.EndpointGroups)
	for upstreamName, spec := range c.upstreamSpecs {
		// find the targetport for our service
		// if targetport is empty, skip this upstream
		targetPort, err := portForUpstream(spec, serviceList)
		if err != nil || targetPort == 0 {
			runtime.HandleError(err)
			continue
		}
		for _, endpoint := range endpointList {
			if spec.ServiceName == endpoint.Name && spec.ServiceNamespace == endpoint.Namespace {
				for _, es := range endpoint.Subsets {
					for _, addr := range es.Addresses {
						// determine whether labels for the owner of this ip (pod) matches the spec
						podLabels, err := getPodLabelsForIp(addr.IP, podList)
						if err != nil {
							// pod not found for ip? what's that about?
							// log it and keep going
							runtime.HandleError(err)
							continue
						}
						if !labels.AreLabelsInWhiteList(spec.Labels, podLabels) {
							continue
						}
						// pod hasn't been assigned address yet
						if addr.IP == "" {
							continue
						}

						m := endpointdiscovery.Endpoint{
							Address: addr.IP,
							Port:    targetPort,
						}
						endpointGroups[upstreamName] = append(endpointGroups[upstreamName], m)
					}
				}
			}
		}
	}
	return endpointGroups, nil
}

func getPodLabelsForIp(ip string, pods []*kubev1.Pod) (map[string]string, error) {
	for _, pod := range pods {
		if pod.Status.PodIP == ip && pod.Status.Phase == kubev1.PodRunning {
			return pod.Labels, nil
		}
	}
	return nil, errors.Errorf("running pod not found with ip %v", ip)
}

func portForUpstream(spec *UpstreamSpec, serviceList []*kubev1resources.Service) (int32, error) {
	for _, svc := range serviceList {
		if spec.ServiceName == svc.Name && spec.ServiceNamespace == svc.Namespace {
			// found the port we want
			if svc.Spec.ExternalName != "" {
				runtime.HandleError(fmt.Errorf("WARNING: external name services are not supported for Kubernetes Endpoint Interface"))
			}
			// if the service only has one port, just assume that's the one we want
			// this way the user doesn't have to specify portname
			if len(svc.Spec.Ports) == 1 {
				return svc.Spec.Ports[0].TargetPort.IntVal, nil
			}
			for _, port := range svc.Spec.Ports {
				if spec.ServicePort == port.Name {
					if port.TargetPort.StrVal != "" {
						runtime.HandleError(fmt.Errorf("WARNING: named ports are not supported for Kubernetes Endpoint Interface"))
					}
					return port.TargetPort.IntVal, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("target port or service not found")
}
