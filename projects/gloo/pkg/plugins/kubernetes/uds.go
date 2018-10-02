package kubernetes

import (
	"fmt"
	"reflect"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/discovery"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (p *plugin) DiscoverUpstreams(watchNamespaces []string, writeNamespace string, opts clients.WatchOpts, discOpts discovery.Opts) (chan v1.UpstreamList, chan error, error) {
	if p.kubeShareFactory == nil {
		p.kubeShareFactory = getInformerFactory(p.kube)
	}

	watch := p.kubeShareFactory.Subscribe()

	opts = opts.WithDefaults()
	upstreamsChan := make(chan v1.UpstreamList)
	errs := make(chan error)
	discoverUpstreams := func() {
		list, err := p.kubeShareFactory.ServicesLister().List(labels.SelectorFromSet(opts.Selector))
		if err != nil {
			errs <- err
			return
		}
		upstreamsChan <- convertServices(watchNamespaces, list, discOpts, writeNamespace)
	}

	go func() {
		defer p.kubeShareFactory.Unsubscribe(watch)
		defer close(upstreamsChan)
		defer close(errs)
		// watch should open up with an initial read
		discoverUpstreams()
		for {
			select {
			case _, ok := <-watch:
				if !ok {
					return
				}
				discoverUpstreams()
			case <-opts.Ctx.Done():
				return
			}
		}
	}()
	return upstreamsChan, errs, nil
}

func convertServices(watchNamespaces []string, list []*kubev1.Service, opts discovery.Opts, writeNamespace string) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, svc := range list {
		if skip(svc, opts) {
			continue
		}

		if !containsString(svc.Namespace, watchNamespaces) {
			continue
		}

		upstreams = append(upstreams, convertService(svc, writeNamespace)...)
	}
	return upstreams
}

func containsString(s string, slice []string) bool {
	for _, s2 := range slice {
		if s2 == s {
			return true
		}
	}
	return false
}

func convertService(svc *kubev1.Service, writeNamespace string) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, port := range svc.Spec.Ports {
		upstreams = append(upstreams, createUpstream(svc.ObjectMeta, port, writeNamespace))
	}
	return upstreams
}

func createUpstream(meta metav1.ObjectMeta, port kubev1.ServicePort, writeNamespace string) *v1.Upstream {
	coremeta := kubeutils.FromKubeMeta(meta)
	coremeta.ResourceVersion = ""
	coremeta.Name = upstreamName(meta.Namespace, meta.Name, port.Port)
	coremeta.Namespace = writeNamespace
	servicePort := port.TargetPort.IntVal
	if servicePort == 0 {
		servicePort = port.Port
	}
	return &v1.Upstream{
		Metadata: coremeta,
		UpstreamSpec: &v1.UpstreamSpec{
			UpstreamType: &v1.UpstreamSpec_Kube{
				Kube: &kubeplugin.UpstreamSpec{
					ServiceName:      meta.Name,
					ServiceNamespace: meta.Namespace,
					ServicePort:      uint32(servicePort),
				},
			},
		},
		DiscoveryMetadata: &v1.DiscoveryMetadata{},
	}
}

func upstreamName(serviceNamespace, serviceName string, servicePort int32) string {
	return fmt.Sprintf("%s-%s-%v", serviceNamespace, serviceName, servicePort)
}

func skip(svc *kubev1.Service, opts discovery.Opts) bool {
	if len(svc.Spec.Selector) == 0 {
		return true
	}
	for _, name := range opts.KubeOpts.IgnoredServices {
		if svc.Name == name {
			return true
		}
	}
	return false
}

func (p *plugin) UpdateUpstream(original, desired *v1.Upstream) (bool, error) {
	originalSpec, ok := original.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Kube)
	if !ok {
		return false, errors.Errorf("internal error: expected *v1.UpstreamSpec_Kube, got %v", reflect.TypeOf(original.UpstreamSpec.UpstreamType).Name())
	}
	desiredSpec, ok := desired.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Kube)
	if !ok {
		return false, errors.Errorf("internal error: expected *v1.UpstreamSpec_Kube, got %v", reflect.TypeOf(original.UpstreamSpec.UpstreamType).Name())
	}
	// copy service spec, we don't want to overwrite that
	desiredSpec.Kube.ServiceSpec = originalSpec.Kube.ServiceSpec
	// copy labels; user may have written them over. cannot be auto-discovered
	desiredSpec.Kube.Selector = originalSpec.Kube.Selector

	if originalSpec.Equal(desiredSpec) {
		return false, nil
	}

	return true, nil
}
