package helpers

import (
	"os"

	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	fakeKubeClientset *fake.Clientset
)

func MustKubeClient() kubernetes.Interface {
	client, err := KubeClient()
	if err != nil {
		log.Fatalf("failed to create kube client: %v", err)
	}
	return client
}

func KubeClient() (kubernetes.Interface, error) {
	if fakeKubeClientset != nil {
		return fakeKubeClientset, nil
	}
	cfg, err := kubeutils.GetConfig("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return nil, errors.Wrapf(err, "getting kube config")
	}
	return kubernetes.NewForConfig(cfg)
}
