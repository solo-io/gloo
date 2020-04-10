package proxylatency_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gogojsonpb "github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/jsonpb"

	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
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
		presentConfig := getConfig(filters[0].HttpFilter)
		Expect(*presentConfig).To(BeEquivalentTo(pl))

	})

})

func getConfig(f *envoyhttp.HttpFilter) *proxylatency.ProxyLatency {
	cfg := f.GetConfig()
	rcfg := new(proxylatency.ProxyLatency)

	buf := &bytes.Buffer{}
	err := (&jsonpb.Marshaler{OrigName: true}).Marshal(buf, cfg)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	err = gogojsonpb.Unmarshal(buf, rcfg)

	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return rcfg
}
