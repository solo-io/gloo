package serviceentry

import (
	"net"
	"strconv"

	"github.com/solo-io/gloo/projects/gloo/pkg/discovery"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pkg/config/schema/gvr"
	kubeclient "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"
	corev1 "k8s.io/api/core/v1"
)

var _ discovery.DiscoveryPlugin = &sePlugin{}

func buildPodEndpoint(
	se *networkingclient.ServiceEntry,
	pod *corev1.Pod,
	writeNamespace string,
	port uint32,
) *v1.Endpoint {
	return &v1.Endpoint{
		Metadata: &core.Metadata{
			Name:        "istio-se-ep-" + se.Name + "-" + se.Namespace + "-pod-" + pod.Name,
			Namespace:   writeNamespace,
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
		},
		Address: pod.Status.PodIP,
		Port:    port,
		// Hostname:    "",
	}
}

func (s *sePlugin) WatchEndpoints(
	writeNamespace string,
	upstreamsToTrack v1.UpstreamList,
	opts clients.WatchOpts,
) (<-chan v1.EndpointList, <-chan error, error) {
	defaultFilter := kclient.Filter{ObjectFilter: s.client.ObjectFilter()}

	// watch service entries everywhere, they're referenced by global hostnames
	serviceEntries := kclient.NewDelayedInformer[*networkingclient.ServiceEntry](s.client,
		gvr.ServiceEntry, kubetypes.StandardInformer, defaultFilter)
	rawServiceEntries := krt.WrapClient(serviceEntries, krt.WithName("ServiceEntries"))
	ServiceEntries := krt.NewCollection(rawServiceEntries, func(ctx krt.HandlerContext, raw *networkingclient.ServiceEntry) *serviceEntry {
		return &serviceEntry{raw}
	})
	seNsIndex := krt.NewNamespaceIndex(ServiceEntries)

	se := networkingclient.ServiceEntry{}
	se.Spec.GetWorkloadSelector().GetLabels()

	// TODO workloadentry

	endpointsChan := make(chan v1.EndpointList)
	errs := make(chan error)

	// index the static upstreams by host:port
	var upstreams []upstream
	for _, us := range upstreamsToTrack {
		upstreams = append(upstreams, upstream{us})
	}

	Upstreams := krt.NewStaticCollection(upstreams)
	staticMeshUpstreams := krt.NewIndex(Upstreams, func(us upstream) []string {
		if !us.GetStatic().GetIsMesh() {
			return nil
		}
		var hostPorts []string
		for _, h := range us.GetStatic().GetHosts() {
			hostPorts = append(hostPorts, net.JoinHostPort(h.GetAddr(), strconv.Itoa(int(h.GetPort()))))
		}
		return hostPorts
	})

	// We generate an Endpoint for each port of a ServiceEntry that selects this pod,
	// when that port is referenced from an Upstream.
	// We can't filter Pod watch by namespace, we're doing a hostname lookup and we don't know
	// which namespaces the ServiceEntries are in, and therefore the namespaces that they're doing selection on.
	pods := kclient.NewFiltered[*corev1.Pod](s.client, kclient.Filter{
		ObjectTransform: kubeclient.StripPodUnusedFields,
	})
	Pods := krt.WrapClient(pods, krt.WithName("Pods"))
	podEndpoints := krt.NewManyCollection(Pods, podEndpoints(
		writeNamespace,
		ServiceEntries, seNsIndex,
		Upstreams, staticMeshUpstreams,
	))

	// TODO endpoints from service
	// weEndpoints := krt.NewManyCollection[*networkingclient.WorkloadEntry, *v1.Endpoint](WorkloadEntries, weEndpoints(
	// 	...
	// ))
	// TODO inline endpoints from serviceentry
	// seEndpoints := krt.NewManyCollection[*networkingclient.ServiceEntry, *v1.Endpoint](ServiceEntries, weEndpoints(
	// 	...
	// ))

	// TODO also handle Upstream_Kube for Waypoint usecase:
	// * ServiceEntry will gen Upstream_Kube so users can use it just like in Istio
	// * We also need to allow Service selects WorkloadEntry for Istio parity

	podEndpoints.RegisterBatch(func(_ []krt.Event[endpoint], _ bool) {
		println("stevenctl: pushing")
		endpointsChan <- endpoints(podEndpoints.List()).Unwrap()
		println("stevenctl: pushed")
	}, true)

	go serviceEntries.Start(opts.Ctx.Done())
	go pods.Start(opts.Ctx.Done())
	// TODO move this client start earlier, confirm the late informer starts are ok
	go s.client.RunAndWait(opts.Ctx.Done())

	return endpointsChan, errs, nil
}

func podEndpoints(
	writeNamespace string,
	ServiceEntries krt.Collection[serviceEntry],
	seNsIndex krt.Index[string, serviceEntry],
	Upstreams krt.Collection[upstream],
	staticMeshUpstreams krt.Index[string, upstream],
) func(krt.HandlerContext, *corev1.Pod) []endpoint {
	return func(ctx krt.HandlerContext, pod *corev1.Pod) []endpoint {
		// find all the ServiceEntries that select us
		selectedBy := krt.Fetch(ctx, ServiceEntries, krt.FilterIndex(seNsIndex, pod.Namespace), krt.FilterSelectsNonEmpty(pod.GetLabels()))
		if len(selectedBy) == 0 {
			return nil
		}
		println("stevenctl: selected by serviceentry: ", pod.Name)

		// only need to generate once for each target port
		outByPort := map[uint32]*v1.Endpoint{}

		seenUs := map[string]*v1.Endpoint{}
		for _, se := range selectedBy {
			for _, h := range se.Spec.Hosts {
				for _, p := range se.Spec.Ports {
					// find the upstreams that have a reference to this hostname/port
					hostPort := net.JoinHostPort(h, strconv.Itoa(int(p.Number)))
					upstreamsWithRef := krt.Fetch(ctx, Upstreams, krt.FilterIndex(staticMeshUpstreams, hostPort))
					if upstreamsWithRef == nil {
						// TODO also handle Upstream_Kube for Waypoint usecase:
						// * ServiceEntry will gen Upstream_Kube so users can use it just like in Istio
						continue
					}

					// init an endpoint instance for this port if necessary
					port := p.GetNumber()
					if tp := p.GetTargetPort(); tp > 0 {
						port = tp
					}
					out := outByPort[port]
					if out == nil {
						out = buildPodEndpoint(se.ServiceEntry, pod, writeNamespace, port)
						outByPort[port] = out
					}

					// add the Endpoint ResourceRefs to the Upstreams that reference the ServiceEntry
					for _, us := range upstreamsWithRef {
						// we only want to reference each upstream once for each port
						k := us.GetMetadata().GetName() + "~" + us.GetMetadata().GetNamespace() + "~" + strconv.Itoa(int(p.Number))
						if _, ok := seenUs[k]; ok {
							continue
						}
						out.Upstreams = append(out.Upstreams, &core.ResourceRef{
							Name:      us.GetMetadata().GetName(),
							Namespace: us.GetMetadata().GetNamespace(),
						})
					}
				}
			}
		}

		var out []endpoint
		for _, ep := range outByPort {
			out = append(out, endpoint{ep})
		}
		return out
	}
}

func getNamespace(obj runtime.Object) string {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return ""
	}
	return metaObj.GetNamespace()
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
