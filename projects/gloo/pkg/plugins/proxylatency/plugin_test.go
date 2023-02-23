package proxylatency_test

import (
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/golang/protobuf/ptypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/proxylatency"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/proxylatency"
)

var _ = Describe("Plugin", func() {

	It("should not add filter when not needed", func() {
		p := NewPlugin()
		listener := new(v1.HttpListener)
		filters, err := p.HttpFilters(plugins.Params{}, listener)
		Expect(err).NotTo(HaveOccurred())
		Expect(filters).To(BeEmpty())
	})
	It("should add filter when needed", func() {
		p := NewPlugin()
		pl := proxylatency.ProxyLatency{
			Request: proxylatency.ProxyLatency_FIRST_INCOMING_FIRST_OUTGOING,
		}
		listener := &v1.HttpListener{
			Options: &v1.HttpListenerOptions{
				ProxyLatency: &pl,
			},
		}
		filters, err := p.HttpFilters(plugins.Params{}, listener)
		Expect(err).NotTo(HaveOccurred())
		Expect(filters).To(HaveLen(1))
		presentConfig := getTypedConfig(filters[0].HttpFilter)
		Expect((*presentConfig).Equal(pl)).To(BeTrue())

	})

})

func getTypedConfig(f *envoyhttp.HttpFilter) *proxylatency.ProxyLatency {
	goTypedConfig := f.GetTypedConfig()
	rcfg := new(proxylatency.ProxyLatency)
	err := ptypes.UnmarshalAny(goTypedConfig, rcfg)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return rcfg
}
