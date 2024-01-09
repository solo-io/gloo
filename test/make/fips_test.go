package make_test

import (
	"fmt"
	"os"
	"path/filepath"

	//. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/testutils"
)

const (
	ExpectStandardCrypto = false
	ExpectBoringCrypto   = true
)

func ValidateCrypto(imageName string, binaryPath string, binaryLocalPath string, expectFips bool) {
	pwd, err := os.Getwd()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can get working directory")

	err = testutils.CopyImageFileToLocal(imageName, binaryPath, binaryLocalPath)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can copy binary from image to local filesystem")

	expectedCrypto := "standard"
	unexpectedCrypto := "boring"

	if expectFips {
		expectedCrypto = "boring"
		unexpectedCrypto = "standard"
	}

	// Expected crypto
	target := fmt.Sprintf("BINARY=%s validate-%s-crypto --ignore-errors", filepath.Join(pwd, binaryLocalPath), unexpectedCrypto)
	testutils.ExpectMakeOutputWithOffset(1, target, ContainSubstring(fmt.Sprintf("validate-%s-crypto] Error 1 (ignored)", unexpectedCrypto)))

	// Fail this one
	target = fmt.Sprintf("BINARY=%s validate-%s-crypto --ignore-errors", filepath.Join(pwd, binaryLocalPath), expectedCrypto)
	testutils.ExpectMakeOutputWithOffset(1, target, And(
		ContainSubstring("goversion -crypto"),
		Not(ContainSubstring("Error 1 (ignored)")),
	))
}

// var _ = Describe("FIPS", func() {

// FContext("validate-%-crypto", Ordered, func() {

// 	const (
// 		standardBinaryLocalPath = "gloo-standard"
// 		fipsBinaryLocalPath     = "gloo-fips"

// 		glooContainerBinaryPath    = "/usr/local/bin/gloo"
// 		StandardSDSImage           = "quay.io/solo-io/sds:1.16.0-beta1"
// 		sdsContainerBinaryPath     = "/usr/local/bin/sds"
// 		sdsStandardBinaryLocalPath = "gloo-standard"
// 	)

// 	AfterAll(func() {
// 		for _, binaryPath := range []string{standardBinaryLocalPath, fipsBinaryLocalPath} {
// 			_ = os.Remove(binaryPath)
// 		}
// 	})

// 	It("successfully validates standard-crypto binary", func() {
// 		ValidateCrypto(StandardSDSImage, sdsContainerBinaryPath, sdsStandardBinaryLocalPath, ExpectStandardCrypto)
// 	})

// It("successfully validates standard-crypto binary", func() {
// 	pwd, err := os.Getwd()
// 	Expect(err).NotTo(HaveOccurred(), "can get working directory")

// 	err = testutils.CopyImageFileToLocal(StandardSDSImage, sdsContainerBinaryPath, sdsStandardBinaryLocalPath)
// 	Expect(err).NotTo(HaveOccurred(), "can copy binary from image to local filesystem")

// 	target := fmt.Sprintf("BINARY=%s validate-boring-crypto --ignore-errors", filepath.Join(pwd, standardBinaryLocalPath))
// 	testutils.ExpectMakeOutputWithOffset(0, target, ContainSubstring("validate-boring-crypto] Error 1 (ignored)"))

// 	target = fmt.Sprintf("BINARY=%s validate-standard-crypto --ignore-errors", filepath.Join(pwd, standardBinaryLocalPath))
// 	testutils.ExpectMakeOutputWithOffset(0, target, And(
// 		ContainSubstring("goversion -crypto"),
// 		Not(ContainSubstring("Error 1 (ignored)")),
// 	))
// })

// It("successfully validates boring-crypto binary", func() {
// 	pwd, err := os.Getwd()
// 	Expect(err).NotTo(HaveOccurred(), "can get working directory")

// 	err = testutils.CopyImageFileToLocal(SdsGlooImage, glooContainerBinaryPath, fipsBinaryLocalPath)
// 	Expect(err).NotTo(HaveOccurred(), "can copy binary from image to local filesystem")

// 	target := fmt.Sprintf("BINARY=%s validate-standard-crypto --ignore-errors", filepath.Join(pwd, fipsBinaryLocalPath))
// 	testutils.ExpectMakeOutputWithOffset(0, target, ContainSubstring("validate-standard-crypto] Error 1 (ignored)"))

// 	target = fmt.Sprintf("BINARY=%s validate-boring-crypto --ignore-errors", filepath.Join(pwd, fipsBinaryLocalPath))
// 	testutils.ExpectMakeOutputWithOffset(0, target, And(
// 		ContainSubstring("goversion -crypto"),
// 		Not(ContainSubstring("Error 1 (ignored)")),
// 	))
// })

// 	})

// })
