package kube_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/solo-io/glue/internal/configwatcher/kube"
	"github.com/solo-io/glue/internal/pkg/kube/storage"
	"github.com/solo-io/glue/pkg/api/types/v1"
	clientset "github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
	. "github.com/solo-io/glue/test/helpers"
)

var _ = Describe("KubeConfigWatcher", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
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
	Describe("controller", func() {
		It("watches kube crds", func() {
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			watcher, err := NewCrdWatcher(masterUrl, kubeconfigPath, time.Second, make(chan struct{}))
			Expect(err).NotTo(HaveOccurred())

			// add a route
			glueClient, err := clientset.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			virtualHost := &crdv1.VirtualHost{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "vhost-",
				},
				Spec: crdv1.DeepCopyVirtualHost(NewTestVirtualHost("vhost-1", NewTestRoute1())),
			}
			createdVirtualHost, err := glueClient.GlueV1().VirtualHosts(namespace).Create(virtualHost)
			Expect(err).NotTo(HaveOccurred())

			// give controller time to register
			time.Sleep(time.Second * 2)

			var expectedVhost v1.VirtualHost
			data, err := json.Marshal(virtualHost.Spec)
			Expect(err).To(BeNil())
			err = json.Unmarshal(data, &expectedVhost)
			Expect(err).To(BeNil())
			expectedVhost.SetStorageRef(storage.CreateStorageRef(namespace, createdVirtualHost.Name))
			select {
			case <-time.After(time.Second * 5):
				Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
			case cfg := <-watcher.Config():
				Expect(len(cfg.VirtualHosts)).To(Equal(1))
				Expect(cfg.VirtualHosts[0]).To(Equal(expectedVhost))
				Expect(len(cfg.VirtualHosts[0].Routes)).To(Equal(1))
				Expect(cfg.VirtualHosts[0].Routes[0]).To(Equal(expectedVhost.Routes[0]))
				Expect(cfg.VirtualHosts[0].GetStorageRef()).To(Equal(storage.CreateStorageRef(namespace, createdVirtualHost.Name)))
			case err := <-watcher.Error():
				Expect(err).To(BeNil())
			}
		})
	})
})
