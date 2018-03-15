package e2e

import (
	"time"

	"os"
	"path/filepath"

	"strings"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	. "github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo/pkg/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespace = RandString(6)
)

var gloo storage.Interface
var kube kubernetes.Interface

var _ = Describe("Kubernetes Deployment", func() {
	BeforeSuite(func() {
		err := SetupKubeForE2eTest(namespace, true, true)
		Must(err)
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl := ""
		Must(err)
		cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
		Must(err)
		gloo, err = crd.NewStorage(cfg, namespace, time.Minute)
		Must(err)
		kube, err = kubernetes.NewForConfig(cfg)
		Must(err)
	})
	AfterSuite(func() {
		TeardownKube(namespace)
	})
})

type curlOpts struct {
	protocol string
	path     string
	method   string
	host     string
	caFile   string
	body     string
	headers  map[string]string
	port     int
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
			log.GreyPrintf("running: %v\nwant %v\nhave: %s", opts, substr, res)
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
	if opts.caFile != "" {
		args = append(args, "--cacert", opts.caFile)
	}
	if opts.body != "" {
		args = append(args, "-H", "Content-Type: application/json")
		args = append(args, "-d", opts.body)
	}
	for h, v := range opts.headers {
		args = append(args, "-H", fmt.Sprintf("%v: %v", h, v))
	}
	port := opts.port
	if port == 0 {
		port = 8080
	}
	protocol := opts.protocol
	if protocol == "" {
		protocol = "http"
	}
	args = append(args, fmt.Sprintf("%v://test-ingress:%v%s", protocol, port, opts.path))
	return TestRunner(args...)
}
