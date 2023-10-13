package make_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/test/testutils"
)

var _ = Describe("FIPS", func() {

	Context("validate-%-crypto", Ordered, func() {

		const (
			standardBinaryLocalPath = "gloo-standard"
			fipsBinaryLocalPath     = "gloo-fips"

			glooContainerBinaryPath = "/usr/local/bin/gloo"
		)

		AfterAll(func() {
			for _, binaryPath := range []string{standardBinaryLocalPath, fipsBinaryLocalPath} {
				_ = os.Remove(binaryPath)
			}
		})

		It("successfully validates standard-crypto binary", func() {
			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred(), "can get working directory")

			err = testutils.CopyImageFileToLocal(StandardGlooImage, glooContainerBinaryPath, standardBinaryLocalPath)
			Expect(err).NotTo(HaveOccurred(), "can copy binary from image to local filesystem")

			target := fmt.Sprintf("BINARY=%s validate-boring-crypto --ignore-errors", filepath.Join(pwd, standardBinaryLocalPath))
			ExpectMakeOutputWithOffset(0, target, ContainSubstring("validate-boring-crypto] Error 1 (ignored)"))

			target = fmt.Sprintf("BINARY=%s validate-standard-crypto --ignore-errors", filepath.Join(pwd, standardBinaryLocalPath))
			ExpectMakeOutputWithOffset(0, target, And(
				ContainSubstring("goversion -crypto"),
				Not(ContainSubstring("Error 1 (ignored)")),
			))
		})

		It("successfully validates boring-crypto binary", func() {
			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred(), "can get working directory")

			err = testutils.CopyImageFileToLocal(FipsGlooImage, glooContainerBinaryPath, fipsBinaryLocalPath)
			Expect(err).NotTo(HaveOccurred(), "can copy binary from image to local filesystem")

			target := fmt.Sprintf("BINARY=%s validate-standard-crypto --ignore-errors", filepath.Join(pwd, fipsBinaryLocalPath))
			ExpectMakeOutputWithOffset(0, target, ContainSubstring("validate-standard-crypto] Error 1 (ignored)"))

			target = fmt.Sprintf("BINARY=%s validate-boring-crypto --ignore-errors", filepath.Join(pwd, fipsBinaryLocalPath))
			ExpectMakeOutputWithOffset(0, target, And(
				ContainSubstring("goversion -crypto"),
				Not(ContainSubstring("Error 1 (ignored)")),
			))
		})

	})

})
