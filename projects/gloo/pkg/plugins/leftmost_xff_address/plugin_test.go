package leftmost_xff_address_test

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/xff_offset"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/leftmost_xff_address"
)

var _ = Describe("Plugin", func() {
	var (
		plugin *Plugin
	)

	Context("with use leftmost xff header", func() {

		It("is true, should include the filter", func() {

			listener := &v1.HttpListener{
				Options: &v1.HttpListenerOptions{
					LeftmostXffAddress: &wrappers.BoolValue{Value: true},
				},
			}

			filters, err := plugin.HttpFilters(plugins.Params{}, listener)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))
			goTypedConfig := filters[0].HttpFilter.GetTypedConfig()
			var filterConfig = SoloXffOffset{}
			err = ptypes.UnmarshalAny(goTypedConfig, &filterConfig)
			Expect(err).NotTo(HaveOccurred())
			// Offset from left side of xff header should be 0, to get the left most address
			Expect(filterConfig.Offset).To(Equal(uint32(0)))
		})

		It("is false, should not include the filter", func() {

			listener := &v1.HttpListener{
				Options: &v1.HttpListenerOptions{
					LeftmostXffAddress: &wrappers.BoolValue{Value: false},
				},
			}

			filters, err := plugin.HttpFilters(plugins.Params{}, listener)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(0))
		})

		It("is not configured, should not include the filter", func() {

			listener := &v1.HttpListener{}

			filters, err := plugin.HttpFilters(plugins.Params{}, listener)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(0))

		})
	})
})
