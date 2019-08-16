package plugins

import (
	"context"
	"plugin"

	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
)

var _ = Describe("Plugin Loader", func() {

	var (
		requiredHeader    = "my-header"
		authorizedRequest = &envoyauthv2.CheckRequest{
			Attributes: &envoyauthv2.AttributeContext{
				Request: &envoyauthv2.AttributeContext_Request{
					Http: &envoyauthv2.AttributeContext_HttpRequest{
						Headers: map[string]string{
							requiredHeader: "value-a",
						},
					},
				},
			},
		}
		requestAuthorizedByFirstPluginOnly = &envoyauthv2.CheckRequest{
			Attributes: &envoyauthv2.AttributeContext{
				Request: &envoyauthv2.AttributeContext_Request{
					Http: &envoyauthv2.AttributeContext_HttpRequest{
						Headers: map[string]string{
							// had the header required by the first plugin, but not a value allowed by the second one
							requiredHeader: "not-allowed",
						},
					},
				},
			},
		}
		unauthorizedRequest = &envoyauthv2.CheckRequest{
			Attributes: &envoyauthv2.AttributeContext{
				Request: &envoyauthv2.AttributeContext_Request{
					Http: &envoyauthv2.AttributeContext_HttpRequest{},
				},
			},
		}
	)

	It("can load plugins", func() {
		ctx := context.Background()

		loader := NewPluginLoader(os.ExpandEnv("$GOPATH/src/github.com/solo-io/solo-projects/test/extauth/plugins"))
		svc, err := loader.Load(ctx, &extauth.PluginAuth{Plugins: []*extauth.AuthPlugin{
			{
				Name: "RequiredHeader",
				// No file name of symbol, check that defaults work correctly
				//PluginFileName:     "RequiredHeader.so",
				//ExportedSymbolName: "RequiredHeader",
				Config: &types.Struct{
					Fields: map[string]*types.Value{
						"RequiredHeader": {
							Kind: &types.Value_StringValue{
								StringValue: requiredHeader,
							},
						},
					},
				},
			},
			{
				Name:               "AllowedHeaderValues",
				PluginFileName:     "CheckHeaderValue.so",
				ExportedSymbolName: "Plugin",
				Config: &types.Struct{
					Fields: map[string]*types.Value{
						"RequiredHeader": {
							Kind: &types.Value_StringValue{
								StringValue: requiredHeader,
							},
						},
						"AllowedValues": {
							Kind: &types.Value_ListValue{
								ListValue: &types.ListValue{
									Values: []*types.Value{
										{
											Kind: &types.Value_StringValue{
												StringValue: "value-a",
											},
										},
										{
											Kind: &types.Value_StringValue{
												StringValue: "value-b",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}})

		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())

		Expect(svc.Start(ctx)).NotTo(HaveOccurred())

		By("send a request that is authorized by both plugins")
		resp, err := svc.Authorize(ctx, authorizedRequest)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())
		Expect(resp.CheckResponse.Status.Code).To(BeEquivalentTo(int32(rpc.OK)))
		// Check headers
		Expect(resp.CheckResponse.GetOkResponse()).NotTo(BeNil())
		// Each plugin appends a header if it accepts the request
		Expect(resp.CheckResponse.GetOkResponse().Headers).To(HaveLen(2))

		By("send a request that is authorized by the first plugins and denied by the second")
		resp, err = svc.Authorize(ctx, requestAuthorizedByFirstPluginOnly)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())
		Expect(resp.CheckResponse.Status.Code).To(BeEquivalentTo(int32(rpc.PERMISSION_DENIED)))

		By("send a request that is denied by both plugins")
		resp, err = svc.Authorize(ctx, unauthorizedRequest)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())
		Expect(resp.CheckResponse.Status.Code).To(BeEquivalentTo(int32(rpc.PERMISSION_DENIED)))
	})

	Describe("unpack function", func() {

		It("can unpack a valid symbol", func() {
			value := 123
			var symbol plugin.Symbol = &value

			i, err := unpack(symbol)
			Expect(err).NotTo(HaveOccurred())

			v, ok := i.(int)
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal(123))
		})

		It("fails when unpacking an invalid symbol", func() {
			var intAsInterface interface{} = 123
			var nonPointerSymbol plugin.Symbol = intAsInterface

			_, err := unpack(nonPointerSymbol)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("reflect: call of reflect.Value.Elem on int Value"))
		})

	})
})
