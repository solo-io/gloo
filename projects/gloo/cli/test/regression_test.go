package cli_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd"
	"github.com/spf13/cobra"
)

var _ = Describe("Regression", func() {

	var (
		app    *cobra.Command
		output *bytes.Buffer
	)

	BeforeEach(func() {
		app = cmd.App("0.1.0")
		output = &bytes.Buffer{}
		app.SetOutput(output)
	})

	It("can run properly without crashing", func() {
		app.SetArgs([]string{})
		err := app.Execute()
		Expect(err).NotTo(HaveOccurred())
	})

})
