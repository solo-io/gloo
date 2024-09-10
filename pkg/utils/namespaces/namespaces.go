package namespaces

import (
	"context"
	"os"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/solo-kit/pkg/api/external/kubernetes/namespace"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_kubernetes "k8s.io/client-go/kubernetes"
)

func AllNamespaces(watchNamespaces []string) bool {
	if len(watchNamespaces) == 0 {
		return true
	}
	if len(watchNamespaces) == 1 && watchNamespaces[0] == "" {
		return true
	}
	return false
}

func ProcessWatchNamespaces(watchNamespaces []string, writeNamespace string) []string {
	if AllNamespaces(watchNamespaces) {
		return watchNamespaces
	}

	var writeNamespaceProvided bool
	for _, ns := range watchNamespaces {
		if ns == writeNamespace {
			writeNamespaceProvided = true
			break
		}
	}
	if !writeNamespaceProvided {
		watchNamespaces = append(watchNamespaces, writeNamespace)
	}

	return watchNamespaces
}

func GetPodNamespace() string {
	if podNamespace := os.Getenv("POD_NAMESPACE"); podNamespace != "" {
		return podNamespace
	}
	return "gloo-system"
}

// NewKubeNamespaceClient returns either an implementation of the KubeNamespaceClient based on the following condition :
// If the method was able to create the kubeclient and it has permission to list namespaces, `namespace.NewNamespaceClient` else `FakeKubeNamespaceWatcher`
func NewKubeNamespaceClient(ctx context.Context) kubernetes.KubeNamespaceClient {
	kubeClient, err := helpers.KubeClientWithKubecontext("")
	if err != nil {
		return &FakeKubeNamespaceWatcher{}
	}

	clientset, ok := kubeClient.(*k8s_kubernetes.Clientset)
	if !ok {
		return &FakeKubeNamespaceWatcher{}
	}

	action := authv1.ResourceAttributes{
		Namespace: "",
		Verb:      "list",
		Resource:  "namespaces",
	}

	selfCheck := authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &action,
		},
	}

	resp, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, &selfCheck, metav1.CreateOptions{})
	if err != nil {
		return &FakeKubeNamespaceWatcher{}
	}

	if resp.Status.Allowed {
		kubeCache, err := cache.NewKubeCoreCache(ctx, kubeClient)
		if err != nil {
			return &FakeKubeNamespaceWatcher{}
		}

		return namespace.NewNamespaceClient(kubeClient, kubeCache)
	}

	return &FakeKubeNamespaceWatcher{}
}
