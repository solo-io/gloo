package serviceentry

import (
	"context"
	"strings"

	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	"github.com/solo-io/gloo/projects/gloo/pkg/discovery"
	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ discovery.DiscoveryPlugin = &sePlugin{}

// TODO informers
// type seNsInformer struct {
// 	serviceEntryInformer cache.SharedIndexInformer
// 	podInformer          cache.SharedIndexInformer
//  // TODO workloadentry
// }
//
// type seInformers struct {
// 	perNamespace map[string]seNsInformer
// }

func (s *sePlugin) listEndpoints(
	ctx context.Context,
	writeNamespace string,
	listOpts metav1.ListOptions,
	watchNamespaces []string,
	upstreams v1.UpstreamList,
) (v1.EndpointList, error) {
	serviceEntries := make(map[string][]*networkingclient.ServiceEntry, len(watchNamespaces))
	pods := make(map[string][]corev1.Pod, len(watchNamespaces))
	for _, ns := range watchNamespaces {
		seList, _ := s.istio.NetworkingV1beta1().ServiceEntries(ns).List(ctx, listOpts)
		for _, se := range seList.Items {
			serviceEntries[se.Namespace] = append(serviceEntries[se.Namespace], se)
		}
		podList, _ := s.kube.CoreV1().Pods(ns).List(ctx, listOpts)
		for _, pod := range podList.Items {
			pods[pod.Namespace] = append(pods[pod.Namespace], pod)
		}
	}

	var out v1.EndpointList
	for _, us := range upstreams {
		kubeUs := us.GetKube()
		if kubeUs == nil {
			// we pre-filtered, this should never happen
			contextutils.LoggerFrom(ctx).DPanicw(
				"upstream was not kube",
				"namespace", us.GetMetadata().GetNamespace(),
				"name", us.GetMetadata().GetName(),
			)
			continue
		}
		if !strings.HasPrefix(us.GetMetadata().Name, UpstreamNamePrefix) {
			continue
		}
		var upstreamServiceEntry *networkingclient.ServiceEntry
		for _, se := range serviceEntries[us.GetMetadata().GetNamespace()] {
			if se.Name == kubeUs.ServiceName && se.Namespace == kubeUs.ServiceNamespace {
				upstreamServiceEntry = se
				break
			}
		}
		if upstreamServiceEntry == nil {
			contextutils.LoggerFrom(ctx).Warn("could not find the associated service entry for this upstream")
			continue
		}
		for _, p := range pods[us.GetMetadata().Namespace] {
			selector := labels.Set(kubeUs.Selector).AsSelector()
			if !selector.Matches(labels.Set(p.Labels)) {
				continue
			}
			out = append(out, buildPodEndpoint(us, upstreamServiceEntry, p))
		}
		// TODO workloadentry
	}
	return out, nil
}

func buildPodEndpoint(us *v1.Upstream, se *networkingclient.ServiceEntry, pod corev1.Pod) *v1.Endpoint {
	port := us.GetKube().GetServicePort()
	for _, p := range se.Spec.GetPorts() {
		if p.GetTargetPort() != 0 && p.GetNumber() == us.GetKube().GetServicePort() {
			port = p.GetTargetPort()
			break
		}
	}
	return &v1.Endpoint{
		Metadata: &core.Metadata{
			Name:        "istio-se-ep-" + se.Name + "-pod-" + pod.Name,
			Namespace:   se.Namespace,
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
		},
		Upstreams: []*core.ResourceRef{{
			// TODO re-use the same endpoint selected by multiple upstreams
			Name:      us.GetMetadata().GetName(),
			Namespace: us.GetMetadata().GetNamespace(),
		}},
		Address: pod.Status.PodIP,
		Port:    port,
		// Hostname:    "",
	}
}

func (s *sePlugin) WatchEndpoints(writeNamespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	watchNamespaces := s.getWatchNamespaces(upstreamsToTrack)
	watchNamespacesSet := sets.NewString(watchNamespaces...)
	// kubeUpstreams := filterMap(upstreamsToTrack, func(e *v1.Upstream) *kubernetes.UpstreamSpec {
	// 	return e.GetKube()
	// })

	// setup watches
	listOpts := buildListOptions(opts.ExpressionSelector, opts.Selector)
	podWatch, err := s.kube.CoreV1().Pods(metav1.NamespaceAll).Watch(opts.Ctx, listOpts)
	if err != nil {
		return nil, nil, err
	}
	seWatch, err := s.istio.NetworkingV1beta1().ServiceEntries(metav1.NamespaceAll).Watch(opts.Ctx, listOpts)
	if err != nil {
		return nil, nil, err
	}
	// TODO workloadentry

	endpointsChan := make(chan v1.EndpointList)
	errs := make(chan error)

	// TODO TODO use informers + listers instead of watch
	// serviceEntry := serviceEntryInformer(s.istio, metav1.NamespaceAll)
	// Pods := serviceEntryInformer(s.istio, metav1.NamespaceAll)
	// updates := make(chan struct{})
	// if err := watchInformers(updates, opts.Ctx.Done()); err != nil {
	// 	return nil, nil, err
	// }

	updateResourceList := func() {
		list, err := s.listEndpoints(opts.Ctx, writeNamespace, listOpts, watchNamespaces, upstreamsToTrack)
		if err != nil {
			errs <- err
			return
		}
		select {
		case <-opts.Ctx.Done():
			return
		case endpointsChan <- list:
		}
	}

	go func() {
		defer seWatch.Stop()
		defer podWatch.Stop()
		defer close(endpointsChan)
		defer close(errs)

		updateResourceList()
		for {
			select {
			case ev, ok := <-seWatch.ResultChan():
				if !ok {
					return
				}
				if !watchNamespacesSet.Has(getNamespace(ev.Object)) {
					continue
				}
				updateResourceList()
			case ev, ok := <-podWatch.ResultChan():
				if !ok {
					return
				}
				if !watchNamespacesSet.Has(getNamespace(ev.Object)) {
					continue
				}
				updateResourceList()
			case <-opts.Ctx.Done():
				return
			}
		}
	}()

	return endpointsChan, errs, nil
}

func getNamespace(obj runtime.Object) string {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return ""
	}
	return metaObj.GetNamespace()
}

func (s *sePlugin) getWatchNamespaces(upstreamsToTrack v1.UpstreamList) []string {
	var namespaces []string
	if settingsutil.IsAllNamespacesFromSettings(s.settings) {
		namespaces = []string{metav1.NamespaceAll}
	} else {
		nsSet := map[string]bool{}
		for _, upstream := range upstreamsToTrack {
			svcNs := upstream.GetKube().GetServiceNamespace()
			// only care about kube upstreams
			if svcNs == "" {
				continue
			}
			nsSet[svcNs] = true
		}
		for ns := range nsSet {
			namespaces = append(namespaces, ns)
		}
	}

	// If there are no upstreams to watch (eg: if discovery is disabled), namespaces remains an empty list.
	// When creating the InformerFactory, by convention, an empty namespace list means watch all namespaces.
	// To ensure that we only watch what we are supposed to, fallback to WatchNamespaces if namespaces is an empty list.
	if len(namespaces) == 0 {
		namespaces = s.settings.GetWatchNamespaces()
	}

	return namespaces
}

// satisfy interface but don't implement; only implement hybrid client for ServiceEntry Upstreams.
func (s *sePlugin) DiscoverUpstreams(watchNamespaces []string, writeNamespace string, opts clients.WatchOpts, discOpts discovery.Opts) (chan v1.UpstreamList, chan error, error) {
	us, errs := make(chan v1.UpstreamList), make(chan error)
	close(us)
	close(errs)
	return us, errs, nil
}

// satisfy interface but don't implement; only implement hybrid client for ServiceEntry Upstreams.
func (s *sePlugin) UpdateUpstream(original *v1.Upstream, desired *v1.Upstream) (bool, error) {
	return false, nil
}
