package configwatcher_test

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

	. "github.com/solo-io/glue/implemented_modules/kube/configwatcher"
	clientset "github.com/solo-io/glue/implemented_modules/kube/configwatcher/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue/implemented_modules/kube/configwatcher/crd/solo.io/v1"
	"github.com/solo-io/glue/pkg/api/types/v1"
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
			route := &crdv1.Route{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "route-",
				},
				Spec: crdv1.DeepCopyRoute(NewTestRoute1()),
			}
			_, err = glueClient.GlueV1().Routes(namespace).Create(route)
			Expect(err).NotTo(HaveOccurred())

			// give controller time to register
			time.Sleep(time.Second * 2)

			var expectedRoute v1.Route
			data, err := json.Marshal(route.Spec)
			Expect(err).To(BeNil())
			err = json.Unmarshal(data, &expectedRoute)
			Expect(err).To(BeNil())
			select {
			case <-time.After(time.Second * 5):
				Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
			case cfg := <-watcher.Config():
				Expect(len(cfg.Routes)).To(Equal(1))
				Expect(cfg.Routes[0]).To(Equal(expectedRoute))
				Expect(cfg.Routes[0].Plugins["auth"]).To(Equal(expectedRoute.Plugins["auth"]))
			case err := <-watcher.Error():
				Expect(err).To(BeNil())
			}
		})
	})
})
