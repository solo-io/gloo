package v1

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("V1Cache", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace      string
		cfg            *rest.Config
		cache          Cache
		artifactClient ArtifactClient
		endpointClient EndpointClient
		proxyClient    ProxyClient
		secretClient   SecretClient
		upstreamClient UpstreamClient
		kube           kubernetes.Interface
	)

	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := services.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

		// Artifact Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		artifactClientFactory := factory.NewResourceClientFactory(&factory.KubeConfigMapClientOpts{
			Clientset: kube,
		})
		artifactClient, err = NewArtifactClient(artifactClientFactory)
		Expect(err).NotTo(HaveOccurred())

		// Endpoint Constructor
		endpointClientFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
			Cache: memory.NewInMemoryResourceCache(),
		})
		endpointClient, err = NewEndpointClient(endpointClientFactory)
		Expect(err).NotTo(HaveOccurred())

		// Proxy Constructor
		proxyClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: ProxyCrd,
			Cfg: cfg,
		})
		proxyClient, err = NewProxyClient(proxyClientFactory)
		Expect(err).NotTo(HaveOccurred())

		// Secret Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		secretClientFactory := factory.NewResourceClientFactory(&factory.KubeSecretClientOpts{
			Clientset: kube,
		})
		secretClient, err = NewSecretClient(secretClientFactory)
		Expect(err).NotTo(HaveOccurred())

		// Upstream Constructor
		upstreamClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd: UpstreamCrd,
			Cfg: cfg,
		})
		upstreamClient, err = NewUpstreamClient(upstreamClientFactory)
		Expect(err).NotTo(HaveOccurred())
		cache = NewCache(artifactClient, endpointClient, proxyClient, secretClient, upstreamClient)
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
		artifact1, err := artifactClient.Write(NewArtifact(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainartifact:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainartifact
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.ArtifactList).To(ContainElement(artifact1))

		artifact2, err := artifactClient.Write(NewArtifact(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ArtifactList).To(ContainElement(artifact1))
			Expect(snap.ArtifactList).To(ContainElement(artifact2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		endpoint1, err := endpointClient.Write(NewEndpoint(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainendpoint:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainendpoint
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.EndpointList).To(ContainElement(endpoint1))

		endpoint2, err := endpointClient.Write(NewEndpoint(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.EndpointList).To(ContainElement(endpoint1))
			Expect(snap.EndpointList).To(ContainElement(endpoint2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		proxy1, err := proxyClient.Write(NewProxy(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainproxy:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainproxy
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.ProxyList).To(ContainElement(proxy1))

		proxy2, err := proxyClient.Write(NewProxy(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ProxyList).To(ContainElement(proxy1))
			Expect(snap.ProxyList).To(ContainElement(proxy2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		secret1, err := secretClient.Write(NewSecret(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainsecret:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainsecret
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.SecretList).To(ContainElement(secret1))

		secret2, err := secretClient.Write(NewSecret(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.SecretList).To(ContainElement(secret1))
			Expect(snap.SecretList).To(ContainElement(secret2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		upstream1, err := upstreamClient.Write(NewUpstream(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

	drainupstream:
		for {
			select {
			case snap = <-snapshots:
			case err := <-errs:
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Millisecond * 500):
				break drainupstream
			case <-time.After(time.Second):
				Fail("expected snapshot before 1 second")
			}
		}
		Expect(snap.UpstreamList).To(ContainElement(upstream1))

		upstream2, err := upstreamClient.Write(NewUpstream(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.UpstreamList).To(ContainElement(upstream1))
			Expect(snap.UpstreamList).To(ContainElement(upstream2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = artifactClient.Delete(artifact2.Metadata.Namespace, artifact2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ArtifactList).To(ContainElement(artifact1))
			Expect(snap.ArtifactList).NotTo(ContainElement(artifact2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = artifactClient.Delete(artifact1.Metadata.Namespace, artifact1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ArtifactList).NotTo(ContainElement(artifact1))
			Expect(snap.ArtifactList).NotTo(ContainElement(artifact2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = endpointClient.Delete(endpoint2.Metadata.Namespace, endpoint2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.EndpointList).To(ContainElement(endpoint1))
			Expect(snap.EndpointList).NotTo(ContainElement(endpoint2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = endpointClient.Delete(endpoint1.Metadata.Namespace, endpoint1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.EndpointList).NotTo(ContainElement(endpoint1))
			Expect(snap.EndpointList).NotTo(ContainElement(endpoint2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = proxyClient.Delete(proxy2.Metadata.Namespace, proxy2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ProxyList).To(ContainElement(proxy1))
			Expect(snap.ProxyList).NotTo(ContainElement(proxy2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = proxyClient.Delete(proxy1.Metadata.Namespace, proxy1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.ProxyList).NotTo(ContainElement(proxy1))
			Expect(snap.ProxyList).NotTo(ContainElement(proxy2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = secretClient.Delete(secret2.Metadata.Namespace, secret2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.SecretList).To(ContainElement(secret1))
			Expect(snap.SecretList).NotTo(ContainElement(secret2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = secretClient.Delete(secret1.Metadata.Namespace, secret1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.SecretList).NotTo(ContainElement(secret1))
			Expect(snap.SecretList).NotTo(ContainElement(secret2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
		err = upstreamClient.Delete(upstream2.Metadata.Namespace, upstream2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.UpstreamList).To(ContainElement(upstream1))
			Expect(snap.UpstreamList).NotTo(ContainElement(upstream2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}

		err = upstreamClient.Delete(upstream1.Metadata.Namespace, upstream1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.UpstreamList).NotTo(ContainElement(upstream1))
			Expect(snap.UpstreamList).NotTo(ContainElement(upstream2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second * 3):
			Fail("expected snapshot before 1 second")
		}
	})
})
