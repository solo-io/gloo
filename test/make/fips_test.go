package make_test

// import (
// 	"os"

// 	. "github.com/onsi/ginkgo/v2"
// 	"github.com/solo-io/gloo/test/testutils"
// )

// var _ = Describe("FIPS", func() {

// 	Context("validate-%-crypto", Ordered, func() {

// 		const (
// 			standardSdsImage           = "quay.io/solo-io/sds:1.16.0-beta1"
// 			fipsSdsImage               = "quay.io/solo-io/sds-fips:1.17.0-beta2-9037"
// 			sdsContainerBinaryPath     = "/usr/local/bin/sds"
// 			sdsStandardBinaryLocalPath = "gloo-standard"
// 			sdsFipsBinaryLocalPath     = "gloo-fips"
// 		)

// 		AfterAll(func() {
// 			for _, binaryPath := range []string{sdsStandardBinaryLocalPath, sdsFipsBinaryLocalPath} {
// 				_ = os.Remove(binaryPath)
// 			}
// 		})

// 		It("successfully validates standard-crypto binary", func() {
// 			testutils.ValidateCrypto(standardSdsImage, sdsContainerBinaryPath, sdsStandardBinaryLocalPath, testutils.ExpectStandardCrypto)
// 		})

// 		It("successfully validates boring-crypto binary", func() {
// 			testutils.ValidateCrypto(fipsSdsImage, sdsContainerBinaryPath, sdsFipsBinaryLocalPath, testutils.ExpectBoringCrypto)
// 		})

// 	})

// })
