package v1alpha3_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/solo-io/solo-kit/projects/supergloo/pkg/api/external/istio/networking/v1alpha3"
)

var _ = FDescribe("VirtualServiceClient.Sk", func() {
	It("works", func() {
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		vsClient, err := NewVirtualServiceClient(&factory.KubeResourceClientFactory{
			Crd:         VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: kube.NewKubeCache(),
		})
		Expect(err).NotTo(HaveOccurred())
		err = vsClient.Register()
		Expect(err).NotTo(HaveOccurred())
		ls, err := vsClient.List("default", clients.ListOpts{})
		ls, err = vsClient.List("default", clients.ListOpts{})
		ls, err = vsClient.List("default", clients.ListOpts{})
		ls, err = vsClient.List("default", clients.ListOpts{})
		ls, err = vsClient.List("default", clients.ListOpts{})
		ls, err = vsClient.List("default", clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		vs := ls[0]
		vs.Metadata.ResourceVersion = ""
		vs.Metadata.Name = "google"
		vs.Gateways = []string{"google.sucks.com"}
		written, err := vsClient.Write(vs, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(written).To(Equal(""))
	})
})
