package make_test

import (
	. "github.com/onsi/ginkgo/v2"
)

const (
	PublishContext = "PUBLISH_CONTEXT"
	UseDigests     = "USE_DIGESTS"
)

// NOTE: These tests are mostly to check that the makefile ifeq directive logic is working as expected.
// It's preferred that if Makefile logic gets anymore complicated, to direct the logic to go code rather than
// increase the complexity of the makefile.
var _ = Describe("Make", func() {
	Context("Determines PUBLISH_CONTEXT from Cloudbuild Env Variables", func() {
		const (
			TestAssetId   = "TEST_ASSET_ID"
			TaggedVersion = "TAGGED_VERSION"
		)

		It("None, when Cloudbuild variables are unset", func() {
			ExpectMakeVarsWithEnvVars([]*EnvVar{
				{TestAssetId, ""},
				{TaggedVersion, ""},
			}, []*MakeVar{
				{PublishContext, "NONE"},
			})
		})

		It("PullRequest, when Cloudbuild TEST_ASSET_ID is set, but TAGGED_VERSION is not", func() {
			ExpectMakeVarsWithEnvVars([]*EnvVar{
				{TestAssetId, "test-asset-id"},
				{TaggedVersion, ""},
			}, []*MakeVar{
				{PublishContext, "PULL_REQUEST"},
			})
		})

		It("Release, when Cloudbuild TAGGED_VERSION and TEST_ASSET_ID are set", func() {
			ExpectMakeVarsWithEnvVars([]*EnvVar{
				{TestAssetId, "test-asset-id"},
				{TaggedVersion, "v1.0.0"},
			}, []*MakeVar{
				{PublishContext, "RELEASE"},
				// We also include a variable which should be set when PUBLISH_CONTEXT is RELEASE
				// This is done to gain further confidence that variables are properly gated by PUBLISH_CONTEXT
				{UseDigests, "--use-digests"},
			})
		})
	})

	Context("Determines Make targets and variables from PUBLISH_CONTEXT", func() {

		When("PUBLISH_CONTEXT is NONE", func() {
			It("does not set USE_DIGESTS variable", func() {
				ExpectMakeVarsWithEnvVars([]*EnvVar{
					{PublishContext, "NONE"},
				}, []*MakeVar{
					{UseDigests, ""},
				})
			})

		})

		When("PUBLISH_CONTEXT is RELEASE", func() {
			It("sets USE_DIGESTS variable", func() {
				ExpectMakeVarsWithEnvVars([]*EnvVar{
					{PublishContext, "RELEASE"},
				}, []*MakeVar{
					{UseDigests, "--use-digests"},
				})
			})

		})

	})
})
