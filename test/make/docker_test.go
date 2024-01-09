package make_test

import (
	"fmt"

	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// dockerImages contains the set of images that we attempt to build and publish during our CI pipeline
var dockerImages = []string{
	"gloo",
	"discovery",
	"gloo-envoy-wrapper",
	"certgen",
	"sds",
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

	// 	It("can push to quay.io for existing repository", func() {
	// 		// This functionality relies on permissions to push to quay.io, which is only enabled in CI
	// 		testutils.ValidateRequirementsAndNotifyGinkgo(testutils.DefinedEnv(testutils.GcloudBuildId))

	// 		err := testutils.DockerTag(StandardGlooImage, validTag)
	// 		Expect(err).NotTo(HaveOccurred(), "can re-tag image locally")

	// 		err = testutils.DockerPush(validTag)
	// 		Expect(err).NotTo(HaveOccurred(), "can push to quay.io for existing repository")
	// 	})

	// 	It("cannot push to quay.io for non-existent repository", func() {
	// 		err := testutils.DockerTag(StandardGlooImage, invalidTag)
	// 		Expect(err).NotTo(HaveOccurred(), "can re-tag image locally")

	// 		err = testutils.DockerPush(invalidTag)
	// 		Expect(err).To(HaveOccurred(), "can NOT push to quay.io for non-existent repository")
	// 	})

	// })

	Context("docker-push", func() {

		const (
			messageFmt = "The push refers to repository [quay.io/solo-io/"
		)

		It("Attempts to push every supported image", func() {
			fmt.Printf("Len(DockerImages): %d\n", len(dockerImages))
			// The point of this test is to ensure that the docker-push target attempts to push every image that we support
			// As a result, we attempt to push images, and then assert that the output contains the expected error message
			// NOTE: At the moment these images do not exist (since they aren't built in the same pipeline as our tests).
			// If that assumption changes, we may need to re-work these tests if the log output changes
			countMatcher := WithTransform(func(output string) int {
				return strings.Count(output, messageFmt)
			}, Equal(len(dockerImages)))

			// If this test fails, it's likely because you added a new image to the list of images that we attempt to build and publish during our CI pipeline
			// If this is true, you will need to request that we configure the Quay repository for this image
			// If you do not, we will hit a failure during our release pipeline:
			// https://github.com/solo-io/solo-projects/issues/5372#issuecomment-1732184633
			ExpectMakeOutputWithOffset(0, "docker-push --ignore-errors", countMatcher)
		})

	})

})
