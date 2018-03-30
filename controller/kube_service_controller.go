package controller

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	kubeplugin "github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/kubecontroller"
)

const (
	kubeSystemNamespace = "kube-system"

	ownerAnnotationKey = "generated_by"
)

type ServiceController struct {
	errors chan error

	serviceLister kubelisters.ServiceLister
	upstreams     storage.Interface
	runFunc       func(stop <-chan struct{})

	generatedBy string
}

func NewServiceController(cfg *rest.Config,
	configStore storage.Interface,
	resyncDuration time.Duration) (*ServiceController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	// attempt to register upstreams if they don't exist
	if err := configStore.V1().Register(); err != nil && !storage.IsAlreadyExists(err) {
		return nil, errors.Wrap(err, "failed to register upstreams")
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, resyncDuration)
	serviceInformer := kubeInformerFactory.Core().V1().Services()

	c := &ServiceController{
		errors: make(chan error),

		serviceLister: serviceInformer.Lister(),
		upstreams:     configStore,
		generatedBy:   "kubernetes-upstream-discovery",
	}

	kubeController := kubecontroller.NewController("gloo-service-discovery", kubeClient,
		kubecontroller.NewLockingSyncHandler(c.syncGlooUpstreamsWithKubeServices),
		serviceInformer.Informer())

	c.runFunc = func(stop <-chan struct{}) {
		go kubeInformerFactory.Start(stop)
		go kubeController.Run(2, stop)
		// refresh every minute
		tick := time.Tick(time.Minute)
		go func() {
			for {
				select {
				case <-tick:
					c.syncGlooUpstreamsWithKubeServices()
				case <-stop:
					return
				}
			}
		}()
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
	var ourUpstreams []*v1.Upstream
	for _, us := range upstreams {
		if us.Metadata != nil && us.Metadata.Annotations[ownerAnnotationKey] == c.generatedBy {
			// our upstream, we supervise it
			ourUpstreams = append(ourUpstreams, us)
		}
	}
	return ourUpstreams, nil
}

func (c *ServiceController) generateDesiredUpstreams() ([]*v1.Upstream, error) {
	serviceList, err := c.serviceLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %v", err)
	}
	var upstreams []*v1.Upstream
	for _, svc := range serviceList {
		// ignore services in the kube-system namespace
		if svc.Namespace == kubeSystemNamespace {
			continue
		}
		// ignore the kubernetes default service
		if svc.Name == "kubernetes" && svc.Namespace == "default" {
			continue
		}

		for _, port := range svc.Spec.Ports {
			// annotate
			annotations := map[string]string{
				ownerAnnotationKey: c.generatedBy,
			}
			for k, v := range svc.Annotations {
				annotations[k] = v
			}
			// copy annotations from the service = happy users
			upstream := &v1.Upstream{
				Name: upstreamName(svc.Namespace, svc.Name, port.Port),
				Type: kubeplugin.UpstreamTypeKube,
				Spec: kubeplugin.EncodeUpstreamSpec(kubeplugin.UpstreamSpec{
					ServiceNamespace: svc.Namespace,
					ServiceName:      svc.Name,
					ServicePort:      port.TargetPort.IntVal,
				}),
				// mark the upstream as ours
				Metadata: &v1.Metadata{
					Annotations: annotations,
				},
			}
			upstreams = append(upstreams, upstream)
		}
	}
	return upstreams, nil
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
		// TODO: think about caring about already exists errors
		// This workaround is necessary because the ingress controller may be running and creating upstreams
		if _, err := c.upstreams.V1().Upstreams().Create(us); err != nil && !storage.IsAlreadyExists(err) {
			log.Debugf("creating upstream %v", us.Name)
			return fmt.Errorf("failed to create upstream crd %s: %v", us.Name, err)
		}
	}
	for _, us := range upstreamsToUpdate {
		log.Debugf("updating upstream %v", us.Name)
		// preserve functions that may have already been discovered
		currentUpstream, err := c.upstreams.V1().Upstreams().Get(us.Name)
		if err != nil {
			return fmt.Errorf("failed to get existing upstream %s: %v", us.Name, err)
		}
		// all we want to do is update the spec and merge the annotations
		currentUpstream.Spec = us.Spec
		if currentUpstream.Metadata == nil {
			currentUpstream.Metadata = &v1.Metadata{}
		}
		currentUpstream.Metadata.Annotations = mergeAnnotations(currentUpstream.Metadata.Annotations, us.Metadata.Annotations)
		if _, err := c.upstreams.V1().Upstreams().Update(currentUpstream); err != nil {
			return fmt.Errorf("failed to update upstream %s: %v", us.Name, err)
		}
	}
	// only remaining are no longer desired, delete em!
	for _, us := range actualUpstreams {
		log.Debugf("deleting upstream %v", us.Name)
		if err := c.upstreams.V1().Upstreams().Delete(us.Name); err != nil {
			return fmt.Errorf("failed to update upstream crd %s: %v", us.Name, err)
		}
	}
	return nil
}

func upstreamName(serviceNamespace, serviceName string, servicePort int32) string {
	return fmt.Sprintf("%s-%s-%v", serviceNamespace, serviceName, servicePort)
}

// get the unique set of funcs between two lists
// if conflict, new wins
func mergeAnnotations(oldAnnotations, newAnnotations map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range oldAnnotations {
		merged[k] = v
	}
	for k, v := range newAnnotations {
		merged[k] = v
	}
	return merged
}
