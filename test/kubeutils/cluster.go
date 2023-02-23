package kubeutils

import (
	"fmt"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/client/config"

	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	"k8s.io/client-go/rest"
)

type ClusterConfig struct {
	ClusterName           string
	KubeContext           string
	RestConfig            *rest.Config
	FederatedClientset    fedv1.Clientset
	MulticlusterClientset multicluster_v1alpha1.Clientset
	GatewayClientset      gatewayv1.Clientset
	GlooClientset         gloov1.Clientset
}

func CreateClusterConfigFromKubeClusterNameEnv(clusterNameEnv string) *ClusterConfig {
	clusterName := os.Getenv(clusterNameEnv)
	Expect(clusterName).NotTo(BeEmpty())
	return CreateClusterConfigFromKubeClusterName(clusterName)
}

func CreateClusterConfigFromKubeClusterName(clusterName string) *ClusterConfig {
	kubeCtx := fmt.Sprintf("kind-%s", clusterName)
	restCfg, err := config.GetConfigWithContext(kubeCtx)
	Expect(err).NotTo(HaveOccurred())

	fedClientset, err := fedv1.NewClientsetFromConfig(restCfg)
	Expect(err).NotTo(HaveOccurred())

	multiclusterClientset, err := multicluster_v1alpha1.NewClientsetFromConfig(restCfg)
	Expect(err).NotTo(HaveOccurred())

	glooClientset, err := gloov1.NewClientsetFromConfig(restCfg)
	Expect(err).NotTo(HaveOccurred())

	gatewayClientset, err := gatewayv1.NewClientsetFromConfig(restCfg)
	Expect(err).NotTo(HaveOccurred())

	return &ClusterConfig{
		ClusterName:           clusterName,
		KubeContext:           kubeCtx,
		RestConfig:            restCfg,
		FederatedClientset:    fedClientset,
		MulticlusterClientset: multiclusterClientset,
		GatewayClientset:      gatewayClientset,
		GlooClientset:         glooClientset,
	}
}
