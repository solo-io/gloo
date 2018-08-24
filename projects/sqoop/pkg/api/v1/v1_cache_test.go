package v1

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

var _ = Describe("V1Cache", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace         string
		cfg               *rest.Config
		cache             Cache
		resolverMapClient ResolverMapClient
	)

	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

		// ResolverMap Constructor
		resolverMapClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: ResolverMapCrd,
			Cfg: cfg,
		})
		resolverMapClient, err = NewResolverMapClient(resolverMapClientFactory)
		Expect(err).NotTo(HaveOccurred())
		cache = NewCache(resolverMapClient)
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
		resolverMap1, err := resolverMapClient.Write(NewResolverMap(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainresolverMap:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainresolverMap
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.ResolverMapList).To(ContainElement(resolverMap1))

		resolverMap2, err := resolverMapClient.Write(NewResolverMap(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ResolverMapList).To(ContainElement(resolverMap1))
			Expect(snap.ResolverMapList).To(ContainElement(resolverMap2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = resolverMapClient.Delete(resolverMap2.Metadata.Namespace, resolverMap2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ResolverMapList).To(ContainElement(resolverMap1))
			Expect(snap.ResolverMapList).NotTo(ContainElement(resolverMap2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = resolverMapClient.Delete(resolverMap1.Metadata.Namespace, resolverMap1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ResolverMapList).NotTo(ContainElement(resolverMap1))
			Expect(snap.ResolverMapList).NotTo(ContainElement(resolverMap2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
	})
})
