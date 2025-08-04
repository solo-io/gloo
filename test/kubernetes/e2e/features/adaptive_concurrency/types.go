package adaptive_concurrency

import (
	"net/http"
	"path/filepath"

	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	e2edefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/skv2/codegen/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	k8sProxySvcName      = "gloo-proxy-gw"
	k8sProxySvcNamespace = "default"
)

var (
	edgeGatewayManifest            = filepath.Join(util.MustGetThisDir(), "testdata", "edge", "edge-gateway.yaml")
	sleepVirtualServiceManifest    = filepath.Join(util.MustGetThisDir(), "testdata", "edge", "sleep-server-vses.yaml")
	gg2SetupManifest               = filepath.Join(util.MustGetThisDir(), "testdata", "gg2", "setup.yaml")
	gg2HttpListenerOptionsManifest = filepath.Join(util.MustGetThisDir(), "testdata", "gg2", "http-listener-options.yaml")

	setupSuite = base.SimpleTestCase{
		Manifests: []string{
			e2edefaults.SleepServerPodManifest,
			e2edefaults.CurlPodManifest,
		},
		Resources: []client.Object{
			e2edefaults.CurlPod, e2edefaults.SleepServerPod, e2edefaults.SleepServerService,
		},
	}

	edgeTestCases = map[string]*base.TestCase{
		"TestAdaptiveConcurrency": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{
					edgeGatewayManifest,
					sleepVirtualServiceManifest,
				},
			},
		},
	}

	gg2TestCases = map[string]*base.TestCase{
		"TestAdaptiveConcurrency": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{
					gg2SetupManifest,
					gg2HttpListenerOptionsManifest,
				},
			},
		},
	}

	okResponse = &testmatchers.HttpResponse{
		StatusCode: http.StatusOK,
	}

	okOrRateLimitedResponse = &testmatchers.HttpResponse{
		StatusCode: []int{http.StatusOK, http.StatusServiceUnavailable},
	}
)
