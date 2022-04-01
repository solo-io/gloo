package locality

import (
	"context"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	k8s_core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	"github.com/solo-io/skv2/contrib/pkg/sets"
	k8s_core_types "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	UnsupportedServiceType = func(svc *k8s_core_types.Service, clusterName string) error {
		return eris.Errorf("Unsupported service type (%s) found for gateway service (%s.%s) on cluster (%s)",
			svc.Spec.Type, svc.GetName(), svc.GetNamespace(), clusterName)
	}
	NoExternallyResolvableIp = func(svc *k8s_core_types.Service, clusterName string) error {
		return eris.Errorf("Service (%s.%s) of type LoadBalancer on cluster (%s) is not yet externally accessible",
			svc.GetName(), svc.GetNamespace(), clusterName)
	}
	NoAvailableIngresses = func(svc *k8s_core_types.Service, clusterName string) error {
		return eris.Errorf("Service (%s.%s) of type LoadBalancer on cluster (%s) has no ingress",
			svc.GetName(), svc.GetNamespace(), clusterName)
	}
	NoScheduledPods = func(svc *k8s_core_types.Service, clusterName string) error {
		return eris.Errorf("no node found for the service's pods. ensure at least one pod has been deployed "+
			"for the (%s.%s) service on cluster (%s)", svc.GetName(), svc.GetNamespace(), clusterName)
	}
	NoActiveAddressesForNode = func(node *k8s_core_types.Node, clusterName string) error {
		return eris.Errorf("no active addresses found for node (%s) on cluster (%s)",
			node.GetName(), clusterName)
	}
)

/*
	ExternalIpFinder takes a list of k8s services and the cluster they are located on, it can then find
	the string address to communicate with the endpoints via an external network
*/
type ExternalIpFinder interface {
	GetExternalIps(
		ctx context.Context,
		svcs []*k8s_core_types.Service,
	) ([]*IngressEndpoint, error)
}

type Port struct {
	Port uint32
	Name string
}
type IngressEndpoint struct {
	Address     string
	Ports       []*Port
	ServiceName string
}

func NewExternalIpFinder(
	clusterName string,
	podClient k8s_core_v1.PodClient,
	nodeClient k8s_core_v1.NodeClient,
) ExternalIpFinder {
	return &externalIpFinderImpl{
		podClient:   podClient,
		nodeClient:  nodeClient,
		clusterName: clusterName,
	}
}

type externalIpFinderImpl struct {
	podClient   k8s_core_v1.PodClient
	nodeClient  k8s_core_v1.NodeClient
	clusterName string
}

func (f *externalIpFinderImpl) GetExternalIps(
	ctx context.Context,
	svcs []*k8s_core_types.Service,
) ([]*IngressEndpoint, error) {
	result := make([]*IngressEndpoint, 0, len(svcs))
	multiErr := multierror.Error{}

	// Pre sort the svcs so that the output is idempotent
	sort.Slice(svcs, func(i, j int) bool {
		return sets.Key(svcs[i]) < sets.Key(svcs[j])
	})

	for _, svc := range svcs {
		switch svc.Spec.Type {
		case k8s_core_types.ServiceTypeLoadBalancer:
			ingress := svc.Status.LoadBalancer.Ingress
			if len(ingress) == 0 {
				multiErr.Errors = append(multiErr.Errors, NoAvailableIngresses(svc, f.clusterName))
				continue
			}

			// depending on the environment, the svc may have either an IP or a Hostname
			// https://istio.io/docs/tasks/traffic-management/ingress/ingress-control/#determining-the-ingress-ip-and-ports
			address := ingress[0].IP
			if address == "" {
				address = ingress[0].Hostname
			}
			if address == "" {
				multiErr.Errors = append(multiErr.Errors, NoExternallyResolvableIp(svc, f.clusterName))
				continue
			}
			ep := &IngressEndpoint{
				Address:     address,
				ServiceName: svc.GetName(),
			}
			for _, port := range svc.Spec.Ports {
				ep.Ports = append(ep.Ports, &Port{
					Port: uint32(port.Port),
					Name: port.Name,
				})
			}
			result = append(result, ep)

		case k8s_core_types.ServiceTypeNodePort:
			address, err := f.getNodeIp(ctx, svc)
			if err != nil {

				multiErr.Errors = append(multiErr.Errors, err)
				continue
			}
			ep := &IngressEndpoint{
				Address:     address,
				ServiceName: svc.GetName(),
			}
			for _, port := range svc.Spec.Ports {
				ep.Ports = append(ep.Ports, &Port{
					Port: uint32(port.NodePort),
					Name: port.Name,
				})
			}
			result = append(result, ep)
		default:
			multiErr.Errors = append(multiErr.Errors, UnsupportedServiceType(svc, f.clusterName))
			continue
		}
	}

	return result, multiErr.ErrorOrNil()
}

func (f *externalIpFinderImpl) getNodeIp(
	ctx context.Context,
	svc *k8s_core_types.Service,
) (string, error) {

	pods, err := f.podClient.ListPod(ctx, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(svc.Spec.Selector),
		Namespace:     svc.Namespace,
	})
	if err != nil {
		return "", err
	}
	var nodeName string
	for _, pod := range pods.Items {
		if pod.Spec.NodeName != "" {
			nodeName = pod.Spec.NodeName
			break
		}
	}
	if nodeName == "" {
		return "", NoScheduledPods(svc, f.clusterName)
	}

	node, err := f.nodeClient.GetNode(ctx, nodeName)
	if err != nil {
		return "", err
	}

	// Default to first address in list
	// TODO: See if this should be specified somehow
	for _, addr := range node.Status.Addresses {
		return addr.Address, nil
	}

	return "", NoActiveAddressesForNode(node, f.clusterName)
}
