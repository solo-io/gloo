package license_test

import (
	"context"
	"os"
	"time"

	"github.com/solo-io/licensing/pkg/defaults"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-projects/pkg/license"
)

var _ = Describe("License", func() {
	It("should ignore spaces", func() {

		km, err := defaults.GetKeyManager()
		Expect(err).NotTo(HaveOccurred())

		key, err := km.KeyGenerator().GenerateKey(context.Background(), time.Now().Add(time.Hour))
		Expect(err).NotTo(HaveOccurred())

		// simulate echo
		os.Setenv("GLOO_LICENSE_KEY", key+"\n")

		Expect(LicenseStatus(context.Background())).NotTo(HaveOccurred())

	})
})
