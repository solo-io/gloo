package debug_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"time"

	cliutil "github.com/solo-io/gloo/pkg/cliutil/install"
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
			out, err := cliutil.KubectlOut(nil, "config", "current-context")
			currentContext := string(out)
			Expect(err).NotTo(HaveOccurred(), err.Error()+", "+currentContext)
			fmt.Println("ARIANA BeforeSuite current-context", currentContext)

			Expect(cliutil.Kubectl(nil, "config", "unset", "current-context")).NotTo(HaveOccurred())

			fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "debug TEST: running debug")
			err = testutils.Glooctl("debug")
			fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "debug TEST: AFTER debug")
			Expect(err).NotTo(HaveOccurred())
		}, &timeout)

		fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "debug TEST: AFTER timeout")
		kubeStateBytes, err := os.ReadFile(kubeStateFile(""))
		Expect(err).NotTo(HaveOccurred(), kubeStateFile("")+" file should be present")
		Expect(kubeStateBytes).NotTo(BeEmpty())
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
				out, err := cliutil.KubectlOut(nil, "config", "current-context")
				currentContext := string(out)
				Expect(err).NotTo(HaveOccurred(), err.Error()+", "+currentContext)
				fmt.Println("ARIANA BeforeSuite current-context", currentContext)

				Expect(cliutil.Kubectl(nil, "config", "unset", "current-context")).NotTo(HaveOccurred())

				fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "custom-dir TEST: running debug")
				err = testutils.Glooctl("debug --directory " + customDir)
				fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "custom-dir TEST: AFTER debug")
				Expect(err).NotTo(HaveOccurred())
			}, &timeout)

			fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "custom-dir TEST: AFTER timeout")
			kubeStateBytes, err := os.ReadFile(kubeStateFile(customDir))
			Expect(err).NotTo(HaveOccurred(), kubeStateFile(customDir)+" file should be present")
			Expect(kubeStateBytes).NotTo(BeEmpty())

			// default dir should not exist
			_, err = os.ReadDir(defaultOutDir)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(os.ErrNotExist))
		})
	})

	It("should error and abort if the user does not consent", func() {
		testutil.ExpectInteractive(func(c *testutil.Console) {
			c.ExpectString("This command will overwrite the \"" + defaultOutDir + "\" directory, if present. Are you sure you want to proceed? [y/N]: ")
			c.SendLine("N")
			c.ExpectEOF()
		}, func() {
			fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "N-debug TEST: running debug")
			err := testutils.Glooctl("debug")
			fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "N-debug TEST: AFTER debug")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Aborting: cannot proceed without overwriting \"" + defaultOutDir + "\" directory"))
		}, nil)

		fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "N-debug TEST: AFTER timeout")
		_, err := os.ReadDir(defaultOutDir)
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(os.ErrNotExist))
	})

})
