package versioned_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/gloo/pkg/log"
	"os"
	"path/filepath"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/solo-io/solo-kit/test/mocks"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1 "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/solo.io/v1"
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
		apiextsClient, err := apiexts.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
		err = mocks.MockCrd.Register(apiextsClient)
		Expect(err).NotTo(HaveOccurred())

		c, err := apiextsClient.ApiextensionsV1beta1().CustomResourceDefinitions().List(v1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(c.Items)).To(BeNumerically(">=", 1))
		var found bool
		for _, i := range c.Items {
			if i.Name == mocks.MockCrd.FullName() {
				found = true
				break
			}
		}
		Expect(found).To(BeTrue())

		mockCrdClient, err := NewForConfig(cfg, mocks.MockCrd)
		Expect(err).NotTo(HaveOccurred())
		name := "foo"
		input := mocks.NewMockResource(name)
		inputCrd := mocks.MockCrd.KubeResource(input)
		created, err := mockCrdClient.ResourcesV1().Resources(namespace).Create(inputCrd)
		Expect(err).NotTo(HaveOccurred())
		Expect(created).NotTo(BeNil())
		Expect(created.Spec).NotTo(BeNil())
		Expect(created.Spec).To(Equal(&crdv1.Spec{
			"data":     "foo",
			"metadata": map[string]interface{}{"name": "foo"},
			"status":   map[string]interface{}{},
		}))
	})
})
