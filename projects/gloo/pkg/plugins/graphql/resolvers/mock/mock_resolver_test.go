package mock_test

import (
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	structpb "github.com/golang/protobuf/ptypes/struct"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/test/matchers"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/mock"
)

var _ = Describe("translates mock resolver", func() {
	getResolverConfig := func(config *v3.TypedExtensionConfig) *v2.StaticResolver {
		mockResolver, err := utils.AnyToMessage(config.GetTypedConfig())
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		ExpectWithOffset(1, mockResolver).NotTo(BeNil())
		ret, ok := mockResolver.(*v2.StaticResolver)
		ExpectWithOffset(1, ok).To(BeTrue())
		return ret
	}

	getStructValueFromJSON := func(JSON string) *structpb.Value {
		msg := &structpb.Value{}
		err := msg.UnmarshalJSON([]byte(JSON))
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		return msg
	}

	Context("sync response", func() {
		It("translates sync response", func() {
			message := `{"response":"test"}`
			mockResolver := &v1beta1.MockResolver{
				Response: &v1beta1.MockResolver_SyncResponse{
					SyncResponse: getStructValueFromJSON(message),
				},
			}
			envoyResolver, err := mock.TranslateMockResolver(mockResolver)
			Expect(err).NotTo(HaveOccurred())
			Expect(envoyResolver.GetName()).To(Equal(mock.StaticResolverTypedExtensionConfigName))
			Expect(getResolverConfig(envoyResolver)).To(matchers.MatchProto(&v2.StaticResolver{
				Response: &v2.StaticResolver_SyncResponse{
					SyncResponse: message,
				},
			}))
		})

		Context("async response", func() {
			It("translates async response", func() {
				message := `{"response":"test"}`
				delayMs := uint32(100)
				mockResolver := &v1beta1.MockResolver{
					Response: &v1beta1.MockResolver_AsyncResponse_{
						AsyncResponse: &v1beta1.MockResolver_AsyncResponse{
							Response: getStructValueFromJSON(message),
							Delay: &duration.Duration{
								Nanos: int32(time.Millisecond.Nanoseconds() * 100),
							},
						},
					},
				}
				envoyResolver, err := mock.TranslateMockResolver(mockResolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(envoyResolver.GetName()).To(Equal(mock.StaticResolverTypedExtensionConfigName))
				Expect(getResolverConfig(envoyResolver)).To(matchers.MatchProto(&v2.StaticResolver{
					Response: &v2.StaticResolver_AsyncResponse_{
						AsyncResponse: &v2.StaticResolver_AsyncResponse{
							Response: message,
							DelayMs:  delayMs,
						},
					},
				}))
			})

			It("translates error response", func() {
				message := `some error!`
				mockResolver := &v1beta1.MockResolver{
					Response: &v1beta1.MockResolver_ErrorResponse{
						ErrorResponse: message,
					},
				}
				envoyResolver, err := mock.TranslateMockResolver(mockResolver)
				Expect(err).NotTo(HaveOccurred())
				Expect(envoyResolver.GetName()).To(Equal(mock.StaticResolverTypedExtensionConfigName))
				Expect(getResolverConfig(envoyResolver)).To(matchers.MatchProto(&v2.StaticResolver{
					Response: &v2.StaticResolver_ErrorResponse{
						ErrorResponse: message,
					},
				}))
			})

		})
	})
})
