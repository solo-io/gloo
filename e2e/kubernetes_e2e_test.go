package e2e

import (
	"time"

	"os"
	"path/filepath"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	. "github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo/pkg/log"
	"k8s.io/client-go/tools/clientcmd"
)

const helloService = "helloservice"
const servicePort = 8080

var _ = Describe("Kubernetes Deployment", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
	)
	BeforeSuite(func() {
		mkb = NewMinikube(true)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
	})
	AfterSuite(func() {
		mkb.Teardown()
	})
	Describe("E2e", func() {
		Describe("updating gloo config", func() {
			var gloo storage.Interface
			BeforeEach(func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Must(err)
				gloo, err = crd.NewStorage(cfg, crd.GlooDefaultNamespace, time.Minute)
				Must(err)
			})
			//It("responds 503 for a route with misconfigured upstream", func() {
			//	curlEventuallyShouldRespond("/broken", "< HTTP/1.1 503", time.Minute*3)
			//})
			Context("update gloo with new rules", func() {
				randomPath := "/" + uuid.New()
				//It("responds 404 before update", func() {
				//	curlEventuallyShouldRespond(randomPath, "< HTTP/1.1 404", time.Minute*15)
				//})
				It("responds 200 after update", func() {
					_, err := gloo.V1().Upstreams().Create(&v1.Upstream{
						Name: helloService,
						Type: kubernetes.UpstreamTypeKube,
						Spec: kubernetes.EncodeUpstreamSpec(kubernetes.UpstreamSpec{
							ServiceNamespace: crd.GlooDefaultNamespace,
							ServiceName:      helloService,
							ServicePort:      fmt.Sprintf("%v", servicePort),
						}),
					})
					Must(err)
					_, err = gloo.V1().VirtualHosts().Create(&v1.VirtualHost{
						Name: "one-route",
						Routes: []*v1.Route{{
							Matcher: &v1.Route_RequestMatcher{
								RequestMatcher: &v1.RequestMatcher{
									Path: &v1.RequestMatcher_PathExact{
										PathExact: randomPath,
									},
									Verbs: []string{"GET"},
								},
							},
							SingleDestination: &v1.Destination{
								DestinationType: &v1.Destination_Upstream{
									Upstream: &v1.UpstreamDestination{
										Name: helloService,
									},
								},
							},
						}},
					})
					Must(err)
					curlEventuallyShouldRespond(randomPath, "< HTTP/1.1 200", time.Minute*5)
				})
			})

		})
	})
})

func curlEventuallyShouldRespond(path, substr string, timeout ...time.Duration) {
	t := time.Second * 20
	if len(timeout) > 0 {
		t = timeout[0]
	}
	// for some useful-ish output
	tick := time.Tick(t / 8)
	Eventually(func() string {
		res, err := curlEnvoy(path)
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

func curlEnvoy(path string) (string, error) {
	return TestRunner("curl", "-v", "http://envoy:8080"+path)
}
