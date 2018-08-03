package v1

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/log"
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
		attributeClient AttributeClient
		roleClient RoleClient
		upstreamClient UpstreamClient
		virtualServiceClient VirtualServiceClient
	)

	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

		attributeClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: AttributeCrd,
			Cfg: cfg,
		})
		attributeClient, err = NewAttributeClient(attributeClientFactory)
		Expect(err).NotTo(HaveOccurred())

		roleClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: RoleCrd,
			Cfg: cfg,
		})
		roleClient, err = NewRoleClient(roleClientFactory)
		Expect(err).NotTo(HaveOccurred())

		upstreamClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: UpstreamCrd,
			Cfg: cfg,
		})
		upstreamClient, err = NewUpstreamClient(upstreamClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: VirtualServiceCrd,
			Cfg: cfg,
		})
		virtualServiceClient, err = NewVirtualServiceClient(virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())
		cache = NewCache(attributeClient, roleClient, upstreamClient, virtualServiceClient)
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
		attribute1, err := attributeClient.Write(NewAttribute(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.AttributeList).To(HaveLen(1))
			Expect(snap.AttributeList).To(ContainElement(attribute1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		attribute2, err := attributeClient.Write(NewAttribute(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.AttributeList).To(HaveLen(2))
			Expect(snap.AttributeList).To(ContainElement(attribute1))
			Expect(snap.AttributeList).To(ContainElement(attribute2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		role1, err := roleClient.Write(NewRole(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.RoleList).To(HaveLen(1))
			Expect(snap.RoleList).To(ContainElement(role1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		role2, err := roleClient.Write(NewRole(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.RoleList).To(HaveLen(2))
			Expect(snap.RoleList).To(ContainElement(role1))
			Expect(snap.RoleList).To(ContainElement(role2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		upstream1, err := upstreamClient.Write(NewUpstream(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.UpstreamList).To(HaveLen(1))
			Expect(snap.UpstreamList).To(ContainElement(upstream1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		upstream2, err := upstreamClient.Write(NewUpstream(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.UpstreamList).To(HaveLen(2))
			Expect(snap.UpstreamList).To(ContainElement(upstream1))
			Expect(snap.UpstreamList).To(ContainElement(upstream2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		virtualService1, err := virtualServiceClient.Write(NewVirtualService(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.VirtualServiceList).To(HaveLen(1))
			Expect(snap.VirtualServiceList).To(ContainElement(virtualService1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		virtualService2, err := virtualServiceClient.Write(NewVirtualService(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.VirtualServiceList).To(HaveLen(2))
			Expect(snap.VirtualServiceList).To(ContainElement(virtualService1))
			Expect(snap.VirtualServiceList).To(ContainElement(virtualService2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		err = attributeClient.Delete(attribute2.Metadata.Namespace, attribute2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.AttributeList).To(HaveLen(1))
			Expect(snap.AttributeList).To(ContainElement(attribute1))
			Expect(snap.AttributeList).NotTo(ContainElement(attribute2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = attributeClient.Delete(attribute1.Metadata.Namespace, attribute1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.AttributeList).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		err = roleClient.Delete(role2.Metadata.Namespace, role2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.RoleList).To(HaveLen(1))
			Expect(snap.RoleList).To(ContainElement(role1))
			Expect(snap.RoleList).NotTo(ContainElement(role2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = roleClient.Delete(role1.Metadata.Namespace, role1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.RoleList).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		err = upstreamClient.Delete(upstream2.Metadata.Namespace, upstream2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.UpstreamList).To(HaveLen(1))
			Expect(snap.UpstreamList).To(ContainElement(upstream1))
			Expect(snap.UpstreamList).NotTo(ContainElement(upstream2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = upstreamClient.Delete(upstream1.Metadata.Namespace, upstream1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.UpstreamList).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		err = virtualServiceClient.Delete(virtualService2.Metadata.Namespace, virtualService2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.VirtualServiceList).To(HaveLen(1))
			Expect(snap.VirtualServiceList).To(ContainElement(virtualService1))
			Expect(snap.VirtualServiceList).NotTo(ContainElement(virtualService2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = virtualServiceClient.Delete(virtualService1.Metadata.Namespace, virtualService1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.VirtualServiceList).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
	})
})
