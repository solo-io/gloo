package chain

import (
	"context"
	"go/parser"

	"github.com/solo-io/gloo/test/matchers"

	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/golang/protobuf/ptypes/wrappers"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/ext-auth-plugins/api"
	chainmocks "github.com/solo-io/solo-projects/projects/extauth/pkg/config/chain/mocks"
)

//go:generate mockgen -destination mocks/auth_service_mock.go -package mocks github.com/solo-io/ext-auth-plugins/api AuthService

var _ = Describe("Plugin Chain", func() {

	// these tests aren't great since they test implementation, not the exposed interface
	Describe("mergeHeaders function", func() {

		var (
			buildHeader = func(key, value string, append bool) *envoycore.HeaderValueOption {
				return &envoycore.HeaderValueOption{
					Header: &envoycore.HeaderValue{
						Key:   key,
						Value: value,
					},
					Append: &wrappers.BoolValue{
						Value: append,
					},
				}
			}

			getFirst = func() *envoyauthv2.OkHttpResponse {
				return &envoyauthv2.OkHttpResponse{
					Headers: []*envoycore.HeaderValueOption{
						buildHeader("a", "foo", true),
						buildHeader("b", "bar", false),
						buildHeader("c", "baz", true),
					},
				}
			}

			getSecond = func() *envoyauthv2.OkHttpResponse {
				return &envoyauthv2.OkHttpResponse{
					Headers: []*envoycore.HeaderValueOption{
						buildHeader("b", "new-b", true),
						buildHeader("c", "new-c", false),
					},
				}
			}
		)

		It("works as expected", func() {
			result := mergeHeaders(getFirst(), getSecond())
			Expect(result.Headers).To(ConsistOf(
				buildHeader("a", "foo", true),
				buildHeader("b", "bar, new-b", true),
				buildHeader("c", "new-c", false),
			))

			result = mergeHeaders(getSecond(), getFirst())
			Expect(result.Headers).To(ConsistOf(
				buildHeader("a", "foo", true),
				buildHeader("b", "bar", false),
				buildHeader("c", "new-c, baz", true),
			))
		})

		It("covers edge cases", func() {
			result := mergeHeaders(getFirst(), nil)
			Expect(result).To(BeEquivalentTo(getFirst()))

			result = mergeHeaders(nil, getSecond())
			Expect(result).To(BeEquivalentTo(getSecond()))

			result = mergeHeaders(nil, nil)
			Expect(result).To(BeNil())
		})
	})

	Describe("plugin chain execution", func() {

		var (
			ctrl          *gomock.Controller
			pluginWrapper authServiceChain
			mockSvc1,
			mockSvc2,
			mockSvc3 *chainmocks.MockAuthService
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())

			mockSvc1 = chainmocks.NewMockAuthService(ctrl)
			mockSvc2 = chainmocks.NewMockAuthService(ctrl)
			mockSvc3 = chainmocks.NewMockAuthService(ctrl)

			pluginWrapper = authServiceChain{}
			err := pluginWrapper.AddAuthService("One", mockSvc1)
			Expect(err).NotTo(HaveOccurred())
			err = pluginWrapper.AddAuthService("Two", mockSvc2)
			Expect(err).NotTo(HaveOccurred())
			err = pluginWrapper.AddAuthService("Three", mockSvc3)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("fails when adding plugins with the same name", func() {
			err := pluginWrapper.AddAuthService("One", mockSvc1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(DuplicateAuthServiceNameError("One").Error()))
		})

		It("panics when adding plugins after Start has been called", func() {
			mockSvc1.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
			mockSvc2.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
			mockSvc3.EXPECT().Start(gomock.Any()).Return(nil).Times(1)

			err := pluginWrapper.Start(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(func() { _ = pluginWrapper.AddAuthService("Four", mockSvc1) }).To(Panic())
		})

		Describe("start functions", func() {

			Context("all functions succeed", func() {

				BeforeEach(func() {
					mockSvc1.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
					mockSvc2.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
					mockSvc3.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
				})

				It("runs all start functions", func() {
					err := pluginWrapper.Start(context.Background())
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("a function fails", func() {

				BeforeEach(func() {
					mockSvc1.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
					mockSvc2.EXPECT().Start(gomock.Any()).Return(errors.New("start failed")).Times(1)
					mockSvc3.EXPECT().Start(gomock.Any()).Return(nil).Times(0)
				})

				It("does not run start function after failure", func() {
					err := pluginWrapper.Start(context.Background())
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("authorize functions", func() {

			Context("all functions succeed", func() {

				BeforeEach(func() {
					mockSvc1.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(api.AuthorizedResponse(), nil).Times(1)
					mockSvc2.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(api.AuthorizedResponse(), nil).Times(1)
					mockSvc3.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(api.AuthorizedResponse(), nil).Times(1)
				})

				It("runs all authorize functions", func() {
					err := pluginWrapper.SetAuthorizer(nil)
					Expect(err).NotTo(HaveOccurred())
					response, err := pluginWrapper.Authorize(context.Background(), &api.AuthorizationRequest{})
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(BeEquivalentTo(api.AuthorizedResponse()))
				})

			})

			Context("a function fails", func() {

				BeforeEach(func() {
					mockSvc1.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
					mockSvc2.EXPECT().Start(gomock.Any()).Return(errors.New("start failed")).Times(1)
					mockSvc3.EXPECT().Start(gomock.Any()).Return(nil).Times(0)
				})

				It("does not run start function after failure", func() {
					err := pluginWrapper.Start(context.Background())
					Expect(err).To(HaveOccurred())
				})
			})

			// somewhat duplicated from "mergeHeaders function" context, but actually tests the exposed API
			// these tests would actually fail if we removed the merge headers logic from the Authorize() impl
			Context("headers", func() {
				var (
					buildHeader = func(key, value string, append bool) *envoycore.HeaderValueOption {
						return &envoycore.HeaderValueOption{
							Header: &envoycore.HeaderValue{
								Key:   key,
								Value: value,
							},
							Append: &wrappers.BoolValue{
								Value: append,
							},
						}
					}

					getFirst = func() *envoyauthv2.OkHttpResponse {
						return &envoyauthv2.OkHttpResponse{
							Headers: []*envoycore.HeaderValueOption{
								buildHeader("a", "foo", true),
								buildHeader("b", "bar", false),
								buildHeader("c", "baz", true),
							},
						}
					}

					getSecond = func() *envoyauthv2.OkHttpResponse {
						return &envoyauthv2.OkHttpResponse{
							Headers: []*envoycore.HeaderValueOption{
								buildHeader("b", "new-b", true),
								buildHeader("c", "new-c", false),
							},
						}
					}

					getThirdDenied = func() *envoyauthv2.DeniedHttpResponse {
						return &envoyauthv2.DeniedHttpResponse{
							Headers: []*envoycore.HeaderValueOption{
								buildHeader("denied", "val-1", true),
								buildHeader("also-denied", "val-2", false),
							},
						}
					}
				)
				BeforeEach(func() {
					f := api.AuthorizedResponse()
					f.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_OkResponse{OkResponse: getFirst()}

					s := api.AuthorizedResponse()
					s.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_OkResponse{OkResponse: getSecond()}

					mockSvc1.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(f, nil).Times(1)
					mockSvc2.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(s, nil).Times(1)
				})

				It("runs all authorize functions and merges headers from successful authentications", func() {
					mockSvc3.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(api.AuthorizedResponse(), nil).Times(1)

					err := pluginWrapper.SetAuthorizer(nil)
					Expect(err).NotTo(HaveOccurred())
					response, err := pluginWrapper.Authorize(context.Background(), &api.AuthorizationRequest{})
					Expect(err).NotTo(HaveOccurred())
					expected := api.AuthorizedResponse()
					expected.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_OkResponse{OkResponse: &envoyauthv2.OkHttpResponse{
						Headers: []*envoycore.HeaderValueOption{
							buildHeader("a", "foo", true),
							buildHeader("b", "bar, new-b", true),
							buildHeader("c", "new-c", false),
						},
					}}

					// non-deterministic order of headers means we can't use BeEquivalentTo on the entire response or else there will be test flakes
					// we assert headers, then nil them out and compare the responses
					Expect(response.CheckResponse.GetOkResponse().GetHeaders()).To(ConsistOf(expected.CheckResponse.GetOkResponse().GetHeaders()))

					r := response.CheckResponse.GetOkResponse()
					r.Headers = nil
					e := expected.CheckResponse.GetOkResponse()
					e.Headers = nil
					Expect(response).To(matchers.BeEquivalentToDiff(expected))
				})

				It("runs all authorize functions and returns headers on denied response", func() {
					t := api.UnauthorizedResponse()
					t.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_DeniedResponse{DeniedResponse: getThirdDenied()}
					mockSvc3.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(t, nil).Times(1)

					err := pluginWrapper.SetAuthorizer(nil)
					Expect(err).NotTo(HaveOccurred())
					response, err := pluginWrapper.Authorize(context.Background(), &api.AuthorizationRequest{})
					Expect(err).NotTo(HaveOccurred())
					expected := api.UnauthorizedResponse()
					expected.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_DeniedResponse{DeniedResponse: &envoyauthv2.DeniedHttpResponse{
						Headers: []*envoycore.HeaderValueOption{
							buildHeader("denied", "val-1", true),
							buildHeader("also-denied", "val-2", false),
						},
					}}

					// non-deterministic order of headers means we can't use BeEquivalentTo on the entire response or else there will be test flakes
					// we assert headers, then nil them out and compare the responses
					Expect(response.CheckResponse.GetDeniedResponse().GetHeaders()).To(ConsistOf(expected.CheckResponse.GetDeniedResponse().GetHeaders()))

					r := response.CheckResponse.GetDeniedResponse()
					r.Headers = nil
					e := expected.CheckResponse.GetDeniedResponse()
					e.Headers = nil
					Expect(response).To(matchers.BeEquivalentToDiff(expected))
				})
			})

			DescribeTable("complex boolean logic",
				func(expr string, expectedResponse, svc1Resp, svc2Resp, svc3Resp *api.AuthorizationResponse) {
					pluginWrapper = authServiceChain{}
					parsedExpr, parseErr := parser.ParseExpr(expr)
					Expect(parseErr).ToNot(HaveOccurred())
					err := pluginWrapper.SetAuthorizer(parsedExpr)
					Expect(err).NotTo(HaveOccurred())
					err = pluginWrapper.AddAuthService("One", mockSvc1)
					Expect(err).NotTo(HaveOccurred())
					err = pluginWrapper.AddAuthService("Two", mockSvc2)
					Expect(err).NotTo(HaveOccurred())
					err = pluginWrapper.AddAuthService("Three", mockSvc3)
					Expect(err).NotTo(HaveOccurred())

					if svc1Resp == nil {
						mockSvc1.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(svc1Resp, nil).Times(0)
					} else {
						mockSvc1.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(svc1Resp, nil).Times(1)
					}
					if svc2Resp == nil {
						mockSvc2.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(svc2Resp, nil).Times(0)
					} else {
						mockSvc2.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(svc2Resp, nil).Times(1)
					}
					if svc3Resp == nil {
						mockSvc3.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(svc3Resp, nil).Times(0)
					} else {
						mockSvc3.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(svc3Resp, nil).Times(1)
					}

					response, err := pluginWrapper.Authorize(context.Background(), &api.AuthorizationRequest{})
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(BeEquivalentTo(expectedResponse))
				},
				Entry("and all - all accepted", "One && Two && Three", api.AuthorizedResponse(), api.AuthorizedResponse(), api.AuthorizedResponse(), api.AuthorizedResponse()),
				Entry("and all - single reject", "One && Two && Three", api.UnauthorizedResponse(), api.AuthorizedResponse(), api.AuthorizedResponse(), api.UnauthorizedResponse()),
				Entry("or all - one accepted", "One || Two || Three", api.AuthorizedResponse(), api.UnauthorizedResponse(), api.UnauthorizedResponse(), api.AuthorizedResponse()),
				Entry("or all - single reject", "One || Two || Three", api.AuthorizedResponse(), api.UnauthorizedResponse(), api.AuthorizedResponse(), nil),
				Entry("short circuits left to right with and", "One && Two || Three", api.AuthorizedResponse(), api.AuthorizedResponse(), api.AuthorizedResponse(), nil),
				Entry("short circuits left to right with or", "One || Two && Three", api.AuthorizedResponse(), api.AuthorizedResponse(), nil, nil),
				Entry("and honors parens", "One || (Two && Three)", api.AuthorizedResponse(), nil, api.AuthorizedResponse(), api.AuthorizedResponse()),
				Entry("and honors parens - w/ short circuit", "One || (Two && Three)", api.AuthorizedResponse(), api.AuthorizedResponse(), api.UnauthorizedResponse(), nil),
				Entry("or honors parens", "One || (Two || Three)", api.AuthorizedResponse(), nil, api.AuthorizedResponse(), nil),
				Entry("not works - flipping unauthorized", "!One && !Two && !Three", api.AuthorizedResponse(), api.UnauthorizedResponse(), api.UnauthorizedResponse(), api.UnauthorizedResponse()),
				Entry("not works - flipping authorized", "!One && !Two && !Three", api.UnauthorizedResponse(), api.AuthorizedResponse(), nil, nil),
			)
		})
	})
})
