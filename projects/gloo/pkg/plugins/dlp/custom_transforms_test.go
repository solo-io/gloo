package dlp

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("dlp plugin", func() {
	Context("custom_transforms", func() {
		It("should fail to compile regex with invalid re2 syntax (e.g. lookaheads)", func() {
			_, err := regexp.Compile(`(?!\D)[0-9]{9}(?=\D|$)`)
			Expect(err).NotTo(BeNil())
		})

		It("should compile all regexes as valid re2", func() {
			for k := range transformMap {
				for _, transform := range GetTransformsFromMap(k) {
					for _, regexAction := range transform.RegexActions {
						_, err := regexp.Compile(regexAction.GetRegex())
						Expect(err).To(BeNil())
					}
				}
			}
		})
	})
})
