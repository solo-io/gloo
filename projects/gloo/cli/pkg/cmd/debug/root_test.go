package debug_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"os"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/testutils"
)

const (
	defaultOutDir = "debug"
	kubeStateFile = defaultOutDir + "/kube-state.log"
)

var _ = Describe("Debug", func() {

	BeforeEach(func() {
		helpers.UseMemoryClients()
	})

	AfterEach(func() {
		Expect(os.RemoveAll(defaultOutDir)).NotTo(HaveOccurred())
	})

	It("should support the top level debug command and should populate the kube-state.log file", func() {
		err := testutils.Glooctl("debug")
		Expect(err).NotTo(HaveOccurred())

		checkDebugKube()
	})

	It("should support the debug kube subcommand and should populate the kube-state.log file", func() {
		err := testutils.Glooctl("debug kube")
		Expect(err).NotTo(HaveOccurred())

		checkDebugKube()
	})

	It("should support the debug gloo subcommand", func() {
		err := testutils.Glooctl("debug gloo")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should support the debug envoy subcommand", func() {
		err := testutils.Glooctl("debug envoy")
		Expect(err).NotTo(HaveOccurred())
	})

})

func checkDebugKube() {
	kubeStateBytes, err := os.ReadFile(kubeStateFile)
	Expect(err).NotTo(HaveOccurred(), kubeStateFile+" file should be present")
	Expect(kubeStateBytes).NotTo(BeEmpty())
}
