package istioconverter_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/glue/internal/istioconverter"
	"github.com/solo-io/glue/pkg/log"
	. "github.com/solo-io/glue/test/helpers"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = XDescribe("Istioconverter", func() {
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
	It("works", func() {
		cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
		Must(err)

		converter, err := NewIstioConverter(kubeconfigPath, "cluster.local", cfg, time.Millisecond, make(chan struct{}))
		Expect(err).NotTo(HaveOccurred())
		log.Printf("%v", converter)
		Must(<-converter.Error())
	})
})
