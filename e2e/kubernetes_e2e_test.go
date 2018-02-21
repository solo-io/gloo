package e2e

import (
	"time"

	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	. "github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo/pkg/log"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterUrl, kubeconfigPath string
	mkb                       *MinikubeInstance
)

var gloo storage.Interface

var _ = Describe("Kubernetes Deployment", func() {
	BeforeSuite(func() {
		mkb = NewMinikube(true)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
		cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
		Must(err)
		gloo, err = crd.NewStorage(cfg, crd.GlooDefaultNamespace, time.Minute)
		Must(err)
	})
	AfterSuite(func() {
		mkb.Teardown()
	})
})

func curlEventuallyShouldRespond(path, method, substr string, timeout ...time.Duration) {
	t := time.Second * 20
	if len(timeout) > 0 {
		t = timeout[0]
	}
	// for some useful-ish output
	tick := time.Tick(t / 8)
	Eventually(func() string {
		res, err := curlEnvoy(path, method)
		if err != nil {
			res = err.Error()
		}
		select {
		default:
			break
		case <-tick:
			log.Printf("curl output: %v", res)
		}
		return res
	}, t).Should(ContainSubstring(substr))
}

func curlEnvoy(path, method string) (string, error) {
	return TestRunner("curl", "-v", "-X"+method, "http://envoy:8080"+path)
}
