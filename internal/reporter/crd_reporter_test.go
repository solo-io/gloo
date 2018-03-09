package reporter_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-testing/helpers"
	. "github.com/solo-io/gloo/internal/reporter"
	"github.com/solo-io/gloo/pkg/log"
)

var _ = Describe("CrdReporter", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		//Skip("This test launches minikube and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.", 1)
		log.Printf("This test launches minikube and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
		rptr                      Interface
	)
	BeforeEach(func() {
		namespace = RandString(8)
		mkb = NewMinikube(false, namespace)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
	})
	AfterEach(func() {
		mkb.Teardown()
	})
	Describe("writereports", func() {
		var (
			glooClient   storage.Interface
			reports      []ConfigObjectReport
			upstreams    []*v1.Upstream
			virtualHosts []*v1.VirtualHost
		)
		Context("writes status reports for cfg crds with 0 errors", func() {
			BeforeEach(func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				glooClient, err = crd.NewStorage(cfg, namespace, time.Second)
				Expect(err).NotTo(HaveOccurred())
				rptr = NewReporter(glooClient)

				testCfg := NewTestConfig()
				upstreams = testCfg.Upstreams
				var storables []v1.ConfigObject
				for _, us := range upstreams {
					_, err := glooClient.V1().Upstreams().Create(us)
					Expect(err).NotTo(HaveOccurred())
					storables = append(storables, us)
				}
				virtualHosts = testCfg.VirtualHosts
				for _, vHost := range virtualHosts {
					_, err := glooClient.V1().VirtualHosts().Create(vHost)
					Expect(err).NotTo(HaveOccurred())
					storables = append(storables, vHost)
				}
				for _, storable := range storables {
					reports = append(reports, ConfigObjectReport{
						CfgObject: storable,
						Err:       nil,
					})
				}
			})

			It("writes an acceptance status for each crd", func() {
				err := rptr.WriteReports(reports)
				Expect(err).NotTo(HaveOccurred())
				updatedUpstreams, err := glooClient.V1().Upstreams().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedUpstreams).To(HaveLen(len(upstreams)))
				for _, updatedUpstream := range updatedUpstreams {
					Expect(updatedUpstream.Status.State).To(Equal(v1.Status_Accepted))
				}
				updatedVhosts, err := glooClient.V1().VirtualHosts().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedVhosts).To(HaveLen(len(upstreams)))
				for _, updatedVhost := range updatedVhosts {
					Expect(updatedVhost.Status.State).To(Equal(v1.Status_Accepted))
				}
			})
		})
		Context("writes status reports for cfg crds with SOME errors", func() {
			BeforeEach(func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				glooClient, err = crd.NewStorage(cfg, namespace, time.Second)
				Expect(err).NotTo(HaveOccurred())
				rptr = NewReporter(glooClient)

				testCfg := NewTestConfig()
				upstreams = testCfg.Upstreams
				var storables []v1.ConfigObject
				for _, us := range upstreams {
					_, err := glooClient.V1().Upstreams().Create(us)
					Expect(err).NotTo(HaveOccurred())
					storables = append(storables, us)
				}
				virtualHosts = testCfg.VirtualHosts
				for _, vHost := range virtualHosts {
					_, err := glooClient.V1().VirtualHosts().Create(vHost)
					Expect(err).NotTo(HaveOccurred())
					storables = append(storables, vHost)
				}
				for _, storable := range storables {
					reports = append(reports, ConfigObjectReport{
						CfgObject: storable,
						Err:       errors.New("oh no an error what did u do!"),
					})
				}
			})

			It("writes an rejected status for each crd", func() {
				err := rptr.WriteReports(reports)
				Expect(err).NotTo(HaveOccurred())
				updatedUpstreams, err := glooClient.V1().Upstreams().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedUpstreams).To(HaveLen(len(upstreams)))
				for _, updatedUpstream := range updatedUpstreams {
					Expect(updatedUpstream.Status.State).To(Equal(v1.Status_Rejected))
				}
				updatedVhosts, err := glooClient.V1().VirtualHosts().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedVhosts).To(HaveLen(len(upstreams)))
				for _, updatedVhost := range updatedVhosts {
					Expect(updatedVhost.Status.State).To(Equal(v1.Status_Rejected))
				}
			})
		})
	})
})
