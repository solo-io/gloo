package locality

import (
	"context"

	k8s_core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	v1sets "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/sets"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewLocalityFinder(nodeClient k8s_core_v1.NodeClient, podClient k8s_core_v1.PodClient) LocalityFinder {
	return &localityFinderImpl{
		nodeClient: nodeClient,
		podClient:  podClient,
	}
}

type localityFinderImpl struct {
	nodeClient k8s_core_v1.NodeClient
	podClient  k8s_core_v1.PodClient
}

func (l *localityFinderImpl) GetRegion(ctx context.Context) (string, error) {
	nodeList, err := l.nodeClient.ListNode(ctx)
	if err != nil {
		return "", err
	}
	for _, node := range nodeList.Items {
		if len(node.Labels) == 0 {
			continue
		}
		// Try stable labels, followed by deprecated tosupport earlier kube versions
		if regionStable, ok := node.Labels[corev1.LabelZoneRegionStable]; ok {
			return regionStable, nil
		} else if regionDeprecated, okDep := node.Labels[corev1.LabelZoneRegion]; okDep {
			return regionDeprecated, nil
		}
	}

	return "", nil
}

func (l *localityFinderImpl) ZonesForDeployment(
	ctx context.Context,
	deployment *appsv1.Deployment,
) ([]string, error) {
	return l.getLocality(ctx, deployment, deployment.Spec.Selector)
}

func (l *localityFinderImpl) ZonesForDaemonSet(
	ctx context.Context,
	daemonSet *appsv1.DaemonSet,
) ([]string, error) {
	return l.getLocality(ctx, daemonSet, daemonSet.Spec.Selector)
}

func (l *localityFinderImpl) getLocality(
	ctx context.Context,
	obj metav1.Object,
	selector *metav1.LabelSelector,
) ([]string, error) {
	labelSelector, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		return nil, err
	}
	podList, err := l.podClient.ListPod(
		ctx,
		client.InNamespace(obj.GetNamespace()),
		client.MatchingLabelsSelector{Selector: labelSelector},
	)
	if err != nil {
		return nil, err
	}

	// Some pods may be scheduled on the same node, so only include each one in the list once
	nodeSet := v1sets.NewNodeSet()
	for _, podIter := range podList.Items {
		pod := podIter
		node, err := l.nodeClient.GetNode(ctx, pod.Spec.NodeName)
		if err != nil {
			return nil, err
		}
		nodeSet.Insert(node)
	}

	zoneSet := sets.NewString()
	for _, node := range nodeSet.List() {
		// If no labels are available, skip this node and try the next
		if len(node.Labels) == 0 {
			continue
		}
		// Try stable labels, followed by deprecated tosupport earlier kube versions
		if zoneStable, ok := node.Labels[corev1.LabelZoneFailureDomainStable]; ok {
			zoneSet.Insert(zoneStable)
		} else if zoneDeprecated, okDep := node.Labels[corev1.LabelZoneFailureDomain]; okDep {
			zoneSet.Insert(zoneDeprecated)
		}
	}

	return zoneSet.List(), nil
}
