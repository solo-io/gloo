package mocks

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
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
		mockDataClient     MockDataClient
	)

	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

		// MockResource Constructor
		mockResourceClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: MockResourceCrd,
			Cfg: cfg,
		})
		mockResourceClient, err = NewMockResourceClient(mockResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())

		// FakeResource Constructor
		fakeResourceClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: FakeResourceCrd,
			Cfg: cfg,
		})
		fakeResourceClient, err = NewFakeResourceClient(fakeResourceClientFactory)
		Expect(err).NotTo(HaveOccurred())

		// MockData Constructor
		mockDataClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: MockDataCrd,
			Cfg: cfg,
		})
		mockDataClient, err = NewMockDataClient(mockDataClientFactory)
		Expect(err).NotTo(HaveOccurred())
		cache = NewCache(mockResourceClient, fakeResourceClient, mockDataClient)
	})
	AfterEach(func() {
		services.TeardownKube(namespace)
	})
	It("tracks snapshots on changes to any resource", func() {
		err := cache.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := cache.Snapshots(namespace, clients.WatchOpts{
			RefreshRate: time.Minute,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *Snapshot
		mockResource1, err := mockResourceClient.Write(NewMockResource(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainmockResource:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainmockResource
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.MockResourceList).To(ContainElement(mockResource1))

		mockResource2, err := mockResourceClient.Write(NewMockResource(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockResourceList).To(ContainElement(mockResource1))
			Expect(snap.MockResourceList).To(ContainElement(mockResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		fakeResource1, err := fakeResourceClient.Write(NewFakeResource(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainfakeResource:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainfakeResource
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.FakeResourceList).To(ContainElement(fakeResource1))

		fakeResource2, err := fakeResourceClient.Write(NewFakeResource(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.FakeResourceList).To(ContainElement(fakeResource1))
			Expect(snap.FakeResourceList).To(ContainElement(fakeResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		mockData1, err := mockDataClient.Write(NewMockData(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainmockData:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainmockData
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.MockDataList).To(ContainElement(mockData1))

		mockData2, err := mockDataClient.Write(NewMockData(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockDataList).To(ContainElement(mockData1))
			Expect(snap.MockDataList).To(ContainElement(mockData2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = mockResourceClient.Delete(mockResource2.Metadata.Namespace, mockResource2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockResourceList).To(ContainElement(mockResource1))
			Expect(snap.MockResourceList).NotTo(ContainElement(mockResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = mockResourceClient.Delete(mockResource1.Metadata.Namespace, mockResource1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockResourceList).NotTo(ContainElement(mockResource1))
			Expect(snap.MockResourceList).NotTo(ContainElement(mockResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = fakeResourceClient.Delete(fakeResource2.Metadata.Namespace, fakeResource2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.FakeResourceList).To(ContainElement(fakeResource1))
			Expect(snap.FakeResourceList).NotTo(ContainElement(fakeResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = fakeResourceClient.Delete(fakeResource1.Metadata.Namespace, fakeResource1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.FakeResourceList).NotTo(ContainElement(fakeResource1))
			Expect(snap.FakeResourceList).NotTo(ContainElement(fakeResource2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = mockDataClient.Delete(mockData2.Metadata.Namespace, mockData2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockDataList).To(ContainElement(mockData1))
			Expect(snap.MockDataList).NotTo(ContainElement(mockData2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = mockDataClient.Delete(mockData1.Metadata.Namespace, mockData1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockDataList).NotTo(ContainElement(mockData1))
			Expect(snap.MockDataList).NotTo(ContainElement(mockData2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
	})
})
