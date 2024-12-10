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

		kubeStateBytes, err := os.ReadFile(kubeStateFile)
		Expect(err).NotTo(HaveOccurred(), kubeStateFile+" file should be present")
		Expect(kubeStateBytes).NotTo(BeEmpty())
	})

})
