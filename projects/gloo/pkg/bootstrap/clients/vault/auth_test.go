package vault

import (
    "context"
    "github.com/avast/retry-go"
    "github.com/golang/mock/gomock"
    "github.com/hashicorp/vault/api"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/onsi/gomega/types"
    "github.com/rotisserie/eris"
    v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
    "github.com/solo-io/gloo/projects/gloo/pkg/bootstrap/clients/vault/mocks"
    "github.com/solo-io/gloo/test/gomega/assertions"
    "go.opencensus.io/stats/view"
)

var metricViews = []*view.View{
    mLastLoginSuccessView,
    mLoginFailuresView,
    mLoginSuccessesView,
    mLastLoginFailureView,
}

var _ = Describe("ClientAuth", func() {

    var (
        ctx    context.Context
        cancel context.CancelFunc

        clientAuth ClientAuth
    )

    BeforeEach(func() {
        ctx, cancel = context.WithCancel(context.Background())

        // The tests below will be responsible for assigning this variable
        // We re-set it here, just to be safe
        clientAuth = nil

        // We should not have any metrics set before running the tests
        // This ensures that we are no leaking metrics between tests
        resetViews()
    })

    AfterEach(func() {
        cancel()
    })

    Context("newStaticTokenAuth", func() {
        // These tests validate the behavior of the staticTokenAuth implementation of the ClientAuth interface

        When("token is empty", func() {

            BeforeEach(func() {
                clientAuth = newStaticTokenAuth("")
            })

            It("login should return an error", func() {
                secret, err := clientAuth.Login(ctx, nil)
                Expect(err).To(MatchError(ErrEmptyToken))
                Expect(secret).To(BeNil())

                assertions.ExpectStatLastValueMatches(mLastLoginFailure, Not(BeZero()))
                assertions.ExpectStatSumMatches(mLoginFailures, Equal(1))
            })

            It("startRenewal should return nil", func() {
                err := clientAuth.StartRenewal(ctx, nil)
                Expect(err).NotTo(HaveOccurred())
            })

        })

        When("token is not empty", func() {

            BeforeEach(func() {
                clientAuth = newStaticTokenAuth("placeholder")
            })

            It("should return a vault.Secret", func() {
                secret, err := clientAuth.Login(ctx, nil)
                Expect(err).NotTo(HaveOccurred())
                Expect(secret).To(Equal(&api.Secret{
                    Auth: &api.SecretAuth{
                        ClientToken: "placeholder",
                    },
                }))
                assertions.ExpectStatLastValueMatches(mLastLoginSuccess, Not(BeZero()))
                assertions.ExpectStatSumMatches(mLoginSuccesses, Equal(1))
            })

            It("startRenewal should return nil", func() {
                err := clientAuth.StartRenewal(ctx, nil)
                Expect(err).NotTo(HaveOccurred())
            })

        })

    })

    Context("newRemoteTokenAuth", func() {
        // These tests validate the behavior of the remoteTokenAuth implementation of the ClientAuth interface

        When("internal auth method always returns an error", func() {

            BeforeEach(func() {
                ctrl := gomock.NewController(GinkgoT())
                internalAuthMethod := mocks.NewMockAuthMethod(ctrl)
                internalAuthMethod.EXPECT().Login(ctx, gomock.Any()).Return(nil, eris.New("mocked error message")).AnyTimes()

                clientAuth = newRemoteTokenAuth(internalAuthMethod, retry.Attempts(3))
            })

            It("should return the error", func() {
                secret, err := clientAuth.Login(ctx, nil)
                Expect(err).To(MatchError("unable to authenticate to vault: mocked error message"))
                Expect(secret).To(BeNil())

                assertions.ExpectStatLastValueMatches(mLastLoginFailure, Not(BeZero()))
                assertions.ExpectStatSumMatches(mLoginFailures, Equal(3))
            })

        })

        When("internal auth method returns an error, and then a success", func() {

            BeforeEach(func() {
                ctrl := gomock.NewController(GinkgoT())
                internalAuthMethod := mocks.NewMockAuthMethod(ctrl)
                internalAuthMethod.EXPECT().Login(ctx, gomock.Any()).Return(nil, eris.New("error")).Times(1)
                internalAuthMethod.EXPECT().Login(ctx, gomock.Any()).Return(&api.Secret{
                    Auth: &api.SecretAuth{
                        ClientToken: "a-client-token",
                    },
                }, nil).Times(1)

                clientAuth = newRemoteTokenAuth(internalAuthMethod, retry.Attempts(2))
            })

            It("should return a secret", func() {
                secret, err := clientAuth.Login(ctx, nil)
                Expect(err).NotTo(HaveOccurred())
                Expect(secret.Auth.ClientToken).To(Equal("a-client-token"))

                assertions.ExpectStatSumMatches(mLoginFailures, Equal(1))
                assertions.ExpectStatSumMatches(mLoginSuccesses, Equal(1))
            })

        })

    })

    Context("newClientAuthForSettings", func() {
        // These tests validate that the constructor maps the Gloo Settings into the appropriate vault.AuthMethod interface
        // it does not test the underlying implementations, as those are handled in the above tests

        DescribeTable("should error on invalid inputs",
            func(vaultSettings *v1.Settings_VaultSecrets, expectedError types.GomegaMatcher) {
                _, err := newClientAuthForSettings(ctx, vaultSettings)
                Expect(err).To(expectedError)
            },
            Entry("partial accessKey / secretAccessKey", &v1.Settings_VaultSecrets{
                AuthMethod: &v1.Settings_VaultSecrets_Aws{
                    Aws: &v1.Settings_VaultAwsAuth{
                        AccessKeyId:     "access-key-id",
                        SecretAccessKey: "",
                    },
                },
            }, HaveOccurred()), // this should be improved to be more explicit
        )

        DescribeTable("should return the correct auth method",
            func(vaultSettings *v1.Settings_VaultSecrets, expectedAuthMethod types.GomegaMatcher) {
                authMethod, err := newClientAuthForSettings(ctx, vaultSettings)
                Expect(err).NotTo(HaveOccurred())
                Expect(authMethod).To(expectedAuthMethod)
            },
            Entry("nil", nil, Not(BeNil())), // this should be improved
        )

    })
})
