package js_test

import (
	_ "embed"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/js"
)

//go:embed data/return_value.js
var returnValueSource string

//go:embed data/run_forever.js
var runForeverSource string

//go:embed data/errors.js
var errorsSource string

var _ = Describe("V8Runner", func() {
	It("should be able to run a program, and return a value", func() {
		input := "Input"
		runner, err := js.NewV8RunnerInputOutput("data/return_value.js", returnValueSource)
		Expect(err).ToNot(HaveOccurred())
		returnValue, err := runner.Run(input)
		Expect(err).ToNot(HaveOccurred())
		Expect(returnValue).To(Equal(fmt.Sprintf("%s Hello World!", input)))
	})

	It("should be able to timeout", func() {
		input := ""
		runner, err := js.NewV8RunnerInputOutput("data/run_forever.js", runForeverSource)
		Expect(err).ToNot(HaveOccurred())
		runner.SetTimeout(time.Millisecond * 100)
		_, err = runner.Run(input)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("timed out"))
		Expect(err.Error()).To(ContainSubstring("100ms"))
	})

	It("should be able to return the error from a failed program", func() {
		input := ""
		runner, err := js.NewV8RunnerInputOutput("data/errors.js", errorsSource)
		Expect(err).ToNot(HaveOccurred())
		_, err = runner.Run(input)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ERROR"))
	})
})
