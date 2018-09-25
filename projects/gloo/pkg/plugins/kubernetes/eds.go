package kubernetes

import (
	"context"
	"fmt"

	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
	kubev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func (p *plugin) WatchEndpoints(writeNamespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	startInformerFactoryOnce(p.kube)
	opts = opts.WithDefaults()

	return newEndpointsWatcher(p.kube, upstreamsToTrack).watch(writeNamespace, opts)
}

type edsWatcher struct {
	kube      kubernetes.Interface
	upstreams map[core.ResourceRef]*kubeplugin.UpstreamSpec
}

func newEndpointsWatcher(kube kubernetes.Interface, upstreams v1.UpstreamList) *edsWatcher {
	upstreamSpecs := make(map[core.ResourceRef]*kubeplugin.UpstreamSpec)
	for _, us := range upstreams {
		kubeUpstream, ok := us.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Kube)
		// only care about kube upstreams
		if !ok {
			continue
		}
		upstreamSpecs[us.Metadata.Ref()] = kubeUpstream.Kube
	}
	return &edsWatcher{
		kube:      kube,
		upstreams: upstreamSpecs,
	}
}

func (c *edsWatcher) List(writeNamespace string, opts clients.ListOpts) (v1.EndpointList, error) {
	endpoints, err := kubePlugin.endpointsLister.List(labels.SelectorFromSet(opts.Selector))
	if err != nil {
		return nil, err
	}

	pods, err := kubePlugin.podsLister.List(labels.SelectorFromSet(opts.Selector))
	if err != nil {
		return nil, err
	}

	return filterEndpoints(opts.Ctx, writeNamespace, endpoints, pods, c.upstreams), nil
}

func (c *edsWatcher) watch(writeNamespace string, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	watch := kubePlugin.subscribe()

	endpointsChan := make(chan v1.EndpointList)
	errs := make(chan error)
	updateResourceList := func() {
		list, err := c.List(writeNamespace, clients.ListOpts{
			Ctx:      opts.Ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		endpointsChan <- list
	}

	go func() {
		defer kubePlugin.unsubscribe(watch)
		defer close(endpointsChan)
		defer close(errs)

		updateResourceList()
		for {
			select {
			case _, ok := <-watch:
				if !ok {
					return
				}
				updateResourceList()
			case <-opts.Ctx.Done():
				return
			}
		}

	}()
	return endpointsChan, errs, nil
}

func filterEndpoints(ctx context.Context, writeNamespace string, kubeEndpoints []*kubev1.Endpoints, pods []*kubev1.Pod, upstreams map[core.ResourceRef]*kubeplugin.UpstreamSpec) v1.EndpointList {
	var endpoints v1.EndpointList

	logger := contextutils.LoggerFrom(contextutils.WithLogger(ctx, "kubernetes_eds"))

	// for each upstream
	for usName, spec := range upstreams {
		// find each matching endpoint
		for _, eps := range kubeEndpoints {
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
					if len(spec.Selector) != 0 {

						// determine whether labels for the owner of this ip (pod) matches the spec
						podLabels, err := getPodLabelsForIp(addr.IP, pods)
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
					}

					hash, _ := hashstructure.Hash([]interface{}{subset, addr}, nil)
					endpointName := fmt.Sprintf("%v-%x", eps.Name, hash)

					ep := createEndpoint(writeNamespace, endpointName, usName, addr.IP, port)
					endpoints = append(endpoints, ep)
				}
			}
		}

	}
	return endpoints
}

func createEndpoint(namespace, name string, upstreamName core.ResourceRef, address string, port uint32) *v1.Endpoint {
	return &v1.Endpoint{
		Metadata: core.Metadata{
			Namespace: namespace,
			Name:      name,
		},
		UpstreamName: upstreamName.Key(), // TODO(yuval-k): should this be Ref?
		Address:      address,
		Port:         port,
	}
}

func getPodLabelsForIp(ip string, pods []*kubev1.Pod) (map[string]string, error) {
	for _, pod := range pods {
		if pod.Status.PodIP == ip && pod.Status.Phase == kubev1.PodRunning {
			return pod.Labels, nil
		}
	}
	return nil, errors.Errorf("running pod not found with ip %v", ip)
}
