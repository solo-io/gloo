package tests

import (
	"os"
	"path/filepath"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/test/services"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/hashicorp/consul/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
)

type ResourceClientTester interface {
	Description() string
	Skip() bool
	Setup(namespace string) factory.ResourceClientFactoryOpts
	Teardown(namespace string)
}

type KubeRcTester struct{
	Crd crd.Crd
}

func (rct *KubeRcTester) Description() string {
	return "kube-crd"
}

func (rct *KubeRcTester) Skip() bool {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return true
	}
	return false
}

func (rct *KubeRcTester) Setup(namespace string) factory.ResourceClientFactoryOpts {
	err := services.SetupKubeForTest(namespace)
	Expect(err).NotTo(HaveOccurred())
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())
	return &factory.KubeResourceClientOpts{
		Crd: rct.Crd,
		Cfg: cfg,
	}
}

func (rct *KubeRcTester) Teardown(namespace string) {
	services.TeardownKube(namespace)
}

type ConsulRcTester struct {
	consulInstance *services.ConsulInstance
	consulFactory  *services.ConsulFactory
}

func (rct *ConsulRcTester) Description() string {
	return "consul-kv"
}

func (rct *ConsulRcTester) Skip() bool {
	if os.Getenv("RUN_CONSUL_TESTS") != "1" {
		log.Printf("This test downloads and runs consul and is disabled by default. To enable, set RUN_CONSUL_TESTS=1 in your env.")
		return true
	}
	return false
}

func (rct *ConsulRcTester) Setup(namespace string) factory.ResourceClientFactoryOpts {
	var err error
	rct.consulFactory, err = services.NewConsulFactory()
	Expect(err).NotTo(HaveOccurred())
	rct.consulInstance, err = rct.consulFactory.NewConsulInstance()
	Expect(err).NotTo(HaveOccurred())
	err = rct.consulInstance.Run()
	Expect(err).NotTo(HaveOccurred())

	consul, err := api.NewClient(api.DefaultConfig())
	Expect(err).NotTo(HaveOccurred())
	return &factory.ConsulResourceClientOpts{
		Consul:  consul,
		RootKey: namespace,
	}
}

func (rct *ConsulRcTester) Teardown(namespace string) {
	rct.consulInstance.Clean()
	rct.consulFactory.Clean()
}
