package kube_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/glue/config/store/kube"
	. "github.com/solo-io/glue/test/helpers"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("KubeStore", func() {
	var (
		mkb                       *MinikubeInstance
		masterUrl, kubeconfigPath string
	)
	BeforeSuite(func() {
		mkb = NewMinikube(false, RandString(6))
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.IP()
		Must(err)
	})
	AfterSuite(func() {
		err := mkb.Teardown()
		Must(err)
	})
	Describe("init the store", func() {
		It("registers crds", func() {
			cache := NewKubeCache()
			err := cache.Init()
			Expect(err).NotTo(HaveOccurred())
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			clientset, err := apiexts.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			crds, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().List(v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(crds.Items).To(HaveLen(1))
			Expect(crds.Items[0].Name).To(Equal("foo"))
		})
	})
})
