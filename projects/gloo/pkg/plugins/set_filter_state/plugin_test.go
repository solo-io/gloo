package set_filter_state_test

import (
	http_set_filter_state_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/set_filter_state/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/set_filter_state"
	. "github.com/solo-io/gloo/projects/gloo/pkg/plugins/set_filter_state"
)

type testCase struct {
	input  *set_filter_state.SetFilterState
	output *http_set_filter_state_v3.Config
}

var testCases = []testCase{
	{},
}

var _ = Describe("SetFilterState Plugin", func() {
	Context("Translation", func() {
		for _, testCase := range testCases {
			It("Translates the filter state", func() {
				Expect(TranslateFilter(testCase.input)).To(Equal(testCase.output))
			})
		}

		// var ()

		// BeforeEach(func() {

		// })

		// Context("HttpFilters", func() {
		// 	var listener *v1.HttpListener

		// 	BeforeEach(func() {
		// 		listener = &v1.HttpListener{}
		// 	})

		// 	It("should return no filters if no config is provided", func() {
		// 		filters, err := p.HttpFilters(params, listener)
		// 		Expect(err).NotTo(HaveOccurred())
		// 		Expect(filters).To(BeEmpty())
		// 	})

		// })
	})
})
