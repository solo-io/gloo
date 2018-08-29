package mocks

import (
	"context"
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
		namespace1         string
		namespace2         string
		cfg                *rest.Config
		cache              Cache
		mockResourceClient MockResourceClient
		fakeResourceClient FakeResourceClient
		mockDataClient     MockDataClient
	)

	BeforeEach(func() {
		namespace1 = helpers.RandString(8)
		namespace2 = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace1)
		Expect(err).NotTo(HaveOccurred())
		err = services.SetupKubeForTest(namespace2)
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
		services.TeardownKube(namespace1)
		services.TeardownKube(namespace2)
	})
	It("tracks snapshots on changes to any resource", func() {
		ctx := context.Background()
		err := cache.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := cache.Snapshots([]string{namespace1, namespace2}, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *Snapshot
		mockResource1a, err := mockResourceClient.Write(NewMockResource(namespace1, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mockResource1b, err := mockResourceClient.Write(NewMockResource(namespace2, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

	drainmockResource:
		for {
			select {
			case snap = <-snapshots:
				// Expect(snap.Mocks.List()).To(ContainElement(mockResource1a))
				// Expect(snap.Mocks.List()).To(ContainElement(mockResource1b))
				_, err1 := snap.Mocks.List().Find(mockResource1a.Metadata.ObjectRef())
				_, err2 := snap.Mocks.List().Find(mockResource1b.Metadata.ObjectRef())
				if err1 == nil && err2 == nil {
					break drainmockResource
				}
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 10):
				nsList1, _ := mockResourceClient.List(namespace1, clients.ListOpts{})
				nsList2, _ := mockResourceClient.List(namespace2, clients.ListOpts{})
				combined := nsList1.ByNamespace()
				combined.Add(nsList2...)
				msg := log.Sprintf("expected final snapshot before 10 seconds.\nexpected %v\nreceived", combined.List(), snap.Mocks.List())
				Fail(msg)
			}
		}

		mockResource2a, err := mockResourceClient.Write(NewMockResource(namespace1, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mockResource2b, err := mockResourceClient.Write(NewMockResource(namespace2, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

	drainmockResource2:
		for {
			select {
			case snap = <-snapshots:
				_, err1 := snap.Mocks.List().Find(mockResource1a.Metadata.ObjectRef())
				_, err2 := snap.Mocks.List().Find(mockResource1b.Metadata.ObjectRef())
				_, err3 := snap.Mocks.List().Find(mockResource2a.Metadata.ObjectRef())
				_, err4 := snap.Mocks.List().Find(mockResource2b.Metadata.ObjectRef())
				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
					break drainmockResource2
				}
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 10):
				nsList1, _ := mockResourceClient.List(namespace1, clients.ListOpts{})
				nsList2, _ := mockResourceClient.List(namespace2, clients.ListOpts{})
				combined := nsList1.ByNamespace()
				combined.Add(nsList2...)
				Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
			}
		}
		fakeResource1a, err := fakeResourceClient.Write(NewFakeResource(namespace1, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		fakeResource1b, err := fakeResourceClient.Write(NewFakeResource(namespace2, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

	drainfakeResource:
		for {
			select {
			case snap = <-snapshots:
				// Expect(snap.Fakes.List()).To(ContainElement(fakeResource1a))
				// Expect(snap.Fakes.List()).To(ContainElement(fakeResource1b))
				_, err1 := snap.Fakes.List().Find(fakeResource1a.Metadata.ObjectRef())
				_, err2 := snap.Fakes.List().Find(fakeResource1b.Metadata.ObjectRef())
				if err1 == nil && err2 == nil {
					break drainfakeResource
				}
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 10):
				nsList1, _ := fakeResourceClient.List(namespace1, clients.ListOpts{})
				nsList2, _ := fakeResourceClient.List(namespace2, clients.ListOpts{})
				combined := nsList1.ByNamespace()
				combined.Add(nsList2...)
				msg := log.Sprintf("expected final snapshot before 10 seconds.\nexpected %v\nreceived", combined.List(), snap.Fakes.List())
				Fail(msg)
			}
		}

		fakeResource2a, err := fakeResourceClient.Write(NewFakeResource(namespace1, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		fakeResource2b, err := fakeResourceClient.Write(NewFakeResource(namespace2, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

	drainfakeResource2:
		for {
			select {
			case snap = <-snapshots:
				_, err1 := snap.Fakes.List().Find(fakeResource1a.Metadata.ObjectRef())
				_, err2 := snap.Fakes.List().Find(fakeResource1b.Metadata.ObjectRef())
				_, err3 := snap.Fakes.List().Find(fakeResource2a.Metadata.ObjectRef())
				_, err4 := snap.Fakes.List().Find(fakeResource2b.Metadata.ObjectRef())
				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
					break drainfakeResource2
				}
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 10):
				nsList1, _ := fakeResourceClient.List(namespace1, clients.ListOpts{})
				nsList2, _ := fakeResourceClient.List(namespace2, clients.ListOpts{})
				combined := nsList1.ByNamespace()
				combined.Add(nsList2...)
				Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
			}
		}
		mockData1a, err := mockDataClient.Write(NewMockData(namespace1, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mockData1b, err := mockDataClient.Write(NewMockData(namespace2, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

	drainmockData:
		for {
			select {
			case snap = <-snapshots:
				// Expect(snap.MockDatas.List()).To(ContainElement(mockData1a))
				// Expect(snap.MockDatas.List()).To(ContainElement(mockData1b))
				_, err1 := snap.MockDatas.List().Find(mockData1a.Metadata.ObjectRef())
				_, err2 := snap.MockDatas.List().Find(mockData1b.Metadata.ObjectRef())
				if err1 == nil && err2 == nil {
					break drainmockData
				}
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 10):
				nsList1, _ := mockDataClient.List(namespace1, clients.ListOpts{})
				nsList2, _ := mockDataClient.List(namespace2, clients.ListOpts{})
				combined := nsList1.ByNamespace()
				combined.Add(nsList2...)
				msg := log.Sprintf("expected final snapshot before 10 seconds.\nexpected %v\nreceived", combined.List(), snap.MockDatas.List())
				Fail(msg)
			}
		}

		mockData2a, err := mockDataClient.Write(NewMockData(namespace1, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mockData2b, err := mockDataClient.Write(NewMockData(namespace2, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

	drainmockData2:
		for {
			select {
			case snap = <-snapshots:
				_, err1 := snap.MockDatas.List().Find(mockData1a.Metadata.ObjectRef())
				_, err2 := snap.MockDatas.List().Find(mockData1b.Metadata.ObjectRef())
				_, err3 := snap.MockDatas.List().Find(mockData2a.Metadata.ObjectRef())
				_, err4 := snap.MockDatas.List().Find(mockData2b.Metadata.ObjectRef())
				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
					break drainmockData2
				}
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 10):
				nsList1, _ := mockDataClient.List(namespace1, clients.ListOpts{})
				nsList2, _ := mockDataClient.List(namespace2, clients.ListOpts{})
				combined := nsList1.ByNamespace()
				combined.Add(nsList2...)
				Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
			}
		}
		err = mockResourceClient.Delete(mockResource2a.Metadata.Namespace, mockResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = mockResourceClient.Delete(mockResource2b.Metadata.Namespace, mockResource2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.Mocks.List()).To(ContainElement(mockResource1a))
			Expect(snap.Mocks.List()).To(ContainElement(mockResource1b))
			Expect(snap.Mocks.List()).NotTo(ContainElement(mockResource2a))
			Expect(snap.Mocks.List()).NotTo(ContainElement(mockResource2b))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = mockResourceClient.Delete(mockResource1a.Metadata.Namespace, mockResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = mockResourceClient.Delete(mockResource1b.Metadata.Namespace, mockResource1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.Mocks.List()).NotTo(ContainElement(mockResource1a))
			Expect(snap.Mocks.List()).NotTo(ContainElement(mockResource1b))
			Expect(snap.Mocks.List()).NotTo(ContainElement(mockResource2a))
			Expect(snap.Mocks.List()).NotTo(ContainElement(mockResource2b))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = fakeResourceClient.Delete(fakeResource2a.Metadata.Namespace, fakeResource2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = fakeResourceClient.Delete(fakeResource2b.Metadata.Namespace, fakeResource2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.Fakes.List()).To(ContainElement(fakeResource1a))
			Expect(snap.Fakes.List()).To(ContainElement(fakeResource1b))
			Expect(snap.Fakes.List()).NotTo(ContainElement(fakeResource2a))
			Expect(snap.Fakes.List()).NotTo(ContainElement(fakeResource2b))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = fakeResourceClient.Delete(fakeResource1a.Metadata.Namespace, fakeResource1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = fakeResourceClient.Delete(fakeResource1b.Metadata.Namespace, fakeResource1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.Fakes.List()).NotTo(ContainElement(fakeResource1a))
			Expect(snap.Fakes.List()).NotTo(ContainElement(fakeResource1b))
			Expect(snap.Fakes.List()).NotTo(ContainElement(fakeResource2a))
			Expect(snap.Fakes.List()).NotTo(ContainElement(fakeResource2b))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = mockDataClient.Delete(mockData2a.Metadata.Namespace, mockData2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = mockDataClient.Delete(mockData2b.Metadata.Namespace, mockData2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockDatas.List()).To(ContainElement(mockData1a))
			Expect(snap.MockDatas.List()).To(ContainElement(mockData1b))
			Expect(snap.MockDatas.List()).NotTo(ContainElement(mockData2a))
			Expect(snap.MockDatas.List()).NotTo(ContainElement(mockData2b))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = mockDataClient.Delete(mockData1a.Metadata.Namespace, mockData1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = mockDataClient.Delete(mockData1b.Metadata.Namespace, mockData1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.MockDatas.List()).NotTo(ContainElement(mockData1a))
			Expect(snap.MockDatas.List()).NotTo(ContainElement(mockData1b))
			Expect(snap.MockDatas.List()).NotTo(ContainElement(mockData2a))
			Expect(snap.MockDatas.List()).NotTo(ContainElement(mockData2b))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
	})
})
