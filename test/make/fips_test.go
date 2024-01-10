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

		// BeforeAll(func() {
		// 	for _, image := range []string{StandardSdsImage} {
		// 		_, err := docker.PullIfNotPresent(context.Background(), image, 3)
		// 		Expect(err).NotTo(HaveOccurred(), "can pull image locally")
		// 	}
		// })

		AfterAll(func() {
			for _, binaryPath := range []string{sdsStandardBinaryLocalPath, sdsFipsBinaryLocalPath} {
				_ = os.Remove(binaryPath)
			}
		})

		FIt("successfully validates standard-crypto binary", func() {
			testutils.ValidateCrypto(StandardSdsImage, sdsContainerBinaryPath, sdsStandardBinaryLocalPath, testutils.ExpectStandardCrypto)
		})

		It("successfully validates boring-crypto binary", func() {
			testutils.ValidateCrypto(FipsSdsImage, sdsContainerBinaryPath, sdsFipsBinaryLocalPath, testutils.ExpectBoringCrypto)
		})

	})

})
