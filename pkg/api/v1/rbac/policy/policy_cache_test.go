package policy

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

var _ = Describe("PolicyCache", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace      string
		cfg            *rest.Config
		cache          Cache
		policyClient   PolicyClient
		identityClient IdentityClient
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
		policyKubeclient, err := versioned.NewForConfig(cfg, PolicyCrd)
		Expect(err).NotTo(HaveOccurred())
		policyClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd:     PolicyCrd,
			Kube:    policyKubeclient,
			ApiExts: apiextsClient,
		})
		policyClient = NewPolicyClient(policyClientFactory)
		identityKubeclient, err := versioned.NewForConfig(cfg, IdentityCrd)
		Expect(err).NotTo(HaveOccurred())
		identityClientFactory := factory.NewResourceClientFactory(&factory.KubeResourceClientOpts{
			Crd:     IdentityCrd,
			Kube:    identityKubeclient,
			ApiExts: apiextsClient,
		})
		identityClient = NewIdentityClient(identityClientFactory)
		cache = NewCache(policyClient, identityClient)
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
		policy1, err := policyClient.Write(NewPolicy(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.PolicyList).To(HaveLen(1))
			Expect(snap.PolicyList).To(ContainElement(policy1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		policy2, err := policyClient.Write(NewPolicy(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.PolicyList).To(HaveLen(2))
			Expect(snap.PolicyList).To(ContainElement(policy1))
			Expect(snap.PolicyList).To(ContainElement(policy2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		identity1, err := identityClient.Write(NewIdentity(namespace, "angela"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.IdentityList).To(HaveLen(1))
			Expect(snap.IdentityList).To(ContainElement(identity1))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		identity2, err := identityClient.Write(NewIdentity(namespace, "lane"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.IdentityList).To(HaveLen(2))
			Expect(snap.IdentityList).To(ContainElement(identity1))
			Expect(snap.IdentityList).To(ContainElement(identity2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		err = policyClient.Delete(policy2.Metadata.Namespace, policy2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.PolicyList).To(HaveLen(1))
			Expect(snap.PolicyList).To(ContainElement(policy1))
			Expect(snap.PolicyList).NotTo(ContainElement(policy2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = policyClient.Delete(policy1.Metadata.Namespace, policy1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.PolicyList).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
		err = identityClient.Delete(identity2.Metadata.Namespace, identity2.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.IdentityList).To(HaveLen(1))
			Expect(snap.IdentityList).To(ContainElement(identity1))
			Expect(snap.IdentityList).NotTo(ContainElement(identity2))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}

		err = identityClient.Delete(identity1.Metadata.Namespace, identity1.Metadata.Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())

		select {
		case snap := <-snapshots:
			Expect(snap.IdentityList).To(HaveLen(0))
		case err := <-errs:
			Expect(err).NotTo(HaveOccurred())
		case <-time.After(time.Second):
			Fail("expected snapshot before 1 second")
		}
	})
})
