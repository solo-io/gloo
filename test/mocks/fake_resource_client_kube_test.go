package mocks

import (
	"os"
	"path/filepath"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/hashicorp/consul/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"time"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/bxcodec/faker"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/onsi/ginkgo/extensions/table"
)

var _ = FDescribe("FakeResourceClient", func() {
	var (
		namespace string
		client    FakeResourceClient

		consulFactory  *services.ConsulFactory
		consulInstance *services.ConsulInstance
		consul         *api.Client
		err            error
	)
	table.DescribeTable("resource client tests",
		func(description string, skip func() bool, setup func(namespace string) factory.ResourceClientFactoryOpts, teardown func(namespace string)) {
			if skip() {
				return
			}
			var _ = Context("with backend "+description, func() {
				var _ = BeforeEach(func() {
					namespace = helpers.RandString(8)
					factoryOpts := setup(namespace)
					client, err = NewFakeResourceClient(factory.NewResourceClientFactory(factoryOpts))
					Expect(err).NotTo(HaveOccurred())
				})
				var _ = AfterEach(func() {
					teardown(namespace)
				})
				var _ = It("CRUDs on resource type FakeResource", func() {
					testClient(namespace, client)
				})
			})
		},

		table.Entry("kube_crd", "kube_crd",
			func() bool {
			if os.Getenv("RUN_KUBE_TESTS") != "1" {
				log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
				return false
			}
			return true
		},
			func(namespace string) factory.ResourceClientFactoryOpts {
				err := services.SetupKubeForTest(namespace)
				Expect(err).NotTo(HaveOccurred())
				kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
				cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				return &factory.KubeResourceClientOpts{
					Crd: FakeResourceCrd,
					Cfg: cfg,
				}
			},
			func(namespace string) {
				services.TeardownKube(namespace)
			}),
		table.Entry("consul_kv", "consul_kv",
			func() bool {
			if os.Getenv("RUN_CONSUL_TESTS") != "1" {
				log.Printf("This test downloads and runs consul and is disabled by default. To enable, set RUN_CONSUL_TESTS=1 in your env.")
				return false
			}
			return true
		},
			func(namespace string) factory.ResourceClientFactoryOpts {
				consulFactory, err = services.NewConsulFactory()
				Expect(err).NotTo(HaveOccurred())
				consulInstance, err = consulFactory.NewConsulInstance()
				Expect(err).NotTo(HaveOccurred())
				err = consulInstance.Run()
				Expect(err).NotTo(HaveOccurred())

				c, err := api.NewClient(api.DefaultConfig())
				Expect(err).NotTo(HaveOccurred())
				consul = c
				return &factory.ConsulResourceClientOpts{
					Consul:  consul,
					RootKey: namespace,
				}
			},
			func(namespace string) {
				consulInstance.Clean()
				consulFactory.Clean()
			}),
	)
})

func testClient(namespace string, client FakeResourceClient) {
	err := client.Register()
	Expect(err).NotTo(HaveOccurred())

	name := "foo"
	input := NewFakeResource(namespace, name)
	input.Metadata.Namespace = namespace
	r1, err := client.Write(input, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	_, err = client.Write(input, clients.WriteOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsExist(err)).To(BeTrue())

	Expect(r1).To(BeAssignableToTypeOf(&FakeResource{}))
	Expect(r1.GetMetadata().Name).To(Equal(name))
	Expect(r1.GetMetadata().Namespace).To(Equal(namespace))
	Expect(r1.GetMetadata().ResourceVersion).NotTo(Equal("7"))
	Expect(r1.Count).To(Equal(input.Count))

	_, err = client.Write(input, clients.WriteOpts{
		OverwriteExisting: true,
	})
	Expect(err).To(HaveOccurred())

	input.Metadata.ResourceVersion = r1.GetMetadata().ResourceVersion
	r1, err = client.Write(input, clients.WriteOpts{
		OverwriteExisting: true,
	})
	Expect(err).NotTo(HaveOccurred())

	read, err := client.Read(namespace, name, clients.ReadOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(read).To(Equal(r1))

	_, err = client.Read("doesntexist", name, clients.ReadOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsNotExist(err)).To(BeTrue())

	name = "boo"
	input = &FakeResource{}

	// ignore return error because interfaces / oneofs mess it up
	faker.FakeData(input)

	input.Metadata = core.Metadata{
		Name:      name,
		Namespace: namespace,
	}

	r2, err := client.Write(input, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	list, err := client.List(namespace, clients.ListOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(list).To(ContainElement(r1))
	Expect(list).To(ContainElement(r2))

	err = client.Delete(namespace, "adsfw", clients.DeleteOpts{})
	Expect(err).To(HaveOccurred())
	Expect(errors.IsNotExist(err)).To(BeTrue())

	err = client.Delete(namespace, "adsfw", clients.DeleteOpts{
		IgnoreNotExist: true,
	})
	Expect(err).NotTo(HaveOccurred())

	err = client.Delete(namespace, r2.GetMetadata().Name, clients.DeleteOpts{})
	Expect(err).NotTo(HaveOccurred())
	list, err = client.List(namespace, clients.ListOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(list).To(ContainElement(r1))
	Expect(list).NotTo(ContainElement(r2))

	w, errs, err := client.Watch(namespace, clients.WatchOpts{
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
		r2, err = client.Write(r2, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		name = "goo"
		input = &FakeResource{}
		// ignore return error because interfaces / oneofs mess it up
		faker.FakeData(input)
		Expect(err).NotTo(HaveOccurred())
		input.Metadata = core.Metadata{
			Name:      name,
			Namespace: namespace,
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
}
