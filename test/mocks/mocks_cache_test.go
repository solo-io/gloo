package mocks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/test/mocks"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"k8s.io/client-go/rest"
	"github.com/solo-io/solo-kit/test/services"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/bxcodec/faker"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"os"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"path/filepath"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("MocksCache", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace string
		cfg       *rest.Config
		cache     Cache
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
		mockResourceKubeclient, err := versioned.NewForConfig(cfg, MockResourceCrd)
		Expect(err).NotTo(HaveOccurred())
		mockResourceClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd:     MockResourceCrd,
			Kube:    mockResourceKubeclient,
			ApiExts: apiextsClient,
		})
		mockResourceClient := NewMockResourceClient(mockResourceClientFactory)
		fakeResourceKubeclient, err := versioned.NewForConfig(cfg, FakeResourceCrd)
		Expect(err).NotTo(HaveOccurred())
		fakeResourceClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd:     FakeResourceCrd,
			Kube:    fakeResourceKubeclient,
			ApiExts: apiextsClient,
		})
		fakeResourceClient := NewFakeResourceClient(fakeResourceClientFactory)
		cache = NewCache(mockResourceClient, fakeResourceClient)
	})
	AfterEach(func() {
		services.TeardownKube(namespace)
	})
	It("CRUDs resources", func() {
		err := cache.Register()
		Expect(err).NotTo(HaveOccurred())

		name := "foo"
		input := NewMockResource(namespace, name)
		input.Metadata.Namespace = namespace
		r1, err := cache.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		_, err = cache.Write(input, clients.WriteOpts{})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsExist(err)).To(BeTrue())

		Expect(r1).To(BeAssignableToTypeOf(&MockResource{}))
		Expect(r1.GetMetadata().Name).To(Equal(name))
		Expect(r1.GetMetadata().Namespace).To(Equal(namespace))
		Expect(r1.GetMetadata().ResourceVersion).NotTo(Equal("7"))
		Expect(r1.Data).To(Equal(input.Data))

		_, err = cache.Write(input, clients.WriteOpts{
			OverwriteExisting: true,
		})
		Expect(err).To(HaveOccurred())

		input.Metadata.ResourceVersion = r1.GetMetadata().ResourceVersion
		r1, err = cache.Write(input, clients.WriteOpts{
			OverwriteExisting: true,
		})
		Expect(err).NotTo(HaveOccurred())

		read, err := cache.Read(name, clients.ReadOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(read).To(Equal(r1))

		_, err = cache.Read(name, clients.ReadOpts{Namespace: "doesntexist"})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		name = "boo"
		input = &MockResource{}
		err = faker.FakeData(input)
		Expect(err).NotTo(HaveOccurred())
		input.Metadata = core.Metadata{
			Name:      name,
			Namespace: namespace,
		}

		r2, err := cache.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		list, err := cache.List(clients.ListOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).To(ContainElement(r2))

		err = cache.Delete("adsfw", clients.DeleteOpts{
			Namespace: namespace,
		})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		err = cache.Delete("adsfw", clients.DeleteOpts{
			IgnoreNotExist: true,
			Namespace:      namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		err = cache.Delete(r2.GetMetadata().Name, clients.DeleteOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		list, err = cache.List(clients.ListOpts{
			Namespace: namespace,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).NotTo(ContainElement(r2))

		w, errs, err := cache.Watch(clients.WatchOpts{
			Namespace:   namespace,
			RefreshRate: time.Hour,
		})
		Expect(err).NotTo(HaveOccurred())

		var r3 resources.Resource
		wait := make(chan struct{})
		go func() {
			defer close(wait)
			defer GinkgoRecover()

			resources.UpdateMetadata(r2, func(meta *core.Metadata) {
				meta.ResourceVersion = ""
			})
			r2, err = cache.Write(r2, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			name = "goo"
			input = &MockResource{}
			err = faker.FakeData(input)
			Expect(err).NotTo(HaveOccurred())
			input.Metadata = core.Metadata{
				Name:      name,
				Namespace: namespace,
			}

			r3, err = cache.Write(input, clients.WriteOpts{})
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
