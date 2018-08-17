package kubernetes

import (
	"fmt"
	"reflect"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/discovery"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
)

func (p *KubePlugin) WatchUpstreams(writeNamespace string, opts clients.WatchOpts, discOpts discovery.Opts) (chan v1.UpstreamList, chan error, error) {
	opts = opts.WithDefaults()
	serviceWatch, err := p.kube.CoreV1().Services("").Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, err
	}
	upstreamsChan := make(chan v1.UpstreamList)
	errs := make(chan error)
	syncUpstreams := func() {
		list, err := p.kube.CoreV1().Services("").List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
		})
		if err != nil {
			errs <- err
			return
		}
		upstreamsChan <- convertServices(list.Items, discOpts, writeNamespace)
	}
	// watch should open up with an initial read
	go syncUpstreams()

	go func() {
		for {
			select {
			case <-time.After(opts.RefreshRate):
				syncUpstreams()
			case event := <-serviceWatch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during pod watch: %v", event)
				default:
					syncUpstreams()
				}
			case <-opts.Ctx.Done():
				return
			}
		}
	}()
	return upstreamsChan, errs, nil
}

func convertServices(list []kubev1.Service, opts discovery.Opts, writeNamespace string) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, svc := range list {
		if skip(svc, opts) {
			continue
		}
		upstreams = append(upstreams, convertService(svc, writeNamespace)...)
	}
	return upstreams
}

func convertService(svc kubev1.Service, writeNamespace string) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, port := range svc.Spec.Ports {
		upstreams = append(upstreams, createUpstream(svc.ObjectMeta, port, writeNamespace))
	}
	return upstreams
}

func createUpstream(meta metav1.ObjectMeta, port kubev1.ServicePort, writeNamespace string) *v1.Upstream {
	coremeta := kubeutils.FromKubeMeta(meta)
	coremeta.Name = upstreamName(meta.Namespace, meta.Name, port.Port)
	coremeta.Namespace = writeNamespace
	return &v1.Upstream{
		Metadata: coremeta,
		UpstreamSpec: &v1.UpstreamSpec{
			UpstreamType: &v1.UpstreamSpec_Kube{
				Kube: &kube.UpstreamSpec{
					ServiceName:      meta.Name,
					ServiceNamespace: meta.Namespace,
					ServicePort:      uint32(port.Port),
				},
			},
		},
		DiscoveryMetadata: &v1.DiscoveryMetadata{},
	}
}

func upstreamName(serviceNamespace, serviceName string, servicePort int32) string {
	return fmt.Sprintf("%s-%s-%v", serviceNamespace, serviceName, servicePort)
}

func skip(svc kubev1.Service, opts discovery.Opts) bool {
	for _, name := range opts.KubeOpts.IgnoredServices {
		if svc.Name == name {
			return true
		}
	}
	return false
}

func (p *KubePlugin) UpdateUpstream(original, desired *v1.Upstream) error {
	originalSpec, ok := original.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Kube)
	if !ok {
		return errors.Errorf("internal error: expected *v1.UpstreamSpec_Kube, got %v", reflect.TypeOf(original.UpstreamSpec.UpstreamType).Name())
	}
	desiredSpec, ok := desired.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Kube)
	if !ok {
		return errors.Errorf("internal error: expected *v1.UpstreamSpec_Kube, got %v", reflect.TypeOf(original.UpstreamSpec.UpstreamType).Name())
	}
	// copy service spec, we don't want to overwrite that
	desiredSpec.Kube.ServiceSpec = originalSpec.Kube.ServiceSpec
	// copy labels; user may have written them over. cannot be auto-discovered
	desiredSpec.Kube.Selector = originalSpec.Kube.Selector

	return nil
}
