package debug_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"time"

	"github.com/solo-io/gloo/pkg/cliutil/testutil"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/testutils"
)

const defaultOutDir = "debug"

var timeout = 20 * time.Second

var kubeStateFile = func(outDir string) string {
	if outDir == "" {
		outDir = defaultOutDir
	}
	return outDir + "/kube-state.log"
}

var _ = Describe("Debug", func() {

	BeforeEach(func() {
		helpers.UseMemoryClients()
	})

	AfterEach(func() {
		Expect(os.RemoveAll(defaultOutDir)).NotTo(HaveOccurred())
	})

	It("should support the top level debug command and should populate the kube-state.log file", func() {

		testutil.ExpectInteractive(func(c *testutil.Console) {
			c.ExpectString("This command will overwrite the \"" + defaultOutDir + "\" directory, if present. Are you sure you want to proceed? [y/N]: ")
			c.SendLine("y")
			c.ExpectEOF()
		}, func() {
			fmt.Println("ARIANA running debug", time.Now())
			err := testutils.Glooctl("debug")
			Expect(err).NotTo(HaveOccurred())

			fmt.Println("ARIANA getting kube state file", time.Now())
			kubeStateBytes, err := os.ReadFile(kubeStateFile(""))
			Expect(err).NotTo(HaveOccurred(), kubeStateFile("")+" file should be present")
			Expect(kubeStateBytes).NotTo(BeEmpty())
		}, &timeout)
	})

	When("a directory is specified", func() {

		const customDir = "custom-dir"

		AfterEach(func() {
			Expect(os.RemoveAll(customDir)).NotTo(HaveOccurred())
		})

		It("should populate specified directory instead", func() {
			testutil.ExpectInteractive(func(c *testutil.Console) {
				c.ExpectString("This command will overwrite the \"" + customDir + "\" directory, if present. Are you sure you want to proceed? [y/N]: ")
				c.SendLine("y")
				c.ExpectEOF()
			}, func() {
				fmt.Println("ARIANA running debug", time.Now())
				err := testutils.Glooctl("debug --directory " + customDir)
				Expect(err).NotTo(HaveOccurred())

				fmt.Println("ARIANA getting kube state file", time.Now())
				kubeStateBytes, err := os.ReadFile(kubeStateFile(customDir))
				Expect(err).NotTo(HaveOccurred(), kubeStateFile(customDir)+" file should be present")
				Expect(kubeStateBytes).NotTo(BeEmpty())

				// default dir should not exist
				_, err = os.ReadDir(defaultOutDir)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(os.ErrNotExist))
			}, &timeout)
		})
	})

	It("should error and abort if the user does not consent", func() {
		testutil.ExpectInteractive(func(c *testutil.Console) {
			c.ExpectString("This command will overwrite the \"" + defaultOutDir + "\" directory, if present. Are you sure you want to proceed? [y/N]: ")
			c.SendLine("N")
			c.ExpectEOF()
		}, func() {
			err := testutils.Glooctl("debug")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Aborting: cannot proceed without overwriting \"" + defaultOutDir + "\" directory"))

			_, err = os.ReadDir(defaultOutDir)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(os.ErrNotExist))
		}, nil)
	})

})
