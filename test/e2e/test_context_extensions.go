package e2e

import (
	"context"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-projects/test/services/extauth"
	"github.com/solo-io/solo-projects/test/services/tap_server"

	"github.com/golang/protobuf/ptypes/duration"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gateway_defaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tap"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"google.golang.org/protobuf/types/known/durationpb"

	ratelimit2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/solo-projects/test/services/ratelimit"

	"github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

// This file contains the set of Extensions that are used to extend the TestContext
// These are used within the implementation of the TestContextWithExtensions

// rateLimitExtension is a TestContext extension that provides a RateLimitInstance
// It is used to isolate the RateLimit specific details from the TestContext
type rateLimitExtension struct {
	rateLimitInstance *ratelimit.Instance
}

// setupDefaults is used to set some default values for the RateLimitServer Settings
// It should be invoked during the BeforeEach phase of the test
func (r *rateLimitExtension) setupDefaults(testContext *TestContextWithExtensions) {
	serverUpstream := r.RateLimitInstance().GetServerUpstream()

	// Define some default values for the RateLimitServer Settings
	testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
		settings.RatelimitServer = &ratelimit2.Settings{
			RatelimitServerRef:      serverUpstream.GetMetadata().Ref(),
			DenyOnFail:              true, // ensures ConsistentlyNotRateLimited() calls will not pass unless server is healthy
			EnableXRatelimitHeaders: false,
			RateLimitBeforeAuth:     false,
			RequestTimeout: &duration.Duration{
				Seconds: 2,
			},
		}
	})

	testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, serverUpstream)
}

// RateLimitInstance returns an Instance of the RateLimit Service
func (r *rateLimitExtension) RateLimitInstance() *ratelimit.Instance {
	return r.rateLimitInstance
}

// runRateLimitService starts running the RateLimit Service
func (r *rateLimitExtension) runRateLimitService(testContext *TestContextWithExtensions) {
	ginkgo.By("RateLimitExtension: Running RateLimit")

	// The EnvoyInstance is similar to the RateLimitService, in that it receives its configuration from Gloo via xDS.
	// As a result, we ensure the RateLimitService uses the same port for connecting to Gloo
	ExpectWithOffset(1, testContext.EnvoyInstance().Port).NotTo(BeZero(), "EnvoyInstance.Port must be set before running RateLimitService")

	go func(testCtx context.Context) {
		defer ginkgo.GinkgoRecover()

		r.rateLimitInstance.RunWithXds(testCtx, testContext.EnvoyInstance().Port)
	}(testContext.Ctx())

	r.rateLimitInstance.EventuallyIsHealthy()
}

type extAuthExtension struct {
	extAuthInstance *extauth.Instance
}

// setupDefaults is used to set some default values for the RateLimitServer Settings
// It should be invoked during the BeforeEach phase of the test
func (r *extAuthExtension) setupDefaults(testContext *TestContextWithExtensions) {
	serverUpstream := r.ExtAuthInstance().GetServerUpstream()

	// Define some default values for the ExtAuth Settings
	testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
		settings.Extauth = &v1.Settings{
			ExtauthzServerRef: serverUpstream.GetMetadata().Ref(),
			FailureModeAllow:  false, // ensures requests will not succeed unless server is healthy
			RequestTimeout: &duration.Duration{
				Seconds: 2,
			},
			RequestBody: &v1.BufferSettings{
				MaxRequestBytes:     0,
				AllowPartialMessage: false,
				PackAsBytes:         false,
			},
			// Required for dynamic metadata emission to work
			TransportApiVersion: v1.Settings_V3,
		}
	})

	testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, serverUpstream)
}

// ExtAuthInstance returns an Instance of the ExtAuth Service
func (r *extAuthExtension) ExtAuthInstance() *extauth.Instance {
	return r.extAuthInstance
}

// runExtAuthService starts running the ExtAuth Service
func (r *extAuthExtension) runExtAuthService(testContext *TestContextWithExtensions) {
	ginkgo.By("ExtAuthExtension: Running ExtAuth")

	// The EnvoyInstance is similar to the ExtAuthService, in that it receives its configuration from Gloo via xDS
	// As a result, we ensure the ExtAuthService uses the same port for connecting to Gloo
	ExpectWithOffset(1, testContext.EnvoyInstance().Port).NotTo(BeZero(), "EnvoyInstance.Port must be set before running ExtAuthService")

	go func(testCtx context.Context) {
		defer ginkgo.GinkgoRecover()

		r.extAuthInstance.RunWithXds(testCtx, testContext.EnvoyInstance().Port)
	}(testContext.Ctx())

	r.extAuthInstance.EventuallyIsHealthy()
}

type tapServerExtension struct {
	tapServerInstance *tap_server.Instance
}

// TapServerInstance returns an Instance of the TapServer Service
func (r *tapServerExtension) TapServerInstance() *tap_server.Instance {
	return r.tapServerInstance
}

// setupDefaults simply configures the gateway with an upstream and a tap sink
func (r *tapServerExtension) setupDefaults(testContext *TestContextWithExtensions) {
	tapUpstream := testContext.TapServerInstance().GetServerUpstream()
	testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, tapUpstream)

	testGateway := gateway_defaults.DefaultGateway(defaults.GlooSystem)
	testGateway.GatewayType = &gatewayv1.Gateway_HttpGateway{
		HttpGateway: &gatewayv1.HttpGateway{
			Options: &gloov1.HttpListenerOptions{
				Tap: &tap.Tap{
					Sinks: []*tap.Sink{{
						SinkType: &tap.Sink_HttpService{
							HttpService: &tap.HttpService{
								TapServer: &core.ResourceRef{
									Name:      tapUpstream.Metadata.Name,
									Namespace: tapUpstream.Metadata.Namespace,
								},
								Timeout: &durationpb.Duration{
									Seconds: 5,
								},
							},
						},
					}},
				},
			},
		},
	}
	testContext.ResourcesToCreate().Gateways = gatewayv1.GatewayList{testGateway}
}

// start the tap server instance
func (r *tapServerExtension) runTapService(testContext *TestContextWithExtensions) {
	ginkgo.By("TapServerExtension: Starting Tap Server")

	go func(testCtx context.Context) {
		defer ginkgo.GinkgoRecover()

		r.tapServerInstance.Run(testCtx)
	}(testContext.Ctx())

	r.tapServerInstance.EventuallyIsHealthy()
}
