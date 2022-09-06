package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/selectors"
)

var _ = Describe("Namespaces", func() {

	var (
		allNamespaces = [][]string{nil, {}, {""}}
	)

	Context("all namespaces", func() {
		It("should consider empty list as all namespaces", func() {
			for _, testCase := range allNamespaces {
				Expect(utils.AllNamespaces(testCase)).To(BeTrue())
			}
		})

		It("should consider list with empty namespace as all namespaces", func() {
			Expect(utils.AllNamespaces([]string{""})).To(BeTrue())
		})

		It("should not consider non empty list as all namespaces", func() {
			Expect(utils.AllNamespaces([]string{"test"})).To(BeFalse())
		})
	})
	Context("process namespaces", func() {

		It("should not change all namespaces", func() {
			for _, testCase := range allNamespaces {
				Expect(utils.ProcessWatchNamespaces(testCase, "test")).To(Equal(testCase))
			}
		})
		It("should add write namespace if not there", func() {
			ns := []string{"ns1"}
			Expect(utils.ProcessWatchNamespaces(ns, "test")).To(Equal([]string{"ns1", "test"}))
		})

		It("should not add write namespace if already there", func() {
			ns := []string{"ns1"}
			Expect(utils.ProcessWatchNamespaces(ns, "ns1")).To(Equal([]string{"ns1"}))
		})

	})
	Context("convert expression selectors", func() {
		var (
			values1 []string
			values2 []string
			se1     selectors.Selector_Expression
			se2     selectors.Selector_Expression
		)
		BeforeEach(func() {
			values1 = []string{"dev"}
			values2 = []string{"foo"}
			se1 = selectors.Selector_Expression{
				Key:      "env",
				Operator: selectors.Selector_Expression_Equals,
				Values:   values1,
			}
			se2 = selectors.Selector_Expression{
				Key:      "meta",
				Operator: selectors.Selector_Expression_Equals,
				Values:   values2,
			}
		})
		It("should convert an expression selector", func() {
			value, err := utils.ConvertExpressionSelectorToString([]*selectors.Selector_Expression{&se1})
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("env=dev"))

			se1.Values = append(se1.Values, "production")
			se1.Operator = selectors.Selector_Expression_In

			value, err = utils.ConvertExpressionSelectorToString([]*selectors.Selector_Expression{&se1})
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("env in (dev,production)"))
		})
		It("should convert multi expression selectors", func() {
			value, err := utils.ConvertExpressionSelectorToString([]*selectors.Selector_Expression{&se1, &se2})
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("env=dev,meta=foo"))

			se1.Values = append(se1.Values, "production")
			se1.Operator = selectors.Selector_Expression_In

			se2.Values = append(se2.Values, "foobar")
			se2.Values = append(se2.Values, "foobarbaz")
			se2.Operator = selectors.Selector_Expression_NotIn

			value, err = utils.ConvertExpressionSelectorToString([]*selectors.Selector_Expression{&se1, &se2})
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("env in (dev,production),meta notin (foo,foobar,foobarbaz)"))
		})
	})

})
