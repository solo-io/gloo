package crd_test

import (
	"os"
	"path/filepath"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/gloo/pkg/storage/crd"
	crdv1 "github.com/solo-io/gloo/pkg/storage/crd/solo.io/v1"
	. "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/pkg/log"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("CrdStorageClient", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		masterUrl, kubeconfigPath string
		namespace                 string
		syncFreq                  = time.Minute
	)
	BeforeEach(func() {
		namespace = RandString(8)
		err := SetupKubeForTest(namespace)
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl = ""
	})
	AfterEach(func() {
		TeardownKube(namespace)
	})
	Describe("New", func() {
		It("creates a new client without error", func() {
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			_, err = NewStorage(cfg, namespace, syncFreq)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("Register", func() {
		It("registers the crds", func() {
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			client, err := NewStorage(cfg, namespace, syncFreq)
			Expect(err).NotTo(HaveOccurred())
			err = client.V1().Register()
			Expect(err).NotTo(HaveOccurred())
			apiextClient, err := apiexts.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			crds, err := apiextClient.ApiextensionsV1beta1().CustomResourceDefinitions().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			for _, crdSchema := range crdv1.KnownCRDs {
				var foundCrd *v1beta1.CustomResourceDefinition
				for _, crd := range crds.Items {
					if crd.Spec.Names.Kind == crdSchema.Kind {
						foundCrd = &crd
						break
					}
				}
				// if crd wasnt found, err
				Expect(foundCrd).NotTo(BeNil())

				Expect(foundCrd.Spec.Version).To(Equal(crdSchema.Version))
				Expect(foundCrd.Spec.Group).To(Equal(crdSchema.Group))
			}
		})
	})
	Describe("upstreams", func() {
		Describe("Create", func() {
			It("creates a crd from the item", func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				client, err := NewStorage(cfg, namespace, syncFreq)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Register()
				Expect(err).NotTo(HaveOccurred())
				upstream := NewTestUpstream1()
				createdUpstream, err := client.V1().Upstreams().Create(upstream)
				Expect(err).NotTo(HaveOccurred())
				upstream.Metadata = createdUpstream.GetMetadata()
				Expect(upstream).To(Equal(createdUpstream))
			})
		})
		Describe("Get", func() {
			It("gets a crd from the name", func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				client, err := NewStorage(cfg, namespace, syncFreq)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Register()
				Expect(err).NotTo(HaveOccurred())
				upstream := NewTestUpstream1()
				_, err = client.V1().Upstreams().Create(upstream)
				Expect(err).NotTo(HaveOccurred())
				created, err := client.V1().Upstreams().Get(upstream.Name)
				Expect(err).NotTo(HaveOccurred())
				upstream.Metadata = created.Metadata
				Expect(created).To(Equal(upstream))
			})
		})
		Describe("Update", func() {
			It("updates a crd from the item", func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				client, err := NewStorage(cfg, namespace, syncFreq)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Register()
				Expect(err).NotTo(HaveOccurred())
				upstream := NewTestUpstream1()
				created, err := client.V1().Upstreams().Create(upstream)
				Expect(err).NotTo(HaveOccurred())
				upstream.Metadata = created.GetMetadata()
				upstream.Type = "something-else"
				upstream.Metadata.Annotations["just_for_this_test"] = "bar"
				updated, err := client.V1().Upstreams().Update(upstream)
				Expect(err).NotTo(HaveOccurred())
				Expect(updated.Metadata.Annotations).To(HaveKey("just_for_this_test"))
				upstream.Metadata = updated.GetMetadata()
				Expect(updated).To(Equal(upstream))
			})
		})
		Describe("Delete", func() {
			It("deletes a crd from the name", func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				client, err := NewStorage(cfg, namespace, syncFreq)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Register()
				Expect(err).NotTo(HaveOccurred())
				upstream := NewTestUpstream1()
				_, err = client.V1().Upstreams().Create(upstream)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Upstreams().Delete(upstream.Name)
				Expect(err).NotTo(HaveOccurred())
				_, err = client.V1().Upstreams().Get(upstream.Name)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("virtualhosts", func() {
		Describe("Create", func() {
			It("creates a crd from the item", func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				client, err := NewStorage(cfg, namespace, syncFreq)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Register()
				Expect(err).NotTo(HaveOccurred())
				virtualhost := NewTestVirtualHost("something", NewTestRoute1())
				createdUpstream, err := client.V1().VirtualHosts().Create(virtualhost)
				Expect(err).NotTo(HaveOccurred())
				virtualhost.Metadata = createdUpstream.GetMetadata()
				Expect(virtualhost).To(Equal(createdUpstream))
			})
		})
		Describe("Get", func() {
			It("gets a crd from the name", func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				client, err := NewStorage(cfg, namespace, syncFreq)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Register()
				Expect(err).NotTo(HaveOccurred())
				vhost := NewTestVirtualHost("something", NewTestRoute1())
				_, err = client.V1().VirtualHosts().Create(vhost)
				Expect(err).NotTo(HaveOccurred())
				created, err := client.V1().VirtualHosts().Get(vhost.Name)
				Expect(err).NotTo(HaveOccurred())
				vhost.Metadata = created.Metadata
				Expect(created).To(Equal(vhost))
			})
		})
		Describe("Update", func() {
			It("updates a crd from the item", func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				client, err := NewStorage(cfg, namespace, syncFreq)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Register()
				Expect(err).NotTo(HaveOccurred())
				vHost := NewTestVirtualHost("something", NewTestRoute1())
				created, err := client.V1().VirtualHosts().Create(vHost)
				Expect(err).NotTo(HaveOccurred())
				// need to set resource ver
				vHost.Metadata = created.GetMetadata()
				vHost.Metadata.Annotations["just_for_this_test"] = "bar"
				updated, err := client.V1().VirtualHosts().Update(vHost)
				Expect(err).NotTo(HaveOccurred())
				Expect(updated.Metadata.Annotations).To(HaveKey("just_for_this_test"))
				vHost.Metadata = updated.GetMetadata()
				Expect(updated).To(Equal(vHost))
			})
		})
		Describe("Delete", func() {
			It("deletes a crd from the name", func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				client, err := NewStorage(cfg, namespace, syncFreq)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().Register()
				Expect(err).NotTo(HaveOccurred())
				vhost := NewTestVirtualHost("something", NewTestRoute1())
				_, err = client.V1().VirtualHosts().Create(vhost)
				Expect(err).NotTo(HaveOccurred())
				err = client.V1().VirtualHosts().Delete(vhost.Name)
				Expect(err).NotTo(HaveOccurred())
				_, err = client.V1().VirtualHosts().Get(vhost.Name)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
