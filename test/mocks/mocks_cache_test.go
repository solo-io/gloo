package mocks

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("MocksCache", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace          string
		cfg                *rest.Config
		cache              Cache
		mockResourceClient MockResourceClient
		fakeResourceClient FakeResourceClient
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
		mockResourceClient = NewMockResourceClient(mockResourceClientFactory)
		fakeResourceKubeclient, err := versioned.NewForConfig(cfg, FakeResourceCrd)
		Expect(err).NotTo(HaveOccurred())
		fakeResourceClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd:     FakeResourceCrd,
			Kube:    fakeResourceKubeclient,
			ApiExts: apiextsClient,
		})
		fakeResourceClient = NewFakeResourceClient(fakeResourceClientFactory)
		cache = NewCache(mockResourceClient, fakeResourceClient)
	})
	AfterEach(func() {
		services.TeardownKube(namespace)
	})
	It("tracks snapshots on changes to any resource", func() {
		err := cache.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := cache.Snapshots(clients.WatchOpts{
			Namespace:   namespace,
			RefreshRate: time.Minute,
		})
		Expect(err).NotTo(HaveOccurred())
		mockResource1, err := mockResourceClient.Write(NewMockResource(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockResourceList).To(HaveLen(1))
			Expect(snap.MockResourceList).To(ContainElement(mockResource1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		mockResource2, err := mockResourceClient.Write(NewMockResource(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockResourceList).To(HaveLen(2))
			Expect(snap.MockResourceList).To(ContainElement(mockResource1))
			Expect(snap.MockResourceList).To(ContainElement(mockResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		fakeResource1, err := fakeResourceClient.Write(NewFakeResource(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.FakeResourceList).To(HaveLen(1))
			Expect(snap.FakeResourceList).To(ContainElement(fakeResource1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		fakeResource2, err := fakeResourceClient.Write(NewFakeResource(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.FakeResourceList).To(HaveLen(2))
			Expect(snap.FakeResourceList).To(ContainElement(fakeResource1))
			Expect(snap.FakeResourceList).To(ContainElement(fakeResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		err = mockResourceClient.Delete(mockResource2.Metadata.Name, clients.DeleteOpts{
			Namespace: mockResource2.Metadata.Namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockResourceList).To(HaveLen(1))
			Expect(snap.MockResourceList).To(ContainElement(mockResource1))
			Expect(snap.MockResourceList).NotTo(ContainElement(mockResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = mockResourceClient.Delete(mockResource1.Metadata.Name, clients.DeleteOpts{
			Namespace: mockResource1.Metadata.Namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockResourceList).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		err = fakeResourceClient.Delete(fakeResource2.Metadata.Name, clients.DeleteOpts{
			Namespace: fakeResource2.Metadata.Namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.FakeResourceList).To(HaveLen(1))
			Expect(snap.FakeResourceList).To(ContainElement(fakeResource1))
			Expect(snap.FakeResourceList).NotTo(ContainElement(fakeResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = fakeResourceClient.Delete(fakeResource1.Metadata.Name, clients.DeleteOpts{
			Namespace: fakeResource1.Metadata.Namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.FakeResourceList).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
	})
})
