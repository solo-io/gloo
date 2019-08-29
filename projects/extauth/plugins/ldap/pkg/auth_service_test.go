package pkg_test

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/go-ldap/ldap"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/errors"
	impl "github.com/solo-io/solo-projects/projects/extauth/plugins/ldap/pkg"
	"github.com/solo-io/solo-projects/projects/extauth/plugins/ldap/pkg/mocks"

	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/solo-io/ext-auth-plugins/api"
)

var _ = Describe("AuthService", func() {

	var (
		ctrl              *gomock.Controller
		clientBuilderMock *mocks.MockClientBuilder
		clientMock        *mocks.MockClient
		authSvc           api.AuthService

		userDnTemplate = "uid=%s,ou=people,dc=solo,dc=io"
		allowedGroups  = []string{
			"cn=developers,ou=groups,dc=solo,dc=io",
			"cn=admins,ou=groups,dc=solo,dc=io",
		}
		config = &impl.Config{
			ServerUrl:      "localhost",
			UserDnTemplate: userDnTemplate,
			AllowedGroups:  allowedGroups,
		}

		requestWithHeaders = func(headers map[string]string) *api.AuthorizationRequest {
			return &api.AuthorizationRequest{
				CheckRequest: &envoyauthv2.CheckRequest{
					Attributes: &envoyauthv2.AttributeContext{
						Request: &envoyauthv2.AttributeContext_Request{
							Http: &envoyauthv2.AttributeContext_HttpRequest{
								Headers: headers,
							},
						},
					},
				},
			}
		}
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		clientMock = mocks.NewMockClient(ctrl)
		clientBuilderMock = mocks.NewMockClientBuilder(ctrl)
		authSvc = impl.NewLdapAuthService(clientBuilderMock, config)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("GetBasicAuthCredentials function", func() {

		It("retrieves credentials from correct header", func() {
			username, pwd, ok := impl.GetBasicAuthCredentials(requestWithHeaders(map[string]string{
				"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user:password"))),
			}))
			Expect(ok).To(BeTrue())
			Expect(username).To(Equal("user"))
			Expect(pwd).To(Equal("password"))
		})

		It("fails when no credentials can be retrieved", func() {
			_, _, ok := impl.GetBasicAuthCredentials(requestWithHeaders(map[string]string{}))
			Expect(ok).To(BeFalse())

			_, _, ok = impl.GetBasicAuthCredentials(requestWithHeaders(map[string]string{
				"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user/password"))),
			}))
			Expect(ok).To(BeFalse())

			_, _, ok = impl.GetBasicAuthCredentials(requestWithHeaders(map[string]string{
				"Auth": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user:password"))),
			}))
			Expect(ok).To(BeFalse())
		})
	})

	Describe("start function", func() {

		BeforeEach(func() {
			clientBuilderMock.EXPECT().Dial("tcp", "localhost").Return(clientMock, nil).Times(1)
			clientMock.EXPECT().Close().Times(1)
		})

		It("test connection on start", func() {
			err := authSvc.Start(context.Background())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("sanitize function", func() {

		It("does not sanitize regular DN values", func() {
			result, wasSanitized := impl.SanitizeLdapDN("john")
			Expect(wasSanitized).To(BeFalse())
			Expect(result).To(Equal("john"))
		})

		//regexp.MustCompile("[,=+<>#;\"]")
		It("sanitizes DN values that contain special characters", func() {
			result, wasSanitized := impl.SanitizeLdapDN(`a,b=c+d<e>f#g;h"i`)
			Expect(wasSanitized).To(BeTrue())
			Expect(result).To(Equal(`a\,b\=c\+d\<e\>f\#g\;h\"i`))
		})
	})

	Describe("authorize function", func() {

		expectStatus := func(resp *api.AuthorizationResponse, expectedCode rpc.Code) {
			ExpectWithOffset(1, resp).NotTo(BeNil())
			ExpectWithOffset(1, resp.CheckResponse).NotTo(BeNil())
			ExpectWithOffset(1, resp.CheckResponse.GetStatus().GetCode()).NotTo(BeNil())
			ExpectWithOffset(1, resp.CheckResponse.GetStatus().GetCode()).To(Equal(int32(expectedCode)))
		}

		Context("basic auth header is wrong", func() {

			It("returns a 401 response", func() {
				response, err := authSvc.Authorize(context.Background(), requestWithHeaders(map[string]string{
					"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("justpassword"))),
				}))
				Expect(err).NotTo(HaveOccurred())
				expectStatus(response, rpc.UNAUTHENTICATED)
			})
		})

		Context("user is unknown", func() {

			BeforeEach(func() {
				clientBuilderMock.EXPECT().Dial("tcp", "localhost").Return(clientMock, nil).Times(1)
				clientMock.EXPECT().Close().Times(1)
				clientMock.EXPECT().Bind("uid=john,ou=people,dc=solo,dc=io", "password").Return(errors.New("BIND failed")).Times(1)
			})

			It("returns a 401 response", func() {
				response, err := authSvc.Authorize(context.Background(), requestWithHeaders(map[string]string{
					"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("john:password"))),
				}))
				Expect(err).NotTo(HaveOccurred())
				expectStatus(response, rpc.UNAUTHENTICATED)
			})
		})

		Context("user is known", func() {

			var membership string

			JustBeforeEach(func() {
				clientBuilderMock.EXPECT().Dial("tcp", "localhost").Return(clientMock, nil).Times(1)
				clientMock.EXPECT().Close().Times(1)
				clientMock.EXPECT().Bind("uid=john,ou=people,dc=solo,dc=io", "password").Return(nil).Times(1)
				clientMock.EXPECT().Search(gomock.Any()).Return(
					&ldap.SearchResult{
						Entries: []*ldap.Entry{
							{
								Attributes: []*ldap.EntryAttribute{
									{
										Name:   impl.MembershipAttribute,
										Values: []string{membership},
									},
								},
							},
						},
					}, nil).Times(1)
			})

			Context("user is not member of any of the required groups", func() {

				BeforeEach(func() {
					membership = "cn=hr,ou=groups,dc=solo,dc=io"
				})

				It("returns a 403 response", func() {
					response, err := authSvc.Authorize(context.Background(), requestWithHeaders(map[string]string{
						"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("john:password"))),
					}))
					Expect(err).NotTo(HaveOccurred())
					expectStatus(response, rpc.PERMISSION_DENIED)
				})
			})

			Context("user is member of one of the required groups", func() {

				BeforeEach(func() {
					membership = "cn=developers,ou=groups,dc=solo,dc=io"
				})

				It("returns an OK response", func() {
					response, err := authSvc.Authorize(context.Background(), requestWithHeaders(map[string]string{
						"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("john:password"))),
					}))
					Expect(err).NotTo(HaveOccurred())
					expectStatus(response, rpc.OK)
				})
			})

		})

	})
})
