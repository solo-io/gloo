package controller_test

import (
	"os"
	"path/filepath"
	"time"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clientset "github.com/solo-io/glue/config/watcher/crd/client/clientset/versioned"
	informers "github.com/solo-io/glue/config/watcher/crd/client/informers/externalversions"
	. "github.com/solo-io/glue/config/watcher/crd/controller"
	"github.com/solo-io/glue/config/watcher/crd/solo.io/v1"
	. "github.com/solo-io/glue/test/helpers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-controller/pkg/signals"
)

var _ = Describe("Controller", func() {
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
		masterUrl, err = mkb.IP()
		Must(err)
	})
	AfterSuite(func() {
		err := mkb.Teardown()
		Must(err)
	})
	Describe("controller", func() {
		It("watches kube resources", func() {

			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			kubeClient, err := kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			glueClient, err := clientset.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			glueInformerFactory := informers.NewSharedInformerFactory(glueClient, time.Millisecond*30)

			controller := NewController(kubeClient, glueClient, glueInformerFactory)

			// set up signals so we handle the first shutdown signal gracefully
			stopCh := signals.SetupSignalHandler()

			go glueInformerFactory.Start(stopCh)

			go func() {
				err = controller.Run(2, stopCh)
				Must(err)
			}()
			// give controller time to register
			time.Sleep(time.Second * 2)

			route := &v1.Route{
				Spec: v1.DeepCopyRoute(*NewTestRoute1()),
			}

			_, err = glueClient.GlueV1().Routes(namespace).Create(route)
			Expect(err).NotTo(HaveOccurred())
			select {
			case <-time.After(time.Second * 5):
				Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
			case cfg := <-controller.Configs():
				Expect(len(cfg.Routes)).To(Equal(1))
				Expect(cfg.Routes[0]).To(Equal(route))
			}
		})
	})
})
