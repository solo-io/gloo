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
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
)

func (p *KubePlugin) WatchEndpoints(opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	opts = opts.WithDefaults()

	// initialize watches
	epWatch, err := p.kube.CoreV1().Endpoints("").Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube eps watch")
	}
	podWatch, err := p.kube.CoreV1().Pods("").Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube pods watch")
	}

	// set up buffers and channels
	endpointsChan := make(chan v1.EndpointList)
	errs := make(chan error)
	upstreamsToTrack := make(map[string]*UpstreamSpec)

	// sync functions
	syncEndpoints := func() {
		endpoints, err := p.kube.CoreV1().Endpoints("").List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
		})
		if err != nil {
			errs <- err
			return
		}
		pods, err := p.kube.CoreV1().Pods("").List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
		})
		if err != nil {
			errs <- err
			return
		}
		endpointsChan <- p.processNewEndpoints(opts.Ctx, endpoints, pods, upstreamsToTrack)
	}

	p.trackUpstreams = func(list v1.UpstreamList) {
		upstreamsToTrack = make(map[string]*UpstreamSpec)
		for _, us := range list {
			kubeUpstream, ok := us.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Kube)
			if !ok {
				continue
			}
			upstreamsToTrack[us.Metadata.Name] = kubeUpstream.Kube
		}
		// need to refresh the eds list for new upstreams
		p.queueResync()
	}

	// watch should open up with an initial read
	go syncEndpoints()

	go func() {
		for {
			select {
			case <-time.After(opts.RefreshRate):
				syncEndpoints()
			case <-p.resync:
				syncEndpoints()
			case event := <-podWatch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during pod watch: %v", event)
				default:
					syncEndpoints()
				}
			case event := <-epWatch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during endpoints watch: %v", event)
				default:
					syncEndpoints()
				}
			case <-opts.Ctx.Done():
				epWatch.Stop()
				podWatch.Stop()
				close(endpointsChan)
				close(errs)
				return
			}
		}
	}()

	return endpointsChan, errs, nil
}

func (p *KubePlugin) processNewEndpoints(ctx context.Context, kubeEndpoints *kubev1.EndpointsList, pods *kubev1.PodList, upstreams map[string]*UpstreamSpec) v1.EndpointList {
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

					ep := createEndpoint(eps.Namespace, endpointName, usName, addr.IP, port)
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

func (p *KubePlugin) TrackUpstreams(list v1.UpstreamList) {
	p.trackUpstreams(list)
}

func (p *KubePlugin) queueResync() {
	go func() {
		p.resync <- struct{}{}
	}()
}
