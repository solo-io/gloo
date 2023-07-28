package clientside_sharding_test

import (
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	v12 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	"github.com/solo-io/gloo/test/helpers"

	"github.com/solo-io/solo-projects/test/kube2e"

	. "github.com/onsi/ginkgo/v2"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	. "github.com/solo-io/gloo/test/kube2e"
)

// This file is largely copied from test/kube2e/gloo-mtls/gloo_mtls_test.go (Oct 2020)

var _ = Describe("Installing gloo with mtls enabled & scaled redis", func() {

	var (
		testContext *kube2e.TestContext
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	It("can route request to upstream", func() {
		testContext.PatchDefaultVirtualService(func(service *v2.VirtualService) *v2.VirtualService {
			return helpers.BuilderFromVirtualService(service).WithRouteOptions(kube2e.DefaultRouteName, &v12.RouteOptions{
				PrefixRewrite: &wrappers.StringValue{
					Value: "/",
				},
			}).Build()
		})

		curlOpts := testContext.DefaultCurlOptsBuilder().WithConnectionTimeout(10).Build()
		testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, GetSimpleTestRunnerHttpResponse(), 1, time.Minute*5)
	})

})
