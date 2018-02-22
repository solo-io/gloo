package e2e

import (
	"time"

	"os"
	"path/filepath"

	"strings"

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
	namespace                 = crd.GlooDefaultNamespace
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
		gloo, err = crd.NewStorage(cfg, namespace, time.Minute)
		Must(err)
	})
	AfterSuite(func() {
		mkb.Teardown()
	})
})

type curlOpts struct {
	path   string
	method string
	host   string
}

func curlEventuallyShouldRespond(opts curlOpts, substr string, timeout ...time.Duration) {
	t := time.Second * 20
	if len(timeout) > 0 {
		t = timeout[0]
	}
	// for some useful-ish output
	tick := time.Tick(t / 8)
	Eventually(func() string {
		res, err := curlEnvoy(opts)
		if err != nil {
			res = err.Error()
		}
		select {
		default:
			break
		case <-tick:
			log.GreyPrintf("curl output: %v", res)
		}
		if strings.Contains(res, substr) {
			log.GreyPrintf("success: %v", res)
		}
		return res
	}, t).Should(ContainSubstring(substr))
}

func curlEnvoy(opts curlOpts) (string, error) {
	args := []string{"curl", "-v"}

	if opts.method != "GET" && opts.method != "" {
		args = append(args, "-X"+opts.method)
	}
	if opts.host != "" {
		args = append(args, "-H", "Host: "+opts.host)
	}
	args = append(args, "http://envoy:8080"+opts.path)
	return TestRunner(args...)
}
