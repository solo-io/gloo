package make_test

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/test/testutils"
)

// Environment Variables which control the value of makefile vars
const (
	PublishContext = "PUBLISH_CONTEXT"

	// TODO: remove this once we are fully off of cloudbuild
	TestAssetId   = "TEST_ASSET_ID"
	TaggedVersion = "TAGGED_VERSION"
)

// Makefile vars
const (
	HelmBucket          = "HELM_BUCKET"
	QuayExpirationLabel = "QUAY_EXPIRATION_LABEL"
)

// NOTE: These tests are mostly to check that the makefile ifeq directive logic is working as expected.
// It's preferred that if Makefile logic gets anymore complicated, to direct the logic to go code rather than
// increase the complexity of the makefile.
var _ = Describe("Make", func() {
	Context("PUBLISH_CONTEXT", func() {
		It("Correctly sets variables when PUBLISH_CONTEXT is unset", func() {
			testutils.ExpectMakeVarsWithEnvVars([]*testutils.EnvVar{
				{PublishContext, ""},
				{TestAssetId, ""},
				{TaggedVersion, ""},
			}, []*testutils.MakeVar{
				{HelmBucket, "gs://solo-public-tagged-helm"},
				{QuayExpirationLabel, "--label quay.expires-after=3w"},
			})
		})

		It("Correctly sets variables when PUBLISH_CONTEXT is RELEASE", func() {
			testutils.ExpectMakeVarsWithEnvVars([]*testutils.EnvVar{
				{PublishContext, "RELEASE"},
				{TestAssetId, ""},
				{TaggedVersion, ""},
			}, []*testutils.MakeVar{
				{HelmBucket, "gs://solo-public-helm"},
				{QuayExpirationLabel, ""},
			})
		})

		It("Correctly sets variables when PUBLISH_CONTEXT is PULL_REQUEST", func() {
			testutils.ExpectMakeVarsWithEnvVars([]*testutils.EnvVar{
				{PublishContext, "PULL_REQUEST"},
				{TestAssetId, ""},
				{TaggedVersion, ""},
			}, []*testutils.MakeVar{
				{HelmBucket, "gs://solo-public-tagged-helm"},
				{QuayExpirationLabel, "--label quay.expires-after=3w"},
			})
		})
	})
})
