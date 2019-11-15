package converter_test

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	. "github.com/solo-io/go-utils/testutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter"
)

var (
	virtualServiceConverter converter.VirtualServiceDetailsConverter
	mockCtrl                *gomock.Controller
	rawGetter               *mocks.MockRawGetter
)

var _ = Describe("VirtualServiceDetailsConverter", func() {
	getVirtualService := func(pluginConfigs map[string]*types.Struct) *gatewayv1.VirtualService {
		return &gatewayv1.VirtualService{
			VirtualHost: &gatewayv1.VirtualHost{
				VirtualHostPlugins: &gloov1.VirtualHostPlugins{
					Extensions: &gloov1.Extensions{
						Configs: pluginConfigs,
					},
				},
			},
		}
	}

	getExpectedRaw := func() *v1.Raw {
		return &v1.Raw{FileName: "fn", Content: "ct"}
	}

	Describe("List", func() {
		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			rawGetter = mocks.NewMockRawGetter(mockCtrl)
			virtualServiceConverter = converter.NewVirtualServiceDetailsConverter(rawGetter)
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("works", func() {
			for _, testCase := range []struct {
				desc            string
				configs         map[string]*types.Struct
				expectedPlugins *v1.Plugins
				expectedRaw     *v1.Raw
			}{
				{
					desc:            "for a nil config",
					configs:         nil,
					expectedPlugins: nil,
					expectedRaw:     getExpectedRaw(),
				},
			} {
				virtualService := getVirtualService(testCase.configs)
				if testCase.expectedRaw != nil {
					rawGetter.EXPECT().
						GetRaw(context.Background(), virtualService, gatewayv1.VirtualServiceCrd).
						Return(testCase.expectedRaw)
				}
				actual := virtualServiceConverter.GetDetails(context.Background(), virtualService)
				expected := &v1.VirtualServiceDetails{
					VirtualService: virtualService,
					Plugins:        testCase.expectedPlugins,
					Raw:            testCase.expectedRaw,
				}
				ExpectEqualProtoMessages(actual, expected, testCase.desc)
			}
		})
	})
})
