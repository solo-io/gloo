package versioned_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/gloo/pkg/log"
	"os"
	"path/filepath"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/solo-io/solo-kit/test/mocks"
	"k8s.io/client-go/rest"
)

var _ = Describe("Clientset", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace string
		cfg       *rest.Config
	)
	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		services.TeardownKube(namespace)
	})
	It("registers, creates, deletes resource implementations", func() {
		)
		mockCrdClient, err := NewForConfig(cfg, mocks.MockCrd)
		Expect(err).NotTo(HaveOccurred())
	})
})
