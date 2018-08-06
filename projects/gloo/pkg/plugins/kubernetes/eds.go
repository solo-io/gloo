package kubernetes

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubewatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type endpointsWatcher struct{
	kube kubernetes.Interface
	upstreams []*v1.Upstream
	resync chan struct{}
}

func newEndpointsWatcher(kube kubernetes.Interface) *endpointsWatcher {
	return &endpointsWatcher{
		kube: kube,
		resync: make(chan struct{}),
	}
}

func (w *endpointsWatcher) TrackUpstreams(upstreams []*v1.Upstream) {
	w.upstreams = upstreams
	go func(){
		w.resync <- struct{}{}
	}()
}

func (w *endpointsWatcher) Watch(namespace string, opts clients.WatchOpts) (<-chan map[edsMapping]*kubev1.Endpoints, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	epWatch, err := w.kube.CoreV1().Endpoints(namespace).Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube eps watch in %v", namespace)
	}
	svcWatch, err := w.kube.CoreV1().Services(namespace).Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube eps watch in %v", namespace)
	}
	podWatch, err := w.kube.CoreV1().Pods(namespace).Watch(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "initiating kube eps watch in %v", namespace)
	}
	endpointsChan := make(chan map[edsMapping]*kubev1.Endpoints)
	errs := make(chan error)
	updateResourceList := func() {
		list, err := w.kube.CoreV1().Endpoints(namespace).List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(opts.Selector).String(),
		})
		if err != nil {
			errs <- err
			return
		}
		endpointsChan <- processNewEndpoints(list)
	}
	// watch should open up with an initial read
	go updateResourceList()

	go func() {
		for {
			select {
			case <-w.resync:
				updateResourceList()
			case event := <-svcWatch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during svc watch: %v", event)
				default:
					updateResourceList()
				}
			case event := <-podWatch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during pod watch: %v", event)
				default:
					updateResourceList()
				}
			case event := <-epWatch.ResultChan():
				switch event.Type {
				case kubewatch.Error:
					errs <- errors.Errorf("error during endpoints watch: %v", event)
				default:
					updateResourceList()
				}
			case <-opts.Ctx.Done():
				epWatch.Stop()
				svcWatch.Stop()
				podWatch.Stop()
				close(endpointsChan)
				close(errs)
				return
			}
		}
	}()

	return endpointsChan, errs, nil
}

type edsMapping struct {
	serviceName string
	servicePort uint32
}

func processNewEndpoints(list *kubev1.EndpointsList) map[edsMapping]*kubev1.Endpoints {
	endpointSet := make(map[edsMapping ]*kubev1.Endpoints)
	for _, eps := range list.Items {
		for _, subset := range eps.Subsets {

		}
	}
	return endpointSet
}
