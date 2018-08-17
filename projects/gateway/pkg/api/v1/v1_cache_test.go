package v1

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
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
		namespace          string
		cfg                *rest.Config
		cache              Cache
		gatewayClient GatewayClient
		virtualServiceClient VirtualServiceClient
	)

	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

		// Gateway Constructor
		gatewayClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: GatewayCrd,
			Cfg: cfg,
		})
		gatewayClient, err = NewGatewayClient(gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		// VirtualService Constructor
		virtualServiceClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: VirtualServiceCrd,
			Cfg: cfg,
		})
		virtualServiceClient, err = NewVirtualServiceClient(virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())
		cache = NewCache(gatewayClient, virtualServiceClient)
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
		gateway1, err := gatewayClient.Write(NewGateway(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	draingateway:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break draingateway
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.GatewayList).To(ContainElement(gateway1))

		gateway2, err := gatewayClient.Write(NewGateway(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.GatewayList).To(ContainElement(gateway1))
			Expect(snap.GatewayList).To(ContainElement(gateway2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		virtualService1, err := virtualServiceClient.Write(NewVirtualService(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainvirtualService:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainvirtualService
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.VirtualServiceList).To(ContainElement(virtualService1))

		virtualService2, err := virtualServiceClient.Write(NewVirtualService(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.VirtualServiceList).To(ContainElement(virtualService1))
			Expect(snap.VirtualServiceList).To(ContainElement(virtualService2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = gatewayClient.Delete(gateway2.Metadata.Namespace, gateway2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.GatewayList).To(ContainElement(gateway1))
			Expect(snap.GatewayList).NotTo(ContainElement(gateway2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = gatewayClient.Delete(gateway1.Metadata.Namespace, gateway1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.GatewayList).NotTo(ContainElement(gateway1))
			Expect(snap.GatewayList).NotTo(ContainElement(gateway2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = virtualServiceClient.Delete(virtualService2.Metadata.Namespace, virtualService2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.VirtualServiceList).To(ContainElement(virtualService1))
			Expect(snap.VirtualServiceList).NotTo(ContainElement(virtualService2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = virtualServiceClient.Delete(virtualService1.Metadata.Namespace, virtualService1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.VirtualServiceList).NotTo(ContainElement(virtualService1))
			Expect(snap.VirtualServiceList).NotTo(ContainElement(virtualService2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
	})
})
