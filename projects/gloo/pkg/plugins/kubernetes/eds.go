package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func (p *KubePlugin) WatchEndpoints(writeNamespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	opts = opts.WithDefaults()

	return newEndpointsWatcher(p.kube, upstreamsToTrack).Watch(writeNamespace, opts)
}

type edsWatcher struct {
	kube      kubernetes.Interface
	upstreams map[string]*kubeplugin.UpstreamSpec
}

func newEndpointsWatcher(kube kubernetes.Interface, upstreams v1.UpstreamList) *edsWatcher {
	upstreamSpecs := make(map[string]*kubeplugin.UpstreamSpec)
	for _, us := range upstreams {
		kubeUpstream, ok := us.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Kube)
		// only care about kube upstreams
		if !ok {
			continue
		}
		upstreamSpecs[us.Metadata.Name] = kubeUpstream.Kube
	}
	return &edsWatcher{
		kube:      kube,
		upstreams: upstreamSpecs,
	}
}

func (c *edsWatcher) List(namespace string, opts clients.ListOpts) (v1.EndpointList, error) {
	endpoints, err := c.kube.CoreV1().Endpoints("").List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, err
	}

	pods, err := c.kube.CoreV1().Pods("").List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, err
	}

	return filterEndpoints(opts.Ctx, namespace, endpoints, pods, c.upstreams), nil
}

func (c *edsWatcher) Watch(namespace string, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	endpointsWatch, err := c.kube.CoreV1().Endpoints("").Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube watch")
	}

	podsWatch, err := c.kube.CoreV1().Pods("").Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube watch")
	}

	resourcesChan := make(chan v1.EndpointList)
	errs := make(chan error)
	updateResourceList := func() {
		list, err := c.List(namespace, clients.ListOpts{
			Ctx:      opts.Ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list
	}
	// watch should open up with an initial read
	go updateResourceList()

	go func() {
		for {
			select {
			case <-time.After(opts.RefreshRate):
				updateResourceList()
			case event := <-endpointsWatch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during watch: %v", event)
				default:
					updateResourceList()
				}
			case event := <-podsWatch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during watch: %v", event)
				default:
					updateResourceList()
				}
			case <-opts.Ctx.Done():
				endpointsWatch.Stop()
				podsWatch.Stop()
				close(resourcesChan)
				close(errs)
				return
			}
		}
	}()

	return resourcesChan, errs, nil
}

func filterEndpoints(ctx context.Context, writeNamespace string, kubeEndpoints *kubev1.EndpointsList, pods *kubev1.PodList, upstreams map[string]*kubeplugin.UpstreamSpec) v1.EndpointList {
	var endpoints v1.EndpointList

	logger := contextutils.LoggerFrom(contextutils.WithLogger(ctx, "kubernetes_eds"))

	// for each upstream
	for usName, spec := range upstreams {
		// find each matching endpoint
		for _, eps := range kubeEndpoints.Items {
			if eps.Namespace != spec.ServiceNamespace || eps.Name != spec.ServiceName {
				continue
			}
			for _, subset := range eps.Subsets {
				var port uint32
				for _, p := range subset.Ports {
					if spec.ServicePort == uint32(p.Port) {
						port = uint32(p.Port)
						break
					}
				}
				if port == 0 {
					logger.Warnf("upstream %v: port %v not found for service %v", usName, spec.ServicePort, spec.ServiceName)
					continue
				}
				for _, addr := range subset.Addresses {
					// determine whether labels for the owner of this ip (pod) matches the spec
					podLabels, err := getPodLabelsForIp(addr.IP, pods.Items)
					if err != nil {
						// pod not found for ip? what's that about?
						logger.Warnf("error for upstream %v service %v: ", usName, spec.ServiceName, err)
						continue
					}
					if !labels.AreLabelsInWhiteList(spec.Selector, podLabels) {
						continue
					}
					// pod hasn't been assigned address yet
					if addr.IP == "" {
						continue
					}
					hash, _ := hashstructure.Hash([]interface{}{subset, addr}, nil)
					endpointName := fmt.Sprintf("%v.%v", eps.Name, hash)

					ep := createEndpoint(writeNamespace, endpointName, usName, addr.IP, port)
					endpoints = append(endpoints, ep)
				}
			}
		}

	}
	return endpoints
}

func createEndpoint(namespace, name, upstreamName, address string, port uint32) *v1.Endpoint {
	return &v1.Endpoint{
		Metadata: core.Metadata{
			Namespace: namespace,
			Name:      name,
		},
		UpstreamName: upstreamName,
		Address:      address,
		Port:         port,
	}
}

func getPodLabelsForIp(ip string, pods []kubev1.Pod) (map[string]string, error) {
	for _, pod := range pods {
		if pod.Status.PodIP == ip && pod.Status.Phase == kubev1.PodRunning {
			return pod.Labels, nil
		}
	}
	return nil, errors.Errorf("running pod not found with ip %v", ip)
}
