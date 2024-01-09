package make_test

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/test/testutils"
)

// dockerImages contains the set of images that we attempt to build and publish during our CI pipeline
var dockerImages = []string{
	"gloo",
	"discovery",
	"gloo-envoy-wrapper",
	"certgen",
	"sds",
	//"sds-fips",
	"ingress",
	"access-logger",
	"kubectl",
}

var _ = Describe("Docker", func() {

	// Context("docker-push-%", func() {
	// 	const (
	// 		validTag   = "quay.io/solo-io/gloo:1.0.0-docker-push-test"
	// 		invalidTag = "quay.io/solo-io/gloo-invalid:1.0.0-docker-push-test"
	// 	)

	// 	testutils.DockerPushTest(validTag, invalidTag, StandardGlooImage)
	// })

	Context("docker-push", func() {

		const (
			messageFmt = "The push refers to repository [quay.io/solo-io/"
		)

		testutils.DockerPushImagesTest(messageFmt, dockerImages)

	})

})
