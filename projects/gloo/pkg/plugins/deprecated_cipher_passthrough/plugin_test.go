package deprecated_cipher_passthrough

import (
	"context"
	"fmt"
	"os"

	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"

	"github.com/solo-io/gloo/pkg/utils/protoutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var _ = Describe("Plugin", func() {

	It("can translate IR to envoy", func() {
		serverNameMap := make(serverNameMap)
		serverNameMap["foo"] = sourceIpCidrMap{noIpRanges: &deprecatedCipherMapping{
			PassthroughCipherSuites:            []uint32{0x003c},
			PassthroughCipherSuitesFilterChain: &envoy_config_listener_v3.FilterChain{Name: "tcp"},
			DefaultFilterChain:                 &envoy_config_listener_v3.FilterChain{Name: "http"},
		}}
		// We have easy to work with intermediate representation, so now we can directly convert it to envoy config.
		m, err := convertIrToEnvoyMatcher(context.Background(), serverNameMap)
		Expect(err).To(BeNil())
		Expect(m.OnNoMatch).To(BeNil())
		Expect(m.GetMatcherTree()).NotTo(BeNil())
	})

	It("should translate server names", func() {
		fcm := []*plugins.ExtendedFilterChain{
			{
				FilterChain: &envoy_config_listener_v3.FilterChain{
					FilterChainMatch: &envoy_config_listener_v3.FilterChainMatch{
						ServerNames: []string{"foo"},
					},
				},
			},
		}
		m, fc, err := ConvertFilterChain(context.Background(), fcm)
		Expect(err).NotTo(HaveOccurred())
		Expect(m).To(BeNil())
		Expect(fc).To(HaveLen(1))
		Expect(fc[0].FilterChainMatch).NotTo(BeNil())
		Expect(fc[0].Name).To(BeEmpty())
	})

	It("should have a valid intermediate representation", func() {
		fcm := []*plugins.ExtendedFilterChain{
			{
				FilterChain: &envoy_config_listener_v3.FilterChain{
					FilterChainMatch: &envoy_config_listener_v3.FilterChainMatch{
						ServerNames: []string{"foo"},
						SourcePrefixRanges: []*envoy_core_v3.CidrRange{
							{
								AddressPrefix: "127.0.0.1",
								PrefixLen:     wrapperspb.UInt32(32),
							},
						},
					},
				},
			},
		}
		serverNames, filterChains, err := filterChainsToMatcherIR(fcm)
		Expect(err).NotTo(HaveOccurred())
		Expect(serverNames).To(HaveLen(1))
		Expect(serverNames["foo"]).To(HaveLen(1))
		expectedSourceIPCidrMapKey := comparableCidrRange{
			AddressPrefix: "127.0.0.1",
			PrefixLen:     32,
		}
		deprecatedCipherMapping := serverNames["foo"][expectedSourceIPCidrMapKey]
		Expect(deprecatedCipherMapping.PassthroughCipherSuites).To(Equal([]uint32(nil)))
		Expect(deprecatedCipherMapping.PassthroughCipherSuitesFilterChain).To(Equal((*envoy_config_listener_v3.FilterChain)(nil)))
		Expect(deprecatedCipherMapping.DefaultFilterChain).To(Equal(filterChains[0]))
	})

	It("should translate server names and source ips", func() {
		fcm := []*plugins.ExtendedFilterChain{
			{
				FilterChain: &envoy_config_listener_v3.FilterChain{
					FilterChainMatch: &envoy_config_listener_v3.FilterChainMatch{
						ServerNames: []string{"foo"},
						SourcePrefixRanges: []*envoy_core_v3.CidrRange{
							{
								AddressPrefix: "127.0.0.1",
								PrefixLen:     wrapperspb.UInt32(32),
							},
						},
					},
				},
			},
		}
		m, fc, err := ConvertFilterChain(context.Background(), fcm)
		Expect(err).NotTo(HaveOccurred())
		Expect(m).To(BeNil())
		Expect(fc).To(HaveLen(1))
		Expect(fc[0].FilterChainMatch).NotTo(BeNil())
		Expect(fc[0].Name).To(BeEmpty())
	})

	It("should translate server and passthrough ciphers with server names", func() {
		fcm := []*plugins.ExtendedFilterChain{
			{
				FilterChain: &envoy_config_listener_v3.FilterChain{
					FilterChainMatch: &envoy_config_listener_v3.FilterChainMatch{
						ServerNames: []string{"foo"},
					},
				},
				PassthroughCipherSuites: []string{"AES128-SHA256"},
			},
		}
		m, fc, err := ConvertFilterChain(context.Background(), fcm)
		Expect(err).NotTo(HaveOccurred())
		Expect(m).NotTo(BeNil())
		Expect(fc).To(HaveLen(1))
		Expect(fc[0].FilterChainMatch).To(BeNil())
		Expect(fc[0].Name).NotTo(BeNil())
	})

	It("should translate server and passthrough ciphers", func() {
		fcm := []*plugins.ExtendedFilterChain{
			{
				FilterChain:             &envoy_config_listener_v3.FilterChain{},
				PassthroughCipherSuites: []string{"AES128-SHA256"},
			},
		}
		m, fc, err := ConvertFilterChain(context.Background(), fcm)
		Expect(err).NotTo(HaveOccurred())
		Expect(m).NotTo(BeNil())
		Expect(fc).To(HaveLen(1))
		Expect(fc[0].FilterChainMatch).To(BeNil())
		Expect(fc[0].Name).NotTo(BeNil())
	})

	It("should translate server and a mix of passthrough ciphers with server names and source ips", func() {
		fcm := []*plugins.ExtendedFilterChain{
			{
				FilterChain: &envoy_config_listener_v3.FilterChain{
					FilterChainMatch: &envoy_config_listener_v3.FilterChainMatch{
						ServerNames: []string{"foo"},
						SourcePrefixRanges: []*envoy_core_v3.CidrRange{
							{
								AddressPrefix: "127.0.0.1",
								PrefixLen:     wrapperspb.UInt32(32),
							},
						},
					},
				},
				PassthroughCipherSuites: []string{"AES128-SHA256", "0x004c"}, // 60 and  76
			},
		}
		m, fc, err := ConvertFilterChain(context.Background(), fcm)
		Expect(err).NotTo(HaveOccurred())
		Expect(m).NotTo(BeNil())
		Expect(fc).To(HaveLen(1))
		Expect(fc[0].FilterChainMatch).To(BeNil())
		Expect(fc[0].Name).NotTo(BeNil())

		Expect(proto2json(m)).To(MatchYAML(testData("servername-srcip-dcp.yaml")))

	})
	It("should translate server and passthrough ciphers with server names and source ips and another filterchain with no ciphers", func() {
		fcm := []*plugins.ExtendedFilterChain{
			{
				FilterChain: &envoy_config_listener_v3.FilterChain{
					FilterChainMatch: &envoy_config_listener_v3.FilterChainMatch{
						ServerNames: []string{"foo"},
						SourcePrefixRanges: []*envoy_core_v3.CidrRange{
							{
								AddressPrefix: "127.0.0.1",
								PrefixLen:     wrapperspb.UInt32(32),
							},
						},
					},
				},
				PassthroughCipherSuites: []string{"AES128-SHA256"},
			}, {
				FilterChain: &envoy_config_listener_v3.FilterChain{
					FilterChainMatch: &envoy_config_listener_v3.FilterChainMatch{
						ServerNames: []string{"foo"},
						SourcePrefixRanges: []*envoy_core_v3.CidrRange{
							{
								AddressPrefix: "127.0.0.1",
								PrefixLen:     wrapperspb.UInt32(32),
							},
						},
					},
				},
			}, {
				FilterChain: &envoy_config_listener_v3.FilterChain{},
			},
		}
		m, fc, err := ConvertFilterChain(context.Background(), fcm)
		Expect(err).NotTo(HaveOccurred())
		Expect(m).NotTo(BeNil())
		Expect(fc).To(HaveLen(3))
		Expect(fc[0].FilterChainMatch).To(BeNil())
		Expect(fc[0].Name).NotTo(BeNil())
		Expect(fc[1].FilterChainMatch).To(BeNil())
		Expect(fc[1].Name).NotTo(BeNil())
		Expect(fc[2].FilterChainMatch).To(BeNil())
		Expect(fc[2].Name).NotTo(BeNil())
		format.MaxLength = 40000
		Expect(proto2json(m)).To(MatchYAML(testData("servername-srcip-dcp-multi.yaml")))

	})

})

func proto2json(m proto.Message) string {
	b, err := protoutils.MarshalBytes(m)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}

func testData(f string) string {
	b, err := os.ReadFile(fmt.Sprintf("testdata/%s", f))
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}
