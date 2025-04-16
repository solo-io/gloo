package usage

import (
	"context"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func calculateNodeResources(nodes []v1.Node) (*NodeResources, error) {
	resources := &NodeResources{}

	for _, node := range nodes {
		// Calculate allocatable (actual capacity)
		cpuAllocatable := node.Status.Allocatable[v1.ResourceCPU]
		memoryAllocatable := node.Status.Allocatable[v1.ResourceMemory]

		// Add to total capacity
		resources.TotalCapacityCPU += cpuAllocatable.MilliValue()
		resources.TotalCapacityMemory += memoryAllocatable.Value()
	}

	return resources, nil
}

func getK8sClusterInfo(opts *Options) (*K8sClusterInfo, error) {
	restCfg, err := kubeutils.GetRestConfigWithKubeContext("")
	if err != nil {
		return nil, err
	}
	kube, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	// Get all pods across all namespaces
	pods, err := kube.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Get all nodes
	nodes, err := kube.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Get all services in the control plane namespace
	services, err := kube.CoreV1().Services(opts.ControlPlaneNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return &K8sClusterInfo{
		Nodes:    nodes.Items,
		Pods:     pods.Items,
		Services: services.Items,
	}, nil
}
