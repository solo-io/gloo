package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utilities", func() {

	It("correctly replaces the default namespace", func() {
		newNamespace := "new-namespace"
		Expect(1).To(Equal(1))
		outs := ReplaceNamespaces(Queries_1731087844762345556, newNamespace)
		for _, out := range outs {
			Expect(out).NotTo(Equal(""))
			Expect(out).To(MatchRegexp(newNamespace))
			Expect(out).NotTo(MatchRegexp(`"namespace":"default"`))
		}
	})
})
