package matchers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/gomega/matchers"
)

var _ = Describe("HaveProtocol", func() {

	It("succeeds if there is no protocol specified", func() {
		httpResponse := &http.Response{
			Proto: matchers.HTTP1Protocol,
		}
		Expect(httpResponse).To(matchers.HaveProtocol(""))
	})

	It("matches the specified protocol", func() {
		httpResponse := &http.Response{
			Proto: matchers.HTTP1Protocol,
		}
		Expect(httpResponse).To(matchers.HaveProtocol(matchers.HTTP1Protocol))
		Expect(httpResponse).To(matchers.HaveHTTP1Protocol())
		Expect(httpResponse).ToNot(matchers.HaveHTTP2Protocol())
	})
})
