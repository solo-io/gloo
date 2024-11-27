package translator_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

var _ = Describe("Utils", func() {

	DescribeTable(
		"IsIpv4Address",
		func(address string, expectedIPv4, expectedPureIPv4 bool, expectedErr error) {
			isIPv4Address, isPureIPv4Address, err := translator.IsIpv4Address(address)

			if expectedErr != nil {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(isIPv4Address).To(Equal(expectedIPv4))
			Expect(isPureIPv4Address).To(Equal(expectedPureIPv4))
		},
		Entry("invalid ip returns original", "invalid", false, false, errors.Errorf("bindAddress invalid is not a valid IP address")),
		Entry("ipv4 returns true", "0.0.0.0", true, true, nil),
		Entry("ipv6 returns false", "::", false, false, nil),
		Entry("ipv4 mapped in ipv6", "::ffff:0.0.0.0", true, false, nil),
	)

})
