package make_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/test/testutils"
)

var _ = Describe("FIPS", func() {

	Context("validate-%-crypto", Ordered, func() {

		const (
			sdsContainerBinaryPath     = "/usr/local/bin/sds"
			sdsStandardBinaryLocalPath = "gloo-standard"
			sdsFipsBinaryLocalPath     = "gloo-fips"
		)

		AfterAll(func() {
			for _, binaryPath := range []string{sdsStandardBinaryLocalPath, sdsFipsBinaryLocalPath} {
				_ = os.Remove(binaryPath)
			}
		})

		It("successfully validates standard-crypto binary", func() {
			testutils.ValidateCrypto(StandardSdsImage, sdsContainerBinaryPath, sdsStandardBinaryLocalPath, testutils.ExpectStandardCrypto)
		})

		It("successfully validates boring-crypto binary", func() {
			testutils.ValidateCrypto(FipsSdsImage, sdsContainerBinaryPath, sdsFipsBinaryLocalPath, testutils.ExpectBoringCrypto)
		})

	})

})
