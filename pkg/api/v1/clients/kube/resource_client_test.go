package kube_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/mocks"
	"github.com/solo-io/solo-kit/test/services"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Base", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace string
		cfg       *rest.Config
		client    *ResourceClient
	)
	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		apiextsClient, err := apiexts.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
		resourceClient, err := versioned.NewForConfig(cfg, mocks.MockCrd)
		Expect(err).NotTo(HaveOccurred())
		client = NewResourceClient(mocks.MockCrd, apiextsClient, resourceClient, &mocks.MockResource{})
	})
	AfterEach(func() {
		services.TeardownKube(namespace)
	})
	It("CRUDs resources", func() {
		err := client.Register()
		Expect(err).NotTo(HaveOccurred())

		name := "foo"
		input := mocks.NewMockResource(name)
		input.Metadata.Namespace = namespace
		r1, err := client.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Write(input, clients.WriteOpts{})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsExist(err)).To(BeTrue())

		Expect(r1).To(BeAssignableToTypeOf(&mocks.MockResource{}))
		Expect(r1.GetMetadata().Name).To(Equal(name))
		Expect(r1.GetMetadata().Namespace).To(Equal(namespace))
		Expect(r1.GetMetadata().ResourceVersion).NotTo(Equal("7"))
		Expect(r1.(*mocks.MockResource).Data).To(Equal(name))

		_, err = client.Write(input, clients.WriteOpts{
			OverwriteExisting: true,
		})
		Expect(err).To(HaveOccurred())

		input.Metadata.ResourceVersion = r1.GetMetadata().ResourceVersion
		r1, err = client.Write(input, clients.WriteOpts{
			OverwriteExisting: true,
		})
		Expect(err).NotTo(HaveOccurred())

		read, err := client.Read(name, clients.ReadOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(read).To(Equal(r1))

		_, err = client.Read(name, clients.ReadOpts{Namespace: "doesntexist"})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		name = "boo"
		input = &mocks.MockResource{
			Data: name,
			Metadata: core.Metadata{
				Name:      name,
				Namespace: namespace,
			},
		}
		r2, err := client.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		list, err := client.List(clients.ListOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).To(ContainElement(r2))

		err = client.Delete("adsfw", clients.DeleteOpts{
			Namespace: namespace,
		})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		err = client.Delete("adsfw", clients.DeleteOpts{
			IgnoreNotExist: true,
			Namespace:      namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		err = client.Delete(r2.GetMetadata().Name, clients.DeleteOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		list, err = client.List(clients.ListOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).NotTo(ContainElement(r2))

		w, errs, err := client.Watch(clients.WatchOpts{
			Namespace:   namespace,
			RefreshRate: time.Hour,
		})
		Expect(err).NotTo(HaveOccurred())

		var r3 resources.Resource
		wait := make(chan struct{})
		go func() {
			defer close(wait)
			defer GinkgoRecover()

			resources.UpdateMetadata(r2, func(meta core.Metadata) core.Metadata {
				meta.ResourceVersion = ""
				return meta
			})
			r2, err = client.Write(r2, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			name = "goo"
			input = &mocks.MockResource{
				Data: name,
				Metadata: core.Metadata{
					Name:      name,
					Namespace: namespace,
				},
			}
			r3, err = client.Write(input, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
		}()
		<-wait

		select {
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case list = <-w:
		case <-time.After(time.Millisecond * 5):
			Fail("expected a message in channel")
		}

	drain:
		for {
			select {
			case list = <-w:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drain
			}
		}

		Expect(list).To(ContainElement(r1))
		Expect(list).To(ContainElement(r2))
		Expect(list).To(ContainElement(r3))
	})
})
