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

	"github.com/solo-io/ext-auth-plugins/api"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
)

var _ = Describe("Plugin Loader", func() {

	var (
		pluginFileDir        = os.ExpandEnv("$GOPATH/src/github.com/solo-io/solo-projects/test/extauth/plugins")
		requiredHeader       = "my-header"
		allowedHeaderValues  = []string{"value-a", "value-b"}
		allowedHeadersPlugin = &extauth.AuthPlugin{
			Name:               "AllowedHeaderValues",
			PluginFileName:     "RequiredHeaderValue.so",
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
											StringValue: allowedHeaderValues[0],
										},
									},
									{
										Kind: &types.Value_StringValue{
											StringValue: allowedHeaderValues[1],
										},
									},
								},
							},
						},
					},
				},
			},
		}
		authorizedRequest = &api.AuthorizationRequest{
			CheckRequest: &envoyauthv2.CheckRequest{
				Attributes: &envoyauthv2.AttributeContext{
					Request: &envoyauthv2.AttributeContext_Request{
						Http: &envoyauthv2.AttributeContext_HttpRequest{
							Headers: map[string]string{
								requiredHeader: allowedHeaderValues[0],
							},
						},
					},
				},
			},
		}
		unauthorizedRequest = &api.AuthorizationRequest{
			CheckRequest: &envoyauthv2.CheckRequest{
				Attributes: &envoyauthv2.AttributeContext{
					Request: &envoyauthv2.AttributeContext_Request{
						Http: &envoyauthv2.AttributeContext_HttpRequest{},
					},
				},
			},
		}
	)

	It("can load a plugin", func() {
		ctx := context.Background()

		loader := NewPluginLoader(pluginFileDir)
		svc, err := loader.LoadAuthPlugin(ctx, allowedHeadersPlugin)

		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())

		Expect(svc.Start(ctx)).NotTo(HaveOccurred())

		By("send a request that is authorized the plugin")
		resp, err := svc.Authorize(ctx, authorizedRequest)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())
		Expect(resp.CheckResponse.Status.Code).To(BeEquivalentTo(int32(rpc.OK)))
		// Check headers
		Expect(resp.CheckResponse.GetOkResponse()).NotTo(BeNil())
		// Each plugin appends a header if it accepts the request
		Expect(resp.CheckResponse.GetOkResponse().Headers).To(HaveLen(1))

		By("send a request that is denied by the plugin")
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
