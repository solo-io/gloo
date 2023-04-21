package utils

import (
	"context"
	"fmt"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Simplified version of the code used by glooctl proxy url.
// Set `local` to true for local testing in kind cluster.
func GetIngressHost(ctx context.Context, ns string) (string, error) {
	kubeClient, err := NewKubeClient()
	if err != nil {
		return "", err
	}

	svc, err := kubeClient.CoreV1().Services(ns).Get(ctx, defaults.GatewayProxyName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	var svcPort *corev1.ServicePort
	switch len(svc.Spec.Ports) {
	case 0:
		return "", eris.New("gateway-proxy service has no ports")
	case 1:
		svcPort = &svc.Spec.Ports[0]
	default:
		for _, p := range svc.Spec.Ports {
			if p.Name == "http" {
				svcPort = &p
				break
			}
		}
		if svcPort == nil {
			return "", eris.New("http port not found on gateway-proxy service")
		}
	}

	var host string
	var port int32
	if len(svc.Status.LoadBalancer.Ingress) == 0 {
		host, err = getNodeIp(ctx, svc, kubeClient)
		if err != nil {
			return "", err
		}
		port = svcPort.NodePort
	} else {
		host = svc.Status.LoadBalancer.Ingress[0].Hostname
		if host == "" {
			host = svc.Status.LoadBalancer.Ingress[0].IP
		}
		port = svcPort.Port
	}

	return fmt.Sprintf("http://%s:%d", host, port), nil
}

func getNodeIp(ctx context.Context, svc *corev1.Service, kube kubernetes.Interface) (string, error) {
	// pick a node where one of our pods is running
	pods, err := kube.CoreV1().Pods(svc.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(svc.Spec.Selector).String(),
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
		return "", eris.Errorf("no node found for %v's pods. ensure at least one pod has been deployed "+
			"for the %v service", svc.Name, svc.Name)
	}

	node, err := kube.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	for _, addr := range node.Status.Addresses {
		return addr.Address, nil
	}

	return "", eris.Errorf("no active addresses found for node %v", node.Name)
}
