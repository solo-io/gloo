package setup

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/bootstrap"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
)

func Setup(bstrp bootstrap.Config, inputResourceOpts factory.ResourceClientFactoryOpts, secretOpts factory.ResourceClientFactoryOpts, artifactOpts factory.ResourceClientFactoryOpts) error {
	inputFactory := factory.NewResourceClientFactory(inputResourceOpts)
	secretFactory := factory.NewResourceClientFactory(secretOpts)
	artifactFactory := factory.NewResourceClientFactory(artifactOpts)
	// endpoints are internal-only, therefore use the in-memory client
	endpointsFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
		Cache: memory.NewInMemoryResourceCache(),
	})

	upstreamClient, err := v1.NewUpstreamClient(inputFactory)
	if err != nil {
		return err
	}

	proxyClient, err := v1.NewProxyClient(inputFactory)
	if err != nil {
		return err
	}

	endpointClient, err := v1.NewEndpointClient(endpointsFactory)
	if err != nil {
		return err
	}

	secretClient, err := v1.NewSecretClient(secretFactory)
	if err != nil {
		return err
	}

	artifactClient, err := v1.NewArtifactClient(artifactFactory)
	if err != nil {
		return err
	}

	// TODO: initialize endpointClient using EDS plugins
	cache := v1.NewCache(artifactClient, endpointClient, proxyClient, secretClient, upstreamClient)
	el := v1.NewEventLoop(cache, &syncer{})
}

type syncer struct{}

func (syncer) Sync(snap *v1.Snapshot) error {

}

type DiscoveryOpts struct {
	KubeOpts struct {
		IgnoredServices []string
	}
}

func SetupDiscovery(bstrp bootstrap.Config, cache v1.Cache, namespace string, opts clients.WatchOpts, discovery DiscoveryOpts) error {
	kube := bstrp.KubeClient()
	serviceWatch, err := kube.CoreV1().Services(namespace).Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return err
	}
	upstreamsChan := make(chan v1.UpstreamList)
	errs := make(chan error)
	syncUpstreams := func() {
		list, err := kube.CoreV1().Services(namespace).List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
		})
		if err != nil {
			errs <- err
			return
		}
		upstreamsChan <- convertServices(list.Items)
	}
	// watch should open up with an initial read
	go syncUpstreams()

	for {
		select {

		case event := <-serviceWatch.ResultChan():
			switch event.Type {
			case kubewatch.Error:
				errs <- errors.Errorf("error during pod watch: %v", event)
			default:
				syncUpstreams()
			}
		case <-opts.Ctx.Done():
			return nil
		}
	}
}

func convertServices(list []kubev1.Service, opts DiscoveryOpts) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, svc := range list {
		if skip(svc, opts) {
			continue
		}
		upstreams = append(upstreams, convertService(svc)...)
	}
	return upstreams
}

func convertService(svc kubev1.Service) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, port := range svc.Spec.Ports {
		upstreams = append(upstreams, createUpstream(svc.ObjectMeta, port))
	}
}

func skip(svc kubev1.Service, opts DiscoveryOpts) bool {
	for _, name := range opts.KubeOpts.IgnoredServices {
		if svc.Name == name {
			return true
		}
	}
	return false
}

func createUpstream(meta metav1.ObjectMeta, port kubev1.ServicePort) *v1.Upstream {
	coremeta := kubeutils.FromKubeMeta(meta)
	coremeta.Name = upstreamName(meta.Namespace, meta.Name, port.Port)
	return &v1.Upstream{
		Metadata:     coremeta,
		UpstreamSpec: &v1.UpstreamSpec{},
	}
}

func upstreamName(serviceNamespace, serviceName string, servicePort int32) string {
	return fmt.Sprintf("%s-%s-%v", serviceNamespace, serviceName, servicePort)
}

/*
start service watch, funnel to a "changes" chan
for {
<change
list services

generate v1.Upstream foreach service

sync desired[] with actual[]

}
*/
