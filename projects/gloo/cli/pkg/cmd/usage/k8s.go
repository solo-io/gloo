package usage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
)

func calculateNodeResources(nodes []corev1.Node) (*NodeResources, error) {
	resources := &NodeResources{}

	for _, node := range nodes {
		// Calculate allocatable (actual capacity)
		cpuAllocatable := node.Status.Allocatable[corev1.ResourceCPU]
		memoryAllocatable := node.Status.Allocatable[corev1.ResourceMemory]

		// Add to total capacity
		resources.TotalCPUCores += cpuAllocatable.Value()
		resources.TotalMemoryGb += memoryAllocatable.Value() / (1024 * 1024 * 1024) // Convert bytes to GB
	}
	resources.Nodes = len(nodes)

	return resources, nil
}

func getK8sClusterInfo(ctx context.Context, opts *Options) (*K8sClusterInfo, error) {
	restCfg, err := kubeutils.GetRestConfigWithKubeContext(opts.Top.KubeContext)
	if err != nil {
		return nil, err
	}
	kube, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	// Get all pods across all namespaces
	pods, err := kube.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Get all nodes
	nodes, err := kube.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	services, err := kube.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return &K8sClusterInfo{
		Nodes:    nodes.Items,
		Pods:     pods.Items,
		Services: services.Items,
	}, nil
}
