package crd_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/glue/internal/pkg/kube/storage"
	"github.com/solo-io/glue/internal/reporter"
	. "github.com/solo-io/glue/internal/reporter/crd"
	"github.com/solo-io/glue/pkg/api/types/v1"
	clientset "github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
	. "github.com/solo-io/glue/test/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("CrdReporter", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
		rptr                      reporter.Interface
	)
	BeforeSuite(func() {
		namespace = RandString(8)
		mkb = NewMinikube(false, namespace)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
	})
	AfterSuite(func() {
		mkb.Teardown()
	})
	Describe("writereports", func() {
		var (
			glueClient   clientset.Interface
			reports      []reporter.ConfigObjectReport
			upstreams    []v1.Upstream
			virtualHosts []v1.VirtualHost
		)
		Context("writes status reports for cfg crds with 0 errors", func() {
			BeforeEach(func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				RegisterCrds(cfg)
				glueClient, err = clientset.NewForConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
				rptr = NewKubeReporter(glueClient)

				testCfg := NewTestConfig()
				upstreams = testCfg.Upstreams
				var storables []v1.StorableConfigObject
				for _, us := range upstreams {
					crdUs := crdv1.UpstreamToCRD(metav1.ObjectMeta{
						Name:      us.Name,
						Namespace: namespace,
					}, us)
					_, err := glueClient.GlueV1().Upstreams(namespace).Create(&crdUs)
					Expect(err).NotTo(HaveOccurred())
					us.SetStorageRef(storage.CreateStorageRef(namespace, us.Name))
					storables = addUpstream(storables, us)
				}
				virtualHosts = testCfg.VirtualHosts
				for _, vHost := range virtualHosts {
					crdVhost := crdv1.VirtualHostToCRD(metav1.ObjectMeta{
						Name:      vHost.Name,
						Namespace: namespace,
					}, vHost)
					_, err := glueClient.GlueV1().VirtualHosts(namespace).Create(&crdVhost)
					Expect(err).NotTo(HaveOccurred())
					vHost.SetStorageRef(storage.CreateStorageRef(namespace, vHost.Name))
					storables = addVirtualHost(storables, vHost)
				}
				for _, storable := range storables {
					reports = append(reports, reporter.ConfigObjectReport{
						CfgObject: storable,
						Err:       nil,
					})
				}
			})

			It("writes an acceptance status for each crd", func() {
				err := rptr.WriteReports(reports)
				Expect(err).NotTo(HaveOccurred())
				updatedUpstreams, err := glueClient.GlueV1().Upstreams(namespace).List(metav1.ListOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedUpstreams.Items).To(HaveLen(len(upstreams)))
				for _, updatedUpstream := range updatedUpstreams.Items {
					Expect(updatedUpstream.Status.State).To(Equal(reporter.ObjectStateAccepted))
				}
				updatedVhosts, err := glueClient.GlueV1().VirtualHosts(namespace).List(metav1.ListOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedVhosts.Items).To(HaveLen(len(upstreams)))
				for _, updatedVhost := range updatedVhosts.Items {
					Expect(updatedVhost.Status.State).To(Equal(reporter.ObjectStateAccepted))
				}
			})
		})
	})
})

func addUpstream(storables []v1.StorableConfigObject, us v1.Upstream) []v1.StorableConfigObject {
	return append(storables, &us)
}

func addVirtualHost(storables []v1.StorableConfigObject, vHost v1.VirtualHost) []v1.StorableConfigObject {
	return append(storables, &vHost)
}
