package debug_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/testutils"
)

var _ = Describe("Debug", func() {

	BeforeEach(func() {
		helpers.UseMemoryClients()
	})

	It("should not require a subcommand after debug", func() {
		err := testutils.Glooctl("debug")
		Expect(err).NotTo(HaveOccurred())
	})

})
