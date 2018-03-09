package configwatcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/gloo-storage/crd"
	. "github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo/pkg/log"
)

var _ = Describe("KubeConfigWatcher", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		//Skip("This test launches minikube and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.", 1)
		log.Printf("This test launches minikube and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
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

			storageClient, err := crd.NewStorage(cfg, namespace, time.Second)
			Expect(err).NotTo(HaveOccurred())

			watcher, err := NewConfigWatcher(storageClient)
			Must(err)
			go func() { watcher.Run(make(chan struct{})) }()

			virtualHost := NewTestVirtualHost("something", NewTestRoute1())
			created, err := storageClient.V1().VirtualHosts().Create(virtualHost)
			Expect(err).NotTo(HaveOccurred())

			// give controller time to register
			time.Sleep(time.Second * 2)

			select {
			case <-time.After(time.Second * 5):
				Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
			case cfg := <-watcher.Config():
				Expect(len(cfg.VirtualHosts)).To(Equal(1))
				Expect(cfg.VirtualHosts[0]).To(Equal(created))
				Expect(len(cfg.VirtualHosts[0].Routes)).To(Equal(1))
				Expect(cfg.VirtualHosts[0].Routes[0]).To(Equal(created.Routes[0]))
			case err := <-watcher.Error():
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})
})
