package kube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type NamespaceClient interface {
	ListNamespaces() ([]string, error)
}

type namespaceClient struct {
	client corev1.CoreV1Interface
}

func (n *namespaceClient) ListNamespaces() ([]string, error) {
	var namespaces []string
	nsList, err := n.client.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Name)
	}
	return namespaces, nil
}

func NewNamespaceClient(client corev1.CoreV1Interface) NamespaceClient {
	return &namespaceClient{
		client: client,
	}
}
