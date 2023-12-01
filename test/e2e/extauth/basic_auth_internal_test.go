package extauth_test

import (
	"net/http"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-projects/test/e2e"
)

// Helper function to return a VirtualService with the given AuthConfig
func createVsWithAuth(testContext *e2e.TestContextWithExtensions, authConfigRef *core.ResourceRef) *v1.VirtualService {
	return helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0]).
		WithVirtualHostOptions(&gloov1.VirtualHostOptions{
			Extauth: &extauth.ExtAuthExtension{
				Spec: &extauth.ExtAuthExtension_ConfigRef{
					ConfigRef: authConfigRef,
				},
			},
		}).
		Build()
}

// Helper function for creating a BasicAuth config based on credentials and encryption type
func createBasicAuthConfigInternalConfig(username, salt, hashedPassword string, encryption *extauth.BasicAuth_EncryptionType) *extauth.AuthConfig {
	return &extauth.AuthConfig{
		Metadata: &core.Metadata{
			Name:      GetBasicAuthExtension().GetConfigRef().Name,
			Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
		},
		Configs: []*extauth.AuthConfig_Config{{
			AuthConfig: &extauth.AuthConfig_Config_BasicAuth{
				BasicAuth: &extauth.BasicAuth{
					Encryption: encryption,
					UserSource: &extauth.BasicAuth_UserList_{
						UserList: &extauth.BasicAuth_UserList{
							Users: map[string]*extauth.BasicAuth_User{
								username: {
									Salt:           salt,
									HashedPassword: hashedPassword,
								},
							},
						},
					},
					Realm: "gloo",
				},
			},
		}},
	}
}

// configure an authconfig with the given credentials and encryption type and a virtual service to use it
func setup(username, salt, hashedPassword string, testContext *e2e.TestContextWithExtensions, encryption *extauth.BasicAuth_EncryptionType) {
	authConfig := createBasicAuthConfigInternalConfig(username, salt, hashedPassword, encryption)
	vsWithAuth := createVsWithAuth(testContext, authConfig.GetMetadata().Ref())
	testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{authConfig}
	testContext.ResourcesToCreate().VirtualServices = v1.VirtualServiceList{vsWithAuth}
}

var _ = Describe("BasicAuthInternal e2e tests for each algorithm", func() {
	type testEntry struct {
		message            string
		hostname           string
		expectedStatusCode int
	}

	// Use the same username and salt for each test
	const (
		username = "user"
		salt     = "0adzfifo"
	)

	var (
		testContext *e2e.TestContextWithExtensions
		// We can use the same test cases for each algorithm, so we define them here
		authTests = []testEntry{
			{
				message:            "should deny request without credentials",
				hostname:           "localhost",
				expectedStatusCode: http.StatusUnauthorized,
			},
			{
				message:            "should allow request with correct user/password",
				hostname:           username + ":password@localhost",
				expectedStatusCode: http.StatusOK,
			},
			{
				message:            "should reject request with incorrect user/password",
				hostname:           "invalid-user:invalid-password@localhost",
				expectedStatusCode: http.StatusUnauthorized,
			},
		}
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	EventuallyRequestWithHostnameReturnsStatus := func(requestBuilder *testutils.HttpRequestBuilder, hostname string, expectedStatusCode int) {
		httpRequestBuilder := testContext.GetHttpRequestBuilder().WithHostname(hostname + "@localhost")
		Eventually(func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(resp).Should(HaveHTTPStatus(expectedStatusCode))
			return resp
		}, "5s", "0.5s")
	}

	runTests := func(entries []testEntry) {
		for _, entry := range entries {
			By(entry.message, func() {
				EventuallyRequestWithHostnameReturnsStatus(
					testContext.GetHttpRequestBuilder(),
					entry.hostname,
					entry.expectedStatusCode)
			})
		}
	}

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	Context("with apr hashing (using before each)", func() {

		const (
			hashedPassword = "14o4fMw/Pm2L34SvyyA2r."
		)

		BeforeEach(func() {
			setup(username, salt, hashedPassword, testContext,
				&extauth.BasicAuth_EncryptionType{
					Algorithm: &extauth.BasicAuth_EncryptionType_Apr_{},
				},
			)
		})

		It("Should run all tests", func() {
			runTests(authTests)
		})
	})

	Context("with sha1 hashing (using before each)", func() {

		const (
			// use mixed case to test case insensitivity
			hashedPassword = "E892ea4b873c7c6f2071C784a8ae6df70a3fbcDB"
		)

		BeforeEach(func() {
			setup(username, salt, hashedPassword, testContext,
				&extauth.BasicAuth_EncryptionType{
					Algorithm: &extauth.BasicAuth_EncryptionType_Sha1_{},
				},
			)
		})

		It("Should run all tests", func() {
			runTests(authTests)
		})

	})
})
