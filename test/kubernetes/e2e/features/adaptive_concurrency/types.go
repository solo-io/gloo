package adaptiveconcurrency

import (
	"net/http"
	"path/filepath"
	"runtime"

	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	e2edefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/skv2/codegen/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	envoyAdminPort = 19000
	envoyStatsPath = "/stats?filter=adaptive_concurrency"
)

var (
	sleepGatewayManifest        = filepath.Join(util.MustGetThisDir(), "testdata", "gateway.yaml")
	sleepVirtualServiceManifest = filepath.Join(util.MustGetThisDir(), "testdata", "sleep-server-vses.yaml")

	setupSuite = base.SimpleTestCase{
		Manifests: []string{
			e2edefaults.SleepServerPodManifest,
			e2edefaults.CurlPodManifest,
		},
		Resources: []client.Object{
			e2edefaults.CurlPod, e2edefaults.SleepServerPod, e2edefaults.SleepServerService,
		},
	}

	testCases = map[string]*base.TestCase{
		"TestAdaptiveConcurrency": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{
					sleepGatewayManifest,
					sleepVirtualServiceManifest,
				},
			},
		},
	}

	okResponse = &testmatchers.HttpResponse{
		StatusCode: http.StatusOK,
	}

	okOrUnavailableResponse = &testmatchers.HttpResponse{
		StatusCode: []int{http.StatusOK, http.StatusServiceUnavailable},
	}

	localClusterDomain = setLocalClusterDomain()
)

func setLocalClusterDomain() string {
	extauthAddr := "localhost"
	if runtime.GOOS == "darwin" {
		extauthAddr = "host.docker.internal"
	}
	return extauthAddr
}
