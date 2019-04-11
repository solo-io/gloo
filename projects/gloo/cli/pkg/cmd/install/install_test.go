package install_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/testutils"
)

var _ = Describe("Install", func() {

	const licenseKey = "--license-key=eyJleHAiOjM4Nzk1MTY3ODYsImlhdCI6MTU1NDk0MDM0OCwiayI6IkJ3ZXZQQSJ9.tbJ9I9AUltZ-iMmHBertugI2YIg1Z8Q0v6anRjc66Jo"

	BeforeEach(func() {

	})

	It("shouldn't get errors for gateway dry run", func() {
		err := testutils.GlooctlEE(fmt.Sprintf("install gateway --file %s --dry-run %s", file, licenseKey))
		Expect(err).NotTo(HaveOccurred())
	})

	It("shouldn't get errors for knative dry run", func() {
		err := testutils.GlooctlEE(fmt.Sprintf("install knative --file %s --dry-run %s", file, licenseKey))
		Expect(err).NotTo(HaveOccurred())
	})

	It("shouldn't get errors for ingress dry run", func() {
		err := testutils.GlooctlEE(fmt.Sprintf("install ingress --file %s --dry-run %s", file, licenseKey))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should error when not overriding helm chart in dev mode", func() {
		err := testutils.GlooctlEE(fmt.Sprintf("install gateway %s", licenseKey))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("you must provide a GlooE Helm chart URI via the 'file' option when running an unreleased version of glooctl"))
	})

	It("should error when not providing file with valid extension", func() {
		err := testutils.GlooctlEE(fmt.Sprintf("install gateway --file foo %s", licenseKey))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unsupported file extension for Helm chart URI: [foo]. Extension must either be .tgz or .tar.gz"))
	})

	It("should error when not providing valid file", func() {
		err := testutils.GlooctlEE(fmt.Sprintf("install gateway --file foo.tgz %s", licenseKey))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("retrieving gloo helm chart archive: opening file"))
	})

	It("should error when not overriding helm chart in dev mode", func() {
		err := testutils.GlooctlEE("install gateway")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("you must provide a valid license key to be able to install GlooE"))
	})

	It("should error when not providing file with valid extension", func() {
		err := testutils.GlooctlEE("install gateway --file foo")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("you must provide a valid license key to be able to install GlooE"))
	})

	It("should error when not providing valid file", func() {
		err := testutils.GlooctlEE("install gateway --file foo.tgz")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("you must provide a valid license key to be able to install GlooE"))
	})
})
